package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestMoneyExact(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want int64
	}{{"0", 0}, {"0.01", 1}, {"0.1", 10}, {"400.00", 40000}, {"99999999.99", 9999999999}} {
		got, e := Money(tt.in)
		if e != nil || got != tt.want {
			t.Errorf("%s = %d %v", tt.in, got, e)
		}
		if back, e := Money(Decimal(got)); e != nil || back != got {
			t.Fatal("round trip lost precision")
		}
	}
	for _, v := range []string{"-1", "NaN", "1e2", "1.001", "1.", ".2", "+1", "9999999999999"} {
		if _, e := Money(v); e == nil {
			t.Errorf("accepted %s", v)
		}
	}
}
func TestPenaltyBoundaries(t *testing.T) {
	start := time.Date(2026, 10, 8, 8, 0, 0, 0, Bangkok)
	for _, tt := range []struct {
		delay time.Duration
		want  int64
	}{{0, 0}, {time.Nanosecond, 30000}, {time.Hour - time.Nanosecond, 30000}, {time.Hour, 40000}, {3 * time.Hour, 40000}, {3*time.Hour + time.Nanosecond, 150000}} {
		check := start.Add(tt.delay)
		if got := Penalty("late", start, &check); got != tt.want {
			t.Errorf("delay %v got %d", tt.delay, got)
		}
	}
	if Penalty("leave", start, nil) != 0 || Penalty("absent", start, nil) != 150000 {
		t.Fatal("leave/absent wrong")
	}
}
func TestPeriodsAndOvernight(t *testing.T) {
	for _, v := range [][2]string{{"2028-02-16", "2028-02-29"}, {"2026-10-01", "2026-10-15"}, {"2026-10-16", "2026-10-31"}} {
		if e := Period(v[0], v[1]); e != nil {
			t.Fatal(e)
		}
	}
	for _, v := range [][2]string{{"2026-02-16", "2026-02-29"}, {"2026-10-01", "2026-10-31"}, {"2026-10-15", "2026-10-31"}, {"2026-10-16", "2026-11-01"}} {
		if e := Period(v[0], v[1]); e == nil {
			t.Fatal("accepted invalid period")
		}
	}
	start, e := Start("2026-10-08", "22:00")
	if e != nil || start.Add(8*time.Hour).Day() != 9 || start.Add(10*time.Hour).Hour() != 8 {
		t.Fatal("overnight cutoff wrong")
	}
}
func TestGPSAndQRBoth(t *testing.T) {
	now := time.Date(2026, 10, 8, 8, 0, 0, 0, Bangkok)
	s := &Service{QRSecret: []byte(strings.Repeat("s", 32)), Now: func() time.Time { return now }}
	id := uuid.NewString()
	token, e := s.signQR(QRClaims{id, now.Unix(), now.Unix() + 60})
	if e != nil || s.checkQR(token, id) != nil {
		t.Fatal("valid QR failed")
	}
	if s.checkQR(token, uuid.NewString()) == nil || s.checkQR(token+"x", id) == nil {
		t.Fatal("QR tamper/area accepted")
	}
	now = now.Add(time.Minute)
	if s.checkQR(token, id) == nil {
		t.Fatal("expired QR accepted")
	}
	lat, lon, accuracy := 13.75, 100.5, 50.0
	in := CheckIn{Latitude: &lat, Longitude: &lon, Accuracy: &accuracy}
	if GPS(in, lat, lon) != nil {
		t.Fatal("same location rejected")
	}
	accuracy = 50.01
	if GPS(in, lat, lon) == nil {
		t.Fatal("poor accuracy accepted")
	}
	accuracy = 1
	lat += .002
	if GPS(in, 13.75, lon) == nil {
		t.Fatal("outside radius accepted")
	}
	in.Latitude = nil
	if GPS(in, 13.75, lon) == nil {
		t.Fatal("missing gps accepted")
	}
}
func TestFileOwnership(t *testing.T) {
	p := auth.Principal{UserID: uuid.New(), Role: auth.RoleWorker}
	if !owned(p, p.UserID.String()+"/"+uuid.NewString()) || owned(p, uuid.NewString()+"/"+uuid.NewString()) || owned(p, p.UserID.String()+"/../secret") {
		t.Fatal("ownership checks wrong")
	}
	if mediaAllowed(p, "application/pdf") || mediaAllowed(p, "image/svg+xml") {
		t.Fatal("worker accepted non-photo")
	}
}
