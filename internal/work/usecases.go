package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// 4A: both the review and optional replacement commit together.
type LeaveReviewInput struct {
	Status   string `json:"status"`
	WorkerID string `json:"replacement_worker_id"`
}

func (s *Service) ReviewLeaveAndReplace(ctx context.Context, id string, in LeaveReviewInput) ([]string, error) {
	if e := validID(id, "request_id"); e != nil {
		return nil, e
	}
	if in.Status != "approved" && in.Status != "rejected" {
		return nil, invalid("status", "approved or rejected required")
	}
	if in.WorkerID != "" {
		if e := validID(in.WorkerID, "replacement_worker_id"); e != nil {
			return nil, e
		}
		if in.Status != "approved" {
			return nil, invalid("replacement_worker_id", "only permitted with approval")
		}
	}
	var ids []string
	err := s.Repo.Transaction(ctx, func(q Query) error {
		if e := s.reviewLeave(ctx, q, id, in.Status); e != nil {
			return e
		}
		if in.WorkerID == "" {
			return nil
		}
		var state, original string
		var plan ScheduleInput
		if e := q.QueryRow(ctx, query64, id).Scan(&state, &plan.AssignmentID, &plan.WorkDate, &plan.Start, &original); e != nil {
			return e
		}
		if original == in.WorkerID {
			return conflict("replacement must be a different worker")
		}
		plan.WorkerIDs = []string{in.WorkerID}
		var e error
		ids, e = s.schedule(ctx, q, plan)
		return e
	})
	if err == nil {
		s.Event("leave_review", id)
		for _, schedule := range ids {
			s.Event("schedule", schedule)
		}
		if in.Status == "approved" && len(ids) == 0 {
			s.Event("replacement_needed", id)
		}
	}
	return ids, err
}

// 5A: use quantities as planning context, not proof of equipment condition.
func (s *Service) InspectRequest(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT r.*,c.project_name,c.contract_no,l.location_name,
 COALESCE((SELECT jsonb_agg(jsonb_build_object('item_id',i.item_id,'equipment_id',i.equipment_id,'equipment_name',e.equipment_name,'required_qty',i.required_qty,'actual_qty',i.actual_qty,'remaining_qty',GREATEST(i.to_buy_qty-COALESCE(i.actual_qty,0),0))) FROM requisition_item i JOIN equipment e USING(equipment_id) WHERE i.requisition_id=r.requisition_id),'[]'::jsonb) AS requested_items,
 COALESCE((SELECT jsonb_agg(jsonb_build_object('requisition_id',b.requisition_id,'assignment_id',b.assignment_id,'status',b.status,'equipment_id',i.equipment_id,'required_qty',i.required_qty,'existing_qty',i.existing_qty,'actual_qty',i.actual_qty,'remaining_qty',GREATEST(i.to_buy_qty-COALESCE(i.actual_qty,0),0))) FROM equipment_requisition b JOIN tor_location_assignment ba ON ba.assignment_id=b.assignment_id JOIN requisition_item i ON i.requisition_id=b.requisition_id WHERE ba.tor_id=a.tor_id AND b.requisition_type='tor_base' AND i.equipment_id IN (SELECT equipment_id FROM requisition_item WHERE requisition_id=r.requisition_id)),'[]'::jsonb) AS tor_requirements,
 COALESCE((SELECT jsonb_agg(jsonb_build_object('requisition_id',other.requisition_id,'assignment_id',other.assignment_id,'status',other.status,'equipment_id',i.equipment_id,'remaining_qty',GREATEST(i.to_buy_qty-COALESCE(i.actual_qty,0),0))) FROM equipment_requisition other JOIN tor_location_assignment oa ON oa.assignment_id=other.assignment_id JOIN requisition_item i ON i.requisition_id=other.requisition_id WHERE oa.tor_id=a.tor_id AND other.requisition_id<>r.requisition_id AND other.requisition_type='additional' AND other.status NOT IN ('completed','rejected') AND i.equipment_id IN (SELECT equipment_id FROM requisition_item WHERE requisition_id=r.requisition_id)),'[]'::jsonb) AS other_pending_requests
 FROM equipment_requisition r JOIN tor_location_assignment a ON a.assignment_id=r.assignment_id JOIN contract_tor c ON c.tor_id=a.tor_id JOIN location l ON l.location_id=a.location_id WHERE r.requisition_id=$1 AND r.requisition_type='additional' AND r.status='pending_survey'`, id)
}

type ProcurementDecision struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

func (s *Service) DecideRequest(ctx context.Context, p auth.Principal, id string, in ProcurementDecision) error {
	if e := validID(id, "requisition_id"); e != nil {
		return e
	}
	if in.Decision != "purchase" && in.Decision != "no_purchase" {
		return invalid("decision", "purchase or no_purchase required")
	}
	if e := text(in.Reason, "reason", 1000); e != nil {
		return e
	}
	var owner string
	err := s.Repo.Transaction(ctx, func(q Query) error {
		var state, kind string
		if e := q.QueryRow(ctx, `SELECT status::text,requisition_type::text,requested_by::text FROM equipment_requisition WHERE requisition_id=$1 FOR UPDATE`, id).Scan(&state, &kind, &owner); e != nil {
			return e
		}
		if state != "pending_survey" || kind != "additional" {
			return conflict("only pending additional requests can be decided")
		}
		var valid bool
		if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1) AND NOT EXISTS(SELECT 1 FROM requisition_item i JOIN equipment e USING(equipment_id) WHERE i.requisition_id=$1 AND (i.required_qty<=0 OR ($2 AND NOT e.is_active)))`, id, in.Decision == "purchase").Scan(&valid); e != nil {
			return e
		}
		if !valid {
			return invalid("items", "positive quantities and active equipment required")
		}
		next := "pending_approval"
		if in.Decision == "no_purchase" {
			next = "rejected"
		}
		_, e := q.Exec(ctx, `UPDATE equipment_requisition SET status=$2::text::requisition_status_enum,reason=CASE WHEN $2::text='pending_approval' THEN reason || E'\n[Purchase decision] ' || $3 ELSE reason END,reviewed_by=CASE WHEN $2::text='rejected' THEN $4::uuid ELSE reviewed_by END,reviewed_at=CASE WHEN $2::text='rejected' THEN $5::timestamptz ELSE reviewed_at END WHERE requisition_id=$1`, id, next, in.Reason, p.UserID, s.Now())
		return e
	})
	return err
}

