package work

import (
	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Explicit field lists keep Supervisor financial data out of the 9A response.
const operationsSummarySQL = `WITH shifts AS (
 SELECT s.schedule_id,at.status FROM work_schedule s
 LEFT JOIN attendance at ON at.schedule_id=s.schedule_id AND at.worker_id=s.worker_id AND at.work_date=s.work_date
 WHERE s.assignment_id=$1 AND s.work_date BETWEEN $2::date AND $3::date AND s.shift_status<>'cancelled'),
 requests AS (SELECT requisition_id,status FROM equipment_requisition WHERE assignment_id=$1 AND (created_at AT TIME ZONE 'Asia/Bangkok')::date BETWEEN $2::date AND $3::date),
 progress AS (SELECT r.status,count(DISTINCT r.requisition_id) AS request_count,COALESCE(sum(i.required_qty),0) AS required_qty,COALESCE(sum(i.existing_qty),0) AS existing_qty,COALESCE(sum(i.actual_qty),0) AS acquired_qty,COALESCE(sum(GREATEST(i.to_buy_qty-COALESCE(i.actual_qty,0),0)),0) AS remaining_qty FROM requests r LEFT JOIN requisition_item i USING(requisition_id) GROUP BY r.status)
 SELECT c.tor_id,c.contract_no,c.project_name,c.status AS contract_status,c.start_date,c.end_date,a.assignment_id,a.required_workers,l.location_name,l.address,
 $2::text AS period_start,$3::text AS period_end,
 (c.status='active' AND $4::date BETWEEN c.start_date AND c.end_date) AS can_continue,
 CASE WHEN c.status='active' AND $4::date BETWEEN c.start_date AND c.end_date THEN '' ELSE 'กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ' END AS continuation_notice,
 (EXISTS(SELECT 1 FROM shifts) OR EXISTS(SELECT 1 FROM requests)) AS has_data,
 CASE WHEN EXISTS(SELECT 1 FROM shifts) OR EXISTS(SELECT 1 FROM requests) THEN '' ELSE 'ยังไม่มีข้อมูลสรุปสำหรับรอบนี้' END AS summary_notice,
 jsonb_build_object('scheduled_shifts',(SELECT count(*) FROM shifts),'on_time',(SELECT count(*) FROM shifts WHERE status='on_time'),'late',(SELECT count(*) FROM shifts WHERE status='late'),'absent',(SELECT count(*) FROM shifts WHERE status='absent'),'leave',(SELECT count(*) FROM shifts WHERE status='leave'),'awaiting_attendance',(SELECT count(*) FROM shifts WHERE status IS NULL)) AS attendance,
 COALESCE((SELECT jsonb_agg(to_jsonb(p) ORDER BY status) FROM progress p),'[]'::jsonb) AS procurement,
 'current attendance by scheduled work_date; current cumulative procurement for requests created within the Bangkok date range'::text AS basis,
 false AS persisted
 FROM tor_location_assignment a JOIN contract_tor c USING(tor_id) JOIN location l USING(location_id) WHERE a.assignment_id=$1`

func (s *Service) OperationsSummary(ctx context.Context, p auth.Principal, id, start, end string) (json.RawMessage, error) {
	if p.Role != auth.RoleAssistant && p.Role != auth.RoleSupervisor {
		return nil, apperr.Forbidden("staff operational summary only")
	}
	if e := validID(id, "assignment_id"); e != nil {
		return nil, e
	}
	if e := validRange(start, end); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, operationsSummarySQL, id, start, end, s.Now().In(Bangkok).Format("2006-01-02"))
}

type OperationsNotificationInput struct {
	Start              string `json:"period_start"`
	End                string `json:"period_end"`
	SupervisorReviewed bool   `json:"supervisor_reviewed"`
}

