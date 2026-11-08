package work

import (
	"context"
)

// 1A: survey initial TOR equipment without writing generated shortage quantities.
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
		if state != "pending_survey" || kind != "tor_base" {
			return conflict("only pending TOR base requests may be surveyed; additional requests use inspection and decision")
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
