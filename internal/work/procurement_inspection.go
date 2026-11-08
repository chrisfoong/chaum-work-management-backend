package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
)

// 5A / 6A: inspect additional requests and decide whether procurement is needed.
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
