package work

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/server"
	"github.com/google/uuid"
)

type loginFixtureVerifier struct{ subject string }

func (v loginFixtureVerifier) Verify(_ context.Context, token string) (string, error) {
	if token != "fixture-token" {
		return "", errors.New("invalid fixture token")
	}
	return v.subject, nil
}

// Exercise the actual /me handlers, DB role resolution and Worker mapping gate.
func TestIntegrationLoginMeAndWorkerMapping(t *testing.T) {
	pool := isolatedPool(t)
	ctx := context.Background()
	id, subject := uuid.New(), "U"+uuid.NewString()
	_, err := pool.Exec(ctx, `INSERT INTO public."USER"(user_id,first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no,is_active)
		VALUES($1,'Login','Fixture','supervisor',$2,$3,'TEST','TEST',true)`, id, fmt.Sprintf("0%09d", time.Now().UnixNano()%1000000000), subject)
	if err != nil {
		t.Fatal("cannot insert isolated fixture")
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM public."USER" WHERE user_id=$1`, id) })
	users := auth.NewUserStore(pool)
	v := loginFixtureVerifier{subject}
	r := server.New(server.Deps{DB: pool,
		WebAuth:    auth.Authenticate(v, users.ByLineID, auth.RoleSupervisor, auth.RoleAssistant),
		WorkerAuth: auth.Authenticate(v, users.ByLineID, auth.RoleWorker),
	})
	Register(r.Web, r.Worker, New(pool, ""))
	check := func(path string, want int, role string) {
		t.Helper()
		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.Header.Set("Authorization", "Bearer fixture-token")
		response := httptest.NewRecorder()
		r.Engine.ServeHTTP(response, request)
		if response.Code != want {
			t.Fatalf("status=%d want=%d", response.Code, want)
		}
		if want == 200 {
			var result struct {
				UserID string `json:"user_id"`
				Role   string `json:"role"`
				Active bool   `json:"is_active"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil || result.UserID != id.String() || result.Role != role || !result.Active {
				t.Fatal("/me did not return this active DB identity and role")
			}
		}
	}
	t.Run("supervisor identity", func(t *testing.T) { check("/api/web/me", 200, "supervisor") })
	t.Run("supervisor denied worker route", func(t *testing.T) { check("/api/liff/me", 403, "") })
	if _, err := pool.Exec(ctx, `UPDATE public."USER" SET role='assistant' WHERE user_id=$1`, id); err != nil {
		t.Fatal("cannot update isolated role")
	}
	t.Run("assistant role read without restart", func(t *testing.T) { check("/api/web/me", 200, "assistant") })
	if _, err := pool.Exec(ctx, `UPDATE public."USER" SET role='worker' WHERE user_id=$1`, id); err != nil {
		t.Fatal("cannot update isolated role")
	}
	t.Run("worker denied web route", func(t *testing.T) { check("/api/web/me", 403, "") })
	t.Run("worker missing mapping", func(t *testing.T) { check("/api/liff/me", 404, "") })
	if _, err := pool.Exec(ctx, `INSERT INTO worker(user_id) VALUES($1)`, id); err != nil {
		t.Fatal("cannot insert isolated Worker")
	}
	t.Run("worker identity", func(t *testing.T) { check("/api/liff/me", 200, "worker") })
	if _, err := pool.Exec(ctx, `INSERT INTO worker(user_id) VALUES($1)`, id); err != nil {
		t.Fatal("cannot insert duplicate isolated Worker")
	}
	t.Run("duplicate Worker rejected", func(t *testing.T) { check("/api/liff/me", 409, "") })
	if _, err := pool.Exec(ctx, `UPDATE public."USER" SET is_active=false WHERE user_id=$1`, id); err != nil {
		t.Fatal("cannot deactivate isolated fixture")
	}
	t.Run("inactive identity rejected", func(t *testing.T) { check("/api/liff/me", 401, "") })
}
