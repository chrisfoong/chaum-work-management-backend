package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

func TestIntegrationAssignmentEquipmentFiltersTORAndOwner(t *testing.T) {
	f := fixtureUsecases(t)
	s := f.service
	ctx := context.Background()
	if _, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"}); e != nil {
		t.Fatal(e)
	}
	in, e := s.Equipment(ctx, "", EquipmentInput{Name: "in-tor-" + uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Equipment(ctx, "", EquipmentInput{Name: "other-tor-" + uuid.NewString()}); e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty) SELECT requisition_id,$2,1 FROM equipment_requisition WHERE assignment_id=$1 AND requisition_type='tor_base'`, f.assignment, in); e != nil {
		t.Fatal(e)
	}
	read := func(limit, offset int) []struct {
		ID string `json:"equipment_id"`
	} {
		data, e := s.AssignmentEquipment(ctx, f.worker, f.assignment, limit, offset)
		if e != nil {
			t.Fatal(e)
		}
		var out []struct {
			ID string `json:"equipment_id"`
		}
		if e = json.Unmarshal(data, &out); e != nil {
			t.Fatal(e)
		}
		return out
	}
	if result := read(100, 0); len(result) != 1 || result[0].ID != in {
		t.Fatalf("wrong TOR equipment: %+v", result)
	}
	if result := read(100, 1); len(result) != 0 {
		t.Fatal("pagination before filtering")
	}
	if _, e = s.AssignmentEquipment(ctx, f.other, f.assignment, 100, 0); e == nil {
		t.Fatal("unassigned worker allowed")
	}
	if _, e = s.AssignmentEquipment(ctx, f.worker, "invalid", 100, 0); e == nil {
		t.Fatal("invalid assignment allowed")
	}
	no := false
	if _, e = s.Equipment(ctx, in, EquipmentInput{Name: "disabled-" + uuid.NewString(), Active: &no}); e != nil {
		t.Fatal(e)
	}
	if result := read(100, 0); len(result) != 0 {
		t.Fatal("inactive equipment leaked")
	}
	if _, e = s.AssignmentEquipment(ctx, f.supervisor, "", 100, 0); e != nil {
		t.Fatal("catalog read regressed", e)
	}
}
