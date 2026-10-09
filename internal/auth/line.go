package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// lineVerifyURL is LINE Login's ID token verification endpoint.
const lineVerifyURL = "https://api.line.me/oauth2/v2.1/verify"

// LineVerifier verifies LIFF ID tokens server-side with LINE.
type LineVerifier struct {
	client    *http.Client
	endpoint  string
	channelID string
}

// NewLineVerifier returns a verifier for tokens issued to channelID.
func NewLineVerifier(channelID string) *LineVerifier {
	return &LineVerifier{
		client:    &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
		endpoint:  lineVerifyURL,
		channelID: channelID,
	}
}

// Verify sends the ID token to LINE and returns the LINE user id (sub).
func (v *LineVerifier) Verify(ctx context.Context, idToken string) (string, error) {
	if v.channelID == "" || idToken == "" {
		return "", errors.New("LINE verifier is not configured or token is empty")
	}
	form := url.Values{"id_token": {idToken}, "client_id": {v.channelID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("build LINE verify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call LINE verify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LINE rejected ID token: status %d", resp.StatusCode)
	}
	var body struct {
		Sub string `json:"sub"`
		Aud string `json:"aud"`
		Iss string `json:"iss"`
		Exp int64  `json:"exp"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return "", fmt.Errorf("decode LINE verify response: %w", err)
	}
	if body.Aud != v.channelID {
		return "", errors.New("LINE ID token audience mismatch")
	}
	if body.Sub == "" {
		return "", errors.New("LINE ID token has no subject")
	}
	if body.Iss != "https://access.line.me" || body.Exp <= time.Now().Unix() {
		return "", errors.New("LINE ID token issuer or expiry invalid")
	}
	return body.Sub, nil
}
