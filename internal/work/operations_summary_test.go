package work

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

type operationsNotifier struct {
	mu       sync.Mutex
	messages []string
	fail     bool
}

func (n *operationsNotifier) Send(_ context.Context, _, msg string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.messages = append(n.messages, msg)
	if n.fail {
		return conflict("fixture failure")
	}
	return nil
}
func (n *operationsNotifier) read() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]string{}, n.messages...)
}

func TestIntegrationOperationsSummary(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	// Pin fixture creation dates; selected-period tests must not depend on wall time.
	if _, e := f.pool.Exec(ctx, `UPDATE equipment_requisition SET created_at='2026-10-04 08:00+07' WHERE assignment_id=$1`, f.assignment); e != nil {
		t.Fatal(e)
	}
	check := func(start, end string, wantData, wantContinue bool, wantAbsent int) {
		t.Helper()
		data, e := s.OperationsSummary(ctx, f.assistant, f.assignment, start, end)
		if e != nil {
			t.Fatal(e)
		}
		var result struct {
			Data         bool   `json:"has_data"`
			Continue     bool   `json:"can_continue"`
			Notice       string `json:"summary_notice"`
			Continuation string `json:"continuation_notice"`
			Attendance   struct {
				Absent int `json:"absent"`
			} `json:"attendance"`
			Persisted bool `json:"persisted"`
		}
		if e = json.Unmarshal(data, &result); e != nil || result.Data != wantData || result.Continue != wantContinue || result.Attendance.Absent != wantAbsent || result.Persisted {
			t.Fatal("invalid operational result", string(data), e)
		}
		if !wantData && result.Notice != "ยังไม่มีข้อมูลสรุปสำหรับรอบนี้" {
			t.Fatal("missing empty-period message")
		}
		if !wantContinue && result.Continuation == "" {
			t.Fatal("missing ended-contract message")
		}
		for _, key := range []string{"base_wage", "net_wage", "daily_wage", "actual_price", "total_deduction", "profit", "bank_account_no", "contract_value", "total_amount"} {
			if strings.Contains(string(data), `"`+key+`"`) {
				t.Fatal("financial field leaked", key)
			}
		}
	}
	t.Run("empty selected period", func(t *testing.T) { check("2026-09-01", "2026-09-15", false, true, 0) })
	ids, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `INSERT INTO attendance(schedule_id,worker_id,work_date,status) VALUES($1,$2,'2026-10-01','absent')`, ids[0], f.workerID); e != nil {
		t.Fatal(e)
	}
	t.Run("current counts", func(t *testing.T) { check("2026-10-01", "2026-10-15", true, true, 1) })
	t.Run("procurement aggregates items without duplicating request count", func(t *testing.T) {
		for _, name := range []string{"summary-spade-" + f.assignment, "summary-rake-" + f.assignment} {
			equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: name})
			if e != nil {
				t.Fatal(e)
			}
			if _, e = f.pool.Exec(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty,existing_qty,actual_qty,actual_price) SELECT requisition_id,$2,5,1,2,999 FROM equipment_requisition WHERE assignment_id=$1`, f.assignment, equipment); e != nil {
				t.Fatal(e)
			}
		}
		data, e := s.OperationsSummary(ctx, f.assistant, f.assignment, "2026-10-01", "2026-10-15")
		var out struct {
			Procurement []struct {
				Count     int `json:"request_count"`
				Required  int `json:"required_qty"`
				Acquired  int `json:"acquired_qty"`
				Remaining int `json:"remaining_qty"`
			} `json:"procurement"`
		}
		if e != nil {
			t.Fatal(e)
		}
		if e = json.Unmarshal(data, &out); e != nil || len(out.Procurement) != 1 {
			t.Fatal("invalid procurement", e)
		}
		v := out.Procurement[0]
		if v.Count != 1 || v.Required != 10 || v.Acquired != 4 || v.Remaining != 4 {
			t.Fatal("incorrect progress quantities", v)
		}
		check("2026-10-01", "2026-10-15", true, true, 1)
	})
	if _, e = f.pool.Exec(ctx, `UPDATE attendance SET status='on_time' WHERE schedule_id=$1`, ids[0]); e != nil {
		t.Fatal(e)
	}
	t.Run("retrospective change updates report", func(t *testing.T) { check("2026-10-01", "2026-10-15", true, true, 0) })
	t.Run("worker and invalid range denied", func(t *testing.T) {
		if _, e := s.OperationsSummary(ctx, f.worker, f.assignment, "2026-10-01", "2026-10-15"); e == nil {
			t.Fatal("Worker read report")
		}
		if _, e := s.OperationsSummary(ctx, f.assistant, f.assignment, "2026-10-15", "2026-10-01"); e == nil {
			t.Fatal("unordered range")
		}
	})
	if _, e = f.pool.Exec(ctx, `UPDATE contract_tor SET status='complete' WHERE tor_id=$1`, f.tor); e != nil {
		t.Fatal(e)
	}
	t.Run("ended report visible scheduling rejected", func(t *testing.T) {
		check("2026-10-01", "2026-10-15", true, false, 0)
		if _, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.otherID}, "2026-10-02", "08:00"}); e == nil {
			t.Fatal("ended contract scheduled")
		}
	})
}

func TestIntegrationManualOperationsNotification(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	// Pin fixture creation dates; selected-period tests must not depend on wall time.
	if _, e := f.pool.Exec(ctx, `UPDATE equipment_requisition SET created_at='2026-10-04 08:00+07' WHERE assignment_id=$1`, f.assignment); e != nil {
		t.Fatal(e)
	}
	n := &operationsNotifier{}
	s.Notify = n
	s.AssistantDashboardURL = "https://assistant.example.test/dashboard"
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	in := OperationsNotificationInput{"2026-10-01", "2026-10-15", true}
	t.Run("empty period cannot send", func(t *testing.T) {
		if _, e := s.NotifyOperationsSummary(ctx, f.supervisor, f.tor, OperationsNotificationInput{"2026-09-01", "2026-09-15", true}); e == nil {
			t.Fatal("empty summary sent")
		}
	})
	// Isolated fixture: an absent shift must be included in processed payroll.
	var schedule string
	if e := f.pool.QueryRow(ctx, `INSERT INTO work_schedule(assignment_id,worker_id,work_date,shift_start_time) VALUES($1,$2,'2026-10-01','08:00') RETURNING schedule_id::text`, f.assignment, f.workerID).Scan(&schedule); e != nil {
		t.Fatal(e)
	}
	if _, e := f.pool.Exec(ctx, `INSERT INTO attendance(schedule_id,worker_id,work_date,status) VALUES($1,$2,'2026-10-01','absent')`, schedule, f.workerID); e != nil {
		t.Fatal(e)
	}
	t.Run("role review payroll and URL guards", func(t *testing.T) {
		if _, e := s.NotifyOperationsSummary(ctx, f.assistant, f.tor, in); e == nil {
			t.Fatal("Assistant sent summary")
		}
		if _, e := s.NotifyOperationsSummary(ctx, f.supervisor, f.tor, OperationsNotificationInput{in.Start, in.End, false}); e == nil {
			t.Fatal("unreviewed summary sent")
		}
		if _, e := s.NotifyOperationsSummary(ctx, f.supervisor, f.tor, in); e == nil {
			t.Fatal("unprocessed payroll accepted")
		}
		s.AssistantDashboardURL = "http://unsafe.example/dashboard"
		if _, e := s.NotifyOperationsSummary(ctx, f.supervisor, f.tor, in); e == nil {
			t.Fatal("unsafe dashboard URL")
		}
		s.AssistantDashboardURL = "https://assistant.example.test/dashboard"
	})
	t.Run("payroll and profit confirmation do not notify Assistant", func(t *testing.T) {
		if _, e := s.Finalize(ctx); e != nil {
			t.Fatal(e)
		}
		if _, e := s.PayrollBatch(ctx, f.supervisor, PayrollBatchInput{in.Start, in.End}); e != nil {
			t.Fatal(e)
		}
		if _, e := s.CloseSummary(ctx, f.tor, in.Start, in.End); e != nil {
			t.Fatal(e)
		}
		// Worker payslip notices remain separate from Assistant operational summaries.
		for _, msg := range n.read() {
			if strings.Contains(msg, "สรุป") || strings.Contains(msg, "continuation") {
				t.Fatal("automatic Assistant summary")
			}
		}
	})
	t.Run("manual message includes scope dates and link without finance", func(t *testing.T) {
		out, e := s.NotifyOperationsSummary(ctx, f.supervisor, f.tor, in)
		if e != nil || out.Status != "accepted" || out.Accepted < 1 || out.Persisted || out.DeliveryVerified {
			t.Fatal(out, e)
		}
		found := false
		for _, msg := range n.read() {
			if strings.Contains(msg, "สรุปผลการดำเนินงาน") {
				found = true
				if !strings.Contains(msg, "period_start=2026-10-01") || !strings.Contains(msg, f.tor) || !strings.Contains(msg, "https://assistant.example.test/dashboard") {
					t.Fatal("missing report link")
				}
				if strings.Contains(msg, "net_wage") || strings.Contains(msg, "profit") {
					t.Fatal("financial notification")
				}
			}
		}
		if !found {
			t.Fatal("manual summary missing")
		}
	})
}
