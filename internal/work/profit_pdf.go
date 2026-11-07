package work

import (
	"bytes"
	"compress/zlib"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

//go:embed assets/Sarabun-Regular.ttf
var profitThaiFont []byte

type profitPDFRow struct {
	ID       string `json:"tor_id"`
	Name     string `json:"project_name"`
	Start    string `json:"period_start"`
	End      string `json:"period_end"`
	Revenue  string `json:"revenue"`
	Labor    string `json:"labor"`
	Material string `json:"material"`
	Profit   string `json:"profit"`
}
type thaiPDF struct {
	cids   map[string]uint32
	gids   []uint16
	face   *font.Face
	shaper shaping.HarfbuzzShaper
	glyphs map[uint32]string
	pages  []*strings.Builder
	page   *strings.Builder
}

func pdfUTF16(s string) string {
	var b bytes.Buffer
	for _, r := range utf16.Encode([]rune(s)) {
		b.WriteByte(byte(r >> 8))
		b.WriteByte(byte(r))
	}
	return strings.ToUpper(hex.EncodeToString(b.Bytes()))
}
func (p *thaiPDF) shape(s string, size float64) shaping.Output {
	text := []rune(s)
	return p.shaper.Shape(shaping.Input{Text: text, RunEnd: len(text), Direction: di.DirectionLTR, Face: p.face, Size: fixed.Int26_6(size * 64), Script: language.Thai, Language: language.NewLanguage("th")})
}
func (p *thaiPDF) width(s string, size float64) float64 {
	out := p.shape(s, size)
	var w float64
	for _, g := range out.Glyphs {
		w += float64(g.Advance) / 64
	}
	return w
}

// Shape Thai marks with OpenType positioning; ActualText preserves selectable
// original Unicode even when several glyphs form a single Thai cluster.
func (p *thaiPDF) text(s string, x, y, size float64, color string) {
	out := p.shape(s, size)
	runes := []rune(s)
	fmt.Fprintf(p.page, "/Span << /ActualText <FEFF%s> >> BDC\nBT /F1 %.2f Tf %s rg\n", pdfUTF16(s), size, color)
	seen := map[int]bool{}
	for _, g := range out.Glyphs {
		index := g.TextIndex()
		value := ""
		if index >= 0 && index < len(runes) && !seen[index] {
			value = pdfUTF16(string(runes[index:min(index+g.RunesCount(), len(runes))]))
			seen[index] = true
		}
		key := fmt.Sprintf("%d:%s", g.GlyphID, value)
		cid, ok := p.cids[key]
		if !ok {
			cid = uint32(len(p.gids) + 1)
			p.cids[key] = cid
			p.gids = append(p.gids, uint16(g.GlyphID))
			p.glyphs[cid] = value
		}
		fmt.Fprintf(p.page, "1 0 0 1 %.3f %.3f Tm <%04X> Tj\n", x+float64(g.XOffset)/64, 842-y+float64(g.YOffset)/64, cid)
		x += float64(g.Advance) / 64
	}
	p.page.WriteString("ET\nEMC\n")
}
func (p *thaiPDF) right(s string, right, y, size float64, color string) {
	p.text(s, right-p.width(s, size), y, size, color)
}
func (p *thaiPDF) rect(x, y, w, h float64, color string) {
	fmt.Fprintf(p.page, "%s rg %.2f %.2f %.2f %.2f re f\n", color, x, 842-y-h, w, h)
}
func (p *thaiPDF) wrap(s string, width, size float64) []string {
	var lines []string
	line := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, line)
			line = ""
			continue
		}
		next := line + string(r)
		if line != "" && p.width(next, size) > width && !unicode.Is(unicode.Mn, r) {
			if at := strings.LastIndex(line, " "); at > 0 {
				lines = append(lines, strings.TrimSpace(line[:at]))
				line = strings.TrimSpace(line[at+1:]) + string(r)
			} else {
				lines = append(lines, strings.TrimSpace(line))
				line = string(r)
			}
		} else {
			line = next
		}
	}
	if line != "" {
		lines = append(lines, strings.TrimSpace(line))
	}
	if len(lines) == 0 {
		return []string{"—"}
	}
	return lines
}
func thaiReportDate(s string) string {
	d, e := time.Parse("2006-01-02", s)
	if e != nil {
		return s
	}
	months := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	return fmt.Sprintf("%d %s %d", d.Day(), months[int(d.Month())-1], d.Year()+543)
}
func signedReportMoney(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	negative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	v, e := Money(s)
	if negative {
		v = -v
	}
	return v, e
}
func reportMoney(v int64) string {
	s := Decimal(v)
	negative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	parts := strings.Split(s, ".")
	whole := parts[0]
	for i := len(whole) - 3; i > 0; i -= 3 {
		whole = whole[:i] + "," + whole[i:]
	}
	if negative {
		whole = "-" + whole
	}
	return whole + "." + parts[1]
}

