package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/db"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
	"chrisfoong/chaum-work-management-backend/internal/notify"
)

func init() {
	gin.SetMode(gin.TestMode)
	httpx.UseJSONFieldNames()
}

var (
	supervisorID = uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assistantID  = uuid.MustParse("b1ffcd00-1d1c-4f09-9c7e-7cc0ce491b22")
)

// fakeStore records calls and returns configured results.
type fakeStore struct {
	dupContractNo     bool
	missingLocation   bool
	dupLocationName   bool
	createContractErr error
	assignmentErr     error
	requisitionErrs   []error // returned in order by CreateInitialRequisition
	assistants        []string
	nextSeq           int

	calls            []string
	requisitionNos   []string
	equipmentCreated []string
}

func (f *fakeStore) call(name string) { f.calls = append(f.calls, name) }

func (f *fakeStore) CheckDuplicateContractNo(context.Context, db.DBTX, string) (bool, error) {
	f.call("checkDuplicateContractNo")
	return f.dupContractNo, nil
}
func (f *fakeStore) ValidateLocation(context.Context, db.DBTX, uuid.UUID) (bool, error) {
	f.call("validateLocation")
	return !f.missingLocation, nil
}
func (f *fakeStore) CheckDuplicateLocationName(context.Context, db.DBTX, string) (bool, error) {
	f.call("checkDuplicateLocationName")
	return f.dupLocationName, nil
}
func (f *fakeStore) CreateContract(context.Context, db.DBTX, uuid.UUID, ContractInfo) (uuid.UUID, error) {
	f.call("createContract")
	return uuid.New(), f.createContractErr
}
func (f *fakeStore) CreateLocation(context.Context, db.DBTX, string, string) (uuid.UUID, error) {
	f.call("createLocation")
	if f.dupLocationName {
		return uuid.Nil, pgUnique(constraintLocationName)
	}
	return uuid.New(), nil
}
func (f *fakeStore) CreateTORLocationAssignment(context.Context, db.DBTX, uuid.UUID, uuid.UUID, int) (uuid.UUID, string, error) {
	f.call("createTORLocationAssignment")
	return uuid.New(), "สวนลุมพินี", f.assignmentErr
}
func (f *fakeStore) NextRequisitionSeq(context.Context, db.DBTX, string) (int, error) {
	f.call("nextRequisitionSeq")
	f.nextSeq++
	return f.nextSeq, nil
}
func (f *fakeStore) CreateInitialRequisition(_ context.Context, _ db.DBTX, _ uuid.UUID, _ uuid.UUID, no string) (uuid.UUID, error) {
	f.call("createInitialRequisition")
	f.requisitionNos = append(f.requisitionNos, no)
	if len(f.requisitionErrs) > 0 {
		err := f.requisitionErrs[0]
		f.requisitionErrs = f.requisitionErrs[1:]
		if err != nil {
			return uuid.Nil, err
		}
	}
	return uuid.New(), nil
}
func (f *fakeStore) CreateEquipment(_ context.Context, _ db.DBTX, name string) error {
	f.call("createEquipment")
	f.equipmentCreated = append(f.equipmentCreated, name)
	return nil
}
func (f *fakeStore) FindEquipmentByName(context.Context, db.DBTX, string) (uuid.UUID, error) {
	f.call("findEquipmentByName")
	return uuid.New(), nil
}
func (f *fakeStore) CreateRequisitionItem(context.Context, db.DBTX, uuid.UUID, uuid.UUID, string, int) (uuid.UUID, error) {
	f.call("createRequisitionItem")
	return uuid.New(), nil
}
func (f *fakeStore) SearchLocations(_ context.Context, _ db.DBTX, query string, limit int) ([]Location, error) {
	f.call(fmt.Sprintf("searchLocations(%q,%d)", query, limit))
	return nil, nil
}
func (f *fakeStore) FindActiveAssistantLineIDs(context.Context, db.DBTX) ([]string, error) {
	return f.assistants, nil
}

func pgUnique(constraint string) error {
	return fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505", ConstraintName: constraint})
}

type failingNotifier struct{ calls int }

func (n *failingNotifier) Send(context.Context, notify.Message) (notify.Status, error) {
	n.calls++
	return notify.StatusFailed, errors.New("LINE unavailable")
}

