package contract

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestIntegrationUnassignedLocationCoordinates(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	prefix := "qa-coordinates-" + uuid.NewString()
	_, err = tx.Exec(ctx, `INSERT INTO location(location_name,address,latitude,longitude) VALUES($1,'mock',13.729864,100.541),($2,'mock',NULL,NULL)`, prefix+"-gps", prefix+"-unset")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := (TORRepository{}).SearchLocations(ctx, tx, prefix, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d locations", len(rows))
	}
	for _, row := range rows {
		if row.LocationName == prefix+"-gps" {
			if row.Latitude == nil || row.Longitude == nil || *row.Latitude != "13.7298640" || *row.Longitude != "100.5410000" {
				t.Fatalf("coordinates lost: %#v", row)
			}
		} else if row.Latitude != nil || row.Longitude != nil {
			t.Fatal("unset coordinates must remain null")
		}
	}
}
