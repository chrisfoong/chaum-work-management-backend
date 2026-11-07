package work

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestIntegrationInvoiceDeductionAllocatesCappedPayrollAcrossProjects(t *testing.T) {
	f := fixtureUsecases(t)
	other := fixtureUsecases(t)
	ctx := context.Background()
	s := f.service
	if _, e := f.pool.Exec(ctx, `UPDATE contract_tor SET contract_value=2000 WHERE tor_id=$1`, f.tor); e != nil {
		t.Fatal(e)
	}
	first, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	second, e := s.Schedule(ctx, ScheduleInput{other.assignment, []string{f.workerID}, "2026-10-03", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	var payroll string
	if e = f.pool.QueryRow(ctx, `INSERT INTO payroll(payroll_slip_no,user_id,managed_by_id,period_start,period_end,base_wage,total_deduction,net_wage) VALUES($1,$2,$3,'2026-10-01','2026-10-15',400,400,0) RETURNING payroll_id::text`, uuid.NewString(), f.worker.UserID, f.supervisor.UserID).Scan(&payroll); e != nil {
		t.Fatal(e)
	}
	for i, id := range []string{first[0], second[0]} {
		date, penalty := "2026-10-02", "300.00"
		if i == 1 {
			date, penalty = "2026-10-03", "1500.00"
		}
		var attendance string
		if e = f.pool.QueryRow(ctx, `INSERT INTO attendance(schedule_id,worker_id,work_date,status) VALUES($1,$2,$3::date,'absent') RETURNING attendance_id::text`, id, f.workerID, date).Scan(&attendance); e != nil {
			t.Fatal(e)
		}
		if _, e = f.pool.Exec(ctx, `INSERT INTO deduction_transaction(worker_id,attendance_id,payroll_id,penalty_amount,reason) VALUES($1,$2,$3,$4::numeric,'test penalty')`, f.workerID, attendance, payroll, penalty); e != nil {
			t.Fatal(e)
		}
	}
	a, e := s.InvoiceDeductions(ctx, f.tor, "2026-10")
	if e != nil || a.Total != "66.67" || a.Pending != 0 || len(a.Lines) != 1 {
		t.Fatal(a, e)
	}
	b, e := s.InvoiceDeductions(ctx, other.tor, "2026-10")
	if e != nil || b.Total != "333.33" {
		t.Fatal(b, e)
	}
	if _, e = s.invoice(ctx, InvoiceInput{f.tor, "2026-10", "10", "0", "test"}, true); e == nil {
		t.Fatal("negative net invoice accepted")
	}
	invoice, e := s.invoice(ctx, InvoiceInput{f.tor, "2026-10", "1000", "0", "auto"}, true)
	if e != nil {
		t.Fatal(e)
	}
	var saved string
	if e = f.pool.QueryRow(ctx, `SELECT exat_deduction_amount::text FROM company_invoice WHERE invoice_id=$1`, invoice).Scan(&saved); e != nil || saved != "66.67" {
		t.Fatal(saved, e)
	}
	if e = s.updateInvoice(ctx, invoice, InvoiceEditInput{Expected: "1100", Deduction: "0", Reason: "auto updated", PreviousExpected: "1000"}, true); e != nil {
		t.Fatal(e)
	}
	if e = f.pool.QueryRow(ctx, `SELECT exat_deduction_amount::text FROM company_invoice WHERE invoice_id=$1`, invoice).Scan(&saved); e != nil || saved != "66.67" {
		t.Fatal(saved, e)
	}
	// A later salary slip outside the selected work period must not be included.
	if _, e = f.pool.Exec(ctx, `INSERT INTO payroll(payroll_slip_no,user_id,managed_by_id,period_start,period_end,base_wage,total_deduction,net_wage) VALUES($1,$2,$3,'2026-09-01','2026-09-15',100,100,0)`, uuid.NewString(), f.worker.UserID, f.supervisor.UserID); e != nil {
		t.Fatal(e)
	}
	a, e = s.InvoiceDeductions(ctx, f.tor, "2026-10")
	if e != nil || a.Total != "66.67" {
		t.Fatal(a, e)
	}
}
func TestIntegrationInvoiceAutoWaitsForPayroll(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	_, e := f.service.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-03", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	preview, e := f.service.InvoiceDeductions(ctx, f.tor, "2026-10")
	if e != nil || preview.Pending != 1 || preview.Total != "0.00" {
		t.Fatal(preview, e)
	}
	if _, e = f.service.invoice(ctx, InvoiceInput{f.tor, "2026-10", "1000", "0", "auto"}, true); e == nil {
		t.Fatal("incomplete payroll silently became zero deduction")
	}
}
