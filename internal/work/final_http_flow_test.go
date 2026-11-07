package work

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/server"
)

// Only the external LINE verifier and object-storage verifier are substitutes.
// HTTP routing, DB identity/roles, transactions and reports use real code/Postgres.
type finalFlowVerifier struct{ subjects map[string]string }

func (v finalFlowVerifier) Verify(_ context.Context, token string) (string, error) {
	if subject, ok := v.subjects[token]; ok {
		return subject, nil
	}
	return "", fmt.Errorf("invalid fixture token")
}
func TestIntegrationFinalThreeRoleHTTPFlow(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	users := auth.NewUserStore(f.pool)
	v := finalFlowVerifier{map[string]string{"supervisor": f.supervisor.UserID.String(), "assistant": f.assistant.UserID.String(), "worker": f.worker.UserID.String(), "other": f.other.UserID.String()}}
	r := server.New(server.Deps{DB: f.pool, WebAuth: auth.Authenticate(v, users.ByLineID, auth.RoleSupervisor, auth.RoleAssistant), WorkerAuth: auth.Authenticate(v, users.ByLineID, auth.RoleWorker)})
	Register(r.Web, r.Worker, s)
	call := func(role, method, path string, body any, want int) *httptest.ResponseRecorder {
		t.Helper()
		b, _ := json.Marshal(body)
		request := httptest.NewRequest(method, path, bytes.NewReader(b))
		request.Header.Set("Authorization", "Bearer "+role)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		r.Engine.ServeHTTP(response, request)
		if response.Code != want {
			t.Fatalf("%s %s as %s: got %d want %d: %s", method, path, role, response.Code, want, response.Body.String())
		}
		t.Logf("%s %s as %s -> %d", method, path, role, want)
		return response
	}
	decode := func(response *httptest.ResponseRecorder) map[string]json.RawMessage {
		t.Helper()
		var out map[string]json.RawMessage
		if e := json.Unmarshal(response.Body.Bytes(), &out); e != nil {
			t.Fatal(e)
		}
		return out
	}
	value := func(response *httptest.ResponseRecorder, key string) string {
		t.Helper()
		var out string
		if e := json.Unmarshal(decode(response)[key], &out); e != nil {
			t.Fatal(e)
		}
		return out
	}
	for _, role := range []string{"supervisor", "assistant", "worker"} {
		path := "/api/web/me"
		if role == "worker" {
			path = "/api/liff/me"
		}
		call(role, "GET", path, nil, 200)
	}
	call("worker", "GET", "/api/web/me", nil, 403)
	call("assistant", "GET", "/api/web/reports/profit?period_start=2026-10-01&period_end=2026-10-31", nil, 403)
	call("supervisor", "POST", "/api/web/schedules", ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"}, 403)
	response := call("assistant", "POST", "/api/web/schedules", ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"}, 200)
	var ids []string
	if e := json.Unmarshal(decode(response)["schedule_ids"], &ids); e != nil || len(ids) != 1 {
		t.Fatal("missing schedule", e)
	}
	id := ids[0]
	qr := value(call("supervisor", "POST", "/api/web/assignments/"+f.assignment+"/qr", nil, 200), "qr_token")
	*f.now = time.Date(2026, 10, 2, 8, 0, 0, 0, Bangkok)
	lat, lon, accuracy := 13.75, 100.5, 10.0
	check := CheckIn{id, &lat, &lon, &accuracy, qr}
	call("other", "POST", "/api/liff/attendance/check-in", check, 404)
	call("worker", "POST", "/api/liff/attendance/check-in", check, 200)
	call("worker", "POST", "/api/liff/attendance/check-in", check, 409)
	out := CheckOut{id, "HTTP E2E completed work", f.worker.UserID.String() + "/test-photo.png"}
	call("worker", "POST", "/api/liff/attendance/check-out", out, 409)
	*f.now = time.Date(2026, 10, 2, 16, 0, 0, 0, Bangkok)
	call("worker", "POST", "/api/liff/attendance/check-out", out, 200)
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	payroll := value(call("supervisor", "POST", "/api/web/payroll", PayrollInput{f.worker.UserID.String(), "2026-10-01", "2026-10-15"}, 200), "payroll_id")
	call("supervisor", "POST", "/api/web/payroll/"+payroll+"/mark-paid", nil, 200)
	call("worker", "GET", "/api/liff/payroll", nil, 200)
	call("assistant", "POST", "/api/web/invoices", InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "900"}, 403)
	call("supervisor", "POST", "/api/web/invoices", InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "1001"}, 400)
	invoice := value(call("supervisor", "POST", "/api/web/invoices", InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "900"}, 200), "invoice_id")
	call("supervisor", "POST", "/api/web/invoices/"+invoice+"/mark-paid", map[string]string{"net_received": "900"}, 200)
	query := "?tor_id=" + f.tor + "&period_start=2026-10-01&period_end=2026-10-31"
	report := call("supervisor", "GET", "/api/web/reports/profit"+query, nil, 200)
	var rows []map[string]any
	if e := json.Unmarshal(decode(report)["data"], &rows); e != nil || len(rows) != 1 {
		t.Fatal("missing report", e)
	}
	if rows[0]["revenue"] != "900.00" || rows[0]["labor"] != "400.00" || rows[0]["profit"] != "500.00" {
		t.Fatalf("report mismatch: %+v", rows[0])
	}
	csvResponse := call("supervisor", "GET", "/api/web/reports/profit.csv"+query, nil, 200)
	records, e := csv.NewReader(strings.NewReader(csvResponse.Body.String())).ReadAll()
	if e != nil || len(records) != 2 || records[1][4] != "900.00" || records[1][5] != "400.00" || records[1][7] != "500.00" {
		t.Fatal("CSV mismatch", e, records)
	}
	pdf := call("supervisor", "GET", "/api/web/reports/profit.pdf"+query, nil, 200)
	if !bytes.HasPrefix(pdf.Body.Bytes(), []byte("%PDF-")) {
		t.Fatal("invalid PDF header")
	}
	if dir := os.Getenv("TEST_ARTIFACT_DIR"); dir != "" {
		if e = os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
		for name, content := range map[string][]byte{"three-role-profit.pdf": pdf.Body.Bytes(), "three-role-profit.csv": csvResponse.Body.Bytes(), "three-role-profit.json": report.Body.Bytes()} {
			if e = os.WriteFile(filepath.Join(dir, name), content, 0600); e != nil {
				t.Fatal(e)
			}
		}
	}
}