// 8A: no requisition FK exists on work_evidence. The explicit description
// marker is a compatibility convention; it does not claim a new DB relation.
type DeliveryInput struct {
	ScheduleID  string   `json:"schedule_id"`
	Description string   `json:"description"`
	Photos      []string `json:"photo_paths"`
}

func deliveryPrefix(id string) string { return "[Chaum delivery " + id + "] " }
func (s *Service) Deliver(ctx context.Context, p auth.Principal, id string, in DeliveryInput) error {
	if e := validID(id, "requisition_id"); e != nil {
		return e
	}
	if e := validID(in.ScheduleID, "schedule_id"); e != nil {
		return e
	}
	if e := text(in.Description, "description", 2000); e != nil {
		return e
	}
	if len(in.Photos) == 0 || len(in.Photos) > 20 {
		return invalid("photo_paths", "1-20 photos required")
	}
	photos := append([]string(nil), in.Photos...)
	sort.Strings(photos)
	for i, path := range photos {
		if i > 0 && path == photos[i-1] {
			return invalid("photo_paths", "duplicate photo")
		}
		if e := s.evidencePhoto(ctx, p, path); e != nil {
			return e
		}
	}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		var state, kind, assignment, owner, number string
		if e := q.QueryRow(ctx, `SELECT status::text,requisition_type::text,assignment_id::text,requested_by::text,requisition_no FROM equipment_requisition WHERE requisition_id=$1 FOR UPDATE`, id).Scan(&state, &kind, &assignment, &owner, &number); e != nil {
			return e
		}
		if state != "completed" || kind != "additional" {
			return conflict("only fully purchased additional requests can be delivered")
		}
		worker, e := s.Worker(ctx, q, owner)
		if e != nil {
			return e
		}
		var valid bool
		if e = q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule WHERE schedule_id=$1 AND worker_id=$2 AND assignment_id=$3 AND shift_status<>'cancelled' AND work_date<=$4::date FOR SHARE)`, in.ScheduleID, worker, assignment, s.Now().In(Bangkok).Format("2006-01-02")).Scan(&valid); e != nil {
			return e
		}
		if !valid {
			return conflict("recipient schedule must match request owner and assignment")
		}
		var complete bool
		if e = q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1) AND NOT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1 AND COALESCE(actual_qty,0)<to_buy_qty)`, id).Scan(&complete); e != nil {
			return e
		}
		if !complete {
			return conflict("delivery requires all requested equipment to be procured")
		}
		var items string
		if e = q.QueryRow(ctx, `SELECT COALESCE(string_agg(e.equipment_name || ' x' || COALESCE(i.actual_qty,0)::text,', ' ORDER BY i.item_id),'') FROM requisition_item i JOIN equipment e USING(equipment_id) WHERE i.requisition_id=$1`, id).Scan(&items); e != nil {
			return e
		}
		description := deliveryPrefix(id) + "request=" + number + "; recipient=" + owner + "; items=" + items + "; " + in.Description
		if len([]rune(description)) > 2000 {
			return invalid("description", "delivery summary plus description exceeds 2000 characters")
		}
		rows, e := q.Query(ctx, `SELECT schedule_id::text,worker_id::text,description,photo_url FROM work_evidence WHERE left(description,length($1))=$1 ORDER BY photo_url`, deliveryPrefix(id))
		if e != nil {
			return e
		}
		var existing []string
		different := false
		for rows.Next() {
			var schedule, recipient, desc, photo string
			if e = rows.Scan(&schedule, &recipient, &desc, &photo); e != nil {
				rows.Close()
				return e
			}
			different = different || schedule != in.ScheduleID || recipient != worker || desc != description
			existing = append(existing, photo)
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			return e
		}
		if len(existing) > 0 {
			if different || strings.Join(existing, "\n") != strings.Join(photos, "\n") {
				return conflict("delivery already recorded with different payload")
			}
			return nil
		}
		for _, photo := range photos {
			if _, e = q.Exec(ctx, query7, in.ScheduleID, worker, description, photo); e != nil {
				return e
			}
		}
		return nil
	})
	return err
}
func (s *Service) Deliveries(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return nil, e
	}
	return s.Repo.List(ctx, s.Repo.Pool, `SELECT evidence_id,schedule_id,worker_id,description,photo_url,submitted_at FROM work_evidence WHERE left(description,length($1))=$1 ORDER BY submitted_at,evidence_id`, deliveryPrefix(id))
}
func (s *Service) evidencePhoto(ctx context.Context, p auth.Principal, path string) error {
	// Staff receipts may be PDF, delivery evidence must be JPEG/PNG.
	photoActor := p
	photoActor.Role = auth.RoleWorker
	return s.evidence(ctx, photoActor, path)
}

