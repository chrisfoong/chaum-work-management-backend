package contract

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/db"
	"chrisfoong/chaum-work-management-backend/internal/notify"
)

// These tests commit real rows into the throwaway database (TEST_DATABASE_URL,
// localhost only); reset it between runs. Skipped without TEST_DATABASE_URL.

// countingService returns a service on pool whose transaction runner counts transactions.
func countingService(pool *pgxpool.Pool, transactions *atomic.Int64) *Service {
	svc := NewService(TORRepository{}, pool, func(ctx context.Context, fn func(db.DBTX) error) error {
		transactions.Add(1)
		return db.PoolTx(pool)(ctx, fn)
	}, notify.Disabled{})
	svc.FileVerifier = func(context.Context, uuid.UUID, string) error { return nil } // Storage is tested separately with an HTTP test server.
	svc.runAsync = func(f func()) { f() }                                            // finish notifications before the pool closes
	return svc
}

func insertCommittedSupervisor(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	suffix := rand.Int64N(1e8)
	err := pool.QueryRow(context.Background(), `
		INSERT INTO public."USER" (first_name, last_name, phone_number, role, line_id, bank_name, bank_account_no)
		VALUES ('ทดสอบ', 'พร้อมกัน', $1, 'supervisor', $2, 'test', '000')
		RETURNING user_id`,
		fmt.Sprintf("09%08d", suffix), fmt.Sprintf("U%032x", suffix)).Scan(&id)
	if err != nil {
		t.Fatalf("insert supervisor: %v", err)
	}
	return id
}

func confirmRequest(projectName string, areas ...Area) ConfirmRequest {
	return ConfirmRequest{
		Contract: ContractInfo{
			ContractNo:    fmt.Sprintf("8%09d", rand.Int64N(1e9)),
			ProjectName:   projectName,
			PartnerAgency: "ทดสอบ", StartDate: "2026-01-01", EndDate: "2026-12-31", ContractValue: "1000.00", ContractFilePath: "test-contract.png",
		},
		Scope: Scope{Areas: areas},
	}
}

func newArea(locationName string) Area {
	return Area{
		NewLocation:     &NewLocation{Name: locationName, Address: "ทดสอบ"},
		RequiredWorkers: 1,
		// Shared names: the equipment upsert also runs concurrently.
		Items: []Item{{EquipmentName: "กรรไกรพร้อมกัน", RequiredQty: 1}, {EquipmentName: "รถเข็นพร้อมกัน", RequiredQty: 2}},
	}
}

