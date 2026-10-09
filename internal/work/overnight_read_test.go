package work

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestIntegrationWorkerReadOngoingOvernightShift(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	ids, err := f.service.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "22:00"})
	if err != nil {
		t.Fatal(err)
	}
	*f.now = time.Date(2026, 10, 2, 1, 0, 0, 0, Bangkok)
	read := func() []map[string]any {
		t.Helper()
		data, err := f.service.FilteredList(ctx, f.worker, "schedules", ListFilter{}, 100, 0)
		if err != nil {
			t.Fatal(err)
		}
		var result []map[string]any
		if err = json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	shifts := read()
	if len(shifts) != 1 || shifts[0]["schedule_id"] != ids[0] {
		t.Fatal("ongoing prior-day shift must remain visible")
	}
	data, err := f.service.FilteredList(ctx, f.other, "schedules", ListFilter{}, 100, 0)
	if err != nil || string(data) != "[]" {
		t.Fatal("ownership must remain enforced", err)
	}
	*f.now = time.Date(2026, 10, 2, 6, 0, 0, 0, Bangkok)
	if len(read()) != 0 {
		t.Fatal("ended prior-day shift must disappear")
	}
}