func (s *Service) Continuation(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "tor_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT c.tor_id,c.contract_no,c.project_name,c.start_date,c.end_date,c.status,
 COALESCE((SELECT jsonb_agg(jsonb_build_object('assignment_id',a.assignment_id,'location_id',a.location_id,'location_name',l.location_name,'required_workers',a.required_workers)) FROM tor_location_assignment a JOIN location l USING(location_id) WHERE a.tor_id=c.tor_id),'[]'::jsonb) AS assignments
 FROM contract_tor c WHERE c.tor_id=$1 AND c.status='active' AND c.end_date >= $2::date`, id, s.Now().In(Bangkok).Format("2006-01-02"))
}

type PayrollBatchInput struct {
	Start string `json:"period_start"`
	End   string `json:"period_end"`
}
type PayrollBatchResult struct {
	UserID    string `json:"user_id"`
	PayrollID string `json:"payroll_id"`
	Existing  bool   `json:"existing"`
}

func (s *Service) PayrollBatch(ctx context.Context, p auth.Principal, in PayrollBatchInput) ([]PayrollBatchResult, error) {
	if e := Period(in.Start, in.End); e != nil {
		return nil, e
	}
	end, _ := Date(in.End)
	if s.Now().Before(end.AddDate(0, 0, 1)) {
		return nil, conflict("payroll period has not closed")
	}
	out := []PayrollBatchResult{}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		rows, e := q.Query(ctx, `SELECT DISTINCT w.user_id::text FROM worker w JOIN work_schedule sc ON sc.worker_id=w.worker_id WHERE sc.work_date BETWEEN $1::date AND $2::date AND sc.shift_status<>'cancelled' ORDER BY w.user_id::text`, in.Start, in.End)
		if e != nil {
			return e
		}
		var users []string
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
			var existing string
			e = q.QueryRow(ctx, `SELECT payroll_id::text FROM payroll WHERE user_id=$1 AND period_start=$2::date AND period_end=$3::date`, user, in.Start, in.End).Scan(&existing)
			if e == nil {
				out = append(out, PayrollBatchResult{user, existing, true})
				continue
			}
			if !isMissing(e) {
				return e
			}
			id, e := s.createPayroll(ctx, q, p, PayrollInput{user, in.Start, in.End})
			if e != nil {
				return e
			}
			out = append(out, PayrollBatchResult{user, id, false})
		}
		return nil
	})
	if err == nil {
		for _, r := range out {
			if !r.Existing {
				s.Event("payroll", r.PayrollID)
			}
		}
	}
	return out, err
}
func (s *Service) PayrollMonth(ctx context.Context, p auth.Principal, month string, limit, offset int) (json.RawMessage, error) {
	if e := Month(month); e != nil {
		return nil, e
	}
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	sql := strings.Replace(query35, "ORDER BY p.period_start", "AND to_char(p.period_start,'YYYY-MM')=$4 ORDER BY p.period_start", 1)
	return s.Repo.List(ctx, s.Repo.Pool, sql, owner, limit, offset, month)
}

func validRange(start, end string) error {
	a, e := Date(start)
	if e != nil {
		return e
	}
	b, e := Date(end)
	if e != nil {
		return e
	}
	if b.Before(a) || b.Sub(a) > 3660*24*time.Hour {
		return invalid("period", "ordered date range of at most ten years required")
	}
	return nil
}
func (s *Service) ProfitRange(ctx context.Context, tor, start, end string) (json.RawMessage, error) {
	if e := validRange(start, end); e != nil {
		return nil, e
	}
	var selected any
	if tor != "" {
		if e := validID(tor, "tor_id"); e != nil {
			return nil, e
		}
		selected = tor
	}
	var duplicates bool
	if e := s.Repo.Row(ctx, s.Repo.Pool, `SELECT EXISTS(SELECT 1 FROM worker GROUP BY user_id HAVING count(*)>1)`).Scan(&duplicates); e != nil {
		return nil, e
	}
	if duplicates {
		return nil, conflict("duplicate Worker/User mappings must be resolved before reporting")
	}
	sql := strings.ReplaceAll(query87, "p.is_paid AND to_char(p.period_start,'YYYY-MM')=$1", "p.is_paid AND p.period_start >= $1::date AND p.period_end <= $2::date")
	sql = strings.ReplaceAll(sql, "billing_month=$1", "billing_month BETWEEN to_char($1::date,'YYYY-MM') AND to_char($2::date,'YYYY-MM')")
	sql = strings.ReplaceAll(sql, "to_char(e.created_at AT TIME ZONE 'Asia/Bangkok','YYYY-MM')=$1", "(e.created_at AT TIME ZONE 'Asia/Bangkok')::date BETWEEN $1::date AND $2::date")
	sql = strings.ReplaceAll(sql, "$1::text AS month", "$1::text AS period_start,$2::text AS period_end")
	sql = strings.ReplaceAll(sql, "ORDER BY c.tor_id", "WHERE ($3::uuid IS NULL OR c.tor_id=$3) ORDER BY c.tor_id")
	return s.Repo.List(ctx, s.Repo.Pool, sql, start, end, selected)
}

func (s *Service) CloseSummary(ctx context.Context, tor, start, end string) (json.RawMessage, error) {
	if e := validID(tor, "tor_id"); e != nil {
		return nil, e
	}
	data, e := s.ProfitRange(ctx, tor, start, end)
	if e != nil {
		return nil, e
	}
	var rows []json.RawMessage
	if e = json.Unmarshal(data, &rows); e != nil {
		return nil, e
	}
	if len(rows) == 0 {
		return nil, conflict("TOR does not exist")
	}
	// 6S confirmation is a snapshot, not contract termination. No close ledger
	// exists; never encode closure by changing contract status.
	return json.Marshal(struct {
		Data        json.RawMessage `json:"data"`
		ConfirmedAt string          `json:"confirmed_at"`
		Persisted   bool            `json:"persisted"`
	}{data, s.Now().Format(time.RFC3339), false})
}

// AssignmentContinuation is a read model for the area selected in 9A.
// Notification history cannot be reconstructed from the unchanged schema.
func (s *Service) AssignmentContinuation(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "assignment_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT c.tor_id,c.contract_no,c.project_name,c.start_date,c.end_date,c.status,a.assignment_id,a.required_workers,l.location_id,l.location_name,l.address,
 (c.status='active' AND $2::date BETWEEN c.start_date AND c.end_date) AS can_continue,
 CASE WHEN c.status='active' AND $2::date BETWEEN c.start_date AND c.end_date THEN '' ELSE 'กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ' END AS notice,
 false AS summary_notification_history_available
 FROM tor_location_assignment a JOIN contract_tor c USING(tor_id) JOIN location l USING(location_id) WHERE a.assignment_id=$1`, id, s.Now().In(Bangkok).Format("2006-01-02"))
}
