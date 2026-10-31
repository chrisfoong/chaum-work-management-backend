package work

import (
	"context"
	"strings"
	"testing"
)

func TestIntegrationWorkerSummaryIsPersonal(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	ids, err := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID, f.otherID}, "2026-10-02", "08:00"})
	if err != nil || len(ids) != 2 {
		t.Fatal(err)
	}
	msg, err := s.WorkerProjectSummaryMessage(ctx, f.worker, f.tor, "2026-10-01", "2026-10-15")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "กะที่ได้รับมอบหมาย: 1 กะ") || !strings.Contains(msg, "01 ต.ค. 2569") {
		t.Fatal("personal counts/date incorrect")
	}
	if strings.Contains(msg, "daily_wage") || strings.Contains(msg, "profit") || strings.Contains(msg, f.other.UserID.String()) {
		t.Fatal("private staff data leaked")
	}
	if _, err = s.WorkerProjectSummaryMessage(ctx, f.supervisor, f.tor, "2026-10-01", "2026-10-15"); err == nil {
		t.Fatal("staff principal accepted as Worker")
	}
	if _, err = s.OperationsSummary(ctx, f.worker, f.assignment, "2026-10-01", "2026-10-15"); err == nil {
		t.Fatal("Worker granted staff 9A access")
	}
	card := (&LinePush{DashboardURL: "https://staff.example/web", WorkerURL: "https://miniapp.line.me/worker"}).messageCard(msg)
	bubble := card["contents"].(map[string]any)
	button := bubble["footer"].(map[string]any)["contents"].([]map[string]any)[0]["action"].(map[string]any)
	if button["uri"] != "https://miniapp.line.me/worker" {
		t.Fatal("summary points Worker to staff page")
	}
}
func TestLeaveCardReadableFormatting(t *testing.T) {
	got := cardLeaveDetails("เลขใบลา: LEAVE-e20a20d3-2195-41c0-b10a-b97f67330ab2\nวันที่ลา: 2026-10-30\nเหตุผลการลา: พักธุระ")
	if !strings.Contains(got, "เลขอ้างอิงย่อ: LEAVE-e20a20d3") || !strings.Contains(got, "30 ต.ค. 2569") || !strings.Contains(got, "พักธุระ") {
		t.Fatal(got)
	}
}
