package work

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// lineRecorder captures asynchronous LINE pushes per recipient line_id.
type lineRecorder struct {
	mu   sync.Mutex
	sent map[string][]string
}

func (r *lineRecorder) Send(_ context.Context, line, msg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sent == nil {
		r.sent = map[string][]string{}
	}
	r.sent[line] = append(r.sent[line], msg)
	return nil
}

// waitFor polls until a message to line containing every fragment arrives.
func (r *lineRecorder) waitFor(t *testing.T, line string, fragments ...string) string {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		r.mu.Lock()
		for _, m := range r.sent[line] {
			ok := true
			for _, f := range fragments {
				ok = ok && strings.Contains(m, f)
			}
			if ok {
				r.mu.Unlock()
				return m
			}
		}
		r.mu.Unlock()
		time.Sleep(20 * time.Millisecond)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	t.Fatalf("no LINE message to %s containing %q; got %q", line, fragments, r.sent[line])
	return ""
}

func (r *lineRecorder) count(line string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.sent[line])
}

// checkIn performs a real GPS+QR check-in; 5W requires an open shift.
func checkIn(t *testing.T, f usecaseFixture, schedule string) {
	t.Helper()
	qr, e := f.service.QR(context.Background(), f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	lat, lon, acc := 13.75, 100.5, 10.0
	if e = f.service.CheckIn(context.Background(), f.worker, CheckIn{schedule, &lat, &lon, &acc, qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
}

func leaveStatus(t *testing.T, f usecaseFixture, id string) string {
	t.Helper()
	var state string
	if e := f.pool.QueryRow(context.Background(), `SELECT status::text FROM leave_request WHERE request_id=$1`, id).Scan(&state); e != nil {
		t.Fatal(e)
	}
	return state
}

func requestState(t *testing.T, f usecaseFixture, id string) (string, string) {
	t.Helper()
	var state, reason string
	if e := f.pool.QueryRow(context.Background(), `SELECT status::text,reason FROM equipment_requisition WHERE requisition_id=$1`, id).Scan(&state, &reason); e != nil {
		t.Fatal(e)
	}
	return state, reason
}

func expenseCount(t *testing.T, f usecaseFixture, id string) int {
	t.Helper()
	var n int
	if e := f.pool.QueryRow(context.Background(), `SELECT count(*) FROM expense_claim WHERE requisition_id=$1`, id).Scan(&n); e != nil {
		t.Fatal(e)
	}
	return n
}

// Test plan section 4: leave and rejection paths, each on its own fixture/day.
func TestIntegrationLeaveAndRejectionPaths(t *testing.T) {
	ctx := context.Background()

	// 2W -> 4A -> 1W/3W
	t.Run("leave approved with replacement", func(t *testing.T) {
		f := fixtureUsecases(t)
		s := f.service
		rec := &lineRecorder{}
		s.Notify = rec
		*f.now = time.Date(2026, 10, 20, 9, 0, 0, 0, Bangkok)
		if _, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-22", "08:00"}); e != nil {
			t.Fatal(e)
		}
		leave, e := s.Leave(ctx, f.worker, LeaveInput{"2026-10-22", "family"})
		if e != nil {
			t.Fatal(e)
		}
		if _, e = s.Leave(ctx, f.worker, LeaveInput{"2026-10-22", "again"}); e == nil {
			t.Fatal("duplicate leave accepted")
		}
		ids, e := s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: "approved", WorkerID: f.otherID})
		if e != nil || len(ids) != 1 {
			t.Fatal(ids, e)
		}
		if leaveStatus(t, f, leave) != "approved" {
			t.Fatal("leave not approved")
		}
		if _, e = s.Leave(ctx, f.worker, LeaveInput{"2026-10-22", "after review"}); e == nil {
			t.Fatal("leave re-submitted after review")
		}
		rec.waitFor(t, f.worker.UserID.String(), "leave approved", "family")
		rec.waitFor(t, f.other.UserID.String(), "schedule details", "08:00")

		// Replacement sees exactly the original slot (1W/3W).
		data, e := s.Schedules(ctx, f.other, 100, 0)
		if e != nil {
			t.Fatal(e)
		}
		var rows []struct {
			ID         string `json:"schedule_id"`
			Assignment string `json:"assignment_id"`
			Date       string `json:"work_date"`
			Start      string `json:"shift_start_time"`
			Status     string `json:"shift_status"`
		}
		if e = json.Unmarshal(data, &rows); e != nil {
			t.Fatal(e)
		}
		found := false
		for _, r := range rows {
			if r.ID == ids[0] {
				found = true
				if r.Assignment != f.assignment || !strings.HasPrefix(r.Date, "2026-10-22") || !strings.HasPrefix(r.Start, "08:00") || r.Status != "scheduled" {
					t.Fatalf("replacement shift wrong: %+v", r)
				}
			}
		}
		if !found {
			t.Fatal("replacement cannot see shift", string(data))
		}
		// The original worker no longer counts toward staffing that day.
		var staffed int
		if e = f.pool.QueryRow(ctx, query49, f.assignment, "2026-10-22").Scan(&staffed); e != nil || staffed != 1 {
			t.Fatal("staffing after replacement", staffed, e)
		}
	})

	// 2W -> 4A (reject)
	t.Run("leave rejected", func(t *testing.T) {
		f := fixtureUsecases(t)
		s := f.service
		rec := &lineRecorder{}
		s.Notify = rec
		*f.now = time.Date(2026, 10, 20, 9, 0, 0, 0, Bangkok)
		if _, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-23", "08:00"}); e != nil {
			t.Fatal(e)
		}
		leave, e := s.Leave(ctx, f.worker, LeaveInput{"2026-10-23", "doctor"})
		if e != nil {
			t.Fatal(e)
		}
		if _, e = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: "rejected", WorkerID: f.otherID, Reason: "หาคนแทนไม่ทัน"}); e == nil {
			t.Fatal("replacement accepted with rejection")
		}
		if _, e = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: "rejected"}); e == nil {
			t.Fatal("rejection without reason accepted")
		}
		if _, e = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: "rejected", Reason: "หาคนแทนไม่ทัน"}); e != nil {
			t.Fatal(e)
		}
		if leaveStatus(t, f, leave) != "rejected" {
			t.Fatal("leave not rejected")
		}
		var storedReason string
		if e = f.pool.QueryRow(ctx, `SELECT reason FROM leave_request WHERE request_id=$1`, leave).Scan(&storedReason); e != nil || !strings.HasSuffix(storedReason, "[Rejection] หาคนแทนไม่ทัน") {
			t.Fatal("leave rejection reason not stored", storedReason, e)
		}
		if _, e = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: "approved"}); e == nil {
			t.Fatal("reviewed leave changed")
		}
		rec.waitFor(t, f.worker.UserID.String(), "leave rejected", "doctor", "เหตุผลที่ไม่อนุมัติ: หาคนแทนไม่ทัน")
		// Worker keeps the shift.
		var state string
		if e = f.pool.QueryRow(ctx, `SELECT shift_status::text FROM work_schedule WHERE worker_id=$1 AND work_date='2026-10-23'`, f.workerID).Scan(&state); e != nil || state != "scheduled" {
			t.Fatal("rejected leave changed shift", state, e)
		}
	})

	// 5W -> 5A/6A -> 6W
	t.Run("no purchase needed", func(t *testing.T) {
		f := fixtureUsecases(t)
		s := f.service
		rec := &lineRecorder{}
		s.Notify = rec
		*f.now = time.Date(2026, 10, 24, 9, 0, 0, 0, Bangkok)
		shift, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-24", "08:00"})
		if e != nil {
			t.Fatal(e)
		}
		checkIn(t, f, shift[0])
		equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
		if e != nil {
			t.Fatal(e)
		}
		req, e := s.Requisition(ctx, f.worker, RequestInput{f.assignment, "broom broken", []RequestItem{{equipment, 1, ""}}})
		if e != nil {
			t.Fatal(e)
		}
		if e = s.DecideRequest(ctx, f.assistant, req, ProcurementDecision{"no_purchase", "spare in storeroom"}); e != nil {
			t.Fatal(e)
		}
		if state, _ := requestState(t, f, req); state != "rejected" {
			t.Fatal("no_purchase state", state)
		}
		result, e := s.RetryNotification(ctx, f.assistant, req, NotificationInput{"no_purchase", "ใช้ไม้กวาดสำรองในห้องเก็บของ"})
		if e != nil || result.Status != "accepted" {
			t.Fatal(result, e)
		}
		rec.waitFor(t, f.worker.UserID.String(), "declined", "ใช้ไม้กวาดสำรองในห้องเก็บของ")
		if expenseCount(t, f, req) != 0 {
			t.Fatal("no_purchase created expense")
		}
		if _, e = s.Fund(ctx, f.supervisor, req, FundInput{"10", uuid.NewString(), f.supervisor.UserID.String() + "/" + uuid.NewString()}); e == nil {
			t.Fatal("funded a no_purchase request")
		}
		if expenseCount(t, f, req) != 0 {
			t.Fatal("expense created after no_purchase")
		}
	})

	// 6A -> 3S (reject) -> 6W
	t.Run("supervisor rejects purchase", func(t *testing.T) {
		f := fixtureUsecases(t)
		s := f.service
		rec := &lineRecorder{}
		s.Notify = rec
		*f.now = time.Date(2026, 10, 25, 9, 0, 0, 0, Bangkok)
		shift, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-25", "08:00"})
		if e != nil {
			t.Fatal(e)
		}
		checkIn(t, f, shift[0])
		equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
		if e != nil {
			t.Fatal(e)
		}
		req, e := s.Requisition(ctx, f.worker, RequestInput{f.assignment, "need gloves", []RequestItem{{equipment, 2, ""}}})
		if e != nil {
			t.Fatal(e)
		}
		var item string
		if e = f.pool.QueryRow(ctx, `SELECT item_id::text FROM requisition_item WHERE requisition_id=$1`, req).Scan(&item); e != nil {
			t.Fatal(e)
		}
		if e = s.DecideRequest(ctx, f.assistant, req, ProcurementDecision{"purchase", "worn out"}); e != nil {
			t.Fatal(e)
		}
		if e = s.ReviewRequest(ctx, f.supervisor, req, ReviewInput{"reject", ""}); e == nil {
			t.Fatal("rejection without reason accepted")
		}
		if e = s.ReviewRequest(ctx, f.supervisor, req, ReviewInput{"reject", "งบเดือนนี้หมดแล้ว"}); e != nil {
			t.Fatal(e)
		}
		state, reason := requestState(t, f, req)
		if state != "rejected" || !strings.HasSuffix(reason, "[Rejection] งบเดือนนี้หมดแล้ว") {
			t.Fatal("rejection not stored", state, reason)
		}
		rec.waitFor(t, f.worker.UserID.String(), "equipment rejected", "งบเดือนนี้หมดแล้ว")
		if e = s.ReviewRequest(ctx, f.supervisor, req, ReviewInput{"approve", ""}); e == nil {
			t.Fatal("rejected request re-approved")
		}
		if _, e = s.Fund(ctx, f.supervisor, req, FundInput{"100", uuid.NewString(), f.supervisor.UserID.String() + "/" + uuid.NewString()}); e == nil {
			t.Fatal("funded a rejected request")
		}
		zero := 0
		if _, e = s.Purchase(ctx, f.assistant, req, PurchaseInput{[]PurchaseItem{{item, 2, "50", &zero}}, f.assistant.UserID.String() + "/" + uuid.NewString()}); e == nil {
			t.Fatal("purchased a rejected request")
		}
		if expenseCount(t, f, req) != 0 {
			t.Fatal("expense created for rejected request")
		}
		if state, _ = requestState(t, f, req); state != "rejected" {
			t.Fatal("rejected state changed", state)
		}
	})
}
