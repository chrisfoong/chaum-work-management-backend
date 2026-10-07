package work

// Exported-schema SQL for dashboard.
const (
	query21 = `SELECT (SELECT count(*) FROM work_schedule sc JOIN worker w ON w.worker_id=sc.worker_id WHERE w.user_id=$1 AND sc.shift_status='scheduled' AND sc.work_date>=($2::timestamptz AT TIME ZONE 'Asia/Bangkok')::date) AS upcoming_schedules,(SELECT count(*) FROM leave_request WHERE user_id=$1 AND status='pending') AS pending_leave,(SELECT COALESCE(sum(net_wage),0)::text FROM payroll WHERE user_id=$1 AND NOT is_paid) AS unpaid_wage`
	query22 = `SELECT (SELECT count(*) FROM contract_tor WHERE status IN ('registered','active')) AS open_contracts,(SELECT count(*) FROM leave_request WHERE status='pending') AS pending_leave,(SELECT count(*) FROM equipment_requisition WHERE status NOT IN ('completed','rejected')) AS open_requisitions`
)
