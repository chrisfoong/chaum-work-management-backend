package db

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
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
