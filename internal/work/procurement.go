package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
)

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
	if e := text(in.Reason, "reason", 1000); e != nil {
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
			if e = s.Repo.Row(ctx, q, query65, worker, in.AssignmentID).Scan(&own); e != nil {
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

type SurveyItem struct {
	ItemID   string `json:"item_id"`
	Existing int    `json:"existing_qty"`
}

func (s *Service) Survey(ctx context.Context, id string, items []SurveyItem) error {
	if e := validID(id, "requisition_id"); e != nil {
		return e
	}
	if len(items) == 0 || len(items) > 100 {
		return invalid("items", "1-100 items required")
	}
	seen := map[string]bool{}
	for _, it := range items {
		if e := validID(it.ItemID, "item_id"); e != nil {
			return e
		}
		if it.Existing < 0 || it.Existing > 1000000 || seen[it.ItemID] {
			return invalid("items", "nonnegative quantity and distinct items required")
		}
		seen[it.ItemID] = true
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		var state, kind string
		if e := s.Repo.Row(ctx, q, query70, id).Scan(&state, &kind); e != nil {
			return e
		}
		if state != "pending_survey" {
			return conflict("survey already submitted")
		}
		var count int
		if e := s.Repo.Row(ctx, q, query71, id).Scan(&count); e != nil {
			return e
		}
		if count != len(items) {
			return invalid("items", "all items must be surveyed")
		}
		for _, it := range items {
			tag, e := s.Repo.Exec(ctx, q, query72, id, it.ItemID, it.Existing)
			if e != nil {
				return e
			}
			if tag.RowsAffected() != 1 {
				return invalid("item_id", "item does not belong to requisition")
			}
		}
		var missing bool
		if e := s.Repo.Row(ctx, q, query73, id).Scan(&missing); e != nil {
			return e
		}
		next := "completed"
		if missing {
			next = "pending_procurement"
			if kind == "additional" {
				next = "pending_approval"
			}
		}
		_, e := s.Repo.Exec(ctx, q, query74, id, next)
		return e
	})
}

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
	return s.Repo.Transaction(ctx, func(q Query) error {
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
}

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
	return expense, e
}

type PurchaseItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"actual_qty"`
	Price    string `json:"actual_price"`
}
type PurchaseInput struct {
	Items   []PurchaseItem `json:"items"`
	Receipt string         `json:"receipt_path"`
}

func (s *Service) Purchase(ctx context.Context, p auth.Principal, id string, in PurchaseInput) (string, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return "", e
	}
	if len(in.Items) == 0 || len(in.Items) > 100 {
		return "", invalid("items", "1-100 items required")
	}
	prices := map[string]int64{}
	total := int64(0)
	for _, it := range in.Items {
		if e := validID(it.ItemID, "item_id"); e != nil {
			return "", e
		}
		price, e := Money(it.Price)
		if e != nil || price <= 0 || it.Quantity <= 0 || it.Quantity > 1000000 || price > 9999999999 {
			return "", invalid("items", "positive quantity and price required")
		}
		if _, ok := prices[it.ItemID]; ok {
			return "", invalid("items", "duplicate item")
		}
		prices[it.ItemID] = price
		total += price * int64(it.Quantity)
	}
	if total > 9999999999 {
		return "", invalid("amount", "total exceeds NUMERIC(10,2)")
	}
	if e := s.evidence(ctx, p, in.Receipt); e != nil {
		return "", e
	}
	var expense string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		var state string
		if e := s.Repo.Row(ctx, q, query81, id).Scan(&state); e != nil {
			return e
		}
		if state != "pending_procurement" {
			return conflict("request is not ready for procurement")
		}
		for _, it := range in.Items {
			var buy int
			var actual *int
			if e := s.Repo.Row(ctx, q, query82, it.ItemID, id).Scan(&buy, &actual); e != nil {
				return e
			}
			if actual != nil {
				return conflict("item already purchased; retries return conflict, not exactly-once success")
			}
			if it.Quantity != buy {
				return invalid("actual_qty", "single purchase must cover the surveyed quantity")
			}
			if _, e := s.Repo.Exec(ctx, q, query83, it.ItemID, it.Quantity, Decimal(prices[it.ItemID])); e != nil {
				return e
			}
		}
		if e := s.Repo.Row(ctx, q, query84, document("EXP"), p.UserID, id, Decimal(total), in.Receipt).Scan(&expense); e != nil {
			return e
		}
		var missing bool
		if e := s.Repo.Row(ctx, q, query85, id).Scan(&missing); e != nil {
			return e
		}
		next := "completed"
		if missing {
			next = "pending_fund"
		}
		_, e := s.Repo.Exec(ctx, q, query86, id, next)
		return e
	})
	if e == nil {
		s.Event("purchase", id)
	}
	return expense, e
}
