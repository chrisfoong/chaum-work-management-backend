package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestIntegrationAtomicAttendanceAndRollback(t *testing.T) {
	pool := isolatedPool(t)
	ctx := context.Background()
	s := New(pool, strings.Repeat("y", 32))
	s.Files = testFiles{}
	now := time.Date(2026, 10, 8, 8, 0, 0, 0, Bangkok)
	s.Now = func() time.Time { return now }
	user := uuid.New()
	worker := uuid.New()
	sup := uuid.New()
	tor := uuid.New()
	loc := uuid.New()
	assignment := uuid.New()
	schedule := uuid.New()
	suffix := time.Now().UnixNano() % 1000000000
	for _, r := range []struct {
		id    uuid.UUID
		role  string
		phone string
	}{{sup, "supervisor", fmt.Sprintf("0%09d", suffix)}, {user, "worker", fmt.Sprintf("0%09d", suffix+1)}} {
		_, e := pool.Exec(ctx, `INSERT INTO public."USER"(user_id,first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no) VALUES($1,'test','atomic',$2,$3,$4,'test','0')`, r.id, r.role, r.phone, r.id.String())
		if e != nil {
			t.Fatal(e)
		}
	}
	for _, step := range []struct {
		sql  string
		args []any
	}{{`INSERT INTO worker(worker_id,user_id) VALUES($1,$2)`, []any{worker, user}}, {`INSERT INTO location(location_id,location_name,latitude,longitude) VALUES($1,$2,13.75,100.5)`, []any{loc, loc.String()}}, {`INSERT INTO contract_tor(tor_id,contract_no,project_name,user_id,partner_agency,start_date,end_date,contract_value) VALUES($1,$2,'atomic',$3,'test','2026-10-01','2026-10-31',1)`, []any{tor, tor.String(), sup}}, {`INSERT INTO tor_location_assignment(assignment_id,tor_id,location_id,required_workers) VALUES($1,$2,$3,1)`, []any{assignment, tor, loc}}, {`INSERT INTO work_schedule(schedule_id,assignment_id,worker_id,work_date,shift_start_time) VALUES($1,$2,$3,'2026-10-08','08:00')`, []any{schedule, assignment, worker}}} {
		if _, e := pool.Exec(ctx, step.sql, step.args...); e != nil {
			t.Fatal(e)
		}
	}
	principal := auth.Principal{UserID: user, Role: auth.RoleWorker}
	qr, e := s.QR(ctx, assignment.String())
	if e != nil {
		t.Fatal(e)
	}
	lat, lon, acc := 13.75, 100.5, 10.0
	in := CheckIn{schedule.String(), &lat, &lon, &acc, qr["qr_token"].(string)}
	var wg sync.WaitGroup
	var outcomes [4]error
	for i := range 4 {
		wg.Add(1)
		go func() { defer wg.Done(); outcomes[i] = s.CheckIn(ctx, principal, in) }()
	}
	wg.Wait()
	success := 0
	for _, e := range outcomes {
		if e == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("check-in successes %d", success)
	}
	// An overlong DB-bound path fails the evidence insert after updating check_out.
	// The whole transaction must roll back, leaving check_out unset.
	now = now.Add(8 * time.Hour)
	if e = s.CheckOut(ctx, principal, CheckOut{schedule.String(), "test", strings.Repeat("x", 300)}); e == nil {
		t.Fatal("overlong evidence unexpectedly saved")
	}
	var checkout *time.Time
	if e = pool.QueryRow(ctx, `SELECT check_out FROM attendance WHERE schedule_id=$1`, schedule).Scan(&checkout); e != nil || checkout != nil {
		t.Fatal("checkout not rolled back", e)
	}
	if e = s.CheckOut(ctx, principal, CheckOut{schedule.String(), "test", user.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
}
