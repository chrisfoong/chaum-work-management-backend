// Package work implements the confirmed MVP workflows. Handlers bind input;
// services enforce policies and transaction boundaries; Repository owns SQL access.
package work

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Bangkok = time.FixedZone("Asia/Bangkok", 7*3600)

type Query interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
type Repository struct{ Pool *pgxpool.Pool }
type Service struct {
	Repo     Repository
	Now      func() time.Time
	QRSecret []byte
	Files    EvidenceStore
	Notify   Notifier
}

func New(pool *pgxpool.Pool, secret string) *Service {
	return &Service{Repo: Repository{pool}, Now: time.Now, QRSecret: []byte(secret)}
}

type EvidenceStore interface {
	Verify(context.Context, auth.Principal, string) error
}
type Notifier interface {
	Send(context.Context, string, string) error
}

func (r Repository) Transaction(ctx context.Context, f func(Query) error) error {
	tx, e := r.Pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(context.WithoutCancel(ctx))
	if e = f(tx); e != nil {
		return e
	}
	return tx.Commit(ctx)
}
func (r Repository) List(ctx context.Context, q Query, sql string, args ...any) (json.RawMessage, error) {
	var out []byte
	e := q.QueryRow(ctx, "SELECT COALESCE(jsonb_agg(to_jsonb(result)), '[]'::jsonb) FROM ("+sql+") result", args...).Scan(&out)
	return json.RawMessage(out), e
}
func lock(ctx context.Context, q Query, key string) error {
	_, e := q.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, key)
	return e
}
func one(ctx context.Context, q Query, sql string, args ...any) (json.RawMessage, error) {
	var out []byte
	e := q.QueryRow(ctx, "SELECT to_jsonb(result) FROM ("+sql+") result", args...).Scan(&out)
	if errors.Is(e, pgx.ErrNoRows) {
		return nil, apperr.NotFound("record not found")
	}
	return out, e
}
func invalid(field, msg string) error {
	return apperr.Validation(apperr.FieldError{Field: field, Message: msg})
}
func conflict(msg string) error { return apperr.Conflict("state_conflict", msg) }
func validID(v, field string) error {
	if _, e := uuid.Parse(v); e != nil {
		return invalid(field, "must be UUID")
	}
	return nil
}
func text(v, field string, max int) error {
	if strings.TrimSpace(v) == "" || len([]rune(v)) > max {
		return invalid(field, fmt.Sprintf("required and at most %d characters", max))
	}
	return nil
}
func document(prefix string) string { return prefix + "-" + uuid.NewString() }

// Money parses nonnegative decimal strings without rounding or binary floating point.
func Money(v string) (int64, error) {
	parts := strings.Split(v, ".")
	if len(parts) > 2 || len(parts[0]) == 0 || len(parts[0]) > 12 {
		return 0, invalid("amount", "invalid decimal")
	}
	for _, part := range parts {
		if len(part) == 0 {
			return 0, invalid("amount", "invalid decimal")
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return 0, invalid("amount", "nonnegative decimal required")
			}
		}
	}
	whole, e := strconv.ParseInt(parts[0], 10, 64)
	if e != nil {
		return 0, e
	}
	frac := int64(0)
	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return 0, invalid("amount", "at most two decimal places")
		}
		f := parts[1]
		if len(f) == 1 {
			f += "0"
		}
		frac, _ = strconv.ParseInt(f, 10, 64)
	}
	return whole*100 + frac, nil
}
func Decimal(c int64) string { return fmt.Sprintf("%d.%02d", c/100, c%100) }
func Date(v string) (time.Time, error) {
	d, e := time.ParseInLocation("2006-01-02", v, Bangkok)
	if e != nil {
		return d, invalid("date", "YYYY-MM-DD required")
	}
	return d, nil
}
func Start(date, clock string) (time.Time, error) {
	if len(clock) == 5 {
		clock += ":00"
	}
	t, e := time.ParseInLocation("2006-01-02 15:04:05", date+" "+clock, Bangkok)
	if e != nil {
		return t, invalid("shift_start_time", "invalid local shift start")
	}
	return t, nil
}
func Period(start, end string) error {
	a, e := Date(start)
	if e != nil {
		return e
	}
	b, e := Date(end)
	if e != nil {
		return e
	}
	last := time.Date(a.Year(), a.Month()+1, 0, 0, 0, 0, 0, Bangkok).Day()
	if a.Year() != b.Year() || a.Month() != b.Month() || !((a.Day() == 1 && b.Day() == 15) || (a.Day() == 16 && b.Day() == last)) {
		return invalid("period", "must be 1-15 or 16-month end")
	}
	return nil
}
func Penalty(status string, start time.Time, check *time.Time) int64 {
	if status == "absent" {
		return 150000
	}
	if status != "late" || check == nil {
		return 0
	}
	d := check.Sub(start)
	switch {
	case d <= 0:
		return 0
	case d < time.Hour:
		return 30000
	case d <= 3*time.Hour:
		return 40000
	default:
		return 150000
	}
}
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	rad := math.Pi / 180
	a := math.Sin((lat2 - lat1) * rad / 2)
	b := math.Sin((lon2 - lon1) * rad / 2)
	h := a*a + math.Cos(lat1*rad)*math.Cos(lat2*rad)*b*b
	if h > 1 {
		h = 1
	}
	return 6371000 * 2 * math.Asin(math.Sqrt(h))
}
func (s *Service) Worker(ctx context.Context, q Query, user string) (string, error) {
	rows, e := q.Query(ctx, `SELECT worker_id::text FROM worker WHERE user_id=$1 ORDER BY worker_id LIMIT 2`, user)
	if e != nil {
		return "", e
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			return "", e
		}
		ids = append(ids, id)
	}
	if e = rows.Err(); e != nil {
		return "", e
	}
	if len(ids) > 1 {
		return "", conflict("multiple workers for user")
	}
	if len(ids) == 0 {
		return "", apperr.NotFound("worker not registered")
	}
	return ids[0], nil
}
func (s *Service) evidence(ctx context.Context, p auth.Principal, path string) error {
	if s.Files == nil {
		return conflict("storage is not configured")
	}
	return s.Files.Verify(ctx, p, path)
}
func (s *Service) notify(ctx context.Context, user, msg string) {
	if s.Notify == nil {
		return
	}
	var line string
	if s.Repo.Pool.QueryRow(ctx, `SELECT line_id FROM public."USER" WHERE user_id=$1 AND is_active`, user).Scan(&line) != nil {
		return
	}
	if e := s.Notify.Send(ctx, line, msg); e != nil {
		slog.Warn("notification failed")
	} else {
		slog.Info("notification accepted by LINE; recipient delivery unverified")
	}
}

func isMissing(e error) bool { return errors.Is(e, pgx.ErrNoRows) }
