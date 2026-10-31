package work

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNotificationCardsPreserveDetailsAndAudience(t *testing.T) {
	sender := LinePush{DashboardURL: "https://app.example/web", WorkerURL: "https://miniapp.line.me/example"}
	cases := []struct{ message, want, link string }{
		{"Chaum: schedule updated. Please open the app to view details.", "ตารางงานอัปเดต", sender.WorkerURL},
		{"Chaum: approval_needed updated. Please open the app to view details.", "รออนุมัติจัดซื้อ", sender.DashboardURL},
		{"Chaum: ส่งมอบอุปกรณ์คำขอ REQ-001\nพลั่ว × 2", "พลั่ว × 2", sender.WorkerURL},
		{"Chaum: equipment request declined. มีอุปกรณ์เพียงพอแล้ว", "มีอุปกรณ์เพียงพอแล้ว", sender.WorkerURL},
		{"สรุปผลการดำเนินงานโครงการ บ้านใหญ่\nhttps://app.example/web?tor_id=123&period_start=2026-10-01", "บ้านใหญ่", "https://app.example/web?tor_id=123&period_start=2026-10-01"},
	}
	for _, tc := range cases {
		card := sender.messageCard(tc.message)
		body, err := json.Marshal(card)
		if err != nil || card["type"] != "flex" || !strings.Contains(string(body), tc.want) {
			t.Fatalf("card lost business details: %s", body)
		}
		bubble := card["contents"].(map[string]any)
		footer := bubble["footer"].(map[string]any)["contents"].([]map[string]any)
		if footer[0]["action"].(map[string]any)["uri"] != tc.link {
			t.Fatal("wrong audience/report link")
		}
	}
}

func TestNotificationCardDoesNotExposeUnsafeURLs(t *testing.T) {
	sender := LinePush{DashboardURL: "http://localhost:3000/web", LogoURL: "https://secret:password@example.com/private"}
	card := sender.messageCard(strings.Repeat("ข้อความ", 200))
	bubble := card["contents"].(map[string]any)
	if bubble["footer"] != nil || bubble["hero"] != nil {
		t.Fatal("unsafe asset or action URL")
	}
	if len([]rune(card["altText"].(string))) > 400 {
		t.Fatal("alt text too long")
	}
}
