package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLINEBoundedRetriesReuseKey(t *testing.T) {
	for _, tc := range []struct {
		name     string
		codes    []int
		accepted bool
	}{
		{"transient recovery", []int{500, 200}, true},
		{"already accepted", []int{500, 409}, true},
		{"permanent rejection", []int{401}, false},
		{"bounded rate limit", []int{429, 429, 429}, false},
		{"plain conflict", []int{409}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			count := 0
			var key, payload string
			endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				if count == 0 {
					key = r.Header.Get("X-Line-Retry-Key")
					payload = string(body)
				}
				if r.Header.Get("X-Line-Retry-Key") != key || key == "" || string(body) != payload {
					t.Error("retry key or payload changed")
				}
				code := tc.codes[count]
				count++
				if tc.accepted && code == 409 {
					w.Header().Set("X-Line-Accepted-Request-Id", "fixture-id")
				}
				w.WriteHeader(code)
			}))
			defer endpoint.Close()
			client := endpoint.Client()
			client.Transport = redirectTransport{client.Transport, endpoint.URL}
			sender := LinePush{Token: "fixture-token", Client: client}
			err := sender.Send(context.Background(), "fixture-recipient", "notice")
			if (err == nil) != tc.accepted || count != len(tc.codes) {
				t.Fatal("unexpected retries", count, err)
			}
		})
	}
}

func TestIntegrationContinuationAndNotificationHTTP(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	t.Run("selected assignment permits active and exposes ended notice", func(t *testing.T) {
		for _, state := range []string{"active", "complete"} {
			if _, e := f.pool.Exec(ctx, `UPDATE contract_tor SET status=$2::text::contract_status_enum WHERE tor_id=$1`, f.tor, state); e != nil {
				t.Fatal(e)
			}
			data, e := s.AssignmentContinuation(ctx, f.assignment)
			if e != nil {
				t.Fatal(e)
			}
			var out struct {
				Can        bool   `json:"can_continue"`
				Notice     string `json:"notice"`
				History    bool   `json:"summary_notification_history_available"`
				Assignment string `json:"assignment_id"`
			}
			if e = json.Unmarshal(data, &out); e != nil || out.Can != (state == "active") || out.Assignment != f.assignment || out.History || (state == "complete" && out.Notice == "") {
				t.Fatal(string(data), e)
			}
		}
	})
	var req string
	e := f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status,reviewed_by,reviewed_at) VALUES($1,$2,$3,'original','additional','rejected',$4,now()) RETURNING requisition_id::text`, uuid.NewString(), f.worker.UserID, f.assignment, f.assistant.UserID).Scan(&req)
	if e != nil {
		t.Fatal(e)
	}
	role := auth.RoleSupervisor
	resolve := func(context.Context, string) (auth.Principal, error) {
		return auth.Principal{UserID: f.assistant.UserID, Role: role}, nil
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	web := r.Group("/api/web", auth.Authenticate(testVerifier{}, resolve, auth.RoleSupervisor, auth.RoleAssistant))
	worker := r.Group("/api/liff", auth.Authenticate(testVerifier{}, resolve, auth.RoleWorker))
	Register(web, worker, s)
	t.Run("supervisor cannot retry assistant notices", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/api/web/notifications/"+req+"/retry", strings.NewReader(`{"kind":"no_purchase","reason":"reuse"}`))
		request.Header.Set("Authorization", "Bearer fixture")
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		if w.Code != 403 {
			t.Fatal(w.Code)
		}
	})
	role = auth.RoleAssistant
	t.Run("assistant gets disabled notification without business writes", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/api/web/notifications/"+req+"/retry", strings.NewReader(`{"kind":"no_purchase","reason":"reuse"}`))
		request.Header.Set("Authorization", "Bearer fixture")
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		if w.Code != 200 || !strings.Contains(w.Body.String(), `"status":"disabled"`) {
			t.Fatal(w.Code, w.Body.String())
		}
	})
}

func TestIntegrationInitialAdditionalPurchaseWithoutFunding(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	var req, item string
	e = f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status,reviewed_by,reviewed_at) VALUES($1,$2,$3,'approved','additional','pending_procurement',$4,now()) RETURNING requisition_id::text`, uuid.NewString(), f.worker.UserID, f.assignment, f.supervisor.UserID).Scan(&req)
	if e != nil {
		t.Fatal(e)
	}
	e = f.pool.QueryRow(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty) VALUES($1,$2,2) RETURNING item_id::text`, req, equipment).Scan(&item)
	if e != nil {
		t.Fatal(e)
	}
	zero := 0
	in := PurchaseInput{Items: []PurchaseItem{{item, 0, "0", &zero}}}
	if _, e = s.Purchase(ctx, f.assistant, req, in); e == nil {
		t.Fatal("7A zero quantity accepted")
	}
	in.Items[0].Quantity = 1
	if _, e = s.Purchase(ctx, f.assistant, req, in); e != nil {
		t.Fatal("first approved purchase incorrectly requires 2S", e)
	}
	one := 1
	in.Items[0].Expected = &one
	if _, e = s.Purchase(ctx, f.assistant, req, in); e == nil {
		t.Fatal("partial purchase bypassed funding state")
	}
}
