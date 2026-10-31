package work

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLINECardsMockTransport(t *testing.T) {
	cases := []struct {
		name, message string
		worker        bool
	}{
		{"schedule", "Chaum: schedule updated. Please open the app to view details.", true},
		{"leave-review", "Chaum: leave_review updated. Please open the app to view details.", true},
		{"leave-request", "Chaum: leave updated. Please open the app to view details.", false},
		{"purchase-approval", "Chaum: approval_needed updated. Please open the app to view details.", false},
		{"purchase-funding", "Chaum: purchase_funding updated. Please open the app to view details.", false},
		{"replacement", "Chaum: replacement_needed updated. Please open the app to view details.", false},
		{"funded", "Chaum: funded updated. Please open the app to view details.", false},
		{"equipment-request", "Chaum: request updated. Please open the app to view details.", false},
		{"equipment-review", "Chaum: request_review updated. Please open the app to view details.", true},
		{"payroll", "Chaum: payroll updated. Please open the app to view details.", true},
		{"delivery", "Chaum: ส่งมอบอุปกรณ์คำขอ REQ-MOCK-001\nพลั่ว × 2\nไม้กวาด × 1", true},
		{"declined", "Chaum: equipment request declined. อุปกรณ์ที่มีอยู่เพียงพอสำหรับพื้นที่แล้ว", true},
		{"rejected", "Chaum: equipment request rejected. กรุณาตรวจสอบจำนวนอุปกรณ์อีกครั้ง", true},
		{"9A-summary", "สรุปผลการดำเนินงานโครงการ บ้านใหญ่ รอบวันที่ 2026-10-01 ถึง 2026-10-15 พร้อมให้ตรวจสอบแล้ว\nhttps://app.example/web?tor_id=mock-tor&period_start=2026-10-01&period_end=2026-10-15", false},
	}
	var captures []map[string]any
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer mock-only" || r.Header.Get("X-Line-Retry-Key") != "mock-"+tc.name {
					t.Error("push headers changed")
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body["to"] != "mock-recipient" {
					t.Error("recipient changed")
				}
				messages := body["messages"].([]any)
				if len(messages) != 1 {
					t.Fatal("expected one bubble")
				}
				card := messages[0].(map[string]any)
				if card["type"] != "flex" || card["altText"] == "" {
					t.Fatal("missing Flex/altText")
				}
				bubble := card["contents"].(map[string]any)
				hero := bubble["hero"].(map[string]any)["contents"].([]any)[0].(map[string]any)
				if hero["url"] != "https://app.example/brand/chaum-logo.png" {
					t.Error("logo URL wrong")
				}
				button := bubble["footer"].(map[string]any)["contents"].([]any)[0].(map[string]any)
				uri := button["action"].(map[string]any)["uri"].(string)
				want := "https://app.example/web"
				if tc.worker {
					want = "https://miniapp.line.me/2011922658-kb0yFlBE"
				}
				if tc.name == "9A-summary" {
					want = strings.Split(tc.message, "\n")[1]
				}
				if uri != want {
					t.Errorf("incorrect action: %s", uri)
				}
				if tc.name == "delivery" || tc.name == "declined" || tc.name == "rejected" {
					detail := bubble["body"].(map[string]any)["contents"].([]any)[2].(map[string]any)["text"].(string)
					if tc.name == "delivery" && (!strings.Contains(detail, "พลั่ว × 2") || !strings.Contains(detail, "ไม้กวาด × 1")) {
						t.Error("equipment quantities lost")
					}
					if tc.name != "delivery" && !strings.Contains(detail, strings.SplitN(tc.message, ". ", 2)[1]) {
						t.Error("reason lost")
					}
				}
				captures = append(captures, map[string]any{"case": tc.name, "transport": "MOCK HTTP — NOT LINE", "request": body})
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			client := server.Client()
			client.Transport = redirectTransport{base: client.Transport, url: server.URL}
			push := LinePush{Token: "mock-only", Client: client, DashboardURL: "https://app.example/web", WorkerURL: "https://miniapp.line.me/2011922658-kb0yFlBE"}
			if err := push.SendKey(context.Background(), "mock-recipient", tc.message, "mock-"+tc.name); err != nil || calls != 1 {
				t.Fatalf("mock send failed: %v calls=%d", err, calls)
			}
		})
	}
	if dir := os.Getenv("CHAUM_MOCK_CARD_ARTIFACTS"); dir != "" {
		data, err := json.MarshalIndent(captures, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, "mock-cards.json"), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
