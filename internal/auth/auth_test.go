package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func init() { gin.SetMode(gin.TestMode) }

type verifierFunc func(ctx context.Context, token string) (string, error)

func (f verifierFunc) Verify(ctx context.Context, token string) (string, error) { return f(ctx, token) }

func TestAuthenticate(t *testing.T) {
	assistant := Principal{UserID: uuid.New(), Role: RoleAssistant}
	worker := Principal{UserID: uuid.New(), Role: RoleWorker}

	verifier := verifierFunc(func(_ context.Context, token string) (string, error) {
		if token == "good" {
			return "subject", nil
		}
		return "", errors.New("bad token")
	})

	tests := []struct {
		name    string
		header  string
		resolve Resolver
		want    int
	}{
		{name: "no header", header: "", want: http.StatusUnauthorized},
		{name: "wrong scheme", header: "Basic good", want: http.StatusUnauthorized},
		{name: "invalid token", header: "Bearer bad", want: http.StatusUnauthorized},
		{
			name:    "unknown or inactive user",
			header:  "Bearer good",
			resolve: func(context.Context, string) (Principal, error) { return Principal{}, ErrUserNotFound },
			want:    http.StatusUnauthorized,
		},
		{
			name:    "lookup failure",
			header:  "Bearer good",
			resolve: func(context.Context, string) (Principal, error) { return Principal{}, errors.New("db down") },
			want:    http.StatusInternalServerError,
		},
		{
			name:    "role not allowed",
			header:  "Bearer good",
			resolve: func(context.Context, string) (Principal, error) { return worker, nil },
			want:    http.StatusForbidden,
		},
		{
			name:    "allowed",
			header:  "bearer good",
			resolve: func(context.Context, string) (Principal, error) { return assistant, nil },
			want:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/x", Authenticate(verifier, tt.resolve, RoleSupervisor, RoleAssistant), func(c *gin.Context) {
				p, ok := PrincipalFrom(c)
				if !ok || p != assistant {
					t.Errorf("principal = %+v, %v", p, ok)
				}
				c.Status(http.StatusOK)
			})
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d (%s)", w.Code, tt.want, w.Body.String())
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name      string
		principal *Principal
		want      int
	}{
		{name: "missing Authenticate", want: http.StatusInternalServerError},
		{name: "assistant blocked from supervisor route", principal: &Principal{Role: RoleAssistant}, want: http.StatusForbidden},
		{name: "supervisor allowed", principal: &Principal{Role: RoleSupervisor}, want: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/x", func(c *gin.Context) {
				if tt.principal != nil {
					c.Set(principalKey, *tt.principal)
				}
			}, RequireRole(RoleSupervisor), func(c *gin.Context) { c.Status(http.StatusOK) })
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}

func TestSupabaseVerifier(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	const issuer = "https://ref.supabase.co/auth/v1"
	v := &SupabaseVerifier{
		keyfunc: func(*jwt.Token) (any, error) { return &key.PublicKey, nil },
		issuer:  issuer,
	}
	sign := func(claims jwt.RegisteredClaims) string {
		s, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	future := jwt.NewNumericDate(time.Now().Add(time.Hour))
	past := jwt.NewNumericDate(time.Now().Add(-time.Hour))
	hs256, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: issuer, Subject: "u", ExpiresAt: future,
	}).SignedString([]byte("legacy-secret"))

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{name: "valid", token: sign(jwt.RegisteredClaims{Issuer: issuer, Subject: "user-1", ExpiresAt: future}), want: "user-1"},
		{name: "expired", token: sign(jwt.RegisteredClaims{Issuer: issuer, Subject: "user-1", ExpiresAt: past}), wantErr: true},
		{name: "no expiry", token: sign(jwt.RegisteredClaims{Issuer: issuer, Subject: "user-1"}), wantErr: true},
		{name: "wrong issuer", token: sign(jwt.RegisteredClaims{Issuer: "https://other/auth/v1", Subject: "user-1", ExpiresAt: future}), wantErr: true},
		{name: "no subject", token: sign(jwt.RegisteredClaims{Issuer: issuer, ExpiresAt: future}), wantErr: true},
		{name: "legacy HS256 rejected", token: hs256, wantErr: true},
		{name: "garbage", token: "not-a-jwt", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Verify(context.Background(), tt.token)
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("got %q, %v; want %q, err=%v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestLineVerifier(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    string
		wantErr bool
	}{
		{name: "valid", status: http.StatusOK, body: `{"sub":"U123","aud":"chan"}`, want: "U123"},
		{name: "rejected by LINE", status: http.StatusBadRequest, body: `{"error":"invalid_request","error_description":"IdToken expired."}`, wantErr: true},
		{name: "audience mismatch", status: http.StatusOK, body: `{"sub":"U123","aud":"other"}`, wantErr: true},
		{name: "no subject", status: http.StatusOK, body: `{"aud":"chan"}`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.FormValue("id_token") != "tok" || r.FormValue("client_id") != "chan" {
					t.Errorf("unexpected request: %s id_token=%q client_id=%q", r.Method, r.FormValue("id_token"), r.FormValue("client_id"))
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			v := NewLineVerifier("chan")
			v.endpoint = srv.URL
			got, err := v.Verify(context.Background(), "tok")
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("got %q, %v; want %q, err=%v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

type fakeRow struct {
	id   uuid.UUID
	role string
	err  error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*dest[0].(*uuid.UUID) = r.id
	*dest[1].(*string) = r.role
	return nil
}

type fakeQuerier struct{ row fakeRow }

func (q fakeQuerier) QueryRow(context.Context, string, ...any) pgx.Row { return q.row }

func TestUserStore(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name    string
		row     fakeRow
		subject string
		byLine  bool
		wantErr error
	}{
		{name: "by user id", row: fakeRow{id: id, role: "assistant"}, subject: id.String()},
		{name: "by line id", row: fakeRow{id: id, role: "worker"}, subject: "U123", byLine: true},
		{name: "subject not a UUID", subject: "nope", wantErr: ErrUserNotFound},
		{name: "no row", row: fakeRow{err: pgx.ErrNoRows}, subject: id.String(), wantErr: ErrUserNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserStore(fakeQuerier{row: tt.row})
			lookup := s.ByUserID
			if tt.byLine {
				lookup = s.ByLineID
			}
			p, err := lookup(context.Background(), tt.subject)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (p.UserID != id || string(p.Role) != tt.row.role) {
				t.Fatalf("principal = %+v", p)
			}
		})
	}
}
