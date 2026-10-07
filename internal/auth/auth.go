// Package auth authenticates requests and checks roles.
//
// Web (Supervisor, Assistant): Supabase Auth access token; users.user_id equals
// the Supabase Auth user id (decision 13, see docs/04).
// Worker (LINE Mini App): LIFF ID token verified with LINE; users.line_id holds
// the LINE user id.
//
// Area/assignment permission derived from USER MANAGES CONTRACT_TOR is
// TODO(decision-12): role checks only until D12 is decided.
package auth

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
)

// Role is a stored USER.role value.
type Role string

const (
	RoleSupervisor Role = "supervisor"
	RoleAssistant  Role = "assistant"
	RoleWorker     Role = "worker"
)

// Principal is the authenticated, active user making the request.
type Principal struct {
	UserID uuid.UUID
	Role   Role
}

// ErrUserNotFound means no active user matches the verified identity.
var ErrUserNotFound = errors.New("no active user for identity")

// Verifier checks a bearer token and returns the subject it vouches for.
type Verifier interface {
	Verify(ctx context.Context, token string) (subject string, err error)
}

// Resolver maps a verified subject to an active user.
type Resolver func(ctx context.Context, subject string) (Principal, error)

const principalKey = "auth.principal"

// PrincipalFrom returns the user set by Authenticate.
func PrincipalFrom(c *gin.Context) (Principal, bool) {
	v, ok := c.Get(principalKey)
	if !ok {
		return Principal{}, false
	}
	p, ok := v.(Principal)
	return p, ok
}

// Authenticate verifies the bearer token, resolves the active user and lets
// through only the allowed roles.
func Authenticate(v Verifier, resolve Resolver, allowed ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			httpx.WriteError(c, apperr.Unauthorized("missing bearer token"))
			return
		}
		subject, err := v.Verify(c.Request.Context(), token)
		if err != nil {
			httpx.WriteError(c, apperr.Unauthorized("invalid or expired token"))
			return
		}
		p, err := resolve(c.Request.Context(), subject)
		if errors.Is(err, ErrUserNotFound) {
			httpx.WriteError(c, apperr.Unauthorized("user is not registered or is inactive"))
			return
		}
		if err != nil {
			httpx.WriteError(c, apperr.Internal(err))
			return
		}
		if !slices.Contains(allowed, p.Role) {
			httpx.WriteError(c, apperr.Forbidden("role is not permitted for this endpoint"))
			return
		}
		c.Set(principalKey, p)
		c.Next()
	}
}

// RequireRole narrows an authenticated group to the given roles.
func RequireRole(allowed ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := PrincipalFrom(c)
		if !ok {
			httpx.WriteError(c, apperr.Internal(errors.New("RequireRole used without Authenticate")))
			return
		}
		if !slices.Contains(allowed, p.Role) {
			httpx.WriteError(c, apperr.Forbidden("role is not permitted for this endpoint"))
			return
		}
		c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	token = strings.TrimSpace(token)
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}
