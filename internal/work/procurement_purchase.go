package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
)

// 2A / 7A: preview and commit purchase rounds using the same transaction path.
type PurchaseItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"actual_qty"`
	Price    string `json:"actual_price"`
	Expected *int   `json:"expected_actual_qty"`
}
type PurchaseInput struct {
	Items   []PurchaseItem `json:"items"`
	Receipt string         `json:"receipt_path"`
}

type PurchaseReviewItem struct {
	ItemID         string `json:"item_id"`
	ActualQuantity int    `json:"new_actual_qty"`
	Remaining      int    `json:"remaining_qty"`
	UnitPrice      string `json:"new_actual_price"`
	Cost           string `json:"round_cost"`
}
type PurchaseReview struct {
	Items      []PurchaseReviewItem `json:"items"`
	Total      string               `json:"round_total"`
	NextStatus string               `json:"next_status"`
	Persisted  bool                 `json:"persisted"`
}

func (s *Service) PreviewPurchase(ctx context.Context, p auth.Principal, id string, in PurchaseInput) (PurchaseReview, error) {
	out := PurchaseReview{Items: []PurchaseReviewItem{}}
	_, e := s.executePurchase(ctx, p, id, in, &out)
	return out, e
}
func (s *Service) Purchase(ctx context.Context, p auth.Principal, id string, in PurchaseInput) (string, error) {
	return s.executePurchase(ctx, p, id, in, nil)
}
func (s *Service) executePurchase(ctx context.Context, p auth.Principal, id string, in PurchaseInput, review *PurchaseReview) (string, error) {
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
		if e != nil || price < 0 || it.Quantity < 0 || it.Quantity > 1000000 || price > 9999999999 {
			return "", invalid("items", "nonnegative quantity and price required")
		}
		if _, ok := prices[it.ItemID]; ok {
			return "", invalid("items", "duplicate item")
		}
		prices[it.ItemID] = price
		total += price * int64(it.Quantity)
	}
	positive := false
	for _, it := range in.Items {
		if it.Quantity > 0 {
			positive = true
		}
		if it.Expected == nil || *it.Expected < 0 {
			return "", invalid("expected_actual_qty", "current cumulative quantity required to prevent replay")
		}
	}
	if total > 9999999999 {
		return "", invalid("amount", "total exceeds NUMERIC(10,2)")
	}
	if total > 0 {
		if e := s.evidencePhoto(ctx, p, in.Receipt); e != nil {
			return "", e
		}
	}
	var expense, nextState string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		var state, kind string
		var approved bool
		if e := s.Repo.Row(ctx, q, query81, id).Scan(&state, &kind, &approved); e != nil {
			return e
		}
		if state != "pending_procurement" || (kind == "additional" && !approved) {
			return conflict("request is not ready for procurement")
		}
		if kind == "additional" {
			if !positive {
				return invalid("items", "additional purchases require at least one positive quantity")
			}
		}
		if kind == "additional" && total > 0 {
			if e := s.evidencePhoto(ctx, p, in.Receipt); e != nil {
				return e
			}
		}
		for _, it := range in.Items {
			var buy int
			var actual int
			var oldPrice string
			if e := s.Repo.Row(ctx, q, query82, it.ItemID, id).Scan(&buy, &actual, &oldPrice); e != nil {
				return e
			}
			if actual != *it.Expected {
				return conflict("purchase quantity changed; reload before confirming")
			}
			if it.Quantity > buy-actual {
				return invalid("actual_qty", "quantity exceeds remaining requirement")
			}
			if review != nil {
				price, e := Money(oldPrice)
				if e != nil {
					return conflict("invalid stored unit price")
				}
				newQty := actual + it.Quantity
				if it.Quantity > 0 {
					price = prices[it.ItemID]
					if kind == "additional" {
						previous, _ := Money(oldPrice)
						numerator := previous*int64(actual) + price*int64(it.Quantity)
						price = (numerator + int64(newQty)/2) / int64(newQty)
					}
				}
				review.Items = append(review.Items, PurchaseReviewItem{it.ItemID, newQty, buy - newQty, Decimal(price), Decimal(prices[it.ItemID] * int64(it.Quantity))})
				continue
			}
			if it.Quantity == 0 {
				continue
			}
			if _, e := s.Repo.Exec(ctx, q, query83, it.ItemID, it.Quantity, Decimal(prices[it.ItemID]), kind == "additional"); e != nil {
				return e
			}
		}
		if total > 0 {
			if e := lock(ctx, q, "purchase-receipt:"+in.Receipt); e != nil {
				return e
			}
			var used bool
			if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM expense_claim WHERE expense_type='actual_expense' AND receipt_photo_url=$1)`, in.Receipt).Scan(&used); e != nil {
				return e
			}
			if used {
				return conflict("receipt already used for a purchase")
			}
			if review == nil {
				if e := s.Repo.Row(ctx, q, query84, document("EXP"), p.UserID, id, Decimal(total), in.Receipt, s.Now()).Scan(&expense); e != nil {
					return e
				}
			}
		}
		if review != nil {
			quantities := map[string]int{}
			for _, it := range in.Items {
				quantities[it.ItemID] = it.Quantity
			}
			rows, e := q.Query(ctx, `SELECT item_id::text,to_buy_qty,COALESCE(actual_qty,0) FROM requisition_item WHERE requisition_id=$1`, id)
			if e != nil {
				return e
			}
			missing := false
			for rows.Next() {
				var item string
				var buy, actual int
				if e = rows.Scan(&item, &buy, &actual); e != nil {
					rows.Close()
					return e
				}
				missing = missing || actual+quantities[item] < buy
			}
			e = rows.Err()
			rows.Close()
			if e != nil {
				return e
			}
			review.Total = Decimal(total)
			review.NextStatus = "completed"
			if missing {
				review.NextStatus = "pending_fund"
			}
			return nil
		}
		var missing bool
		if e := s.Repo.Row(ctx, q, query85, id).Scan(&missing); e != nil {
			return e
		}
		next := "completed"
		if missing {
			next = "pending_fund"
		}
		nextState = next
		_, e := s.Repo.Exec(ctx, q, query86, id, next)
		return e
	})
	if e == nil {
		if nextState == "pending_fund" {
			s.Event("purchase_funding", id)
		}
	}
	return expense, e
}

// 5A: use quantities as planning context, not proof of equipment condition.
