package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStorageDoesNotAcceptOtherOwnersOrRedirect(t *testing.T) {
	calls := 0
	p := auth.Principal{UserID: uuid.New(), Role: auth.RoleWorker}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing storage authentication")
		}
		http.Redirect(w, r, "https://example.invalid/leak", 302)
	}))
	defer server.Close()
	st := NewStorage(server.URL, "test-key", "private")
	if e := st.Verify(context.Background(), p, uuid.NewString()+"/"+uuid.NewString()); e == nil || calls != 0 {
		t.Fatal("other owner's object was requested")
	}
	if e := st.Verify(context.Background(), p, p.UserID.String()+"/"+uuid.NewString()); e == nil || calls != 1 {
		t.Fatal("redirect accepted")
	}
}
func TestPushUsesAuthenticatedPostAndRejectsFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("incorrect push request")
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"to":"U-test"`) {
			t.Error("recipient missing")
		}
		w.WriteHeader(429)
	}))
	defer server.Close()
	client := server.Client()
	client.Transport = redirectTransport{base: client.Transport, url: server.URL}
	push := LinePush{Token: "test-token", Client: client}
	if e := push.Send(context.Background(), "U-test", "notice"); e == nil {
		t.Fatal("LINE failure reported as sent")
	}
}

type redirectTransport struct {
	base http.RoundTripper
	url  string
}

func (t redirectTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	copy := r.Clone(r.Context())
	target, _ := http.NewRequest(r.Method, t.url, nil)
	copy.URL = target.URL
	return t.base.RoundTrip(copy)
}
