package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http/httptest"
	"strings"
	"testing"
)

type testVerifier struct{}

func (testVerifier) Verify(context.Context, string) (string, error) { return "verified", nil }
func TestMutationRoleBoundaries(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.NewString()
	tests := []struct {
		method, path string
		role         auth.Role
	}{{"POST", "/users", auth.RoleAssistant}, {"PATCH", "/users/" + id, auth.RoleAssistant}, {"POST", "/leave-requests/" + id + "/review", auth.RoleSupervisor}, {"POST", "/payroll", auth.RoleAssistant}, {"POST", "/payroll/" + id + "/mark-paid", auth.RoleAssistant}, {"POST", "/invoices", auth.RoleAssistant}, {"POST", "/requisitions/" + id + "/review", auth.RoleAssistant}, {"POST", "/requisitions/" + id + "/fund-transfers", auth.RoleAssistant}, {"POST", "/schedules", auth.RoleSupervisor}, {"POST", "/leave-requests/" + id + "/replacement", auth.RoleSupervisor}, {"POST", "/requisitions/" + id + "/purchase", auth.RoleSupervisor}}
	tests = append(tests, struct {
		method, path string
		role         auth.Role
	}{"GET", "/workers/available?work_date=2026-10-08", auth.RoleSupervisor}, struct {
		method, path string
		role         auth.Role
	}{"GET", "/leave-requests/" + id + "/candidates", auth.RoleSupervisor}, struct {
		method, path string
		role         auth.Role
	}{"GET", "/requisitions/" + id + "/delivery-schedules", auth.RoleSupervisor})
	tests = append(tests, struct {
		method, path string
		role         auth.Role
	}{"POST", "/requisitions/" + id + "/decision", auth.RoleSupervisor}, struct {
		method, path string
		role         auth.Role
	}{"POST", "/requisitions/" + id + "/delivery", auth.RoleSupervisor}, struct {
		method, path string
		role         auth.Role
	}{"POST", "/payroll/batch", auth.RoleAssistant}, struct {
		method, path string
		role         auth.Role
	}{"POST", "/reports/profit/confirm", auth.RoleAssistant})
	tests = append(tests, struct {
		method, path string
		role         auth.Role
	}{"POST", "/contracts/" + id + "/operations-summary/notify", auth.RoleAssistant}, struct {
		method, path string
		role         auth.Role
	}{"GET", "/assignments/" + id + "/operations-summary?period_start=2026-10-01&period_end=2026-10-15", auth.RoleSupervisor})
	for _, tt := range tests {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			r := gin.New()
			resolve := func(context.Context, string) (auth.Principal, error) {
				return auth.Principal{UserID: uuid.New(), Role: tt.role}, nil
			}
			web := r.Group("/api/web", auth.Authenticate(testVerifier{}, resolve, auth.RoleSupervisor, auth.RoleAssistant))
			worker := r.Group("/api/liff", auth.Authenticate(testVerifier{}, resolve, auth.RoleWorker))
			Register(web, worker, &Service{})
			request := httptest.NewRequest(tt.method, "/api/web"+tt.path, strings.NewReader("{}"))
			request.Header.Set("Authorization", "Bearer test")
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			if w.Code != 403 {
				t.Fatalf("expected 403 before any DB access, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
func TestRoutesRejectMissingIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	resolve := func(context.Context, string) (auth.Principal, error) { return auth.Principal{}, auth.ErrUserNotFound }
	web := r.Group("/api/web", auth.Authenticate(testVerifier{}, resolve, auth.RoleSupervisor, auth.RoleAssistant))
	worker := r.Group("/api/liff", auth.Authenticate(testVerifier{}, resolve, auth.RoleWorker))
	Register(web, worker, &Service{})
	for _, route := range r.Routes() {
		path := strings.ReplaceAll(route.Path, ":id", uuid.NewString())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(route.Method, path, strings.NewReader("{}")))
		if w.Code != 401 {
			t.Errorf("%s %s: %d", route.Method, path, w.Code)
		}
	}
}
