package db

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeTx struct {
	pgx.Tx     // unused methods panic if called
	committed  bool
	rolledBack bool
	commitErr  error
}

func (f *fakeTx) Commit(context.Context) error   { f.committed = true; return f.commitErr }
func (f *fakeTx) Rollback(context.Context) error { f.rolledBack = true; return nil }

type fakeBeginner struct {
	tx  *fakeTx
	err error
}

func (b fakeBeginner) Begin(context.Context) (pgx.Tx, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.tx, nil
}

func TestWithTx(t *testing.T) {
	errFn := errors.New("fn failed")
	errCommit := errors.New("commit failed")
	errBegin := errors.New("begin failed")

	tests := []struct {
		name         string
		beginErr     error
		commitErr    error
		fnErr        error
		wantErr      error
		wantCommit   bool
		wantRollback bool
	}{
		{name: "success commits", wantCommit: true},
		{name: "fn error rolls back", fnErr: errFn, wantErr: errFn, wantRollback: true},
		{name: "commit error reported", commitErr: errCommit, wantErr: errCommit, wantCommit: true, wantRollback: true},
		{name: "begin error", beginErr: errBegin, wantErr: errBegin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &fakeTx{commitErr: tt.commitErr}
			err := WithTx(context.Background(), fakeBeginner{tx: tx, err: tt.beginErr}, func(pgx.Tx) error {
				return tt.fnErr
			})
			if !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && err != nil) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tx.committed != tt.wantCommit || tx.rolledBack != tt.wantRollback {
				t.Fatalf("committed=%v rolledBack=%v, want %v/%v", tx.committed, tx.rolledBack, tt.wantCommit, tt.wantRollback)
			}
		})
	}
}

func TestWithTxPanicRollsBackAndRepanics(t *testing.T) {
	tx := &fakeTx{}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic to propagate")
		}
		if !tx.rolledBack || tx.committed {
			t.Fatalf("committed=%v rolledBack=%v", tx.committed, tx.rolledBack)
		}
	}()
	_ = WithTx(context.Background(), fakeBeginner{tx: tx}, func(pgx.Tx) error { panic("boom") })
}

func TestViolations(t *testing.T) {
	unique := fmt.Errorf("insert: %w", &pgconn.PgError{Code: "23505", ConstraintName: "uq_contract_tor_contract_no"})
	fk := fmt.Errorf("insert: %w", &pgconn.PgError{Code: "23503", ConstraintName: "tor_location_assignment_location_id_fkey"})
	other := &pgconn.PgError{Code: "23514", ConstraintName: "some_check"}

	tests := []struct {
		name       string
		err        error
		wantUnique string
		wantFK     string
	}{
		{name: "unique", err: unique, wantUnique: "uq_contract_tor_contract_no"},
		{name: "foreign key", err: fk, wantFK: "tor_location_assignment_location_id_fkey"},
		{name: "other pg error", err: other},
		{name: "plain error", err: errors.New("boom")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, uok := UniqueViolation(tt.err)
			f, fok := ForeignKeyViolation(tt.err)
			if u != tt.wantUnique || uok != (tt.wantUnique != "") {
				t.Fatalf("UniqueViolation = %q, %v", u, uok)
			}
			if f != tt.wantFK || fok != (tt.wantFK != "") {
				t.Fatalf("ForeignKeyViolation = %q, %v", f, fok)
			}
		})
	}
}

func TestPoolTxPassesTx(t *testing.T) {
	tx := &fakeTx{}
	run := PoolTx(fakeBeginner{tx: tx})
	var got DBTX
	if err := run(context.Background(), func(q DBTX) error { got = q; return nil }); err != nil {
		t.Fatal(err)
	}
	if got != tx || !tx.committed {
		t.Fatalf("fn got %v, committed=%v", got, tx.committed)
	}
}
