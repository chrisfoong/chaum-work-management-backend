package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestIntegrationClosedContractHiddenFromWorker(t *testing.T) {
	for _, state := range []string{"complete", "cancelled"} {
		t.Run(state, func(t *testing.T) {
			f := fixtureUsecases(t)
			ctx := context.Background()
			_, e := f.service.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
			if e != nil {
				t.Fatal(e)
			}
			if e = f.service.ContractState(ctx, f.tor, state); e != nil {
				t.Fatal(e)
			}
			data, e := f.service.FilteredList(ctx, f.worker, "schedules", ListFilter{}, 100, 0)
			var rows []any
			if e != nil || json.Unmarshal(data, &rows) != nil || len(rows) != 0 {
				t.Fatal(string(data), e)
			}
			data, e = f.service.FilteredList(ctx, f.assistant, "schedules", ListFilter{TorID: f.tor}, 100, 0)
			if e != nil || json.Unmarshal(data, &rows) != nil || len(rows) != 1 {
				t.Fatal("staff history lost", string(data), e)
			}
		})
	}
}
func TestIntegrationDeleteOnlyUnusedArea(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	lat, lon := 13.75, 100.5
	loc, e := f.service.Location(ctx, "", LocationInput{Name: uuid.NewString(), Latitude: &lat, Longitude: &lon})
	if e != nil {
		t.Fatal(e)
	}
	eq, e := f.service.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	id, e := f.service.SaveAssignment(ctx, f.supervisor, f.tor, "", AssignmentEditInput{LocationID: loc, RequiredWorkers: 1, Items: []ScopeItemEdit{{EquipmentID: eq, Required: 1}}})
	if e != nil {
		t.Fatal(e)
	}
	if e = f.service.DeleteAssignment(ctx, f.assistant, id); e == nil {
		t.Fatal("assistant deleted area")
	}
	if e = f.service.DeleteAssignment(ctx, f.supervisor, id); e != nil {
		t.Fatal(e)
	}
	var exists bool
	if e = f.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM location WHERE location_id=$1)`, loc).Scan(&exists); e != nil || !exists {
		t.Fatal("master location removed", e)
	}
	*f.now = time.Date(2026, 10, 1, 8, 0, 0, 0, Bangkok)
	if e = f.service.DeleteAssignment(ctx, f.supervisor, f.assignment); e == nil {
		t.Fatal("completed procurement removed")
	}
}
