package contract

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"chrisfoong/chaum-work-management-backend/internal/db"
)

// These tests need a throwaway LOCAL PostgreSQL with migrations/0001_init.up.sql
// applied, given as TEST_DATABASE_URL. Without it they are skipped and the 1S SQL
// stays unverified. Each test runs in a transaction that is rolled back.

func testConn(t *testing.T) db.DBTX {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set: integration test skipped (SQL unverified)")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal("TEST_DATABASE_URL is not a valid connection string")
	}
	switch cfg.ConnConfig.Host {
	case "localhost", "127.0.0.1", "::1":
	default:
		if len(cfg.ConnConfig.Host) == 0 || cfg.ConnConfig.Host[0] != '/' { // unix socket dirs are local
			t.Fatal("refusing to run: TEST_DATABASE_URL must point to a local throwaway database")
		}
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal("cannot connect to the test database")
	}
	t.Cleanup(pool.Close)
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal("cannot begin a transaction on the test database")
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	return tx
}

func insertSupervisor(t *testing.T, q db.DBTX) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := q.QueryRow(context.Background(), `
		INSERT INTO users (first_name, last_name, phone_number, role, line_id, bank_name, bank_account_no)
		VALUES ('ทดสอบ', 'ระบบ', '0800000001', 'supervisor', 'U00000000000000000000000000000001', 'test', '000')
		RETURNING user_id`).Scan(&id)
	if err != nil {
		t.Fatalf("insert supervisor: %v", err)
	}
	return id
}

func TestIntegrationDuplicateContractNoIsUniqueViolation(t *testing.T) {
	q := testConn(t)
	ctx := context.Background()
	repo := TORRepository{}
	userID := insertSupervisor(t, q)
	c := validInfo()
	if _, err := repo.CreateContract(ctx, q, userID, c); err != nil {
		t.Fatal(err)
	}
	_, err := repo.CreateContract(ctx, q, userID, c)
	if name, ok := db.UniqueViolation(err); !ok || name != constraintContractNo {
		t.Fatalf("err = %v, want unique violation on %s", err, constraintContractNo)
	}
}

func TestIntegrationMissingLocationIsForeignKeyViolation(t *testing.T) {
	q := testConn(t)
	ctx := context.Background()
	repo := TORRepository{}
	torID, err := repo.CreateContract(ctx, q, insertSupervisor(t, q), validInfo())
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = repo.CreateTORLocationAssignment(ctx, q, torID, uuid.New(), 2)
	if name, ok := db.ForeignKeyViolation(err); !ok || name != constraintAssignmentLoc {
		t.Fatalf("err = %v, want foreign key violation on %s", err, constraintAssignmentLoc)
	}
}

func TestIntegrationDuplicateLocationNameIsUniqueViolation(t *testing.T) {
	q := testConn(t)
	ctx := context.Background()
	repo := TORRepository{}
	if _, err := repo.CreateLocation(ctx, q, "สวนทดสอบ", "ที่อยู่"); err != nil {
		t.Fatal(err)
	}
	_, err := repo.CreateLocation(ctx, q, "สวนทดสอบ", "อื่น")
	if name, ok := db.UniqueViolation(err); !ok || name != constraintLocationName {
		t.Fatalf("err = %v, want unique violation on %s", err, constraintLocationName)
	}
}

func TestIntegrationEquipmentFindOrCreate(t *testing.T) {
	q := testConn(t)
	ctx := context.Background()
	repo := TORRepository{}
	var ids [2]uuid.UUID
	for i := range ids {
		if err := repo.CreateEquipment(ctx, q, "เครื่องตัดหญ้าทดสอบ"); err != nil {
			t.Fatal(err)
		}
		id, err := repo.FindEquipmentByName(ctx, q, "เครื่องตัดหญ้าทดสอบ")
		if err != nil {
			t.Fatal(err)
		}
		ids[i] = id
	}
	if ids[0] != ids[1] {
		t.Fatalf("second find-or-create returned a different id: %v", ids)
	}
}

func TestIntegrationRequisitionNumbering(t *testing.T) {
	q := testConn(t)
	ctx := context.Background()
	repo := TORRepository{}
	userID := insertSupervisor(t, q)
	torID, err := repo.CreateContract(ctx, q, userID, validInfo())
	if err != nil {
		t.Fatal(err)
	}
	locID, err := repo.CreateLocation(ctx, q, "สวนเลขที่", "x")
	if err != nil {
		t.Fatal(err)
	}
	assignmentID, name, err := repo.CreateTORLocationAssignment(ctx, q, torID, locID, 1)
	if err != nil || name != "สวนเลขที่" {
		t.Fatalf("assignment: %v %q", err, name)
	}

	const prefix = "REQ-29990101"
	next, err := repo.NextRequisitionSeq(ctx, q, prefix)
	if err != nil || next != 1 {
		t.Fatalf("first next = %d, %v", next, err)
	}
	if _, err := repo.CreateInitialRequisition(ctx, q, assignmentID, userID, prefix+"-007"); err != nil {
		t.Fatal(err)
	}
	if next, err = repo.NextRequisitionSeq(ctx, q, prefix); err != nil || next != 8 {
		t.Fatalf("next after 007 = %d, %v", next, err)
	}
	_, err = repo.CreateInitialRequisition(ctx, q, assignmentID, userID, prefix+"-007")
	if name, ok := db.UniqueViolation(err); !ok || name != constraintRequisitionNo {
		t.Fatalf("err = %v, want unique violation on %s", err, constraintRequisitionNo)
	}
}

// compile-time check that a transaction satisfies DBTX.
var _ db.DBTX = pgx.Tx(nil)