func count(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), sql, args...).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// TestIntegrationConcurrentConfirms releases 6 valid confirms and 1 failing confirm
// at once (start barrier), for 5 rounds, all on the same Bangkok day. All 30 valid
// confirms must succeed with distinct requisition numbers; the failing ones must
// leave no rows behind.
func TestIntegrationConcurrentConfirms(t *testing.T) {
	const (
		goroutines = 6
		rounds     = 5
		itemsEach  = 2
	)
	pool := testPool(t)
	ctx := context.Background()
	supervisorID := insertCommittedSupervisor(t, pool)

	// A committed location whose name the failing confirm reuses in its 2nd area.
	existingName := "สวนมีอยู่แล้ว-" + uuid.NewString()
	if _, err := (TORRepository{}).CreateLocation(ctx, pool, existingName, "ทดสอบ"); err != nil {
		t.Fatalf("create existing location: %v", err)
	}

	prefix := "REQ-" + time.Now().In(bangkok).Format("20060102")
	dayReqsBefore := count(t, pool, `SELECT count(*) FROM equipment_requisition WHERE requisition_no LIKE $1 || '-%'`, prefix)

	var transactions atomic.Int64
	svc := countingService(pool, &transactions)

	var (
		goodNumbers   []string
		goodContracts []string
		failContracts []string
		failLocations []string
	)
	seen := map[string]bool{}
	for round := range rounds {
		start := make(chan struct{})
		var wg sync.WaitGroup
		goodOut := make([]ConfirmedContract, goroutines)
		goodErr := make([]error, goroutines)
		for g := range goroutines {
			req := confirmRequest(fmt.Sprintf("โครงการพร้อมกัน %d-%d", round, g), newArea("สวนพร้อมกัน-"+uuid.NewString()))
			goodContracts = append(goodContracts, req.Contract.ContractNo)
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				goodOut[g], goodErr[g] = svc.ConfirmContract(ctx, supervisorID, req)
			}()
		}
		// The failing confirm: contract + 1st area are inserted, then the 2nd area's
		// location name already exists, so the whole transaction must roll back.
		firstName := "สวนต้องไม่เหลือ-" + uuid.NewString()
		failReq := confirmRequest(fmt.Sprintf("โครงการล้มเหลว %d", round), newArea(firstName), newArea(existingName))
		failContracts = append(failContracts, failReq.Contract.ContractNo)
		failLocations = append(failLocations, firstName)
		var failErr error
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, failErr = svc.ConfirmContract(ctx, supervisorID, failReq)
		}()

		close(start)
		wg.Wait()

		for g := range goroutines {
			if goodErr[g] != nil {
				t.Errorf("round %d goroutine %d: %v", round, g, goodErr[g])
				continue
			}
			for _, a := range goodOut[g].Areas {
				no := a.Requisition.RequisitionNo
				if !strings.HasPrefix(no, prefix+"-") || seen[no] {
					t.Errorf("round %d goroutine %d: bad or duplicate number %q", round, g, no)
				}
				seen[no] = true
				goodNumbers = append(goodNumbers, no)
			}
		}
		var ae *apperr.Error
		if !errors.As(failErr, &ae) || ae.Code != "duplicate_location_name" {
			t.Errorf("round %d failing confirm: err = %v, want duplicate_location_name", round, failErr)
		}
	}

	if got := len(goodNumbers); got != goroutines*rounds {
		t.Errorf("successful confirms = %d, want %d", got, goroutines*rounds)
	}

	// No orphans from the failing confirms.
	if n := count(t, pool, `SELECT count(*) FROM contract_tor WHERE contract_no = ANY($1)`, failContracts); n != 0 {
		t.Errorf("orphan contracts from failed confirms: %d", n)
	}
	if n := count(t, pool, `SELECT count(*) FROM location WHERE location_name = ANY($1)`, failLocations); n != 0 {
		t.Errorf("orphan locations from failed confirms: %d", n)
	}
	if n := count(t, pool, `SELECT count(*) FROM equipment_requisition WHERE requisition_no LIKE $1 || '-%'`, prefix); n != dayReqsBefore+len(goodNumbers) {
		t.Errorf("requisitions for the day = %d, want %d (before %d + %d new)", n, dayReqsBefore+len(goodNumbers), dayReqsBefore, len(goodNumbers))
	}
	if n := count(t, pool, `
		SELECT count(*) FROM requisition_item ri
		JOIN equipment_requisition er ON er.requisition_id = ri.requisition_id
		WHERE er.requisition_no = ANY($1)`, goodNumbers); n != len(goodNumbers)*itemsEach {
		t.Errorf("items for new requisitions = %d, want %d", n, len(goodNumbers)*itemsEach)
	}
	if n := count(t, pool, `
		SELECT count(*) FROM tor_location_assignment tla
		JOIN contract_tor ct ON ct.tor_id = tla.tor_id
		WHERE ct.contract_no = ANY($1)
		  AND NOT EXISTS (SELECT 1 FROM equipment_requisition er WHERE er.assignment_id = tla.assignment_id)`, goodContracts); n != 0 {
		t.Errorf("assignments without a requisition: %d", n)
	}

	confirms := int64(goroutines*rounds + rounds)
	retries := transactions.Load() - confirms
	t.Logf("valid confirms ok: %d of %d; failing confirms: %d; transactions: %d; retries: %d",
		len(goodNumbers), goroutines*rounds, rounds, transactions.Load(), retries)
	if retries == 0 {
		t.Log("no retry happened: the advisory lock serialised numbering")
	}
}

// TestIntegrationLockReleasedOnCancel: while another transaction holds today's
// numbering lock, a confirm whose context expires fails with requisition_no_busy
// and leaves no rows; once the holder rolls back, a confirm succeeds.
func TestIntegrationLockReleasedOnCancel(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	supervisorID := insertCommittedSupervisor(t, pool)
	var transactions atomic.Int64
	svc := countingService(pool, &transactions)

	today := time.Now().In(bangkok)
	day := int32(today.Year()*10000 + int(today.Month())*100 + today.Day())
	holder, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := (TORRepository{}).LockRequisitionNumbering(ctx, holder, day); err != nil {
		t.Fatal(err)
	}

	blocked := confirmRequest("โครงการรอล็อก", newArea("สวนรอล็อก-"+uuid.NewString()))
	shortCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	_, err = svc.ConfirmContract(shortCtx, supervisorID, blocked)
	cancel()
	var ae *apperr.Error
	if !errors.As(err, &ae) || ae.Code != "requisition_no_busy" {
		t.Errorf("blocked confirm: err = %v, want requisition_no_busy", err)
	}
	if n := count(t, pool, `SELECT count(*) FROM contract_tor WHERE contract_no = $1`, blocked.Contract.ContractNo); n != 0 {
		t.Errorf("blocked confirm left %d contract rows", n)
	}

	if err := holder.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ConfirmContract(ctx, supervisorID, confirmRequest("โครงการหลังปล่อยล็อก", newArea("สวนหลังปล่อยล็อก-"+uuid.NewString()))); err != nil {
		t.Errorf("confirm after the lock was released: %v", err)
	}
}
