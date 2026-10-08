package work

// Exported-schema SQL for people.
const (
	query41 = `INSERT INTO public."USER"(first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no,daily_wage) VALUES($1,$2,$3,$4,$5,$6,$7,$8::numeric) RETURNING user_id::text`
	query42 = `INSERT INTO worker(user_id) VALUES($1)`
	query43 = `SELECT user_id,first_name,last_name,role,is_active FROM public."USER" WHERE user_id=$1`
	query44 = `UPDATE public."USER" SET first_name=$2,last_name=$3,daily_wage=$4::numeric,is_active=COALESCE($5,is_active),updated_at=now() WHERE user_id=$1`
	query45 = `UPDATE worker SET is_available=$2 WHERE worker_id=$1`
	query46 = `SELECT user_id,first_name,last_name,role,phone_number,bank_name,bank_account_no,daily_wage::text,is_active FROM public."USER" WHERE user_id=$1`
	query47 = `SELECT u.user_id,u.first_name,u.last_name,u.role,u.is_active,u.daily_wage::text,w.worker_id,w.is_available FROM public."USER" u LEFT JOIN worker w ON w.user_id=u.user_id ORDER BY u.user_id,w.worker_id LIMIT $1 OFFSET $2`
	query48 = `SELECT a.required_workers,c.start_date::text,c.end_date::text,c.status::text FROM tor_location_assignment a JOIN contract_tor c ON c.tor_id=a.tor_id WHERE a.assignment_id=$1 FOR UPDATE OF a`
	query49 = `SELECT count(*) FROM work_schedule s JOIN worker w ON w.worker_id=s.worker_id WHERE s.assignment_id=$1 AND s.work_date=$2::date AND s.shift_status<>'cancelled' AND NOT EXISTS(SELECT 1 FROM leave_request l WHERE l.user_id=w.user_id AND l.leave_date=s.work_date AND l.status='approved')`
	query50 = `SELECT w.user_id::text,w.is_available,u.is_active FROM worker w JOIN public."USER" u ON u.user_id=w.user_id WHERE w.worker_id=$1 AND u.role='worker' FOR UPDATE OF w`
	query51 = `SELECT EXISTS(SELECT 1 FROM payroll WHERE user_id=$1 AND $2::date BETWEEN period_start AND period_end)`
	query52 = `INSERT INTO work_schedule(assignment_id,worker_id,work_date,shift_start_time) VALUES($1,$2,$3::date,$4::time) RETURNING schedule_id::text`
	query53 = `SELECT ` + frontendScheduleFields + frontendScheduleJoins + ` WHERE ($1::uuid IS NULL OR (w.user_id=$1 AND s.work_date >= $4::date)) ORDER BY s.work_date ASC,s.shift_start_time,s.schedule_id LIMIT $2 OFFSET $3`
	query54 = `SELECT count(*) FROM leave_request WHERE user_id=$1 AND leave_date=$2::date`
	query55 = `SELECT shift_start_time::text FROM work_schedule WHERE worker_id=$1 AND work_date=$2::date AND shift_status='scheduled' FOR UPDATE`
	query56 = `INSERT INTO leave_request(leave_no,user_id,leave_date,reason,is_advance_notice) VALUES($1,$2,$3::date,$4,$5) RETURNING request_id::text`
	query57 = `SELECT l.*,u.first_name,u.last_name FROM leave_request l JOIN public."USER" u ON u.user_id=l.user_id WHERE ($1::uuid IS NULL OR l.user_id=$1) ORDER BY l.created_at DESC,l.request_id LIMIT $2 OFFSET $3`
	query58 = `SELECT user_id::text,leave_date::text,is_advance_notice FROM leave_request WHERE request_id=$1 FOR UPDATE`
	query59 = `SELECT EXISTS(SELECT 1 FROM payroll WHERE user_id=$1 AND $2::date BETWEEN period_start AND period_end)`
	query60 = `SELECT schedule_id::text FROM work_schedule WHERE worker_id=$1 AND work_date=$2::date FOR UPDATE`
	query61 = `SELECT EXISTS(SELECT 1 FROM attendance WHERE schedule_id=$1 AND check_in IS NOT NULL)`
	query62 = `UPDATE leave_request SET status=$2 WHERE request_id=$1 AND status='pending'`
	query63 = `INSERT INTO attendance(schedule_id,worker_id,work_date,status) VALUES($1,$2,$3::date,'leave') ON CONFLICT(schedule_id) DO UPDATE SET status='leave' WHERE attendance.check_in IS NULL`
	query64 = `SELECT l.status::text,s.assignment_id::text,s.work_date::text,s.shift_start_time::text,w.worker_id::text FROM leave_request l JOIN worker w ON w.user_id=l.user_id JOIN work_schedule s ON s.worker_id=w.worker_id AND s.work_date=l.leave_date WHERE l.request_id=$1`
)
