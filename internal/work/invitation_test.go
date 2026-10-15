package work

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type invitationVerifier struct{}

func (invitationVerifier) Verify(_ context.Context, token string) (string, error) {
	if token != "valid" {
		return "", fmt.Errorf("invalid token")
	}
	return "U" + strings.Repeat("a", 32), nil
}

func TestInvitationTokenExpiryAndStorage(t *testing.T) {
	now := time.Date(2026, 10, 14, 12, 0, 0, 0, time.UTC)
	token, digest, expires, err := newInvitation(now)
	if err != nil {
		t.Fatal(err)
	}
	if len(digest) > 50 || len(token) != 54 || !strings.HasPrefix(digest, invitationPrefix) || strings.Contains(digest, token) {
		t.Fatal("invalid digest storage")
	}
	if !expires.Equal(now.Add(24 * time.Hour)) {
		t.Fatal("incorrect expiry")
	}
	if _, err = invitationExpiry(token, now); err != nil {
		t.Fatal(err)
	}
	if _, err = invitationExpiry(token, expires); err == nil {
		t.Fatal("expired token accepted")
	}
	for _, bad := range []string{"", "x", token + "=", strings.Repeat("x", 54)} {
		if _, err = invitationExpiry(bad, now); err == nil {
			t.Fatal("invalid token accepted")
		}
	}
	payload, _ := base64.RawURLEncoding.DecodeString(token)
	payload[20] ^= 1
	if invitationDigest(base64.RawURLEncoding.EncodeToString(payload)) == digest {
		t.Fatal("tampered token retained identity")
	}
}
func TestInvitationHTTPRejectsForgedIdentityAndRole(t *testing.T) {
	binding.EnableDecoderDisallowUnknownFields = true
	t.Cleanup(func() { binding.EnableDecoderDisallowUnknownFields = false })
	r := gin.New()
	RegisterOnboarding(r, invitationVerifier{}, invitationVerifier{}, New(nil, ""))
	for _, tc := range []struct {
		token, body string
		status      int
	}{{"", `{}`, 401}, {"invalid", `{}`, 401}, {"valid", `{"token":"x","role":"supervisor"}`, 400}, {"valid", `{"token":"x","line_id":"forged"}`, 400}} {
		req := httptest.NewRequest(http.MethodPost, "/api/liff/invitation/claim", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.status {
			t.Fatalf("got %d want %d", w.Code, tc.status)
		}
	}
}
func TestInvitationCreationAndRenewalRequireSupervisor(t *testing.T) {
	s := New(nil, "")
	for _, role := range []auth.Role{auth.RoleWorker, auth.RoleAssistant} {
		p := auth.Principal{UserID: uuid.New(), Role: role}
		if _, err := s.CreateInvitation(context.Background(), p, InvitationUserInput{}); err == nil {
			t.Fatal("creation accepted")
		}
		if _, err := s.RenewInvitation(context.Background(), p, uuid.NewString()); err == nil {
			t.Fatal("renewal accepted")
		}
	}
}

func TestIntegrationInvitationUsesExistingSchema(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	create := func(role string) gin.H {
		t.Helper()
		id := uuid.New()
		in := InvitationUserInput{"Invited", "Person", role, fmt.Sprintf("0%09d", uint64(id.ID())%1000000000), "TEST", "0", "450.50"}
		result, err := s.CreateInvitation(ctx, f.supervisor, in)
		if err != nil {
			t.Fatal(err)
		}
		return result.(gin.H)
	}
	worker := create("worker")
	id := worker["user_id"].(string)
	token := worker["token"].(string)
	var active bool
	var marker string
	var count int
	if err := f.pool.QueryRow(ctx, `SELECT line_id,is_active FROM public."USER" WHERE user_id=$1`, id).Scan(&marker, &active); err != nil {
		t.Fatal(err)
	}
	if active || !strings.HasPrefix(marker, invitationPrefix) || len(marker) > 50 {
		t.Fatal("pending user not safely stored")
	}
	f.pool.QueryRow(ctx, `SELECT count(*) FROM worker WHERE user_id=$1`, id).Scan(&count)
	if count != 0 {
		t.Fatal("pending Worker has a mapping")
	}
	if err := s.UpdateUser(ctx, id, UserUpdate{FirstName: "Invited", LastName: "Person", DailyWage: "450.50", Active: boolPtr(true)}); err == nil {
		t.Fatal("manual activation bypassed LINE binding")
	}
	if _, err := s.PreviewInvitation(ctx, "web", token); err == nil {
		t.Fatal("wrong platform accepted")
	}
	subject := "U" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := s.ClaimInvitation(ctx, subject, "web", token); err == nil {
		t.Fatal("wrong channel claim accepted")
	}
	renewed, err := s.RenewInvitation(ctx, f.supervisor, id)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.ClaimInvitation(ctx, subject, "worker", token); err == nil {
		t.Fatal("replaced invitation accepted")
	}
	token = renewed.(gin.H)["token"].(string)
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.ClaimInvitation(ctx, "U"+strings.ReplaceAll(uuid.NewString(), "-", ""), "worker", token)
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("single-use claim produced %d successes", success)
	}
	f.pool.QueryRow(ctx, `SELECT line_id,is_active FROM public."USER" WHERE user_id=$1`, id).Scan(&subject, &active)
	f.pool.QueryRow(ctx, `SELECT count(*) FROM worker WHERE user_id=$1`, id).Scan(&count)
	if !active || count != 1 || !strings.HasPrefix(subject, "U") {
		t.Fatal("incorrect activation or Worker mapping")
	}
	if _, err = s.ClaimInvitation(ctx, subject, "worker", token); err == nil {
		t.Fatal("replay accepted")
	}
	if _, err = s.RenewInvitation(ctx, f.supervisor, id); err == nil {
		t.Fatal("linked user invitation renewed")
	}
	other := create("worker")
	if _, err = s.ClaimInvitation(ctx, subject, "worker", other["token"].(string)); err == nil {
		t.Fatal("LINE identity reused for another account")
	}
	f.pool.Exec(ctx, `UPDATE public."USER" SET is_active=false WHERE user_id=$1`, id)
	state, err := s.OnboardingState(ctx, subject)
	if err != nil || state.(gin.H)["status"] != "inactive" {
		t.Fatal("suspended account not distinguished")
	}
	for _, role := range []string{"assistant", "supervisor"} {
		invited := create(role)
		line := "U" + strings.ReplaceAll(uuid.NewString(), "-", "")
		if _, err = s.ClaimInvitation(ctx, line, "web", invited["token"].(string)); err != nil {
			t.Fatal(err)
		}
		f.pool.QueryRow(ctx, `SELECT count(*) FROM worker WHERE user_id=$1`, invited["user_id"]).Scan(&count)
		if count != 0 {
			t.Fatal("staff gained Worker mapping")
		}
		user, err := auth.NewUserStore(f.pool).ByLineID(ctx, line)
		if err != nil || string(user.Role) != role {
			t.Fatal("server role not retained")
		}
	}
	data, err := s.Users(ctx, 100, 0)
	if err != nil {
		t.Fatal(err)
	}
	var rows []map[string]any
	if err = json.Unmarshal(data, &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if _, ok := row["line_id"]; ok {
			t.Fatal("invitation digest leaked in list")
		}
	}
}

