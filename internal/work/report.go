package work

import (
	"context"
	"encoding/json"
	"strings"
	"time"
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

func validRange(start, end string) error {
	a, e := Date(start)
	if e != nil {
		return e
	}
	b, e := Date(end)
	if e != nil {
		return e
	}
	if b.Before(a) || b.Sub(a) > 3660*24*time.Hour {
		return invalid("period", "ordered date range of at most ten years required")
	}
	return nil
}
func (s *Service) ProfitRange(ctx context.Context, tor, start, end string) (json.RawMessage, error) {
	if e := validRange(start, end); e != nil {
		return nil, e
	}
	var selected any
	if tor != "" {
		if e := validID(tor, "tor_id"); e != nil {
			return nil, e
		}
		selected = tor
	}
	var duplicates bool
	if e := s.Repo.Row(ctx, s.Repo.Pool, `SELECT EXISTS(SELECT 1 FROM worker GROUP BY user_id HAVING count(*)>1)`).Scan(&duplicates); e != nil {
		return nil, e
	}
	if duplicates {
		return nil, conflict("duplicate Worker/User mappings must be resolved before reporting")
	}
	sql := strings.ReplaceAll(query87, "p.is_paid AND to_char(p.period_start,'YYYY-MM')=$1", "p.is_paid AND p.period_start >= $1::date AND p.period_end <= $2::date")
	sql = strings.ReplaceAll(sql, "billing_month=$1", "billing_month BETWEEN to_char($1::date,'YYYY-MM') AND to_char($2::date,'YYYY-MM')")
	sql = strings.ReplaceAll(sql, "to_char(e.created_at AT TIME ZONE 'Asia/Bangkok','YYYY-MM')=$1", "(e.created_at AT TIME ZONE 'Asia/Bangkok')::date BETWEEN $1::date AND $2::date")
	sql = strings.ReplaceAll(sql, "$1::text AS month", "$1::text AS period_start,$2::text AS period_end")
	sql = strings.ReplaceAll(sql, "ORDER BY c.tor_id", "WHERE ($3::uuid IS NULL OR c.tor_id=$3) ORDER BY c.tor_id")
	return s.Repo.List(ctx, s.Repo.Pool, sql, start, end, selected)
}

func (s *Service) CloseSummary(ctx context.Context, tor, start, end string) (json.RawMessage, error) {
	if e := validID(tor, "tor_id"); e != nil {
		return nil, e
	}
	data, e := s.ProfitRange(ctx, tor, start, end)
	if e != nil {
		return nil, e
	}
	var rows []json.RawMessage
	if e = json.Unmarshal(data, &rows); e != nil {
		return nil, e
	}
	if len(rows) == 0 {
		return nil, conflict("TOR does not exist")
	}
	// 6S confirmation is a snapshot, not contract termination. No close ledger
	// exists; never encode closure by changing contract status.
	return json.Marshal(struct {
		Data        json.RawMessage `json:"data"`
		ConfirmedAt string          `json:"confirmed_at"`
		Persisted   bool            `json:"persisted"`
	}{data, s.Now().Format(time.RFC3339), false})
}