const pdfGreen = "0.10 0.34 0.26"
const pdfInk = "0.13 0.20 0.18"
const pdfMuted = "0.40 0.47 0.44"
const pdfPale = "0.94 0.97 0.95"
const pdfRed = "0.69 0.18 0.18"

func profitPDF(data json.RawMessage) ([]byte, error) {
	var rows []profitPDFRow
	if e := json.Unmarshal(data, &rows); e != nil {
		return nil, e
	}
	face, e := font.ParseTTF(bytes.NewReader(profitThaiFont))
	if e != nil {
		return nil, e
	}
	p := &thaiPDF{face: face, glyphs: map[uint32]string{}, cids: map[string]uint32{}}
	var totals [4]int64
	amounts := make([][4]int64, len(rows))
	for i, r := range rows {
		for j, value := range []string{r.Revenue, r.Labor, r.Material, r.Profit} {
			n, err := signedReportMoney(value)
			if err != nil {
				return nil, err
			}
			amounts[i][j] = n
			totals[j] += n
		}
	}
	// Use the same active-first order as the website, retaining every project.
	type entry struct {
		row    profitPDFRow
		amount [4]int64
	}
	entries := make([]entry, len(rows))
	for i := range rows {
		entries[i] = entry{rows[i], amounts[i]}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		active := func(a [4]int64) bool { return a[0] != 0 || a[1] != 0 || a[2] != 0 }
		return active(entries[i].amount) && !active(entries[j].amount)
	})
	period := "ตามช่วงวันที่ที่เลือก"
	if len(rows) > 0 {
		period = thaiReportDate(rows[0].Start) + " – " + thaiReportDate(rows[0].End)
	}
	headers := []string{"โครงการ", "รายรับจริง", "ค่าแรงที่จ่าย", "จัดซื้อจริง", "กำไรสุทธิ"}
	ends := []float64{241, 319, 397, 475, 553}
	newPage := func(first bool) float64 {
		p.page = &strings.Builder{}
		p.pages = append(p.pages, p.page)
		p.rect(42, 35, 34, 4, pdfGreen)
		p.text("ชะอุ่ม", 42, 73, 25, pdfGreen)
		p.right("CHAUM RESOURCE MANAGEMENT", 553, 55, 10, pdfGreen)
		p.right("รายงานการเงินโครงการ", 553, 73, 10, pdfMuted)
		p.rect(42, 91, 511, 1, "0.83 0.88 0.85")
		p.text("สรุปต้นทุนและกำไร", 42, 126, 24, pdfInk)
		p.text("ช่วงวันที่ "+period, 42, 151, 12, pdfMuted)
		y := 178.0
		if first {
			labels := []string{"รายรับจริง", "ต้นทุนค่าแรง", "ต้นทุนจัดซื้อ", "กำไรสุทธิ"}
			for i, label := range labels {
				x := 42 + float64(i)*130
				p.rect(x, y, 121, 79, pdfPale)
				p.text(label, x+12, y+24, 11, pdfMuted)
				color := pdfGreen
				if totals[i] < 0 {
					color = pdfRed
				}
				value := reportMoney(totals[i])
				size := 19.0
				for p.width(value, size) > 99 && size > 10 {
					size--
				}
				p.text(value, x+12, y+53, size, color)
				p.text("บาท", x+12, y+69, 9, pdfMuted)
			}
			y = 294
			p.text(fmt.Sprintf("รายละเอียด %d โครงการ", len(rows)), 42, 282, 12, pdfInk)
		}
		p.rect(42, y, 511, 31, pdfGreen)
		p.text(headers[0], 52, y+20, 11, "1 1 1")
		for j := 1; j < 5; j++ {
			p.right(headers[j], ends[j]-9, y+20, 10, "1 1 1")
		}
		return y + 31
	}
	y := newPage(true)
	if len(entries) == 0 {
		p.text("ไม่พบข้อมูลโครงการในช่วงวันที่เลือก", 52, y+31, 12, pdfMuted)
		y += 66
	}
	for i, entry := range entries {
		name := entry.row.Name
		if name == "" {
			name = entry.row.ID
		}
		lines := p.wrap(name, 179, 11)
		// Limit each page fragment so arbitrarily long project names cannot overflow.
		for len(lines) > 0 {
			capacity := int((709 - y - 14) / 15)
			if capacity < 1 {
				y = newPage(false)
				capacity = int((709 - y - 14) / 15)
			}
			take := min(capacity, len(lines))
			height := float64(take)*15 + 14
			if i%2 == 0 {
				p.rect(42, y, 511, height, "0.97 0.98 0.97")
			}
			for k, line := range lines[:take] {
				p.text(line, 52, y+20+float64(k)*15, 11, pdfInk)
			}
			// A continued name repeats the amounts; it remains one logical row.
			for j, n := range entry.amount {
				color := pdfInk
				if j == 3 {
					color = pdfGreen
					if n < 0 {
						color = pdfRed
					}
				}
				v := reportMoney(n)
				size := 11.0
				for p.width(v, size) > 60 && size > 7 {
					size -= 0.5
				}
				p.right(v, ends[j+1]-9, y+20, size, color)
			}
			p.rect(42, y+height-0.6, 511, 0.6, "0.87 0.91 0.88")
			y += height
			lines = lines[take:]
		}
	}
	if y > 640 {
		y = newPage(false)
	}
	p.rect(42, y, 511, 34, pdfPale)
	p.text("รวมทั้งหมด", 52, y+22, 12, pdfGreen)
	for j, n := range totals {
		color := pdfGreen
		if n < 0 {
			color = pdfRed
		}
		v := reportMoney(n)
		size := 12.0
		for p.width(v, size) > 60 && size > 7 {
			size -= 0.5
		}
		p.right(v, ends[j+1]-9, y+22, size, color)
	}
	y += 59
	p.text("หลักเกณฑ์ของรายงาน", 42, y, 11, pdfGreen)
	y += 19
	notes := []string{"รายรับ: ยอดที่ยืนยันรับเงินจริงตามเดือนวางบิล · ค่าแรง: เฉพาะสลิปที่ยืนยันจ่ายแล้ว", "ค่าแรงแบ่งตามวันทำงานของแต่ละโครงการ จึงเป็นการจัดสรรต้นทุนโดยประมาณ", "จัดซื้อ: ค่าใช้จ่ายจริงตามวันที่ซื้อ ไม่รวมยอดโอนเงินสำรองซ้ำเป็นต้นทุน", "ยอด 0 คือไม่มีรายการเข้าเงื่อนไขในช่วงที่เลือก กำไรติดลบได้เมื่อมีต้นทุนมากกว่ารายรับ"}
	for _, note := range notes {
		for _, line := range p.wrap(note, 511, 10) {
			p.text(line, 42, y, 10, pdfMuted)
			y += 16
		}
	}
	for i, page := range p.pages {
		p.page = page
		p.rect(42, 791, 511, 0.6, "0.83 0.88 0.85")
		p.text("ชะอุ่ม · รายงานสรุปต้นทุนและกำไร (หน่วย: บาท)", 42, 811, 9, pdfMuted)
		p.right(fmt.Sprintf("หน้า %d / %d", i+1, len(p.pages)), 553, 811, 9, pdfMuted)
	}
	return p.bytes(), nil
}
func pdfStream(data []byte, extra string) []byte {
	var b bytes.Buffer
	z := zlib.NewWriter(&b)
	_, _ = z.Write(data)
	_ = z.Close()
	return append(append([]byte(fmt.Sprintf("<< /Length %d /Filter /FlateDecode %s >>\nstream\n", b.Len(), extra)), b.Bytes()...), []byte("\nendstream")...)
}
func (p *thaiPDF) bytes() []byte {
	objects := make([][]byte, 8)
	objects[0] = []byte("<< /Type /Catalog /Pages 2 0 R /Lang (th-TH) >>")
	objects[2] = pdfStream(profitThaiFont, fmt.Sprintf("/Length1 %d", len(profitThaiFont)))
	objects[3] = []byte("<< /Type /FontDescriptor /FontName /Sarabun /Flags 4 /FontBBox [-1000 -500 2000 1500] /ItalicAngle 0 /Ascent 1000 /Descent -300 /CapHeight 700 /StemV 80 /FontFile2 3 0 R >>")
	objects[4] = []byte("<< /Type /Font /Subtype /CIDFontType2 /BaseFont /Sarabun /CIDSystemInfo << /Registry (Adobe) /Ordering (Identity) /Supplement 0 >> /FontDescriptor 4 0 R /DW 1000 /CIDToGIDMap 8 0 R >>")
	var cmap strings.Builder
	cmap.WriteString("/CIDInit /ProcSet findresource begin 12 dict begin begincmap /CIDSystemInfo << /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def /CMapName /SarabunUnicode def /CMapType 2 def 1 begincodespacerange <0000> <FFFF> endcodespacerange\n")
	ids := make([]int, 0, len(p.glyphs))
	for id := range p.glyphs {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for start := 0; start < len(ids); start += 100 {
		end := min(start+100, len(ids))
		fmt.Fprintf(&cmap, "%d beginbfchar\n", end-start)
		for _, id := range ids[start:end] {
			fmt.Fprintf(&cmap, "<%04X> <%s>\n", id, p.glyphs[uint32(id)])
		}
		cmap.WriteString("endbfchar\n")
	}
	cmap.WriteString("endcmap CMapName currentdict /CMap defineresource pop end end")
	objects[5] = pdfStream([]byte(cmap.String()), "")
	objects[6] = []byte("<< /Type /Font /Subtype /Type0 /BaseFont /Sarabun /Encoding /Identity-H /DescendantFonts [5 0 R] /ToUnicode 6 0 R >>")
	mapping := []byte{0, 0}
	for _, gid := range p.gids {
		mapping = append(mapping, byte(gid>>8), byte(gid))
	}
	objects[7] = pdfStream(mapping, "")
	var kids []string
	for _, page := range p.pages {
		id := len(objects) + 1
		kids = append(kids, fmt.Sprintf("%d 0 R", id))
		objects = append(objects, []byte(fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 7 0 R >> >> /Contents %d 0 R >>", id+1)), pdfStream([]byte(page.String()), ""))
	}
	objects[1] = []byte(fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>", len(p.pages), strings.Join(kids, " ")))
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")
	offsets := []int{0}
	for i, o := range objects {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n", i+1)
		out.Write(o)
		out.WriteString("\nendobj\n")
	}
	xref := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(offsets))
	for _, n := range offsets[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", n)
	}
	out.WriteString("trailer\n<< /Size " + strconv.Itoa(len(offsets)) + " /Root 1 0 R >>\nstartxref\n" + strconv.Itoa(xref) + "\n%%EOF\n")
	return out.Bytes()
}
