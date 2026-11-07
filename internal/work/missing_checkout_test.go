package work

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestIntegrationMissingCheckoutBlocksPayroll(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	s := f.service
	ids, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	*f.now = time.Date(2026, 10, 1, 8, 1, 0, 0, Bangkok)
	qr, e := s.QR(ctx, f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	lat, lon, acc := 13.75, 100.5, 10.0
	if e = s.CheckIn(ctx, f.worker, CheckIn{ids[0], &lat, &lon, &acc, qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	in := PayrollInput{f.worker.UserID.String(), "2026-10-01", "2026-10-15"}
	preview, e := s.PreviewPayroll(ctx, in)
	if e != nil || preview.AttendanceComplete || preview.PendingCheckouts != 1 || preview.Workdays != 0 {
		t.Fatal(preview, e)
	}
	if _, e = s.Payroll(ctx, f.supervisor, in); e == nil {
		t.Fatal("unclosed work paid")
	}
	if _, e = s.PayrollBatch(ctx, f.supervisor, PayrollBatchInput{in.Start, in.End}); e == nil {
		t.Fatal("batch accepted unclosed work")
	}
	data, e := s.FilteredList(ctx, f.supervisor, "attendance", ListFilter{TorID: f.tor}, 100, 0)
	if e != nil || !strings.Contains(string(data), `"checkout_review_required": true`) {
		t.Fatal(string(data), e)
	}
	var count int
	if e = f.pool.QueryRow(ctx, `SELECT count(*) FROM payroll WHERE user_id=$1`, f.worker.UserID).Scan(&count); e != nil || count != 0 {
		t.Fatal("payroll persisted", count, e)
	}
	if e = s.CheckOut(ctx, f.worker, CheckOut{ids[0], "finished", f.worker.UserID.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	preview, e = s.PreviewPayroll(ctx, in)
	if e != nil || !preview.AttendanceComplete || preview.PendingCheckouts != 0 || preview.Workdays != 1 {
		t.Fatal(preview, e)
	}
	if _, e = s.Payroll(ctx, f.supervisor, in); e != nil {
		t.Fatal("completed work not payable", e)
	}
}
