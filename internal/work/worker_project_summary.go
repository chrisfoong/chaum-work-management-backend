package work

import (
	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"fmt"
	"strings"
)

// WorkerProjectSummaryMessage is a personal work summary, separate from staff
// report 9A. The caller's identity scopes every count; no coworker or financial
// information is included. This does not grant access to OperationsSummary.
func (s *Service) WorkerProjectSummaryMessage(ctx context.Context, p auth.Principal, tor, start, end string) (string, error) {
	if p.Role != auth.RoleWorker {
		return "", apperr.Forbidden("Worker personal summary only")
	}
	if err := validID(tor, "tor_id"); err != nil {
		return "", err
	}
	if err := validRange(start, end); err != nil {
		return "", err
	}
	worker, err := s.Worker(ctx, s.Repo.Pool, p.UserID.String())
	if err != nil {
		return "", err
	}
	var project string
	var total, finished, working, onTime, late, absent, leave, waiting int
	err = s.Repo.Pool.QueryRow(ctx, `SELECT c.project_name,count(sc.schedule_id),
 count(*) FILTER(WHERE at.check_out IS NOT NULL),
 count(*) FILTER(WHERE at.check_in IS NOT NULL AND at.check_out IS NULL),
 count(*) FILTER(WHERE at.status='on_time'),count(*) FILTER(WHERE at.status='late'),
 count(*) FILTER(WHERE at.status='absent'),count(*) FILTER(WHERE at.status='leave'),
 count(*) FILTER(WHERE at.schedule_id IS NULL)
 FROM contract_tor c JOIN tor_location_assignment a USING(tor_id)
 JOIN work_schedule sc USING(assignment_id)
 LEFT JOIN attendance at ON at.schedule_id=sc.schedule_id AND at.worker_id=sc.worker_id AND at.work_date=sc.work_date
 WHERE c.tor_id=$1 AND sc.worker_id=$2 AND sc.work_date BETWEEN $3::date AND $4::date AND sc.shift_status<>'cancelled'
 GROUP BY c.project_name`, tor, worker, start, end).Scan(&project, &total, &finished, &working, &onTime, &late, &absent, &leave, &waiting)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Chaum: worker project summary.\nโครงการ: %s\nช่วงวันที่: %s – %s\n\nสรุปเฉพาะงานของคุณ\nกะที่ได้รับมอบหมาย: %d กะ\nเช็กเอาท์แล้ว: %d กะ\nเช็กอินแล้วยังไม่เช็กเอาท์: %d กะ\nเข้างานตรงเวลา: %d กะ\nเข้างานสาย: %d กะ\nขาดงาน: %d กะ\nบันทึกสถานะลา: %d กะ\nยังไม่มีบันทึกเข้างาน: %d กะ\n\nข้อมูลล่าสุด ณ เวลาส่ง ไม่ใช่รายงานปิดโครงการ", strings.Join(strings.Fields(project), " "), cardThaiDate(start), cardThaiDate(end), total, finished, working, onTime, late, absent, leave, waiting), nil
}
