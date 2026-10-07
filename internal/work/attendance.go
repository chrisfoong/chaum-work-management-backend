package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"
	"time"
)

type QRClaims struct {
	Assignment string `json:"assignment_id"`
	Issued     int64  `json:"iat"`
	Expiry     int64  `json:"exp"`
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
	now := s.Now().Unix()
	if c.Assignment != assignment || c.Issued > now || c.Expiry <= now || c.Expiry-c.Issued != 60 {
		return invalid("qr_token", "expired QR or wrong assignment")
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
	now := s.Now()
	token, e := s.signQR(QRClaims{id, now.Unix(), now.Unix() + 60})
	return map[string]any{"qr_token": token, "expires_at": now.Add(time.Minute)}, e
}

type CheckIn struct {
	ScheduleID string   `json:"schedule_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Accuracy   *float64 `json:"accuracy_m"`
	QR         string   `json:"qr_token"`
}

func GPS(in CheckIn, lat, lon float64) error {
	if in.Latitude == nil || in.Longitude == nil || in.Accuracy == nil {
		return invalid("gps", "coordinates and accuracy required")
	}
	a, b, c := *in.Latitude, *in.Longitude, *in.Accuracy
	if math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(c) || math.IsInf(a, 0) || math.IsInf(b, 0) || math.IsInf(c, 0) || a < -90 || a > 90 || b < -180 || b > 180 || c < 0 || c > 50 {
		return invalid("gps", "invalid coordinates or accuracy over 50m")
	}
	if Distance(a, b, lat, lon) > 200 {
		return invalid("gps", "outside 200m radius")
	}
	return nil
}
func (s *Service) CheckIn(ctx context.Context, p auth.Principal, in CheckIn) error {
	if e := validID(in.ScheduleID, "schedule_id"); e != nil {
		return e
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		worker, e := s.Worker(ctx, q, p.UserID.String())
		if e != nil {
			return e
		}
		var assignment, date, clock, state string
		var lat, lon *float64
		if e = s.Repo.Row(ctx, q, query2, in.ScheduleID, worker).Scan(&assignment, &date, &clock, &state, &lat, &lon); e != nil {
			return e
		}
		start, e := Start(date, clock)
		if e != nil {
			return e
		}
		now := s.Now()
		if state != "scheduled" || now.Before(start) || !now.Before(start.Add(8*time.Hour)) {
			return conflict("shift is not open for check-in")
		}
		var leave bool
		if e = s.Repo.Row(ctx, q, query3, p.UserID, date).Scan(&leave); e != nil {
			return e
		}
		if leave {
			return conflict("approved leave cannot check in")
		}
		if lat == nil || lon == nil {
			return conflict("location coordinates missing")
		}
		if e = s.checkQR(in.QR, assignment); e != nil {
			return e
		}
		if e = GPS(in, *lat, *lon); e != nil {
			return e
		}
		status := "on_time"
		if now.After(start) {
			status = "late"
		}
		tag, e := s.Repo.Exec(ctx, q, query4, in.ScheduleID, worker, date, now, status)
		if e != nil {
			return e
		}
		if tag.RowsAffected() != 1 {
			return conflict("check-in already exists or attendance inconsistent")
		}
		return nil
	})
}

type CheckOut struct {
	ScheduleID  string `json:"schedule_id"`
	Description string `json:"description"`
	Photo       string `json:"photo_path"`
}

func (s *Service) CheckOut(ctx context.Context, p auth.Principal, in CheckOut) error {
	if e := validID(in.ScheduleID, "schedule_id"); e != nil {
		return e
	}
	if e := text(in.Description, "description", 1000); e != nil {
		return e
	}
	if e := s.evidence(ctx, p, in.Photo); e != nil {
		return e
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		worker, e := s.Worker(ctx, q, p.UserID.String())
		if e != nil {
			return e
		}
		var status string
		if e = s.Repo.Row(ctx, q, query5, in.ScheduleID, worker).Scan(&status); e != nil {
			return e
		}
		if status != "scheduled" {
			return conflict("shift not open")
		}
		tag, e := s.Repo.Exec(ctx, q, query6, in.ScheduleID, worker, s.Now())
		if e != nil {
			return e
		}
		if tag.RowsAffected() != 1 {
			return conflict("check-in missing or already checked out")
		}
		if _, e = s.Repo.Exec(ctx, q, query7, in.ScheduleID, worker, in.Description, in.Photo); e != nil {
			return e
		}
		_, e = s.Repo.Exec(ctx, q, query8, in.ScheduleID)
		return e
	})
}
func (s *Service) Attendance(ctx context.Context, p auth.Principal, limit, offset int) (json.RawMessage, error) {
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return s.Repo.List(ctx, s.Repo.Pool, query9, owner, limit, offset)
}
func (s *Service) Finalize(ctx context.Context) (int64, error) {
	var affected int64
	e := s.Repo.Transaction(ctx, func(q Query) error {
		tag, e := s.Repo.Exec(ctx, q, query10, s.Now())
		if e == nil {
			affected = tag.RowsAffected()
		}
		return e
	})
	return affected, e
}
