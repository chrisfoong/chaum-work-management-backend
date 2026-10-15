package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestIntegrationCheckoutWaitsForScheduledOvernightEnd(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	ids, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "20:00"})
	if e != nil {
		t.Fatal(e)
	}
	*f.now = time.Date(2026, 10, 1, 20, 7, 0, 0, Bangkok)
	qr, e := s.QR(ctx, f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	lat, lon, accuracy := 13.75, 100.5, 10.0
	if e = s.CheckIn(ctx, f.worker, CheckIn{ids[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
	in := CheckOut{ids[0], "finished", f.worker.UserID.String() + "/" + uuid.NewString()}
	for _, early := range []time.Time{time.Date(2026, 10, 1, 23, 0, 0, 0, Bangkok), time.Date(2026, 10, 2, 3, 59, 59, 999999999, Bangkok)} {
		*f.now = early
		if e = s.CheckOut(ctx, f.worker, in); e == nil {
			t.Fatal("early checkout accepted")
		}
		var untouched bool
		if e = f.pool.QueryRow(ctx, `SELECT a.check_out IS NULL AND s.shift_status='scheduled' AND NOT EXISTS(SELECT 1 FROM work_evidence WHERE schedule_id=s.schedule_id) FROM attendance a JOIN work_schedule s ON s.schedule_id=a.schedule_id WHERE s.schedule_id=$1`, ids[0]).Scan(&untouched); e != nil || !untouched {
			t.Fatal("early checkout changed records", e)
		}
	}
	*f.now = time.Date(2026, 10, 2, 4, 0, 0, 0, Bangkok)
	data, e := s.FilteredList(ctx, f.worker, "schedules", ListFilter{}, 100, 0)
	if e != nil {
		t.Fatal(e)
	}
	var open []map[string]any
	if json.Unmarshal(data, &open) != nil || len(open) != 1 {
		t.Fatal("unclosed overnight shift disappeared at checkout boundary", string(data))
	}
	if e = s.CheckOut(ctx, f.worker, in); e != nil {
		t.Fatal("checkout at exact scheduled end rejected", e)
	}
	if e = s.CheckOut(ctx, f.worker, in); e == nil {
		t.Fatal("duplicate checkout accepted")
	}
	var saved time.Time
	if e = f.pool.QueryRow(ctx, `SELECT check_out FROM attendance WHERE schedule_id=$1`, ids[0]).Scan(&saved); e != nil || !saved.Equal(*f.now) {
		t.Fatal("wrong checkout timestamp", e)
	}
}
