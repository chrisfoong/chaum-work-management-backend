package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"sort"
	"strings"
)

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
