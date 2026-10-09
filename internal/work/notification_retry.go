package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"time"
)

type NotificationInput struct {
	Kind   string `json:"kind"`
	Reason string `json:"reason,omitempty"`
}
type NotificationResult struct {
	Status           string `json:"status"`
	Accepted         int    `json:"accepted_recipients"`
	Failed           int    `json:"failed_recipients"`
	DeliveryVerified bool   `json:"delivery_verified"`
	Persisted        bool   `json:"persisted"`
}

// RetryNotification reads committed business rows only. It never replays a
// purchase, leave review or delivery write. There is no durable outbox.
func (s *Service) RetryNotification(ctx context.Context, p auth.Principal, id string, in NotificationInput) (NotificationResult, error) {
	out := NotificationResult{Status: "failed"}
	if e := validID(id, "id"); e != nil {
		return out, e
	}
	var guard, recipients string
	message := "Chaum: " + in.Kind + " updated. Please open the app to view details."
	switch in.Kind {
	case "no_purchase":
		if p.Role != auth.RoleAssistant {
			return out, conflict("Assistant notification required")
		}
		if e := text(in.Reason, "reason", 1000); e != nil {
			return out, e
		}
		guard = `SELECT EXISTS(SELECT 1 FROM equipment_requisition WHERE requisition_id=$1 AND requisition_type='additional' AND status='rejected' AND reviewed_by=$2)`
		recipients = query26
		message = "Chaum: equipment request declined. " + in.Reason
	case "delivery":
		if p.Role != auth.RoleAssistant {
			return out, conflict("Assistant notification required")
		}
		guard = `SELECT EXISTS(SELECT 1 FROM equipment_requisition r WHERE requisition_id=$1 AND requisition_type='additional' AND status='completed' AND EXISTS(SELECT 1 FROM work_evidence WHERE left(description,length('[Chaum delivery ' || r.requisition_id::text || '] '))='[Chaum delivery ' || r.requisition_id::text || '] ') AND $2::uuid IS NOT NULL)`
		recipients = query26
	case "approval_needed", "purchase_funding":
		if p.Role != auth.RoleAssistant {
			return out, conflict("Assistant notification required")
		}
		state := "pending_approval"
		if in.Kind == "purchase_funding" {
			state = "pending_fund"
		}
		guard = `SELECT EXISTS(SELECT 1 FROM equipment_requisition WHERE requisition_id=$1 AND status='` + state + `' AND $2::uuid IS NOT NULL)`
		recipients = `SELECT user_id::text FROM public."USER" WHERE role='supervisor' AND is_active AND $1::uuid IS NOT NULL`
	case "leave_review", "replacement_needed":
		if p.Role != auth.RoleAssistant {
			return out, conflict("Assistant notification required")
		}
		guard = `SELECT EXISTS(SELECT 1 FROM leave_request WHERE request_id=$1 AND status IN ('approved','rejected') AND $2::uuid IS NOT NULL)`
		recipients = query25
		if in.Kind == "replacement_needed" {
			guard = `SELECT EXISTS(SELECT 1 FROM leave_request WHERE request_id=$1 AND status='approved' AND $2::uuid IS NOT NULL)`
			recipients = `SELECT user_id::text FROM public."USER" WHERE role='supervisor' AND is_active AND $1::uuid IS NOT NULL`
		}
	default:
		return out, invalid("kind", "supported kinds: no_purchase, delivery, approval_needed, purchase_funding, leave_review, replacement_needed")
	}
	var allowed bool
	if e := s.Repo.Pool.QueryRow(ctx, guard, id, p.UserID).Scan(&allowed); e != nil {
		return out, e
	}
	if !allowed {
		return out, conflict("notification does not match committed business state")
	}
	if in.Kind == "delivery" {
		if e := s.Repo.Pool.QueryRow(ctx, `SELECT 'Chaum: ส่งมอบอุปกรณ์คำขอ ' || r.requisition_no || E'\n' || COALESCE(string_agg(e.equipment_name || ' × ' || COALESCE(i.actual_qty,0)::text,E'\n' ORDER BY i.item_id),'') FROM equipment_requisition r JOIN requisition_item i USING(requisition_id) JOIN equipment e USING(equipment_id) WHERE r.requisition_id=$1 GROUP BY r.requisition_no`, id).Scan(&message); e != nil {
			return out, e
		}
	}
	rows, e := s.Repo.Pool.Query(ctx, recipients, id)
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
		if s.notifyKey(ctx, user, message, in.Kind+":"+id) {
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
	if s.Notify == nil {
		out.Status = "disabled"
	}
	return out, nil
}
