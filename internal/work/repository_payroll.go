package work

// Exported-schema SQL for payroll.
const (
	query29 = `SELECT daily_wage::text FROM public."USER" WHERE user_id=$1 AND role='worker'`
	query30 = `SELECT s.schedule_id::text,s.work_date::text,s.shift_start_time::text,a.attendance_id::text,a.status::text,a.check_in,EXISTS(SELECT 1 FROM leave_request l WHERE l.user_id=$1 AND l.leave_date=s.work_date AND l.status='approved') FROM work_schedule s LEFT JOIN attendance a ON a.schedule_id=s.schedule_id AND a.worker_id=s.worker_id AND a.work_date=s.work_date WHERE s.worker_id=$2 AND s.work_date BETWEEN $3::date AND $4::date AND s.shift_status<>'cancelled' ORDER BY s.work_date,s.schedule_id`
	query31 = `SELECT EXISTS(SELECT 1 FROM payroll WHERE user_id=$1 AND period_start<=$3::date AND period_end>=$2::date)`
	query32 = `INSERT INTO payroll(payroll_slip_no,user_id,managed_by_id,period_start,period_end,base_wage,total_deduction,net_wage) VALUES($1,$2,$3,$4::date,$5::date,$6::numeric,$7::numeric,$8::numeric) RETURNING payroll_id::text`
	query33 = `SELECT count(*) FROM deduction_transaction WHERE attendance_id=$1`
	query34 = `INSERT INTO deduction_transaction(worker_id,attendance_id,payroll_id,penalty_amount,reason) VALUES($1,$2,$3,$4::numeric,$5)`
	query35 = `SELECT p.payroll_id,p.payroll_slip_no,p.user_id,p.managed_by_id,p.period_start,p.period_end,p.base_wage::text,p.total_deduction::text,p.net_wage::text,p.is_paid,COALESCE((SELECT jsonb_agg(jsonb_build_object('work_date',a.work_date,'penalty_amount',d.penalty_amount::text,'reason',d.reason)) FROM deduction_transaction d LEFT JOIN attendance a ON a.attendance_id=d.attendance_id WHERE d.payroll_id=p.payroll_id),'[]'::jsonb) AS calculated_penalties FROM payroll p WHERE ($1::uuid IS NULL OR p.user_id=$1) ORDER BY p.period_start DESC,p.payroll_id LIMIT $2 OFFSET $3`
	query36 = `UPDATE payroll SET is_paid=true WHERE payroll_id=$1`
	query37 = `SELECT EXISTS(SELECT 1 FROM company_invoice WHERE tor_id=$1 AND billing_month=$2)`
	query38 = `INSERT INTO company_invoice(invoice_no,tor_id,billing_month,expected_amount,exat_deduction_amount,deduction_reason) VALUES($1,$2,$3,$4::numeric,$5::numeric,$6) RETURNING invoice_id::text`
	query39 = `SELECT status::text,net_received::text FROM company_invoice WHERE invoice_id=$1 FOR UPDATE`
	query40 = `UPDATE company_invoice SET status='paid',net_received=$2::numeric WHERE invoice_id=$1`
)
