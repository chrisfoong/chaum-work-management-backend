package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestFrontendFiltersRejectInvalidInputs(t *testing.T) {
	for _, f := range []ListFilter{{Status: "pending_supervisor"}, {AssignmentID: "invalid"}, {From: "2026-10-10", To: "2026-10-01"}, {Kind: "anything"}, {From: "not-a-date"}} {
		if e := f.validate("requisitions"); e == nil {
			t.Errorf("invalid filter accepted: %+v", f)
		}
	}
}
func TestNotificationSplitsLongUnicodeWithoutLosingContent(t *testing.T) {
	text := strings.Repeat("🌿ไทย", 2000)
	parts := notificationParts(text)
	if strings.Join(parts, "") != text || len(parts) < 2 {
		t.Fatal("notification corrupted")
	}
	for _, part := range parts {
		units := 0
		for _, r := range part {
			units++
			if r > 0xffff {
				units++
			}
		}
		if units > 5000 {
			t.Fatal("LINE message exceeds character limit")
		}
	}
}
func TestIntegrationFrontendReadModelsAndDeliveryGuards(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	schedules, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "20:00"})
	if e != nil {
		t.Fatal(e)
	}
	future, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	data, e := s.FilteredList(ctx, f.worker, "schedules", ListFilter{From: "2026-10-01", To: "2026-10-01"}, 100, 0)
	if e != nil {
		t.Fatal(e)
	}
	var shifts []struct {
		Project string `json:"project_name"`
		End     string `json:"shift_end_at"`
	}
	if e = json.Unmarshal(data, &shifts); e != nil || len(shifts) != 1 || shifts[0].Project != "usecases" || !strings.HasPrefix(shifts[0].End, "2026-10-02T") {
		t.Fatal("overnight schedule read model", string(data), e)
	}
	for _, kind := range []string{"leave", "attendance", "requisitions"} {
		if _, e = s.FilteredList(ctx, f.worker, kind, ListFilter{TorID: f.tor}, 100, 0); e != nil {
			t.Fatal(kind, e)
		}
	}
	if _, e = s.Assignments(ctx, f.tor, 100, 0); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Contracts(ctx, 100, 0); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Dashboard(ctx, f.assistant); e != nil {
		t.Fatal(e)
	}
	leave, e := s.Leave(ctx, f.worker, LeaveInput{"2026-10-01", "need leave"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.LeaveDetail(ctx, leave); e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `UPDATE public."USER" SET first_name='000-candidate' WHERE user_id=$1`, f.other.UserID); e != nil {
		t.Fatal(e)
	}
	candidates, e := s.LeaveCandidates(ctx, leave, 100, 0)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(candidates), f.workerID) || !strings.Contains(string(candidates), f.otherID) {
		t.Fatal("candidates did not exclude leaver / include free worker", e)
	}
	if _, e = s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.otherID}, "2026-10-01", "08:00"}); e != nil {
		t.Fatal(e)
	}
	candidates, e = s.LeaveCandidates(ctx, leave, 100, 0)
	if e != nil || strings.Contains(string(candidates), f.otherID) {
		t.Fatal("busy candidate included", e)
	}
	var request, item string
	equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	if e = f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'delivery','additional','completed') RETURNING requisition_id::text`, uuid.NewString(), f.worker.UserID, f.assignment).Scan(&request); e != nil {
		t.Fatal(e)
	}
	if e = f.pool.QueryRow(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty,actual_qty,actual_price) VALUES($1,$2,1,1,1) RETURNING item_id::text`, request, equipment).Scan(&item); e != nil {
		t.Fatal(e)
	}
	input := DeliveryInput{future[0], "received", []string{f.assistant.UserID.String() + "/" + uuid.NewString()}}
	if e = s.Deliver(ctx, f.assistant, request, input); e == nil {
		t.Fatal("future handover accepted")
	}
	choices, e := s.DeliverySchedules(ctx, request, 100, 0)
	if e != nil || strings.Contains(string(choices), future[0]) || !strings.Contains(string(choices), schedules[0]) {
		t.Fatal("handover picker", e)
	}
	input.ScheduleID = schedules[0]
	if _, e = f.pool.Exec(ctx, `UPDATE requisition_item SET actual_qty=0 WHERE item_id=$1`, item); e != nil {
		t.Fatal(e)
	}
	if e = s.Deliver(ctx, f.assistant, request, input); e == nil {
		t.Fatal("incomplete quantities accepted despite completed status")
	}
	if _, e = f.pool.Exec(ctx, `UPDATE requisition_item SET actual_qty=1 WHERE item_id=$1`, item); e != nil {
		t.Fatal(e)
	}
	if e = s.Deliver(ctx, f.assistant, request, input); e != nil {
		t.Fatal(e)
	}
	if e = s.Deliver(ctx, f.assistant, request, input); e != nil {
		t.Fatal("retry", e)
	}
	other, e := s.FilteredList(ctx, f.other, "requisitions", ListFilter{Status: "completed", TorID: f.tor}, 100, 0)
	if e != nil || string(other) != "[]" {
		t.Fatal("worker ownership not preserved", e)
	}
	filtered, e := s.FilteredList(ctx, f.assistant, "requisitions", ListFilter{Status: "completed", Kind: "additional", TorID: f.tor}, 1, 0)
	if e != nil || !strings.Contains(string(filtered), request) {
		t.Fatal("filter must apply before paging", e)
	}
	// A duplicated mapping must conflict, even before a candidate list is paged.
	var duplicate string
	if e = f.pool.QueryRow(ctx, `INSERT INTO worker(user_id) VALUES($1) RETURNING worker_id::text`, f.other.UserID).Scan(&duplicate); e != nil {
		t.Fatal(e)
	}
	defer f.pool.Exec(ctx, `DELETE FROM worker WHERE worker_id=$1`, duplicate)
	if _, e = s.AvailableWorkers(ctx, "2026-10-03", "", 1, 0); e == nil {
		t.Fatal("duplicate mapping silently hidden by paging")
	}
}
