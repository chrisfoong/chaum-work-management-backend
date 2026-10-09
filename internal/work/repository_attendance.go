package work

// Exported-schema SQL for attendance.
const (
	query1  = `SELECT EXISTS(SELECT 1 FROM tor_location_assignment a JOIN contract_tor c ON c.tor_id=a.tor_id WHERE a.assignment_id=$1 AND c.status NOT IN ('cancelled','complete'))`
	query2  = `SELECT s.assignment_id::text,s.work_date::text,s.shift_start_time::text,s.shift_status::text,l.latitude::float8,l.longitude::float8 FROM work_schedule s JOIN tor_location_assignment a ON a.assignment_id=s.assignment_id JOIN location l ON l.location_id=a.location_id WHERE s.schedule_id=$1 AND s.worker_id=$2 FOR UPDATE OF s`
	query3  = `SELECT EXISTS(SELECT 1 FROM leave_request WHERE user_id=$1 AND leave_date=$2::date AND status='approved')`
	query4  = `INSERT INTO attendance(schedule_id,worker_id,work_date,check_in,status) VALUES($1,$2,$3::date,$4,$5) ON CONFLICT(schedule_id) DO UPDATE SET check_in=EXCLUDED.check_in,status=EXCLUDED.status WHERE attendance.check_in IS NULL AND attendance.check_out IS NULL AND attendance.worker_id=EXCLUDED.worker_id AND attendance.work_date=EXCLUDED.work_date AND attendance.status IS DISTINCT FROM 'leave'`
	query5  = `SELECT shift_status::text FROM work_schedule WHERE schedule_id=$1 AND worker_id=$2 FOR UPDATE`
	query6  = `UPDATE attendance SET check_out=$3 WHERE schedule_id=$1 AND worker_id=$2 AND check_in IS NOT NULL AND check_out IS NULL AND check_in<=$3`
	query7  = `INSERT INTO work_evidence(schedule_id,worker_id,description,photo_url) VALUES($1,$2,$3,$4)`
	query8  = `UPDATE work_schedule SET shift_status='completed' WHERE schedule_id=$1`
	query9  = `SELECT a.*,s.assignment_id FROM attendance a JOIN work_schedule s ON s.schedule_id=a.schedule_id JOIN worker w ON w.worker_id=a.worker_id WHERE ($1::uuid IS NULL OR w.user_id=$1) ORDER BY a.work_date DESC,a.attendance_id LIMIT $2 OFFSET $3`
	query10 = `INSERT INTO attendance(schedule_id,worker_id,work_date,status) SELECT s.schedule_id,s.worker_id,s.work_date,CASE WHEN EXISTS(SELECT 1 FROM leave_request l WHERE l.user_id=w.user_id AND l.leave_date=s.work_date AND l.status='approved' AND l.is_advance_notice) THEN 'leave'::attendance_status_enum ELSE 'absent'::attendance_status_enum END FROM work_schedule s JOIN worker w ON w.worker_id=s.worker_id WHERE s.shift_status='scheduled' AND ((s.work_date+s.shift_start_time) AT TIME ZONE 'Asia/Bangkok')+interval '10 hours'<=$1 ON CONFLICT(schedule_id) DO UPDATE SET status=EXCLUDED.status WHERE attendance.check_in IS NULL AND attendance.worker_id=EXCLUDED.worker_id AND attendance.work_date=EXCLUDED.work_date AND attendance.status IS DISTINCT FROM EXCLUDED.status`
)
