package work

import (
	"context"
	"strings"
	"testing"
	"time"
)

type leaveNoticeCapture struct{ messages chan string }

func (n *leaveNoticeCapture) Send(_ context.Context, _ string, msg string) error {
	n.messages <- msg
	return nil
}
func TestIntegrationLeaveDecisionFlexNotification(t *testing.T) {
	for _, decision := range []string{"approved", "rejected"} {
		t.Run(decision, func(t *testing.T) {
			f := fixtureUsecases(t)
			ctx := context.Background()
			s := f.service
			if _, err := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-02", "08:00"}); err != nil {
				t.Fatal(err)
			}
			id, err := s.Leave(ctx, f.worker, LeaveInput{"2026-10-02", "ทดสอบใบลา พักธุระส่วนตัว"})
			if err != nil {
				t.Fatal(err)
			}
			if _, err = s.Leave(ctx, f.worker, LeaveInput{"2026-10-02", "ซ้ำ"}); err == nil {
				t.Fatal("duplicate leave accepted")
			}
			capture := &leaveNoticeCapture{messages: make(chan string, 10)}
			s.Notify = capture
			if err = s.ReviewLeave(ctx, id, decision); err != nil {
				t.Fatal(err)
			}
			var message string
			select {
			case message = <-capture.messages:
			case <-time.After(5 * time.Second):
				t.Fatal("committed review emitted no notification")
			}
			if !strings.HasPrefix(message, "Chaum: leave "+decision+".\n") || !strings.Contains(message, "2026-10-02") || !strings.Contains(message, "พักธุระส่วนตัว") {
				t.Fatal("decision details missing")
			}
			sender := LinePush{WorkerURL: "https://miniapp.line.me/2011922658-kb0yFlBE"}
			card := sender.messageCard(message)
			bubble := card["contents"].(map[string]any)
			body := bubble["body"].(map[string]any)["contents"].([]map[string]any)
			want := "ใบลาได้รับอนุมัติแล้ว"
			if decision == "rejected" {
				want = "ใบลาไม่ได้รับอนุมัติ"
			}
			if body[1]["text"] != want {
				t.Fatal("incorrect approval title")
			}
			retry, err := s.leaveReviewMessage(ctx, id)
			if err != nil || retry != message {
				t.Fatal("resend differs from committed state")
			}
			if err = s.ReviewLeave(ctx, id, decision); err == nil {
				t.Fatal("duplicate review accepted")
			}
			select {
			case <-capture.messages:
				t.Fatal("duplicate review sent another notification")
			case <-time.After(50 * time.Millisecond):
			}
		})
	}
}

func TestLeaveFlexRequiresCommittedDecision(t *testing.T) {
	if _, err := formatLeaveReviewMessage("L-1", "2026-10-30", "ทดสอบ", "pending"); err == nil {
		t.Fatal("pending announced as reviewed")
	}
}
