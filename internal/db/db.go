// Package db owns the PostgreSQL connection pool and the transaction helper.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool connects to PostgreSQL and checks the connection.
func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// Beginner starts a transaction. *pgxpool.Pool satisfies it.
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// WithTx runs fn inside one transaction. It commits when fn returns nil and
// rolls back when fn returns an error or panics.
func WithTx(ctx context.Context, db Beginner, fn func(tx pgx.Tx) error) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = errors.Join(err, fmt.Errorf("rollback: %w", rbErr))
			}
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// DBTX runs statements. *pgxpool.Pool and pgx.Tx satisfy it, so repositories
// work the same inside and outside a transaction.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TxRunner runs fn inside one transaction.
type TxRunner func(ctx context.Context, fn func(q DBTX) error) error

// PoolTx returns a TxRunner that opens transactions on b.
func PoolTx(b Beginner) TxRunner {
	return func(ctx context.Context, fn func(q DBTX) error) error {
		return WithTx(ctx, b, func(tx pgx.Tx) error { return fn(tx) })
	}
}

const (
	codeUniqueViolation     = "23505"
	codeForeignKeyViolation = "23503"
)

// UniqueViolation reports whether err is a unique-constraint violation and
// returns the violated constraint's name.
func UniqueViolation(err error) (constraint string, ok bool) {
	return pgViolation(err, codeUniqueViolation)
}

// ForeignKeyViolation reports whether err is a foreign-key violation and
// returns the violated constraint's name.
func ForeignKeyViolation(err error) (constraint string, ok bool) {
	return pgViolation(err, codeForeignKeyViolation)
}

func pgViolation(err error, code string) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == code {
		return pgErr.ConstraintName, true
	}
	return "", false
}
