package work

import (
	"context"
)

// 4S: finalize attendance independently of payroll processing.
func (s *Service) Finalize(ctx context.Context) (int64, error) {
	var affected int64
	e := s.Repo.Transaction(ctx, func(q Query) error {
		tag, e := s.Repo.Exec(ctx, q, query10, s.Now())
		if e == nil {
			affected = tag.RowsAffected()
		}
		return e
	})
	return affected, e
}
