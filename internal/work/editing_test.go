package work

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestContractEditValidatesBeforeWriting(t *testing.T) {
	valid := ContractDetailsInput{ProjectName: "โครงการ", Agency: "หน่วยงาน", Start: "2026-10-01", End: "2026-10-31", Value: "1000.00", ExpectedUpdatedAt: "2026-10-01T00:00:00Z"}
	if err := validateContractEdit(valid); err != nil {
		t.Fatal(err)
	}
	for _, modify := range []func(*ContractDetailsInput){func(v *ContractDetailsInput) { v.End = "2026-09-01" }, func(v *ContractDetailsInput) { v.Value = "-1" }, func(v *ContractDetailsInput) { v.ExpectedUpdatedAt = "" }, func(v *ContractDetailsInput) { v.ProjectName = "" }} {
		v := valid
		modify(&v)
		if validateContractEdit(v) == nil {
			t.Fatal("invalid contract edit accepted")
		}
	}
}

func TestScopeEditRejectsDuplicatesAndInvalidQuantities(t *testing.T) {
	id := uuid.NewString()
	v := AssignmentEditInput{LocationID: uuid.NewString(), RequiredWorkers: 1, Items: []ScopeItemEdit{{EquipmentID: id, Required: 2}}}
	if err := validateAssignmentEdit(v); err != nil {
		t.Fatal(err)
	}
	v.Items = append(v.Items, ScopeItemEdit{EquipmentID: id, Required: 1})
	if validateAssignmentEdit(v) == nil {
		t.Fatal("duplicate equipment accepted")
	}
	v.Items = v.Items[:1]
	v.Items[0].Required = 0
	if validateAssignmentEdit(v) == nil {
		t.Fatal("zero TOR requirement accepted")
	}
	v.Items[0].Required = 1
	v.RequiredWorkers = 0
	if validateAssignmentEdit(v) == nil {
		t.Fatal("zero minimum staffing accepted")
	}
}

func TestEditingEndpointsRejectAssistant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	web := r.Group("/api/web", auth.Authenticate(qrTestVerifier{}, func(context.Context, string) (auth.Principal, error) {
		return auth.Principal{UserID: uuid.New(), Role: auth.RoleAssistant}, nil
	}, auth.RoleSupervisor, auth.RoleAssistant))
	Register(web, r.Group("/api/liff"), &Service{})
	id := uuid.NewString()
	for _, route := range []struct{ method, path string }{{"PATCH", "contracts/" + id + "/details"}, {"POST", "contracts/" + id + "/assignments"}, {"GET", "assignments/" + id + "/scope"}, {"PATCH", "assignments/" + id + "/scope"}, {"PATCH", "invoices/" + id}, {"PATCH", "users/" + id}, {"POST", "payroll/" + id + "/mark-paid"}} {
		req := httptest.NewRequest(route.method, "/api/web/"+route.path, strings.NewReader(`{}`))
		req.Header.Set("Authorization", "Bearer test-only")
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		r.ServeHTTP(response, req)
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s %s returned %d", route.method, route.path, response.Code)
		}
	}
}

func TestIntegrationContractInvoiceAndPeopleEdits(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	data, err := s.Contract(ctx, f.tor)
	if err != nil {
		t.Fatal(err)
	}
	var detail struct {
		Updated string `json:"updated_at"`
	}
	if err = json.Unmarshal(data, &detail); err != nil {
		t.Fatal(err)
	}
	edit := ContractDetailsInput{ProjectName: "Corrected", Agency: "Agency", Start: "2026-10-01", End: "2026-10-31", Value: "1000.00", ExpectedUpdatedAt: detail.Updated}
	if err = s.UpdateContractDetails(ctx, f.supervisor, f.tor, edit); err != nil {
		t.Fatal(err)
	}
	if err = s.UpdateContractDetails(ctx, f.supervisor, f.tor, edit); err == nil {
		t.Fatal("stale contract form overwrote newer edit")
	}
	invoice, err := s.Invoice(ctx, InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "1000", Deduction: "0"})
	if err != nil {
		t.Fatal(err)
	}
	change := InvoiceEditInput{Expected: "1100", Deduction: "100", PreviousExpected: "1000"}
	if err = s.UpdateInvoice(ctx, invoice, change); err != nil {
		t.Fatal(err)
	}
	if err = s.UpdateInvoice(ctx, invoice, change); err == nil {
		t.Fatal("stale invoice form accepted")
	}
	if err = s.Receive(ctx, invoice, "1000"); err != nil {
		t.Fatal(err)
	}
	change.PreviousExpected = "1100"
	if err = s.UpdateInvoice(ctx, invoice, change); err == nil {
		t.Fatal("paid invoice overwritten")
	}
	phone, bank, account := fmt.Sprintf("0%09d", uint64(uuid.New().ID())%1000000000), "Updated bank", "0123456789"
	inactive := false
	if err = s.UpdateUser(ctx, f.worker.UserID.String(), UserUpdate{FirstName: "Updated", LastName: "Worker", DailyWage: "450", Active: &inactive, Phone: &phone, BankName: &bank, BankAccount: &account}); err != nil {
		t.Fatal(err)
	}
	var active, available bool
	var gotPhone, gotBank, gotAccount string
	if err = f.pool.QueryRow(ctx, `SELECT u.is_active,w.is_available,u.phone_number,u.bank_name,u.bank_account_no FROM public."USER" u JOIN worker w ON w.user_id=u.user_id WHERE u.user_id=$1`, f.worker.UserID).Scan(&active, &available, &gotPhone, &gotBank, &gotAccount); err != nil {
		t.Fatal(err)
	}
	if active || available || gotPhone != phone || gotBank != bank || gotAccount != account {
		t.Fatal("user edit or automatic availability guard failed")
	}
}

