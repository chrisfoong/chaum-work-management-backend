package work

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Printed area codes are stable for the assignment and signing key. Authorization,
// contract validity, shift ownership and GPS are checked afresh at every check-in.
const permanentQRVersion = 2

type QRClaims struct {
	Assignment string `json:"assignment_id"`
	Version    int    `json:"v"`
}

func (s *Service) signQR(c QRClaims) (string, error) {
	if len(s.QRSecret) < 32 {
		return "", conflict("QR signing secret is not configured")
	}
	body, _ := json.Marshal(c)
	payload := base64.RawURLEncoding.EncodeToString(body)
	mac := hmac.New(sha256.New, s.QRSecret)
	mac.Write([]byte(payload))
	return payload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
func (s *Service) checkQR(token, assignment string) error {
	if len(s.QRSecret) < 32 {
		return conflict("QR signing secret is not configured")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 || len(token) > 2048 {
		return invalid("qr_token", "invalid QR")
	}
	signature, e := base64.RawURLEncoding.DecodeString(parts[1])
	if e != nil {
		return invalid("qr_token", "invalid QR")
	}
	mac := hmac.New(sha256.New, s.QRSecret)
	mac.Write([]byte(parts[0]))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return invalid("qr_token", "invalid signature")
	}
	body, e := base64.RawURLEncoding.DecodeString(parts[0])
	if e != nil {
		return invalid("qr_token", "invalid QR")
	}
	var c QRClaims
	if json.Unmarshal(body, &c) != nil {
		return invalid("qr_token", "invalid QR")
	}
	if c.Version != permanentQRVersion || c.Assignment != assignment {
		return invalid("qr_token", "unsupported QR or wrong assignment; scan the current area sign")
	}
	return nil
}
func (s *Service) QR(ctx context.Context, id string) (map[string]any, error) {
	if e := validID(id, "assignment_id"); e != nil {
		return nil, e
	}
	var exists bool
	if e := s.Repo.Row(ctx, s.Repo.Pool, query1, id).Scan(&exists); e != nil {
		return nil, e
	}
	if !exists {
		return nil, conflict("assignment not available")
	}
	token, e := s.signQR(QRClaims{Assignment: id, Version: permanentQRVersion})
	return map[string]any{"qr_token": token, "assignment_id": id, "permanent": true}, e
}
