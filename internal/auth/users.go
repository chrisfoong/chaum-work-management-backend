package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Querier runs a single-row query. *pgxpool.Pool and pgx.Tx satisfy it.
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// UserStore resolves verified identities to active USER rows.
type UserStore struct {
	db Querier
}

// NewUserStore returns a UserStore backed by db.
func NewUserStore(db Querier) *UserStore {
	return &UserStore{db: db}
}

// ByUserID resolves a Supabase Auth subject (the user's UUID).
func (s *UserStore) ByUserID(ctx context.Context, subject string) (Principal, error) {
	id, err := uuid.Parse(subject)
	if err != nil {
		return Principal{}, ErrUserNotFound
	}
	return s.one(ctx, `SELECT user_id, role FROM public."USER" WHERE user_id = $1 AND is_active`, id)
}

// ByLineID resolves a verified LINE user id.
func (s *UserStore) ByLineID(ctx context.Context, subject string) (Principal, error) {
	return s.one(ctx, `SELECT user_id, role FROM public."USER" WHERE line_id = $1 AND is_active`, subject)
}

func (s *UserStore) one(ctx context.Context, sql string, arg any) (Principal, error) {
	var (
		p    Principal
		role string
	)
	err := s.db.QueryRow(ctx, sql, arg).Scan(&p.UserID, &role)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrUserNotFound
	}
	if err != nil {
		return Principal{}, fmt.Errorf("look up user: %w", err)
	}
	p.Role = Role(role)
	return p, nil
}
