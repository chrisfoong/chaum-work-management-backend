package work

import (
	"context"
)

type InvoiceDeductionLine struct {
	Date   string `json:"work_date"`
	Worker string `json:"worker_name"`
	Reason string `json:"reason"`
	Amount string `json:"amount"`
}
type InvoiceDeductionPreview struct {
	ContractValue   string                 `json:"contract_value"`
	BilledAmount    string                 `json:"billed_amount"`
	RemainingAmount string                 `json:"remaining_amount"`
	TorID           string                 `json:"tor_id"`
	Month           string                 `json:"billing_month"`
	Total           string                 `json:"deduction_amount"`
	Pending         int                    `json:"pending_payroll"`
	Lines           []InvoiceDeductionLine `json:"lines"`
}

// Allocate only the deduction actually applied to each payslip. Largest-remainder
// allocation preserves every satang when a capped slip spans projects/months.
const invoiceDeductionEntries = `WITH targets AS (
 SELECT DISTINCT p.payroll_id,p.total_deduction
 FROM payroll p JOIN worker w ON w.user_id=p.user_id
 JOIN work_schedule s ON s.worker_id=w.worker_id
 JOIN tor_location_assignment ta ON ta.assignment_id=s.assignment_id
 WHERE ta.tor_id=$1 AND to_char(s.work_date,'YYYY-MM')=$2 AND s.shift_status<>'cancelled' AND s.work_date BETWEEN p.period_start AND p.period_end
), weights AS (
 SELECT d.deduction_id,d.reason,d.penalty_amount*100 AS weight,
 p.total_deduction*100 AS budget,d.payroll_id,
 sum(d.penalty_amount*100) OVER(PARTITION BY d.payroll_id) AS denominator,
 a.work_date,ta.tor_id,u.first_name||' '||u.last_name AS worker_name
 FROM targets p JOIN deduction_transaction d USING(payroll_id)
 JOIN worker w ON w.worker_id=d.worker_id JOIN public."USER" u ON u.user_id=w.user_id
 LEFT JOIN attendance a ON a.attendance_id=d.attendance_id
 LEFT JOIN work_schedule s ON s.schedule_id=a.schedule_id
 LEFT JOIN tor_location_assignment ta ON ta.assignment_id=s.assignment_id
), shares AS (
 SELECT *,CASE WHEN denominator>0 THEN floor(budget*weight/denominator) ELSE 0 END AS base,
 CASE WHEN denominator>0 THEN mod(budget*weight,denominator) ELSE 0 END AS remainder FROM weights
), ranked AS (
 SELECT *,sum(base) OVER(PARTITION BY payroll_id) AS allocated,
 row_number() OVER(PARTITION BY payroll_id ORDER BY remainder DESC,deduction_id) AS position FROM shares
)
 SELECT COALESCE(work_date::text,''),COALESCE(tor_id::text,''),worker_name,reason,
 round((base+CASE WHEN position<=budget-allocated THEN 1 ELSE 0 END)/100,2)::text
 FROM ranked WHERE budget>0 ORDER BY work_date,deduction_id`