// This manual command has no ledger, snapshot or read receipt. Review of 6S
// is an explicit Supervisor attestation because no confirmation record exists.
func (s *Service) NotifyOperationsSummary(ctx context.Context, p auth.Principal, tor string, in OperationsNotificationInput) (NotificationResult, error) {
	out := NotificationResult{Status: "failed"}
	if p.Role != auth.RoleSupervisor {
		return out, apperr.Forbidden("Supervisor must send the summary")
	}
	if e := validID(tor, "tor_id"); e != nil {
		return out, e
	}
	if e := validRange(in.Start, in.End); e != nil {
		return out, e
	}
	if !in.SupervisorReviewed {
		return out, invalid("supervisor_reviewed", "Supervisor must review 5S and 6S before sending")
	}
	end, _ := Date(in.End)
	if s.Now().Before(end.AddDate(0, 0, 1)) {
		return out, conflict("summary period has not closed")
	}
	dashboard, e := url.Parse(s.AssistantDashboardURL)
	if e != nil || dashboard.Scheme != "https" || dashboard.Host == "" || dashboard.User != nil || dashboard.Fragment != "" {
		return out, conflict("ASSISTANT_DASHBOARD_URL must be a configured HTTPS dashboard URL")
	}
	var project string
	if e = s.Repo.Pool.QueryRow(ctx, `SELECT project_name FROM contract_tor WHERE tor_id=$1`, tor).Scan(&project); e != nil {
		return out, e
	}
	var hasData, missingPayroll bool
	e = s.Repo.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule sc JOIN tor_location_assignment a USING(assignment_id) WHERE a.tor_id=$1 AND sc.work_date BETWEEN $2::date AND $3::date AND sc.shift_status<>'cancelled') OR EXISTS(SELECT 1 FROM equipment_requisition r JOIN tor_location_assignment a USING(assignment_id) WHERE a.tor_id=$1 AND (r.created_at AT TIME ZONE 'Asia/Bangkok')::date BETWEEN $2::date AND $3::date),
 EXISTS(SELECT 1 FROM work_schedule sc JOIN tor_location_assignment a USING(assignment_id) JOIN worker w USING(worker_id) WHERE a.tor_id=$1 AND sc.work_date BETWEEN $2::date AND $3::date AND sc.shift_status<>'cancelled' AND NOT EXISTS(SELECT 1 FROM payroll p WHERE p.user_id=w.user_id AND sc.work_date BETWEEN p.period_start AND p.period_end))`, tor, in.Start, in.End).Scan(&hasData, &missingPayroll)
	if e != nil {
		return out, e
	}
	if !hasData {
		return out, conflict("ยังไม่มีข้อมูลสรุปสำหรับรอบนี้")
	}
	if missingPayroll {
		return out, conflict("process 5S payroll for scheduled workers before sending")
	}
	// Recompute 6S through its existing validated SQL, never expose financials.
	if _, e = s.ProfitRange(ctx, tor, in.Start, in.End); e != nil {
		return out, e
	}
	if s.Notify == nil {
		out.Status = "disabled"
		return out, nil
	}
	params := dashboard.Query()
	params.Set("tor_id", tor)
	params.Set("period_start", in.Start)
	params.Set("period_end", in.End)
	dashboard.RawQuery = params.Encode()
	project = strings.Join(strings.Fields(project), " ")
	message := fmt.Sprintf("สรุปผลการดำเนินงานโครงการ %s รอบวันที่ %s ถึง %s พร้อมให้ตรวจสอบแล้ว\n%s", project, in.Start, in.End, dashboard.String())
	rows, e := s.Repo.Pool.Query(ctx, `SELECT user_id::text FROM public."USER" WHERE role='assistant' AND is_active ORDER BY user_id`)
	if e != nil {
		return out, e
	}
	var users []string
	for rows.Next() {
		var user string
		if e = rows.Scan(&user); e != nil {
			rows.Close()
			return out, e
		}
		users = append(users, user)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return out, e
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for _, user := range users {
		if s.notifyKey(ctx, user, message, "operations-summary:"+tor+":"+in.Start+":"+in.End) {
			out.Accepted++
		} else {
			out.Failed++
		}
	}
	if out.Accepted > 0 && out.Failed == 0 {
		out.Status = "accepted"
	}
	if out.Accepted > 0 && out.Failed > 0 {
		out.Status = "partial"
	}
	return out, nil
}