func TestIntegrationScopeEditsStopAfterSurveyOrScheduling(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	lat, lon := 13.75, 100.5
	location, err := s.Location(ctx, "", LocationInput{Name: uuid.NewString(), Latitude: &lat, Longitude: &lon})
	if err != nil {
		t.Fatal(err)
	}
	equipment, err := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if err != nil {
		t.Fatal(err)
	}
	in := AssignmentEditInput{LocationID: location, RequiredWorkers: 1, Items: []ScopeItemEdit{{EquipmentID: equipment, Required: 2}}}
	area, err := s.SaveAssignment(ctx, f.supervisor, f.tor, "", in)
	if err != nil {
		t.Fatal(err)
	}
	expectedItems := []ScopeItemEdit{{EquipmentID: equipment, Required: 2}}
	in.ExpectedItems = &expectedItems
	previous := 1
	in.ExpectedWorkers = &previous
	in.RequiredWorkers = 2
	in.Items[0].Required = 3
	if _, err = s.SaveAssignment(ctx, f.supervisor, f.tor, area, in); err != nil {
		t.Fatal(err)
	}
	if _, err = s.SaveAssignment(ctx, f.supervisor, f.tor, area, in); err == nil {
		t.Fatal("stale minimum staffing accepted")
	}
	previous = 2
	if _, err = f.pool.Exec(ctx, `UPDATE equipment_requisition SET status='pending_procurement' WHERE assignment_id=$1`, area); err != nil {
		t.Fatal(err)
	}
	if _, err = s.SaveAssignment(ctx, f.supervisor, f.tor, area, in); err == nil {
		t.Fatal("scope changed after survey")
	}
}

func TestIntegrationPayrollRecordsConfirmerAndRejectsReplay(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	ids, err := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if err != nil {
		t.Fatal(err)
	}
	qr, err := s.QR(ctx, f.assignment)
	if err != nil {
		t.Fatal(err)
	}
	lat, lon, accuracy := 13.75, 100.5, 10.0
	if err = s.CheckIn(ctx, f.worker, CheckIn{ids[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); err != nil {
		t.Fatal(err)
	}
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	if _, err = s.Finalize(ctx); err != nil {
		t.Fatal(err)
	}
	slip, err := s.Payroll(ctx, f.supervisor, PayrollInput{UserID: f.worker.UserID.String(), Start: "2026-10-01", End: "2026-10-15"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.pool.Exec(ctx, `UPDATE public."USER" SET role='supervisor' WHERE user_id=$1`, f.assistant.UserID); err != nil {
		t.Fatal(err)
	}
	confirmer := auth.Principal{UserID: f.assistant.UserID, Role: auth.RoleSupervisor}
	if err = s.Pay(ctx, confirmer, slip); err != nil {
		t.Fatal(err)
	}
	if err = s.Pay(ctx, confirmer, slip); err == nil {
		t.Fatal("repeat payment accepted")
	}
	var manager string
	var paid bool
	if err = f.pool.QueryRow(ctx, `SELECT managed_by_id::text,is_paid FROM payroll WHERE payroll_id=$1`, slip).Scan(&manager, &paid); err != nil || !paid || manager != confirmer.UserID.String() {
		t.Fatal("payment confirmer not recorded", err)
	}
}

type supervisorCount int

func (c supervisorCount) Scan(dest ...any) error { *dest[0].(*int) = int(c); return nil }

type supervisorCountQuery int

func (c supervisorCountQuery) QueryRow(context.Context, string, ...any) pgx.Row {
	return supervisorCount(c)
}
func (c supervisorCountQuery) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	panic("guard must not write")
}
func (c supervisorCountQuery) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected query")
}
func TestLastSupervisorGuard(t *testing.T) {
	inactive := false
	active := true
	if protectLastSupervisor(context.Background(), supervisorCountQuery(1), "supervisor", true, &inactive) == nil {
		t.Fatal("last active Supervisor suspension accepted")
	}
	if err := protectLastSupervisor(context.Background(), supervisorCountQuery(2), "supervisor", true, &inactive); err != nil {
		t.Fatal("second Supervisor may be suspended", err)
	}
	if err := protectLastSupervisor(context.Background(), nil, "worker", true, &inactive); err != nil {
		t.Fatal(err)
	}
	if err := protectLastSupervisor(context.Background(), nil, "supervisor", false, &active); err != nil {
		t.Fatal(err)
	}
}
