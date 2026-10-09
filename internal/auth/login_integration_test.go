package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This composes the real LINE verifier, real USER lookup and role middleware.
// Only LINE's HTTP response is simulated; PostgreSQL is isolated and real.
func TestIntegrationLINELoginPipeline(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("isolated PostgreSQL not configured; login integration unverified")
	}
	if os.Getenv("TEST_DATABASE_ISOLATED") != "yes" {
		t.Fatal("isolated test marker required")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || (cfg.ConnConfig.Host != "127.0.0.1" && cfg.ConnConfig.Host != "localhost") {
		t.Fatal("isolated loopback PostgreSQL required")
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal("isolated PostgreSQL unavailable")
	}
	t.Cleanup(pool.Close)
	// Roll back all fixtures, including role/active changes, after this test.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal("cannot begin isolated test transaction")
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	id := uuid.New()
	subject := "U" + uuid.NewString()
	_, err = tx.Exec(ctx, `INSERT INTO public."USER"(user_id,first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no,is_active)
		VALUES($1,'Login','Fixture','supervisor',$2,$3,'TEST','TEST',true)`, id, fmt.Sprintf("0%09d", time.Now().UnixNano()%1000000000), subject)
	if err != nil {
		t.Fatal("cannot insert isolated login fixture")
	}
	line := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, channel := r.FormValue("id_token"), r.FormValue("client_id")
		if (token != "web-fixture" || channel != "web") && (token != "worker-fixture" || channel != "worker") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"sub": subject, "aud": channel, "iss": "https://access.line.me", "exp": time.Now().Add(time.Hour).Unix()})
	}))
	defer line.Close()
	web, worker := NewLineVerifier("web"), NewLineVerifier("worker")
	web.endpoint, worker.endpoint = line.URL, line.URL
	users := NewUserStore(tx)
	router := gin.New()
	identity := func(c *gin.Context) {
		p, ok := PrincipalFrom(c)
		if !ok || p.UserID != id {
			t.Error("identity was not resolved from verified LINE subject")
		}
		c.JSON(200, gin.H{"user_id": p.UserID, "role": p.Role})
	}
	router.GET("/api/web/me", Authenticate(web, users.ByLineID, RoleSupervisor, RoleAssistant), identity)
	router.GET("/api/liff/me", Authenticate(worker, users.ByLineID, RoleWorker), identity)
	for _, tc := range []struct {
		name, role, path, token string
		active, registered      bool
		want                    int
	}{
		{"supervisor web", "supervisor", "/api/web/me", "web-fixture", true, true, 200},
		{"assistant web", "assistant", "/api/web/me", "web-fixture", true, true, 200},
		{"worker channel", "worker", "/api/liff/me", "worker-fixture", true, true, 200},
		{"web token rejected by worker channel", "supervisor", "/api/liff/me", "web-fixture", true, true, 401},
		{"worker token rejected by web channel", "worker", "/api/web/me", "worker-fixture", true, true, 401},
		{"supervisor with valid worker token forbidden", "supervisor", "/api/liff/me", "worker-fixture", true, true, 403},
		{"worker with valid web token forbidden", "worker", "/api/web/me", "web-fixture", true, true, 403},
		{"inactive supervisor", "supervisor", "/api/web/me", "web-fixture", false, true, 401},
		{"unknown LINE subject", "supervisor", "/api/web/me", "web-fixture", true, false, 401},
		{"missing bearer", "supervisor", "/api/web/me", "", true, true, 401},
		{"invalid bearer", "supervisor", "/api/web/me", "invalid-fixture", true, true, 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			storedSubject := subject
			if !tc.registered {
				storedSubject = "unrelated-fixture"
			}
			if _, err := tx.Exec(ctx, `UPDATE public."USER" SET role=$2,is_active=$3,line_id=$4 WHERE user_id=$1`, id, tc.role, tc.active, storedSubject); err != nil {
				t.Fatal("cannot update isolated login fixture")
			}
			request := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.token != "" {
				request.Header.Set("Authorization", "Bearer "+tc.token)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tc.want {
				t.Fatalf("status = %d, want %d", response.Code, tc.want)
			}
			if tc.want == 200 {
				var result struct {
					Role string `json:"role"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil || result.Role != tc.role {
					t.Fatal("response role differs from database")
				}
			}
		})
	}
}
