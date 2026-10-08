package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"sync"
	"testing"
	"time"
)

type usecaseFixture struct {
	pool                                 *pgxpool.Pool
	service                              *Service
	assistant, supervisor, worker, other auth.Principal
	assignment, tor, workerID, otherID   string
	now                                  *time.Time
}

func fixtureUsecases(t *testing.T) usecaseFixture {
	t.Helper()
	pool := isolatedPool(t)
	ctx := context.Background()
	s := New(pool, strings.Repeat("s", 32))
	s.Files = testFiles{}
	now := time.Date(2026, 10, 1, 8, 0, 0, 0, Bangkok)
	s.Now = func() time.Time { return now }
	makeUser := func(role auth.Role) auth.Principal {
		id := uuid.New()
		phone := fmt.Sprintf("0%09d", uint64(id.ID())%1000000000)
		_, e := pool.Exec(ctx, `INSERT INTO public."USER"(user_id,first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no,daily_wage) VALUES($1,'fixture','usecase',$2,$3,$4,'test','0',400)`, id, string(role), phone, id.String())
		if e != nil {
			t.Fatal(e)
		}
		if role == auth.RoleWorker {
			if _, e = pool.Exec(ctx, `INSERT INTO worker(user_id) VALUES($1)`, id); e != nil {
				t.Fatal(e)
			}
		}
		return auth.Principal{UserID: id, Role: role}
	}
	f := usecaseFixture{pool: pool, service: s, supervisor: makeUser(auth.RoleSupervisor), assistant: makeUser(auth.RoleAssistant), worker: makeUser(auth.RoleWorker), other: makeUser(auth.RoleWorker), now: &now}
	f.workerID, _ = s.Worker(ctx, pool, f.worker.UserID.String())
	f.otherID, _ = s.Worker(ctx, pool, f.other.UserID.String())
	lat, lon := 13.75, 100.5
	loc, e := s.Location(ctx, "", LocationInput{Name: uuid.NewString(), Latitude: &lat, Longitude: &lon})
	if e != nil {
		t.Fatal(e)
	}
	if e = pool.QueryRow(ctx, `INSERT INTO contract_tor(contract_no,project_name,user_id,partner_agency,start_date,end_date,contract_value,status) VALUES($1,'usecases',$2,'test','2026-10-01','2026-10-31',1000,'active') RETURNING tor_id::text`, uuid.NewString(), f.supervisor.UserID).Scan(&f.tor); e != nil {
		t.Fatal(e)
	}
	if e = pool.QueryRow(ctx, `INSERT INTO tor_location_assignment(tor_id,location_id,required_workers) VALUES($1,$2,1) RETURNING assignment_id::text`, f.tor, loc).Scan(&f.assignment); e != nil {
		t.Fatal(e)
	}
	if _, e = pool.Exec(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'fixture','tor_base','completed')`, uuid.NewString(), f.supervisor.UserID, f.assignment); e != nil {
		t.Fatal(e)
	}
	return f
}
func TestIntegrationEmergencyLeaveAndAtomicReplacement(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	*f.now = time.Date(2026, 10, 4, 7, 0, 0, 0, Bangkok)
	original, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-04", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	leave, e := s.Leave(ctx, f.worker, LeaveInput{"2026-10-04", "emergency"})
	if e != nil {
		t.Fatal(e)
	}
	var advance bool
	if e = f.pool.QueryRow(ctx, `SELECT is_advance_notice FROM leave_request WHERE request_id=$1`, leave).Scan(&advance); e != nil || advance {
		t.Fatal("same-day must be emergency", e)
	}
	if _, e = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{"approved", uuid.NewString()}); e == nil {
		t.Fatal("invalid replacement accepted")
	}
	var state string
	if e = f.pool.QueryRow(ctx, `SELECT status::text FROM leave_request WHERE request_id=$1`, leave).Scan(&state); e != nil || state != "pending" {
		t.Fatal("failed replacement did not roll back approval", e)
	}
	ids, e := s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{"approved", f.otherID})
	if e != nil || len(ids) != 1 {
		t.Fatal(e)
	}
	*f.now = time.Date(2026, 10, 4, 18, 0, 0, 0, Bangkok)
	if _, e = s.Finalize(ctx); e != nil {
		t.Fatal(e)
	}
	if e = f.pool.QueryRow(ctx, `SELECT status::text FROM attendance WHERE schedule_id=$1`, original[0]).Scan(&state); e != nil || state != "absent" {
		t.Fatal("4S emergency absence classification wrong", state, e)
	}
	out, e := s.PreviewPayroll(ctx, PayrollInput{f.worker.UserID.String(), "2026-10-01", "2026-10-15"})
	if e != nil || out.Calculated != "1500.00" || out.Base != "0.00" || out.Net != "0.00" {
		t.Fatal(out, e)
	}
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	batch, e := s.PayrollBatch(ctx, f.supervisor, PayrollBatchInput{"2026-10-01", "2026-10-15"})
	if e != nil || len(batch) < 2 {
		t.Fatal("batch failed", batch, e)
	}
	retry, e := s.PayrollBatch(ctx, f.supervisor, PayrollBatchInput{"2026-10-01", "2026-10-15"})
	if e != nil || len(retry) != len(batch) {
		t.Fatal("batch retry duplicated payroll", retry, e)
	}
	for _, result := range retry {
		if !result.Existing {
			t.Fatal("batch replay created another payroll", result)
		}
	}
}
func TestIntegrationPurchaseRoundsDeliveryAndPaidReport(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	schedules, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	request := RequestInput{f.assignment, "need additional", []RequestItem{{equipment, 3, ""}}}
	if _, e = s.Requisition(ctx, f.worker, request); e == nil {
		t.Fatal("request without check-in allowed")
	}
	lat, lon, accuracy := 13.75, 100.5, 10.0
	qr, e := s.QR(ctx, f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.CheckIn(ctx, f.worker, CheckIn{schedules[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
	id, e := s.Requisition(ctx, f.worker, request)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.InspectRequest(ctx, id); e != nil {
		t.Fatal(e)
	}
	var item string
	if e = f.pool.QueryRow(ctx, `SELECT item_id::text FROM requisition_item WHERE requisition_id=$1`, id).Scan(&item); e != nil {
		t.Fatal(e)
	}
	if e = s.Survey(ctx, id, []SurveyItem{{item, 0}}); e == nil {
		t.Fatal("additional request bypassed decision")
	}
	if e = s.DecideRequest(ctx, f.assistant, id, ProcurementDecision{"purchase", "TOR reviewed"}); e != nil {
		t.Fatal(e)
	}
	if e = s.ReviewRequest(ctx, f.supervisor, id, ReviewInput{"approve", ""}); e != nil {
		t.Fatal(e)
	}
	fund := func() {
		t.Helper()
		if _, e = s.Fund(ctx, f.supervisor, id, FundInput{"100", uuid.NewString(), f.supervisor.UserID.String() + "/" + uuid.NewString()}); e != nil {
			t.Fatal(e)
		}
	}
	fund()
	zero, one := 0, 1
	first := PurchaseInput{[]PurchaseItem{{item, 1, "10.00", &zero}}, f.assistant.UserID.String() + "/" + uuid.NewString()}
	var wg sync.WaitGroup
	var outcomes [2]error
	for i := range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); _, outcomes[i] = s.Purchase(ctx, f.assistant, id, first) }()
	}
	wg.Wait()
	if (outcomes[0] == nil) == (outcomes[1] == nil) {
		t.Fatal("concurrent confirmation must apply once", outcomes)
	}
	var state string
	if e = f.pool.QueryRow(ctx, `SELECT status::text FROM equipment_requisition WHERE requisition_id=$1`, id).Scan(&state); e != nil || state != "pending_fund" {
		t.Fatal("partial status wrong", state, e)
	}
	fund()
	second := PurchaseInput{[]PurchaseItem{{item, 2, "20.00", &one}}, f.assistant.UserID.String() + "/" + uuid.NewString()}
	if _, e = s.Purchase(ctx, f.assistant, id, second); e != nil {
		t.Fatal(e)
	}
	var quantity int
	var price string
	if e = f.pool.QueryRow(ctx, `SELECT actual_qty,actual_price::text FROM requisition_item WHERE item_id=$1`, item).Scan(&quantity, &price); e != nil || quantity != 3 || price != "16.67" {
		t.Fatal("weighted average wrong", quantity, price, e)
	}
	if _, e = s.Purchase(ctx, f.assistant, id, second); e == nil {
		t.Fatal("repeat purchase applied twice")
	}
	delivery := DeliveryInput{schedules[0], "received by requesting worker", []string{f.assistant.UserID.String() + "/" + uuid.NewString(), f.assistant.UserID.String() + "/" + uuid.NewString()}}
	badDelivery := DeliveryInput{schedules[0], delivery.Description, []string{delivery.Photos[0], strings.Repeat("x", 300)}}
	if e = s.Deliver(ctx, f.assistant, id, badDelivery); e == nil {
		t.Fatal("invalid second evidence photo committed")
	}
	empty, e := s.Deliveries(ctx, id)
	if e != nil || string(empty) != "[]" {
		t.Fatal("delivery did not roll back all photos", e)
	}
	if e = s.Deliver(ctx, f.assistant, id, delivery); e != nil {
		t.Fatal(e)
	}
	if e = s.Deliver(ctx, f.assistant, id, delivery); e != nil {
		t.Fatal("matching delivery retry failed", e)
	}
	delivered, e := s.Deliveries(ctx, id)
	if e != nil {
		t.Fatal(e)
	}
	var evidence []json.RawMessage
	if e = json.Unmarshal(delivered, &evidence); e != nil || len(evidence) != 2 {
		t.Fatal("delivery evidence duplicated", e)
	}
	delivery.Description = "changed"
	if e = s.Deliver(ctx, f.assistant, id, delivery); e == nil {
		t.Fatal("delivery overwritten")
	}
	if _, e = f.pool.Exec(ctx, `INSERT INTO payroll(payroll_slip_no,user_id,managed_by_id,period_start,period_end,base_wage,total_deduction,net_wage) VALUES($1,$2,$3,'2026-10-01','2026-10-15',400,0,400)`, uuid.NewString(), f.worker.UserID, f.supervisor.UserID); e != nil {
		t.Fatal(e)
	}
	assertLabor := func(want string) {
		t.Helper()
		data, e := s.ProfitRange(ctx, f.tor, "2026-10-01", "2026-10-31")
		if e != nil {
			t.Fatal(e)
		}
		var rows []struct {
			Labor    string `json:"labor"`
			Material string `json:"material"`
		}
		if e = json.Unmarshal(data, &rows); e != nil || len(rows) != 1 || rows[0].Labor != want || rows[0].Material != "50.00" {
			t.Fatal("paid report mapping wrong", string(data), e)
		}
		if _, e = profitPDF(data); e != nil {
			t.Fatal(e)
		}
	}
	assertLabor("0.00")
	if _, e = f.pool.Exec(ctx, `UPDATE payroll SET is_paid=true WHERE user_id=$1`, f.worker.UserID); e != nil {
		t.Fatal(e)
	}
	assertLabor("400.00")
	noBuy, e := s.Requisition(ctx, f.worker, request)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.DecideRequest(ctx, f.assistant, noBuy, ProcurementDecision{"no_purchase", "reuse existing stock"}); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Purchase(ctx, f.assistant, noBuy, first); e == nil {
		t.Fatal("no-purchase decision still purchased")
	}
	// 2A keeps the most recent unit price; a zero-cost round has no expense.
	var base, baseItem string
	if e = f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type) VALUES($1,$2,$3,'base','tor_base') RETURNING requisition_id::text`, uuid.NewString(), f.supervisor.UserID, f.assignment).Scan(&base); e != nil {
		t.Fatal(e)
	}
	if e = f.pool.QueryRow(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty) VALUES($1,$2,2) RETURNING item_id::text`, base, equipment).Scan(&baseItem); e != nil {
		t.Fatal(e)
	}
	if e = s.Survey(ctx, base, []SurveyItem{{baseItem, 0}}); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Purchase(ctx, f.assistant, base, PurchaseInput{[]PurchaseItem{{baseItem, 1, "10.00", &zero}}, f.assistant.UserID.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Fund(ctx, f.supervisor, base, FundInput{"10", uuid.NewString(), f.supervisor.UserID.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	if expense, e := s.Purchase(ctx, f.assistant, base, PurchaseInput{[]PurchaseItem{{baseItem, 1, "0.00", &one}}, ""}); e != nil || expense != "" {
		t.Fatal("zero-cost round created expense", expense, e)
	}
	if e = f.pool.QueryRow(ctx, `SELECT actual_qty,actual_price::text FROM requisition_item WHERE item_id=$1`, baseItem).Scan(&quantity, &price); e != nil || quantity != 2 || price != "0.00" {
		t.Fatal("2A must use last unit price", quantity, price, e)
	}
	*f.now = f.now.Add(8 * time.Hour)
	if e = s.CheckOut(ctx, f.worker, CheckOut{schedules[0], "finished", f.worker.UserID.String() + "/" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Requisition(ctx, f.worker, request); e == nil {
		t.Fatal("equipment request after checkout accepted")
	}
}
