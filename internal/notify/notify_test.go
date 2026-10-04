package notify

import (
	"context"
	"testing"
)

func TestNoopNeverReportsSent(t *testing.T) {
	var n Noop
	msg := Message{LineUserID: "U1", Text: "hello"}

	status, err := n.Send(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusSkipped {
		t.Fatalf("status = %q, want %q", status, StatusSkipped)
	}
	if got := n.Messages(); len(got) != 1 || got[0] != msg {
		t.Fatalf("messages = %+v", got)
	}
}
