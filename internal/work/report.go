package work

import (
	"context"
	"encoding/json"
)

func (s *Service) Profit(ctx context.Context, month string) (json.RawMessage, error) {
	if e := Month(month); e != nil {
		return nil, e
	}
	var duplicates bool
	if e := s.Repo.Row(ctx, s.Repo.Pool, `SELECT EXISTS(SELECT 1 FROM worker GROUP BY user_id HAVING count(*)>1)`).Scan(&duplicates); e != nil {
		return nil, e
	}
	if duplicates {
		return nil, conflict("duplicate Worker/User mappings must be resolved before reporting")
	}
	return s.Repo.List(ctx, s.Repo.Pool, query87, month)
}
func (s *Service) Invoices(ctx context.Context, limit, offset int) (json.RawMessage, error) {
	return s.Repo.List(ctx, s.Repo.Pool, query88, limit, offset)
}
