package contract

import (
	"context"
	"github.com/google/uuid"
	"math"
	"testing"
)

func TestNewLocationCoordinateValidation(t *testing.T) {
	for _, tc := range []struct {
		name     string
		lat, lon *float64
		bad      bool
	}{
		{"legacy no coordinates", nil, nil, false},
		{"valid", ptrCoord(13.8141661), ptrCoord(100.6357345), false},
		{"zero valid", ptrCoord(0), ptrCoord(0), false},
		{"missing longitude", ptrCoord(13), nil, true},
		{"bad latitude", ptrCoord(91), ptrCoord(100), true},
		{"bad longitude", ptrCoord(13), ptrCoord(181), true},
		{"nonfinite", ptrCoord(math.NaN()), ptrCoord(100), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scope := Scope{Areas: []Area{{NewLocation: &NewLocation{Name: "park", Address: "address", Latitude: tc.lat, Longitude: tc.lon}, RequiredWorkers: 1, Items: []Item{{EquipmentName: "tool", RequiredQty: 1}}}}}
			normalized := scope.normalized()
			if normalized.Areas[0].NewLocation.Latitude != tc.lat || normalized.Areas[0].NewLocation.Longitude != tc.lon {
				t.Fatal("normalization dropped coordinates")
			}
			if got := len(ValidateScopeFormat(normalized)) > 0; got != tc.bad {
				t.Fatalf("invalid=%v want=%v", got, tc.bad)
			}
		})
	}
}
func ptrCoord(v float64) *float64 { return &v }

func TestIntegrationCreateTORLocationCoordinates(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	tx, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer tx.Rollback(ctx)
	lat, lon := 13.8141661, 100.6357345
	id, e := (TORRepository{}).CreateLocation(ctx, tx, "coordinate-"+uuid.NewString(), "address", &lat, &lon)
	if e != nil {
		t.Fatal(e)
	}
	var savedLat, savedLon float64
	if e = tx.QueryRow(ctx, `SELECT latitude::float8,longitude::float8 FROM location WHERE location_id=$1`, id).Scan(&savedLat, &savedLon); e != nil {
		t.Fatal(e)
	}
	if math.Abs(savedLat-lat) > 0.00000005 || math.Abs(savedLon-lon) > 0.00000005 {
		t.Fatal("coordinates not persisted")
	}
}
