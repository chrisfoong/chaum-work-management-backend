package work

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type ProfitConfirmation struct {
	TorID string `json:"tor_id"`
	Start string `json:"period_start"`
	End   string `json:"period_end"`
}

func payrollList(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		if month := c.Query("period_month"); month != "" {
			data, e := s.PayrollMonth(c.Request.Context(), actor(c), month, limit, offset)
			respond(c, gin.H{"data": data, "limit": limit, "offset": offset}, e)
			return
		}
		data, e := s.Payrolls(c.Request.Context(), actor(c), limit, offset)
		respond(c, gin.H{"data": data, "limit": limit, "offset": offset}, e)
	}
}
func (s *Service) profitRequest(c *gin.Context) (json.RawMessage, error) {
	start, end := c.Query("period_start"), c.Query("period_end")
	if start == "" && end == "" {
		month := c.Query("month")
		if e := Month(month); e != nil {
			return nil, e
		}
		d, _ := time.Parse("2006-01", month)
		start = d.Format("2006-01-02")
		end = d.AddDate(0, 1, -1).Format("2006-01-02")
	}
	return s.ProfitRange(c.Request.Context(), c.Query("tor_id"), start, end)
}
func (s *Service) ExportProfitPDF(c *gin.Context) {
	data, e := s.profitRequest(c)
	if e != nil {
		respond(c, nil, e)
		return
	}
	pdf, e := profitPDF(data)
	if e != nil {
		respond(c, nil, e)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="chaum-profit.pdf"`)
	c.Data(200, "application/pdf", pdf)
}

// A dependency-free financial PDF. UUIDs and decimal amounts are ASCII and
// render with the standard PDF font. The JSON/CSV response retains Thai names.
func profitPDF(data json.RawMessage) ([]byte, error) {
	var records []map[string]any
	if e := json.Unmarshal(data, &records); e != nil {
		return nil, e
	}
	lines := []string{"Chaum profit report", "Basis: paid invoices, paid payroll allocated by workdays, actual expenses.", "Payroll allocation is estimated. This report is not a persisted closing ledger.", ""}
	for _, r := range records {
		value := func(key string) string { v, _ := r[key].(string); return v }
		lines = append(lines, "TOR: "+value("tor_id"), "Period: "+value("period_start")+" to "+value("period_end"), "Revenue THB: "+value("revenue")+"    Paid labor THB: "+value("labor"), "Material THB: "+value("material")+"    Profit THB: "+value("profit"), "")
	}
	if len(records) == 0 {
		lines = append(lines, "No matching contracts.")
	}
	pages := (len(lines) + 43) / 44
	objects := []string{"", ""}
	fontID := 3
	objects = append(objects, `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`)
	var kids []string
	for page := 0; page < pages; page++ {
		pageID := len(objects) + 1
		contentID := pageID + 1
		kids = append(kids, fmt.Sprintf("%d 0 R", pageID))
		var stream strings.Builder
		stream.WriteString("BT /F1 10 Tf 48 790 Td 15 TL\n")
		for _, line := range lines[page*44 : min((page+1)*44, len(lines))] {
			line = strings.NewReplacer("\\", "\\\\", "(", "\\(", ")", "\\)").Replace(line)
			stream.WriteString("(" + line + ") Tj T*\n")
		}
		stream.WriteString("ET\n")
		content := stream.String()
		objects = append(objects, fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 %d 0 R >> >> /Contents %d 0 R >>", fontID, contentID), fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", len(content), content))
	}
	objects[0] = "<< /Type /Catalog /Pages 2 0 R >>"
	objects[1] = fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>", pages, strings.Join(kids, " "))
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n")
	offsets := []int{0}
	for i, object := range objects {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xref := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(offsets))
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", offset)
	}
	out.WriteString("trailer\n<< /Size " + strconv.Itoa(len(offsets)) + " /Root 1 0 R >>\nstartxref\n" + strconv.Itoa(xref) + "\n%%EOF\n")
	return out.Bytes(), nil
}
