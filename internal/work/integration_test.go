package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/schema"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type testFiles struct{}

func (testFiles) Verify(context.Context, auth.Principal, string) error { return nil }
func isolatedPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	v := os.Getenv("TEST_DATABASE_URL")
	if v == "" {
		t.Skip("isolated PostgreSQL not configured; integration unverified")
	}
	if os.Getenv("TEST_DATABASE_ISOLATED") != "yes" {
		t.Fatal("isolated test marker required")
	}
	cfg, e := pgxpool.ParseConfig(v)
	if e != nil {
		t.Fatal("invalid test DSN")
	}
	if cfg.ConnConfig.Host != "127.0.0.1" && cfg.ConnConfig.Host != "localhost" {
		t.Fatal("local isolated database required")
	}
	pool, e := pgxpool.NewWithConfig(context.Background(), cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	if e = schema.Check(context.Background(), pool); e != nil {
		t.Fatal(e)
	}
	return pool
}
func TestIntegrationMVP(t *testing.T) {
	pool := isolatedPool(t)
	ctx := context.Background()
	s := New(pool, strings.Repeat("x", 32))
	s.Files = testFiles{}
	now := time.Date(2026, 10, 1, 8, 0, 0, 0, Bangkok)
	s.Now = func() time.Time { return now }
	create := func(role string) auth.Principal {
		u := uuid.New()
		phone := fmt.Sprintf("0%09d", time.Now().UnixNano()%1000000000)
		out, e := s.CreateUser(ctx, auth.Principal{}, UserInput{FirstName: "Test", LastName: "MVP", Role: role, Phone: phone, LineID: u.String(), BankName: "TEST", BankAccount: "0", DailyWage: "400"})
		if e != nil {
			t.Fatal(e)
		}
		var got struct {
			ID uuid.UUID `json:"user_id"`
		}
		if e = json.Unmarshal(out, &got); e != nil {
			t.Fatal(e)
		}
		return auth.Principal{UserID: got.ID, Role: auth.Role(role)}
	}
	sup := create("supervisor")
	worker := create("worker")
	replacement := create("worker")
	w1, e := s.Worker(ctx, pool, worker.UserID.String())
	if e != nil {
		t.Fatal(e)
	}
	w2, e := s.Worker(ctx, pool, replacement.UserID.String())
	if e != nil {
		t.Fatal(e)
	}
	lat, lon := 13.75, 100.5
	loc, e := s.Location(ctx, "", LocationInput{Name: uuid.NewString(), Address: "test", Latitude: &lat, Longitude: &lon})
	if e != nil {
		t.Fatal(e)
	}
	var tor, assignment string
	if e = pool.QueryRow(ctx, `INSERT INTO contract_tor(contract_no,project_name,user_id,partner_agency,start_date,end_date,contract_value,status) VALUES($1,'test',$2,'test','2026-10-01','2026-10-31',1000,'active') RETURNING tor_id::text`, uuid.NewString(), sup.UserID).Scan(&tor); e != nil {
		t.Fatal(e)
	}
	if e = pool.QueryRow(ctx, `INSERT INTO tor_location_assignment(tor_id,location_id,required_workers) VALUES($1,$2,1) RETURNING assignment_id::text`, tor, loc).Scan(&assignment); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'ready fixture','tor_base','completed')`, uuid.NewString(), sup.UserID, assignment); e != nil {
		t.Fatal(e)
	}
	ids, e := s.Schedule(ctx, ScheduleInput{assignment, []string{w1}, "2026-10-04", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	leave, e := s.Leave(ctx, worker, LeaveInput{"2026-10-04", "planned"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Leave(ctx, worker, LeaveInput{"2026-10-04", "duplicate"}); e == nil {
		t.Fatal("duplicate leave allowed")
	}
	if e = s.ReviewLeave(ctx, leave, "approved"); e != nil {
		t.Fatal(e)
	}
	sub, e := s.Replacement(ctx, leave, w2)
	if e != nil {
		t.Fatal(e)
	}
	now = time.Date(2026, 10, 4, 8, 0, 0, 0, Bangkok)
	qr, e := s.QR(ctx, assignment)
	if e != nil {
		t.Fatal(e)
	}
	accuracy := 10.0
	check := CheckIn{ScheduleID: sub[0], Latitude: &lat, Longitude: &lon, Accuracy: &accuracy, QR: qr["qr_token"].(string)}
	if e = s.CheckIn(ctx, worker, check); e == nil {
		t.Fatal("cross-owner check-in allowed")
	}
	if e = s.CheckIn(ctx, replacement, check); e != nil {
		t.Fatal(e)
	}
	if e = s.CheckIn(ctx, replacement, check); e == nil {
		t.Fatal("duplicate check-in allowed")
	}
	now = now.Add(8 * time.Hour)
	if e = s.CheckOut(ctx, replacement, CheckOut{sub[0], "finished", replacement.UserID.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	if e = s.CheckOut(ctx, replacement, CheckOut{sub[0], "again", "fake"}); e == nil {
		t.Fatal("duplicate checkout allowed")
	}
	now = time.Date(2026, 10, 5, 8, 0, 0, 0, Bangkok)
	absent, e := s.Schedule(ctx, ScheduleInput{assignment, []string{w2}, "2026-10-05", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	now = time.Date(2026, 10, 6, 8, 0, 0, 0, Bangkok)
	if _, e = s.Finalize(ctx); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Finalize(ctx); e != nil {
		t.Fatal(e)
	}
	var status string
	if e = pool.QueryRow(ctx, `SELECT status::text FROM attendance WHERE schedule_id=$1`, absent[0]).Scan(&status); e != nil || status != "absent" {
		t.Fatal("absent job failed", e)
	}
	if e = pool.QueryRow(ctx, `SELECT status::text FROM attendance WHERE schedule_id=$1`, ids[0]).Scan(&status); e != nil || status != "leave" {
		t.Fatal("leave incorrectly classified", e)
	}
	now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	period := PayrollInput{replacement.UserID.String(), "2026-10-01", "2026-10-15"}
	preview, e := s.PreviewPayroll(ctx, period)
	if e != nil || preview.Base != "400.00" || preview.Calculated != "1500.00" || preview.Applied != "400.00" || preview.Net != "0.00" {
		t.Fatalf("bad payroll cap: %+v %v", preview, e)
	}
	var outcomes [2]error
	var payroll [2]string
	var wg sync.WaitGroup
	for i := range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); payroll[i], outcomes[i] = s.Payroll(ctx, sup, period) }()
	}
	wg.Wait()
	success := 0
	var slip string
	for i, e := range outcomes {
		if e == nil {
			success++
			slip = payroll[i]
		}
	}
	if success != 1 {
		t.Fatalf("concurrent payroll successes %d: %v", success, outcomes)
	}
	if e = s.Pay(ctx, slip); e != nil {
		t.Fatal(e)
	}
	if e = s.ReviewLeave(ctx, leave, "rejected"); e == nil {
		t.Fatal("reviewed leave changed")
	}
	unpaid, e := s.PreviewPayroll(ctx, PayrollInput{worker.UserID.String(), period.Start, period.End})
	if e != nil || unpaid.Base != "0.00" || unpaid.Calculated != "0.00" {
		t.Fatal("approved leave must be unpaid without penalty", e)
	}
	equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	open, e := s.Schedule(ctx, ScheduleInput{assignment, []string{w2}, "2026-10-16", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	qr, e = s.QR(ctx, assignment)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.CheckIn(ctx, replacement, CheckIn{ScheduleID: open[0], Latitude: &lat, Longitude: &lon, Accuracy: &accuracy, QR: qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
	req, e := s.Requisition(ctx, replacement, RequestInput{assignment, "extra", []RequestItem{{equipment, 2, "test"}}})
	if e != nil {
		t.Fatal(e)
	}
	var item string
	if e = pool.QueryRow(ctx, `SELECT item_id::text FROM requisition_item WHERE requisition_id=$1`, req).Scan(&item); e != nil {
		t.Fatal(e)
	}
	if e = s.DecideRequest(ctx, sup, req, ProcurementDecision{"purchase", "needed"}); e != nil {
		t.Fatal(e)
	}
	if e = s.ReviewRequest(ctx, sup, req, ReviewInput{"approve", ""}); e != nil {
		t.Fatal(e)
	}
	fund := FundInput{"100", uuid.NewString(), sup.UserID.String() + "/" + uuid.NewString()}
	first, e := s.Fund(ctx, sup, req, fund)
	if e != nil {
		t.Fatal(e)
	}
	second, e := s.Fund(ctx, sup, req, fund)
	if e != nil || first != second {
		t.Fatal("fund retry duplicated", e)
	}
	fund.Amount = "101"
	if _, e = s.Fund(ctx, sup, req, fund); e == nil {
		t.Fatal("fund retry payload changed")
	}
	zero := 0
	purchase := PurchaseInput{[]PurchaseItem{{item, 2, "50", &zero}}, replacement.UserID.String() + "/" + uuid.NewString()}
	if _, e = s.Purchase(ctx, replacement, req, purchase); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Purchase(ctx, replacement, req, purchase); e == nil {
		t.Fatal("purchase overwritten")
	}
	invoice, e := s.Invoice(ctx, InvoiceInput{TorID: tor, Month: "2026-10", Expected: "1000"})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Receive(ctx, invoice, "900"); e != nil {
		t.Fatal(e)
	}
	report, e := s.Profit(ctx, "2026-10")
	if e != nil {
		t.Fatal(e)
	}
	var reports []struct {
		Tor      string `json:"tor_id"`
		Revenue  string `json:"revenue"`
		Material string `json:"material"`
		Profit   string `json:"profit"`
	}
	if e = json.Unmarshal(report, &reports); e != nil {
		t.Fatal(e)
	}
	found := false
	for _, r := range reports {
		if r.Tor == tor {
			found = true
			a, _ := Money(r.Revenue)
			b, _ := Money(r.Material)
			c, _ := Money(r.Profit)
			if a != 90000 || b != 10000 || c != 80000 {
				t.Fatalf("fund counted twice or revenue wrong: %+v", r)
			}
		}
	}
	if !found {
		t.Fatal("TOR absent from report")
	}
	// Duplicate physical worker rows must be reported, not silently selected.
	var duplicateID string
	if e = pool.QueryRow(ctx, `INSERT INTO worker(user_id) VALUES($1) RETURNING worker_id::text`, worker.UserID).Scan(&duplicateID); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM worker WHERE worker_id=$1`, duplicateID) })
	if _, e = s.Worker(ctx, pool, worker.UserID.String()); e == nil {
		t.Fatal("duplicate worker accepted")
	}
}
