package work

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type qrTestVerifier struct{}

func (qrTestVerifier) Verify(context.Context, string) (string, error) { return "test-subject", nil }

func TestQRIssuanceRejectsAssistant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	web := r.Group("/api/web", auth.Authenticate(qrTestVerifier{}, func(context.Context, string) (auth.Principal, error) {
		return auth.Principal{UserID: uuid.New(), Role: auth.RoleAssistant}, nil
	}, auth.RoleSupervisor, auth.RoleAssistant))
	// Nil database deliberately catches a role guard that allows the request
	// through to issuance before authorization.
	Register(web, r.Group("/api/liff"), &Service{})
	req := httptest.NewRequest(http.MethodPost, "/api/web/assignments/"+uuid.NewString()+"/qr", nil)
	req.Header.Set("Authorization", "Bearer test-only")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	if response.Code != http.StatusForbidden {
		t.Fatalf("assistant issuance returned %d", response.Code)
	}
}

func TestIntegrationPermanentQRRechecksContractAtCheckIn(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	shifts, err := f.service.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"})
	if err != nil {
		t.Fatal(err)
	}
	qr, err := f.service.QR(ctx, f.assignment)
	if err != nil {
		t.Fatal(err)
	}
	if qr["permanent"] != true || qr["expires_at"] != nil {
		t.Fatal("permanent response contract incorrect")
	}
	*f.now = time.Date(2026, 10, 2, 8, 0, 0, 0, Bangkok)
	reopened, err := f.service.QR(ctx, f.assignment)
	if err != nil || reopened["qr_token"] != qr["qr_token"] {
		t.Fatal("printed token changed overnight", err)
	}
	lat, lon, accuracy := 13.75, 100.5, 10.0
	in := CheckIn{shifts[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}
	if _, err := f.pool.Exec(ctx, `UPDATE contract_tor SET status='complete' WHERE tor_id=$1`, f.tor); err != nil {
		t.Fatal(err)
	}
	if err := f.service.CheckIn(ctx, f.worker, in); err == nil {
		t.Fatal("old sign accepted for completed contract")
	}
	var count int
	if err := f.pool.QueryRow(ctx, `SELECT count(*) FROM attendance WHERE schedule_id=$1`, shifts[0]).Scan(&count); err != nil || count != 0 {
		t.Fatal("blocked check-in wrote attendance", err)
	}
	if _, err := f.pool.Exec(ctx, `UPDATE contract_tor SET status='active' WHERE tor_id=$1`, f.tor); err != nil {
		t.Fatal(err)
	}
	if err := f.service.CheckIn(ctx, f.worker, in); err != nil {
		t.Fatal("valid old printed code failed", err)
	}
	if err := f.service.CheckIn(ctx, f.worker, in); err == nil {
		t.Fatal("duplicate check-in accepted")
	}
}