type testEnv struct {
	router   *gin.Engine
	store    *fakeStore
	txCount  int
	notifier notify.Notifier
}

func newEnv(store *fakeStore, notifier notify.Notifier) *testEnv {
	env := &testEnv{store: store, notifier: notifier}
	svc := NewService(store, nil, func(_ context.Context, fn func(db.DBTX) error) error {
		env.txCount++
		return fn(nil)
	}, notifier)
	svc.now = func() time.Time { return time.Date(2026, 10, 4, 18, 0, 0, 0, time.UTC) } // 2026-10-05 01:00 in Bangkok
	svc.runAsync = func(f func()) { f() }

	verifier := verifierFunc(func(_ context.Context, token string) (string, error) { return token, nil })
	resolve := func(_ context.Context, subject string) (auth.Principal, error) {
		switch subject {
		case "supervisor":
			return auth.Principal{UserID: supervisorID, Role: auth.RoleSupervisor}, nil
		case "assistant":
			return auth.Principal{UserID: assistantID, Role: auth.RoleAssistant}, nil
		}
		return auth.Principal{}, auth.ErrUserNotFound
	}
	r := gin.New()
	web := r.Group("/api/web", auth.Authenticate(verifier, resolve, auth.RoleSupervisor, auth.RoleAssistant))
	RegisterRoutes(web, NewHandler(svc))
	env.router = r
	return env
}

type verifierFunc func(ctx context.Context, token string) (string, error)

func (f verifierFunc) Verify(ctx context.Context, token string) (string, error) { return f(ctx, token) }

func (e *testEnv) do(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return w
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) httpx.ErrorDetail {
	t.Helper()
	var body httpx.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("not an error body: %s", w.Body.String())
	}
	return body.Error
}

