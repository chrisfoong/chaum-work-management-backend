package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// SupabaseVerifier verifies Supabase Auth access tokens signed with the
// project's asymmetric signing keys (JWKS). Legacy HS256 tokens are rejected.
type SupabaseVerifier struct {
	keyfunc jwt.Keyfunc
	issuer  string
}

// NewSupabaseVerifier loads and keeps refreshing signing keys from
// <supabaseURL>/auth/v1/.well-known/jwks.json until ctx ends.
func NewSupabaseVerifier(ctx context.Context, supabaseURL string) (*SupabaseVerifier, error) {
	k, err := keyfunc.NewDefaultCtx(ctx, []string{supabaseURL + "/auth/v1/.well-known/jwks.json"})
	if err != nil {
		return nil, fmt.Errorf("load supabase jwks: %w", err)
	}
	return &SupabaseVerifier{keyfunc: k.Keyfunc, issuer: supabaseURL + "/auth/v1"}, nil
}

// Verify checks signature, issuer and expiry, and returns the subject (user id).
func (v *SupabaseVerifier) Verify(_ context.Context, token string) (string, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, v.keyfunc,
		jwt.WithValidMethods([]string{"ES256", "RS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", errors.New("token has no subject")
	}
	return claims.Subject, nil
}