func boolPtr(v bool) *bool { return &v }

func TestIntegrationCancelInvitationAndReissue(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	create := func() gin.H {
		t.Helper()
		id := uuid.New()
		out, err := s.CreateInvitation(ctx, f.supervisor, InvitationUserInput{"Cancel", "Invite", "worker", fmt.Sprintf("0%09d", uint64(id.ID())%1000000000), "TEST", "0", "400"})
		if err != nil {
			t.Fatal(err)
		}
		return out.(gin.H)
	}
	invited := create()
	id, token := invited["user_id"].(string), invited["token"].(string)
	if _, err := s.CancelInvitation(ctx, f.assistant, id); err == nil {
		t.Fatal("Assistant cancelled invitation")
	}
	if _, err := s.CancelInvitation(ctx, f.worker, id); err == nil {
		t.Fatal("Worker cancelled invitation")
	}
	if _, err := s.CancelInvitation(ctx, f.supervisor, id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CancelInvitation(ctx, f.supervisor, id); err != nil {
		t.Fatal("repeat cancellation must be idempotent")
	}
	var marker string
	var active bool
	if err := f.pool.QueryRow(ctx, `SELECT line_id,is_active FROM public."USER" WHERE user_id=$1`, id).Scan(&marker, &active); err != nil {
		t.Fatal(err)
	}
	if active || !strings.HasPrefix(marker, cancelledInvitationPrefix) || len(marker) > 50 {
		t.Fatal("incorrect cancellation storage")
	}
	line := "U" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := s.PreviewInvitation(ctx, "worker", token); err == nil {
		t.Fatal("cancelled link preview accepted")
	}
	if _, err := s.ClaimInvitation(ctx, line, "worker", token); err == nil {
		t.Fatal("cancelled link claim accepted")
	}
	if err := s.UpdateUser(ctx, id, UserUpdate{FirstName: "Cancel", LastName: "Invite", DailyWage: "400", Active: boolPtr(true)}); err == nil {
		t.Fatal("cancelled account manually activated")
	}
	fresh, err := s.RenewInvitation(ctx, f.supervisor, id)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.ClaimInvitation(ctx, line, "worker", fresh.(gin.H)["token"].(string)); err != nil {
		t.Fatal(err)
	}
	if _, err = s.CancelInvitation(ctx, f.supervisor, id); err == nil {
		t.Fatal("linked account cancelled")
	}
	if err = f.pool.QueryRow(ctx, `SELECT is_active FROM public."USER" WHERE user_id=$1`, id).Scan(&active); err != nil || !active {
		t.Fatal("linked account deactivated")
	}

	raced := create()
	raceID, raceToken := raced["user_id"].(string), raced["token"].(string)
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, err := s.CancelInvitation(ctx, f.supervisor, raceID)
		results <- err
	}()
	go func() {
		defer wg.Done()
		<-start
		_, err := s.ClaimInvitation(ctx, "U"+strings.ReplaceAll(uuid.NewString(), "-", ""), "worker", raceToken)
		results <- err
	}()
	close(start)
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatal("cancel and claim both succeeded or neither succeeded")
	}
	f.pool.QueryRow(ctx, `SELECT line_id,is_active FROM public."USER" WHERE user_id=$1`, raceID).Scan(&marker, &active)
	if active && strings.HasPrefix(marker, cancelledInvitationPrefix) {
		t.Fatal("cancellation corrupted active account")
	}
	if !active && !strings.HasPrefix(marker, cancelledInvitationPrefix) {
		t.Fatal("unclaimed invitation not cancelled")
	}
}
