package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
)

// 3S: Supervisor approval or rejection of additional procurement.
type ReviewInput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

func (s *Service) ReviewRequest(ctx context.Context, p auth.Principal, id string, in ReviewInput) error {
	if e := validID(id, "requisition_id"); e != nil {
		return e
	}
	if in.Decision != "approve" && in.Decision != "reject" {
		return invalid("decision", "approve or reject required")
	}
	if in.Decision == "reject" {
		if e := text(in.Reason, "reason", 1000); e != nil {
			return e
		}
	}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		var state string
		if e := s.Repo.Row(ctx, q, query75, id).Scan(&state); e != nil {
			return e
		}
		if state != "pending_approval" {
			return conflict("request is not pending approval")
		}
		next := "pending_procurement"
		if in.Decision == "reject" {
			next = "rejected"
		}
		_, e := s.Repo.Exec(ctx, q, query76, id, next, p.UserID, s.Now(), in.Reason)
		return e
	})
	if err == nil {
		s.Event("request_review", id)
	}
	return err
}
