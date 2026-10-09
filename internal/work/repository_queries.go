package work

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SQL is owned by the repository. Services choose operations and transaction order.
func (r Repository) Row(ctx context.Context, q Query, sql string, args ...any) pgx.Row {
	return q.QueryRow(ctx, sql, args...)
}
func (r Repository) Rows(ctx context.Context, q Query, sql string, args ...any) (pgx.Rows, error) {
	return q.Query(ctx, sql, args...)
}
func (r Repository) Exec(ctx context.Context, q Query, sql string, args ...any) (pgconn.CommandTag, error) {
	return q.Exec(ctx, sql, args...)
}
