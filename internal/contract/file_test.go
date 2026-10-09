package contract

import (
	"chrisfoong/chaum-work-management-backend/internal/db"
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestContractFileRequiredAndVerifiedBeforeWrites(t *testing.T) {
	calls := 0
	s := NewService(&fakeStore{}, nil, func(context.Context, func(db.DBTX) error) error { calls++; return nil }, nil)
	req := confirmBody()
	req.Contract.ContractFilePath = ""
	if _, e := s.ConfirmContract(context.Background(), supervisorID, req); e == nil || calls != 0 {
		t.Fatal("missing file accepted / transaction started")
	}
	req = confirmBody()
	if _, e := s.ConfirmContract(context.Background(), supervisorID, req); e == nil || calls != 0 {
		t.Fatal("storage configuration bypassed")
	}
	s.FileVerifier = func(_ context.Context, user uuid.UUID, path string) error {
		if user != supervisorID || path != req.Contract.ContractFilePath {
			t.Fatal("wrong file identity")
		}
		return errors.New("not PNG")
	}
	if _, e := s.ConfirmContract(context.Background(), supervisorID, req); e == nil || calls != 0 {
		t.Fatal("unverified file committed")
	}
}
