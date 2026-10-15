package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
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
		Start   string `json:"shift_start_at"`
	}
	if e = json.Unmarshal(data, &shifts); e != nil || len(shifts) != 1 || shifts[0].Project != "usecases" {
		t.Fatal("overnight schedule read model", string(data), e)
	}
	start, startErr := time.Parse(time.RFC3339, shifts[0].Start)
	end, endErr := time.Parse(time.RFC3339, shifts[0].End)
	if startErr != nil || endErr != nil || end.Sub(start) != 8*time.Hour || start.In(Bangkok).Format("2006-01-02 15:04") != "2026-10-01 20:00" || end.In(Bangkok).Format("2006-01-02 15:04") != "2026-10-02 04:00" {
		t.Fatal("schedule must cross midnight in Bangkok and last exactly eight hours")
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

func TestIntegrationDashboardShowsOnlyCurrentWork(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	read := func() struct {
		Open    int `json:"open_contracts"`
		Pending []struct {
			ID     string `json:"requisition_id"`
			Status string `json:"status"`
		} `json:"pending_requisitions"`
		Contracts []struct {
			ID     string `json:"tor_id"`
			Number string `json:"contract_no"`
		} `json:"contracts"`
	} {
		t.Helper()
		data, err := f.service.Dashboard(ctx, f.supervisor)
		if err != nil {
			t.Fatal(err)
		}
		var result struct {
			Open    int `json:"open_contracts"`
			Pending []struct {
				ID     string `json:"requisition_id"`
				Status string `json:"status"`
			} `json:"pending_requisitions"`
			Contracts []struct {
				ID     string `json:"tor_id"`
				Number string `json:"contract_no"`
			} `json:"contracts"`
		}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	before := read()
	found := false
	for _, c := range before.Contracts {
		if c.ID == f.tor {
			found = c.Number != ""
		}
	}
	if !found {
		t.Fatal("dashboard must include the actual contract number")
	}
	if _, err := f.pool.Exec(ctx, `UPDATE contract_tor SET start_date='2026-11-01',end_date='2026-11-30' WHERE tor_id=$1`, f.tor); err != nil {
		t.Fatal(err)
	}
	if after := read(); after.Open != before.Open-1 {
		t.Fatal("future contract counted as operating", before.Open, after.Open)
	}
	if _, err := f.pool.Exec(ctx, `UPDATE contract_tor SET start_date='2026-09-01',end_date='2026-09-30' WHERE tor_id=$1`, f.tor); err != nil {
		t.Fatal(err)
	}
	if after := read(); after.Open != before.Open-1 {
		t.Fatal("expired contract counted as operating")
	}
	if _, err := f.pool.Exec(ctx, `UPDATE contract_tor SET status='registered',start_date='2026-10-01',end_date='2026-10-31' WHERE tor_id=$1`, f.tor); err != nil {
		t.Fatal(err)
	}
	if after := read(); after.Open != before.Open-1 {
		t.Fatal("registered contract counted as operating")
	}
	var pendingID string
	if err := f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'dashboard','tor_base','pending_survey') RETURNING requisition_id::text`, uuid.NewString(), f.supervisor.UserID, f.assignment).Scan(&pendingID); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := f.pool.Exec(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'finished','tor_base','completed')`, uuid.NewString(), f.supervisor.UserID, f.assignment); err != nil {
			t.Fatal(err)
		}
	}
	after := read()
	found = false
	for _, r := range after.Pending {
		if r.Status == "completed" || r.Status == "rejected" {
			t.Fatal("terminal request shown as pending")
		}
		if r.ID == pendingID {
			found = true
		}
	}
	if !found {
		t.Fatal("newer completed requests must not hide pending work")
	}
}
