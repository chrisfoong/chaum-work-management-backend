package work

// Exported-schema SQL for report.
const (
	query87 = `
 WITH days AS (
 SELECT p.payroll_id,p.net_wage,a.tor_id,count(*) AS days
 FROM payroll p JOIN worker w ON w.user_id=p.user_id
 JOIN work_schedule sc ON sc.worker_id=w.worker_id AND sc.work_date BETWEEN p.period_start AND p.period_end
 JOIN attendance at ON at.schedule_id=sc.schedule_id AND at.worker_id=sc.worker_id AND at.work_date=sc.work_date
 JOIN tor_location_assignment a ON a.assignment_id=sc.assignment_id
 WHERE p.is_paid AND to_char(p.period_start,'YYYY-MM')=$1 AND at.status IN ('on_time','late') AND sc.shift_status<>'cancelled'
 AND NOT EXISTS(SELECT 1 FROM leave_request l WHERE l.user_id=w.user_id AND l.leave_date=sc.work_date AND l.status='approved' AND l.is_advance_notice)
 GROUP BY p.payroll_id,p.net_wage,a.tor_id),
 shares AS (SELECT *,net_wage*100*days/sum(days) OVER(PARTITION BY payroll_id) AS exact FROM days),
 rounded AS(SELECT *,floor(exact) AS cents,row_number() OVER(PARTITION BY payroll_id ORDER BY exact-floor(exact) DESC,tor_id) AS rank FROM shares),
 allocated AS(SELECT payroll_id,days,tor_id,(cents+CASE WHEN rank<=net_wage*100-sum(cents) OVER(PARTITION BY payroll_id) THEN 1 ELSE 0 END)/100 AS amount FROM rounded),
 labor AS(SELECT tor_id,sum(amount) AS amount,jsonb_agg(jsonb_build_object('payroll_id',payroll_id,'workdays',days,'allocated_net_wage',round(amount,2)::text) ORDER BY payroll_id) AS details FROM allocated GROUP BY tor_id),
 revenue AS(SELECT tor_id,sum(net_received) AS amount FROM company_invoice WHERE status='paid' AND billing_month=$1 GROUP BY tor_id),
 material AS(SELECT a.tor_id,sum(e.total_amount) AS amount,jsonb_agg(jsonb_build_object('expense_id',e.expense_id,'requisition_id',e.requisition_id,'amount',e.total_amount::text,'created_at',e.created_at) ORDER BY e.created_at,e.expense_id) AS details FROM expense_claim e JOIN equipment_requisition r ON r.requisition_id=e.requisition_id JOIN tor_location_assignment a ON a.assignment_id=r.assignment_id WHERE e.expense_type='actual_expense' AND to_char(e.created_at AT TIME ZONE 'Asia/Bangkok','YYYY-MM')=$1 GROUP BY a.tor_id)
 SELECT c.tor_id,c.project_name,$1::text AS month,COALESCE(r.amount,0)::text AS revenue,round(COALESCE(l.amount,0),2)::text AS labor,COALESCE(m.amount,0)::text AS material,round(COALESCE(r.amount,0)-COALESCE(l.amount,0)-COALESCE(m.amount,0),2)::text AS profit,COALESCE(l.details,'[]'::jsonb) AS labor_details,COALESCE(m.details,'[]'::jsonb) AS material_details,true AS estimated,'paid invoices by billing_month; payroll by workdays; actual expenses by purchase month'::text AS basis
 FROM contract_tor c LEFT JOIN revenue r ON r.tor_id=c.tor_id LEFT JOIN labor l ON l.tor_id=c.tor_id LEFT JOIN material m ON m.tor_id=c.tor_id ORDER BY c.tor_id`
	query88 = `SELECT invoice_id,invoice_no,tor_id,billing_month,expected_amount::text,net_received::text,exat_deduction_amount::text,deduction_reason,status FROM company_invoice ORDER BY billing_month DESC,invoice_id LIMIT $1 OFFSET $2`
)
