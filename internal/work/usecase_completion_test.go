package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"strings"
	"testing"
)

type completionNotifier struct {
	fail     bool
	messages []string
}

func (n *completionNotifier) Send(_ context.Context, _, msg string) error {
	n.messages = append(n.messages, msg)
	if n.fail {
		return errors.New("fixture unavailable")
	}
	return nil
}

func TestWorkerEquipmentReasonLimit(t *testing.T) {
	s := &Service{}
	_, e := s.Requisition(context.Background(), auth.Principal{}, RequestInput{AssignmentID: uuid.NewString(), Reason: strings.Repeat("ก", 501)})
	if e == nil {
		t.Fatal("501-character reason reached database")
	}
}

func TestIntegrationUsecaseCompletion(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	shifts, e := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if e != nil {
		t.Fatal(e)
	}
	_, e = f.pool.Exec(ctx, `UPDATE location SET address='fixture address' WHERE location_id=(SELECT location_id FROM tor_location_assignment WHERE assignment_id=$1)`, f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	t.Run("worker schedule address", func(t *testing.T) {
		data, e := s.FilteredList(ctx, f.worker, "schedules", ListFilter{}, 100, 0)
		if e != nil || !strings.Contains(string(data), `"address": "fixture address"`) {
			t.Fatal("address missing", e)
		}
	})
	qr, e := s.QR(ctx, f.assignment)
	if e != nil {
		t.Fatal(e)
	}
	lat, lon, accuracy := 13.75, 100.5, 10.0
	if e = s.CheckIn(ctx, f.worker, CheckIn{shifts[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); e != nil {
		t.Fatal(e)
	}
	equipment, e := s.Equipment(ctx, "", EquipmentInput{Name: uuid.NewString()})
	if e != nil {
		t.Fatal(e)
	}
	request := RequestInput{f.assignment, "original Worker reason", []RequestItem{{equipment, 1, ""}}}
	declined, e := s.Requisition(ctx, f.worker, request)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `UPDATE equipment SET is_active=false WHERE equipment_id=$1`, equipment); e != nil {
		t.Fatal(e)
	}
	t.Run("purchase rejects inactive equipment", func(t *testing.T) {
		if e := s.DecideRequest(ctx, f.assistant, declined, ProcurementDecision{"purchase", "buy"}); e == nil {
			t.Fatal("inactive purchase accepted")
		}
	})
	t.Run("no purchase records reviewer and preserves reason", func(t *testing.T) {
		if e := s.DecideRequest(ctx, f.assistant, declined, ProcurementDecision{"no_purchase", "reuse"}); e != nil {
			t.Fatal(e)
		}
		var reviewer, reason string
		var timestamp bool
		e := f.pool.QueryRow(ctx, `SELECT reviewed_by::text,reason,reviewed_at IS NOT NULL FROM equipment_requisition WHERE requisition_id=$1`, declined).Scan(&reviewer, &reason, &timestamp)
		if e != nil || reviewer != f.assistant.UserID.String() || reason != request.Reason || !timestamp {
			t.Fatal("review metadata mismatch", e)
		}
		if _, e = s.InspectRequest(ctx, declined); e == nil {
			t.Fatal("non-pending inspection accepted")
		}
	})
	n := &completionNotifier{fail: true}
	s.Notify = n
	t.Run("notification failure and resend do not alter decision", func(t *testing.T) {
		result, e := s.RetryNotification(ctx, f.assistant, declined, NotificationInput{"no_purchase", "reuse"})
		if e != nil || result.Status != "failed" || result.Failed != 1 {
			t.Fatal(result, e)
		}
		n.fail = false
		result, e = s.RetryNotification(ctx, f.assistant, declined, NotificationInput{"no_purchase", "reuse"})
		if e != nil || result.Status != "accepted" || result.Accepted != 1 || result.DeliveryVerified || result.Persisted {
			t.Fatal(result, e)
		}
		var state, reason string
		if e = f.pool.QueryRow(ctx, `SELECT status::text,reason FROM equipment_requisition WHERE requisition_id=$1`, declined).Scan(&state, &reason); e != nil || state != "rejected" || reason != request.Reason {
			t.Fatal("resend rewrote business data", e)
		}
	})
	t.Run("notification role and state guards", func(t *testing.T) {
		for _, kind := range []string{"delivery", "approval_needed", "purchase_funding", "bogus"} {
			if _, e := s.RetryNotification(ctx, f.assistant, declined, NotificationInput{Kind: kind}); e == nil {
				t.Fatal("wrong state accepted", kind)
			}
		}
		if _, e := s.RetryNotification(ctx, f.supervisor, declined, NotificationInput{"no_purchase", "reuse"}); e == nil {
			t.Fatal("wrong role accepted")
		}
		if _, e := s.RetryNotification(ctx, f.other, declined, NotificationInput{"no_purchase", "reuse"}); e == nil {
			t.Fatal("wrong actor accepted")
		}
	})
	s.Notify = nil
	if _, e = f.pool.Exec(ctx, `UPDATE equipment SET is_active=true WHERE equipment_id=$1`, equipment); e != nil {
		t.Fatal(e)
	}
	handed, e := s.Requisition(ctx, f.worker, request)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `UPDATE equipment_requisition SET status='completed' WHERE requisition_id=$1`, handed); e != nil {
		t.Fatal(e)
	}
	if _, e = f.pool.Exec(ctx, `UPDATE requisition_item SET actual_qty=1,actual_price=0 WHERE requisition_id=$1`, handed); e != nil {
		t.Fatal(e)
	}
	t.Run("no delivery notification before evidence", func(t *testing.T) {
		if _, e := s.RetryNotification(ctx, f.assistant, handed, NotificationInput{Kind: "delivery"}); e == nil {
			t.Fatal("notification before handover")
		}
	})
	if e = s.Deliver(ctx, f.assistant, handed, DeliveryInput{shifts[0], "handover", []string{f.assistant.UserID.String() + "/" + uuid.NewString()}}); e != nil {
		t.Fatal(e)
	}
	t.Run("delivery resend retains single evidence and actual item summary", func(t *testing.T) {
		s.Notify = n
		n.fail = true
		result, e := s.RetryNotification(ctx, f.assistant, handed, NotificationInput{Kind: "delivery"})
		if e != nil || result.Status != "failed" {
			t.Fatal(result, e)
		}
		n.fail = false
		result, e = s.RetryNotification(ctx, f.assistant, handed, NotificationInput{Kind: "delivery"})
		if e != nil || result.Status != "accepted" {
			t.Fatal(result, e)
		}
		evidence, e := s.Deliveries(ctx, handed)
		var rows []json.RawMessage
		if e != nil || json.Unmarshal(evidence, &rows) != nil || len(rows) != 1 {
			t.Fatal("resend duplicated evidence", e)
		}
		if !strings.Contains(n.messages[len(n.messages)-1], "× 1") {
			t.Fatal("actual quantities missing")
		}
		s.Notify = nil
		result, e = s.RetryNotification(ctx, f.assistant, handed, NotificationInput{Kind: "delivery"})
		if e != nil || result.Status != "disabled" {
			t.Fatal(result, e)
		}
	})
	t.Run("base zero acquisition escalates without expense", func(t *testing.T) {
		var base, item string
		e := f.pool.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type) VALUES($1,$2,$3,'base','tor_base') RETURNING requisition_id::text`, uuid.NewString(), f.supervisor.UserID, f.assignment).Scan(&base)
		if e != nil {
			t.Fatal(e)
		}
		e = f.pool.QueryRow(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty) VALUES($1,$2,2) RETURNING item_id::text`, base, equipment).Scan(&item)
		if e != nil {
			t.Fatal(e)
		}
		if e = s.Survey(ctx, base, []SurveyItem{{item, 0}}); e != nil {
			t.Fatal(e)
		}
		zero := 0
		in := PurchaseInput{Items: []PurchaseItem{{item, 0, "0", &zero}}}
		expense, e := s.Purchase(ctx, f.assistant, base, in)
		if e != nil || expense != "" {
			t.Fatal(expense, e)
		}
		var qty, expenses int
		var state string
		e = f.pool.QueryRow(ctx, `SELECT r.status::text,COALESCE(i.actual_qty,0),(SELECT count(*) FROM expense_claim WHERE requisition_id=r.requisition_id) FROM equipment_requisition r JOIN requisition_item i USING(requisition_id) WHERE r.requisition_id=$1`, base).Scan(&state, &qty, &expenses)
		if e != nil || state != "pending_fund" || qty != 0 || expenses != 0 {
			t.Fatal("zero acquisition incorrectly persisted", e)
		}
		if _, e = s.Purchase(ctx, f.assistant, base, in); e == nil {
			t.Fatal("zero-round replay accepted")
		}
	})
}
