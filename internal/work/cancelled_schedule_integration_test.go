package work

import (
	"context"
	"strings"
	"testing"
)

func TestIntegrationReassignUnusedCancelledContract(t *testing.T) {
	for _, recorded := range []bool{false, true} {
		t.Run(map[bool]string{false: "unused", true: "leave_history"}[recorded], func(t *testing.T) {
			f := fixtureUsecases(t)
			target := fixtureUsecases(t)
			ctx := context.Background()
			in := ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-04", "08:00"}
			ids, err := f.service.Schedule(ctx, in)
			if err != nil {
				t.Fatal(err)
			}
			if recorded {
				if _, err = f.service.Leave(ctx, f.worker, LeaveInput{"2026-10-04", "personal"}); err != nil {
					t.Fatal(err)
				}
			}
			if err = f.service.ContractState(ctx, f.tor, "cancelled"); err != nil {
				t.Fatal(err)
			}
			candidates, err := f.service.AvailableWorkers(ctx, "2026-10-04", "", 1000, 0)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(candidates), f.workerID) == recorded {
				t.Fatal("incorrect availability", string(candidates))
			}
			in.AssignmentID = target.assignment
			next, err := f.service.Schedule(ctx, in)
			if recorded {
				if err == nil {
					t.Fatal("replaced a slot with leave history")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(next) != 1 || next[0] != ids[0] {
				t.Fatal("must reuse unique worker/day slot", next, ids)
			}
			// A second assignment must not overwrite the newly active contract.
			if _, err = f.service.Schedule(ctx, in); err == nil {
				t.Fatal("duplicate active assignment accepted")
			}
			var assignment string
			if err = f.pool.QueryRow(ctx, `SELECT assignment_id::text FROM work_schedule WHERE schedule_id=$1`, ids[0]).Scan(&assignment); err != nil || assignment != target.assignment {
				t.Fatal(assignment, err)
			}
		})
	}
}
