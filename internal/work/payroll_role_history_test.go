package work

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationPayrollRetainsScheduledWorkAfterRoleChange(t *testing.T) {
	for _, role := range []string{"assistant", "supervisor"} {
		t.Run(role, func(t *testing.T) {
			f := fixtureUsecases(t)
			ctx := context.Background()
			*f.now = time.Date(2026, 10, 4, 8, 0, 0, 0, Bangkok)
			ids, err := f.service.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-04", "08:00"})
			if err != nil {
				t.Fatal(err)
			}
			qr, err := f.service.QR(ctx, f.assignment)
			if err != nil {
				t.Fatal(err)
			}
			lat, lon, accuracy := 13.75, 100.5, 10.0
			if err = f.service.CheckIn(ctx, f.worker, CheckIn{ids[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); err != nil {
				t.Fatal(err)
			}
			if _, err = f.pool.Exec(ctx, `UPDATE public."USER" SET role=$1 WHERE user_id=$2`, role, f.worker.UserID); err != nil {
				t.Fatal(err)
			}
			*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
			if _, err = f.service.Finalize(ctx); err != nil {
				t.Fatal(err)
			}
			input := PayrollBatchInput{"2026-10-01", "2026-10-15"}
			batch, err := f.service.PayrollBatch(ctx, f.supervisor, input)
			if err != nil {
				t.Fatal("historical scheduled worker omitted", err)
			}
			var own PayrollBatchResult
			for _, result := range batch {
				if result.UserID == f.worker.UserID.String() {
					own = result
				}
			}
			if own.PayrollID == "" {
				t.Fatal("historical worker absent from batch")
			}
			var base, net string
			if err = f.pool.QueryRow(ctx, `SELECT base_wage::text,net_wage::text FROM payroll WHERE payroll_id=$1`, own.PayrollID).Scan(&base, &net); err != nil || base != "400.00" || net != "400.00" {
				t.Fatalf("historical work wage lost: %s %s %v", base, net, err)
			}
			again, err := f.service.PayrollBatch(ctx, f.supervisor, input)
			if err != nil {
				t.Fatal(err)
			}
			matched := 0
			for _, result := range again {
				if result.UserID == f.worker.UserID.String() {
					matched++
					if !result.Existing || result.PayrollID != own.PayrollID {
						t.Fatal("retry duplicated payroll")
					}
				}
			}
			if matched != 1 {
				t.Fatal("retry must include exactly one historical worker")
			}
			if _, err = f.service.PreviewPayroll(ctx, PayrollInput{f.supervisor.UserID.String(), input.Start, input.End}); err == nil {
				t.Fatal("non-worker without mapping accepted")
			}
		})
	}
}