func (s *Service) InvoiceDeductions(ctx context.Context, tor, month string) (InvoiceDeductionPreview, error) {
	if e := validID(tor, "tor_id"); e != nil {
		return InvoiceDeductionPreview{}, e
	}
	if e := Month(month); e != nil {
		return InvoiceDeductionPreview{}, e
	}
	return s.invoiceDeductions(ctx, s.Repo.Pool, tor, month)
}
func (s *Service) invoiceDeductions(ctx context.Context, q Query, tor, month string) (InvoiceDeductionPreview, error) {
	out := InvoiceDeductionPreview{TorID: tor, Month: month, Total: "0.00", Lines: []InvoiceDeductionLine{}}
	if e := q.QueryRow(ctx, `SELECT c.contract_value::text,COALESCE(sum(i.expected_amount),0)::text,(c.contract_value-COALESCE(sum(i.expected_amount),0))::text FROM contract_tor c LEFT JOIN company_invoice i USING(tor_id) WHERE c.tor_id=$1 GROUP BY c.tor_id`, tor).Scan(&out.ContractValue, &out.BilledAmount, &out.RemainingAmount); e != nil {
		return out, e
	}
	var exists bool
	if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contract_tor WHERE tor_id=$1)`, tor).Scan(&exists); e != nil {
		return out, e
	}
	if !exists {
		return out, conflict("ไม่พบโครงการ")
	}
	if e := q.QueryRow(ctx, `SELECT count(*) FROM work_schedule s JOIN tor_location_assignment ta USING(assignment_id) JOIN worker w USING(worker_id) WHERE ta.tor_id=$1 AND to_char(s.work_date,'YYYY-MM')=$2 AND s.shift_status<>'cancelled' AND NOT EXISTS(SELECT 1 FROM payroll p WHERE p.user_id=w.user_id AND s.work_date BETWEEN p.period_start AND p.period_end)`, tor, month).Scan(&out.Pending); e != nil {
		return out, e
	}
	var invalid bool
	if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM payroll p JOIN worker w ON w.user_id=p.user_id JOIN work_schedule s ON s.worker_id=w.worker_id JOIN tor_location_assignment ta USING(assignment_id) WHERE ta.tor_id=$1 AND to_char(s.work_date,'YYYY-MM')=$2 AND s.shift_status<>'cancelled' AND p.total_deduction>0 AND s.work_date BETWEEN p.period_start AND p.period_end AND NOT EXISTS(SELECT 1 FROM deduction_transaction d WHERE d.payroll_id=p.payroll_id AND d.penalty_amount>0))`, tor, month).Scan(&invalid); e != nil {
		return out, e
	}
	if invalid {
		return out, conflict("สลิปมียอดหักแต่ไม่มีประวัติเงินหัก กรุณาตรวจสอบข้อมูลค่าจ้าง")
	}
	rows, e := q.Query(ctx, invoiceDeductionEntries, tor, month)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	var total int64
	for rows.Next() {
		var line InvoiceDeductionLine
		var source string
		if e = rows.Scan(&line.Date, &source, &line.Worker, &line.Reason, &line.Amount); e != nil {
			return out, e
		}
		if line.Date == "" || source == "" {
			return out, conflict("ประวัติเงินหักไม่มีการเข้างานอ้างอิง กรุณาตรวจสอบค่าจ้าง")
		}
		if source != tor || len(line.Date) < 7 || line.Date[:7] != month {
			continue
		}
		amount, e := Money(line.Amount)
		if e != nil {
			return out, e
		}
		total += amount
		out.Lines = append(out.Lines, line)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	out.Total = Decimal(total)
	return out, nil
}
func (s *Service) lockInvoicePayroll(ctx context.Context, q Query, tor, month string) error {
	var locked string
	if e := q.QueryRow(ctx, `SELECT tor_id::text FROM contract_tor WHERE tor_id=$1 FOR UPDATE`, tor).Scan(&locked); e != nil {
		return e
	}
	rows, e := q.Query(ctx, `SELECT DISTINCT w.user_id::text FROM worker w JOIN work_schedule sc USING(worker_id) JOIN tor_location_assignment ta USING(assignment_id) WHERE ta.tor_id=$1 AND to_char(sc.work_date,'YYYY-MM')=$2 ORDER BY w.user_id::text`, tor, month)
	if e != nil {
		return e
	}
	users := []string{}
	for rows.Next() {
		var user string
		if e = rows.Scan(&user); e != nil {
			rows.Close()
			return e
		}
		users = append(users, user)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, user := range users {
		if _, e = s.Worker(ctx, q, user); e != nil {
			return e
		}
		if e = lock(ctx, q, "payroll:"+user); e != nil {
			return e
		}
	}
	return nil
}
