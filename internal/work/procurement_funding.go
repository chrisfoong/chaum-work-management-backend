package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
)

// 2S: record funding evidence and resume procurement.
type FundInput struct {
	Amount    string `json:"amount"`
	Reference string `json:"transfer_ref_no"`
	Receipt   string `json:"receipt_path"`
}

func (s *Service) Fund(ctx context.Context, p auth.Principal, id string, in FundInput) (string, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return "", e
	}
	amount, e := Money(in.Amount)
	if e != nil || amount <= 0 || amount > 9999999999 {
		return "", invalid("amount", "positive amount within NUMERIC(10,2) required")
	}
	if e = text(in.Reference, "transfer_ref_no", 100); e != nil {
		return "", e
	}
	if e = s.evidence(ctx, p, in.Receipt); e != nil {
		return "", e
	}
	var expense string
	e = s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "fund:"+in.Reference); e != nil {
			return e
		}
		var existingRequest, existingAmount, existingReceipt, user string
		e := s.Repo.Row(ctx, q, query77, in.Reference).Scan(&expense, &existingRequest, &existingAmount, &existingReceipt, &user)
		if e == nil {
			if existingRequest != id || existingAmount != Decimal(amount) || existingReceipt != in.Receipt || user != p.UserID.String() {
				return conflict("transfer reference reused with different payload")
			}
			return nil
		}
		if !isMissing(e) {
			return e
		}
		var state string
		if e = s.Repo.Row(ctx, q, query78, id).Scan(&state); e != nil {
			return e
		}
		if state != "pending_fund" && state != "pending_procurement" {
			return conflict("request cannot receive procurement funding")
		}
		if e = s.Repo.Row(ctx, q, query79, document("EXP"), p.UserID, id, Decimal(amount), in.Reference, in.Receipt).Scan(&expense); e != nil {
			return e
		}
		_, e = s.Repo.Exec(ctx, q, query80, id)
		return e
	})
	if e == nil {
		s.Event("funded", id)
	}
	return expense, e
}
