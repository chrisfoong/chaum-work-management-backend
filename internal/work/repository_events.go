package work

// Exported-schema SQL for events.
const (
	query23 = `SELECT w.user_id::text FROM work_schedule sc JOIN worker w ON w.worker_id=sc.worker_id WHERE sc.schedule_id=$1`
	query24 = `SELECT user_id::text FROM public."USER" WHERE role='assistant' AND is_active AND $1::uuid IS NOT NULL`
	query25 = `SELECT user_id::text FROM leave_request WHERE request_id=$1`
	query26 = `SELECT requested_by::text FROM equipment_requisition WHERE requisition_id=$1`
	query27 = `SELECT user_id::text FROM public."USER" WHERE role='assistant' AND is_active AND $1::uuid IS NOT NULL`
	query28 = `SELECT user_id::text FROM payroll WHERE payroll_id=$1`
)
