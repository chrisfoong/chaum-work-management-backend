package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
)

// 5W: worker equipment requests and ownership-filtered reads.
type RequestItem struct {
	EquipmentID string `json:"equipment_id"`
	Required    int    `json:"required_qty"`
	Remark      string `json:"remark"`
}
type RequestInput struct {
	AssignmentID string        `json:"assignment_id"`
	Reason       string        `json:"reason"`
	Items        []RequestItem `json:"items"`
}

func (s *Service) Requisition(ctx context.Context, p auth.Principal, in RequestInput) (string, error) {
	if e := validID(in.AssignmentID, "assignment_id"); e != nil {
		return "", e
	}
	if e := text(in.Reason, "reason", 500); e != nil {
		return "", e
	}
	if len(in.Items) == 0 || len(in.Items) > 100 {
		return "", invalid("items", "1-100 items required")
	}
	seen := map[string]bool{}
	for _, it := range in.Items {
		if e := validID(it.EquipmentID, "equipment_id"); e != nil {
			return "", e
		}
		if it.Required <= 0 || it.Required > 1000000 || seen[it.EquipmentID] {
			return "", invalid("items", "positive quantity and distinct equipment required")
		}
		seen[it.EquipmentID] = true
	}
	var id string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if p.Role == auth.RoleWorker {
			worker, e := s.Worker(ctx, q, p.UserID.String())
			if e != nil {
				return e
			}
			var own bool
			if e = s.Repo.Row(ctx, q, query65, worker, in.AssignmentID, s.Now()).Scan(&own); e != nil {
				return e
			}
			if !own {
				return conflict("assignment is not assigned to worker")
			}
		}
		if e := s.Repo.Row(ctx, q, query66, document("REQ"), p.UserID, in.AssignmentID, in.Reason).Scan(&id); e != nil {
			return e
		}
		for _, it := range in.Items {
			tag, e := s.Repo.Exec(ctx, q, query67, id, it.EquipmentID, it.Required, it.Remark)
			if e != nil {
				return e
			}
			if tag.RowsAffected() != 1 {
				return invalid("equipment_id", "active equipment required")
			}
		}
		return nil
	})
	if e == nil {
		s.Event("request", id)
	}
	return id, e
}
func (s *Service) Requisitions(ctx context.Context, p auth.Principal, limit, offset int) (json.RawMessage, error) {
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return s.Repo.List(ctx, s.Repo.Pool, query68, owner, limit, offset)
}
func (s *Service) RequisitionDetail(ctx context.Context, p auth.Principal, id string) (json.RawMessage, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return nil, e
	}
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return one(ctx, s.Repo.Pool, query69, id, owner)
}
