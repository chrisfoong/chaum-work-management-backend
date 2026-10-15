package work

import (
	"context"
	"encoding/json"
	"testing"
)

func TestIntegrationOnboardingRoutesVerifiedDatabaseRole(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	for _, tc := range []struct{ subject, platform string }{
		{f.worker.UserID.String(), "worker"}, {f.assistant.UserID.String(), "web"}, {f.supervisor.UserID.String(), "web"},
	} {
		result, e := f.service.OnboardingState(ctx, tc.subject)
		if e != nil {
			t.Fatal(e)
		}
		data, _ := json.Marshal(result)
		var state map[string]any
		json.Unmarshal(data, &state)
		if state["status"] != "registered" || state["platform"] != tc.platform {
			t.Fatalf("wrong platform: %s", data)
		}
	}
	if _, e := f.pool.Exec(ctx, `UPDATE public."USER" SET is_active=false WHERE user_id=$1`, f.worker.UserID); e != nil {
		t.Fatal(e)
	}
	result, e := f.service.OnboardingState(ctx, f.worker.UserID.String())
	if e != nil {
		t.Fatal(e)
	}
	data, _ := json.Marshal(result)
	var state map[string]any
	json.Unmarshal(data, &state)
	if state["status"] != "inactive" || state["platform"] != nil {
		t.Fatal("inactive account acquired destination", string(data))
	}
}