func TestRoleAccess(t *testing.T) {
	routes := []struct{ method, path string }{
		{http.MethodPost, "/api/web/contracts/info"},
		{http.MethodPost, "/api/web/contracts/scope"},
		{http.MethodPost, "/api/web/contracts/confirm"},
		{http.MethodGet, "/api/web/locations?q=x"},
	}
	for _, rt := range routes {
		t.Run(rt.path, func(t *testing.T) {
			env := newEnv(&fakeStore{}, &notify.Noop{})
			if w := env.do(t, rt.method, rt.path, "", nil); w.Code != http.StatusUnauthorized {
				t.Fatalf("no token: %d", w.Code)
			}
			w := env.do(t, rt.method, rt.path, "assistant", nil)
			if w.Code != http.StatusForbidden || decodeError(t, w).Code != "forbidden" {
				t.Fatalf("assistant: %d %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestSubmitContractInfo(t *testing.T) {
	tests := []struct {
		name       string
		store      fakeStore
		body       any
		wantStatus int
		wantCode   string
		wantFields []string
	}{
		{name: "valid", body: validInfo(), wantStatus: http.StatusOK},
		{name: "invalid fields named", body: ContractInfo{ContractNo: "12"}, wantStatus: http.StatusBadRequest, wantCode: "validation_failed",
			wantFields: []string{"contract_no", "project_name", "partner_agency", "start_date", "end_date", "contract_value"}},
		{name: "duplicate contract no", store: fakeStore{dupContractNo: true}, body: validInfo(), wantStatus: http.StatusConflict,
			wantCode: "duplicate_contract_no", wantFields: []string{"contract_no"}},
		{name: "contract value as JSON number rejected", body: map[string]any{"contract_no": "6700001148", "contract_value": 1500.5},
			wantStatus: http.StatusBadRequest, wantCode: "validation_failed", wantFields: []string{"body"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			env := newEnv(&store, &notify.Noop{})
			w := env.do(t, http.MethodPost, "/api/web/contracts/info", "supervisor", tt.body)
			assertResponse(t, w, tt.wantStatus, tt.wantCode, tt.wantFields)
			if len(store.calls) > 0 && tt.wantCode == "validation_failed" {
				t.Fatalf("database called after a format error: %v", store.calls)
			}
		})
	}
}

func TestSubmitScopeData(t *testing.T) {
	tests := []struct {
		name       string
		store      fakeStore
		body       any
		wantStatus int
		wantCode   string
		wantFields []string
		wantCalls  []string
	}{
		{name: "valid: validateLocation before checkDuplicateLocationName", body: validScope(), wantStatus: http.StatusOK,
			wantCalls: []string{"validateLocation", "checkDuplicateLocationName"}},
		{name: "existing location missing", store: fakeStore{missingLocation: true}, body: validScope(),
			wantStatus: http.StatusBadRequest, wantCode: "validation_failed", wantFields: []string{"areas[1].location_id"},
			wantCalls: []string{"validateLocation"}},
		{name: "new location name taken", store: fakeStore{dupLocationName: true}, body: validScope(),
			wantStatus: http.StatusConflict, wantCode: "duplicate_location_name", wantFields: []string{"areas[0].new_location.name"},
			wantCalls: []string{"validateLocation", "checkDuplicateLocationName"}},
		{name: "format error stops before database", body: Scope{}, wantStatus: http.StatusBadRequest,
			wantCode: "validation_failed", wantFields: []string{"areas"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			env := newEnv(&store, &notify.Noop{})
			w := env.do(t, http.MethodPost, "/api/web/contracts/scope", "supervisor", tt.body)
			assertResponse(t, w, tt.wantStatus, tt.wantCode, tt.wantFields)
			if !reflect.DeepEqual(store.calls, tt.wantCalls) && !(len(store.calls) == 0 && len(tt.wantCalls) == 0) {
				t.Fatalf("calls = %v, want %v", store.calls, tt.wantCalls)
			}
		})
	}
}

func confirmBody() ConfirmRequest {
	return ConfirmRequest{Contract: validInfo(), Scope: validScope()}
}

func TestConfirmContractSuccess(t *testing.T) {
	store := &fakeStore{assistants: []string{"U1", "U2"}}
	noop := &notify.Noop{}
	env := newEnv(store, noop)

	w := env.do(t, http.MethodPost, "/api/web/contracts/confirm", "supervisor", confirmBody())
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var out ConfirmedContract
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Status != "registered" || out.ContractFileURL != nil || len(out.Areas) != 2 {
		t.Fatalf("unexpected body: %+v", out)
	}
	// One requisition per area, numbered per Bangkok day (UTC 18:00 on 4 Oct = 5 Oct in Bangkok).
	if got := []string{out.Areas[0].Requisition.RequisitionNo, out.Areas[1].Requisition.RequisitionNo}; !reflect.DeepEqual(got, []string{"REQ-20261005-001", "REQ-20261005-002"}) {
		t.Fatalf("requisition numbers = %v", got)
	}
	if out.Areas[0].Requisition.Status != "pending_survey" || out.Areas[0].Requisition.RequisitionType != "tor_base" || len(out.Areas[0].Requisition.Items) != 2 {
		t.Fatalf("area 0 requisition = %+v", out.Areas[0].Requisition)
	}
	wantOrder := []string{
		"createContract", "nextRequisitionSeq",
		"createLocation", "createTORLocationAssignment", "createInitialRequisition",
		"createEquipment", "findEquipmentByName", "createRequisitionItem",
		"createEquipment", "findEquipmentByName", "createRequisitionItem",
		"createTORLocationAssignment", "createInitialRequisition",
		"createEquipment", "findEquipmentByName", "createRequisitionItem",
	}
	if !reflect.DeepEqual(store.calls, wantOrder) {
		t.Fatalf("calls = %v", store.calls)
	}
	if env.txCount != 1 {
		t.Fatalf("transactions = %d", env.txCount)
	}
	if msgs := noop.Messages(); len(msgs) != 2 || msgs[0].LineUserID != "U1" {
		t.Fatalf("notifications = %+v", msgs)
	}
}

func TestConfirmContractErrors(t *testing.T) {
	tests := []struct {
		name       string
		store      fakeStore
		body       ConfirmRequest
		wantStatus int
		wantCode   string
		wantFields []string
	}{
		{name: "format errors prefixed by body section",
			body: ConfirmRequest{Contract: ContractInfo{ContractNo: "1"}, Scope: Scope{}}, wantStatus: http.StatusBadRequest, wantCode: "validation_failed",
			wantFields: []string{"contract.contract_no", "contract.project_name", "contract.partner_agency", "contract.start_date", "contract.end_date", "contract.contract_value", "scope.areas"}},
		{name: "duplicate contract no from unique key", store: fakeStore{createContractErr: pgUnique(constraintContractNo)}, body: confirmBody(),
			wantStatus: http.StatusConflict, wantCode: "duplicate_contract_no", wantFields: []string{"contract.contract_no"}},
		{name: "duplicate location name from unique key", store: fakeStore{dupLocationName: true}, body: confirmBody(),
			wantStatus: http.StatusConflict, wantCode: "duplicate_location_name", wantFields: []string{"scope.areas[0].new_location.name"}},
		{name: "missing location from foreign key", store: fakeStore{assignmentErr: fmt.Errorf("x: %w", &pgconn.PgError{Code: "23503", ConstraintName: constraintAssignmentLoc})},
			body: confirmBody(), wantStatus: http.StatusBadRequest, wantCode: "validation_failed", wantFields: []string{"scope.areas[0].location_id"}},
		{name: "unknown database error hidden", store: fakeStore{createContractErr: errors.New("connection reset")}, body: confirmBody(),
			wantStatus: http.StatusInternalServerError, wantCode: "internal"},
		{name: "requisition number busy after retries", store: fakeStore{requisitionErrs: []error{
			pgUnique(constraintRequisitionNo), pgUnique(constraintRequisitionNo), pgUnique(constraintRequisitionNo)}},
			body: confirmBody(), wantStatus: http.StatusConflict, wantCode: "requisition_no_busy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			noop := &notify.Noop{}
			env := newEnv(&store, noop)
			w := env.do(t, http.MethodPost, "/api/web/contracts/confirm", "supervisor", tt.body)
			assertResponse(t, w, tt.wantStatus, tt.wantCode, tt.wantFields)
			if len(noop.Messages()) != 0 {
				t.Fatal("notification sent for a failed confirm")
			}
		})
	}
}

func TestConfirmContractRetriesTakenRequisitionNo(t *testing.T) {
	store := &fakeStore{requisitionErrs: []error{pgUnique(constraintRequisitionNo)}}
	env := newEnv(store, &notify.Noop{})
	w := env.do(t, http.MethodPost, "/api/web/contracts/confirm", "supervisor", confirmBody())
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	if env.txCount != 2 {
		t.Fatalf("transactions = %d, want 2 (one retry)", env.txCount)
	}
	// The retry re-reads the next number inside the new transaction.
	if store.requisitionNos[0] != "REQ-20261005-001" || store.requisitionNos[1] != "REQ-20261005-002" {
		t.Fatalf("numbers tried = %v", store.requisitionNos)
	}
}

func TestConfirmContractNotificationFailureKeepsContract(t *testing.T) {
	notifier := &failingNotifier{}
	env := newEnv(&fakeStore{assistants: []string{"U1"}}, notifier)
	w := env.do(t, http.MethodPost, "/api/web/contracts/confirm", "supervisor", confirmBody())
	if w.Code != http.StatusCreated || notifier.calls != 1 {
		t.Fatalf("status %d, sends %d", w.Code, notifier.calls)
	}
}

func TestSearchLocations(t *testing.T) {
	store := &fakeStore{}
	env := newEnv(store, &notify.Noop{})
	w := env.do(t, http.MethodGet, "/api/web/locations?q=%E0%B8%AA%E0%B8%A7%E0%B8%99", "supervisor", nil)
	if w.Code != http.StatusOK || w.Body.String() != `{"data":[]}` {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if !reflect.DeepEqual(store.calls, []string{`searchLocations("สวน",20)`}) {
		t.Fatalf("calls = %v", store.calls)
	}

	long := make([]byte, 0, 300)
	for range 256 {
		long = append(long, 'a')
	}
	w = env.do(t, http.MethodGet, "/api/web/locations?q="+string(long), "supervisor", nil)
	assertResponse(t, w, http.StatusBadRequest, "validation_failed", []string{"q"})
}

func assertResponse(t *testing.T, w *httptest.ResponseRecorder, status int, code string, fields []string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status = %d, want %d: %s", w.Code, status, w.Body.String())
	}
	if code == "" {
		return
	}
	e := decodeError(t, w)
	if e.Code != code {
		t.Fatalf("code = %q, want %q", e.Code, code)
	}
	got := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		got[i] = f.Field
	}
	assertFields(t, got, fields)
}

func fieldNames(errs []apperr.FieldError) []string {
	out := make([]string, len(errs))
	for i, e := range errs {
		out[i] = e.Field
	}
	return out
}

func assertFields(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fields = %v, want %v", got, want)
	}
}
