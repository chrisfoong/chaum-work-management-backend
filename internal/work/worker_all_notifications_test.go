package work

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type workerFlowCapture struct {
	line string
	ch   chan string
}

func (n *workerFlowCapture) Send(_ context.Context, to, msg string) error {
	if to == n.line {
		n.ch <- msg
	}
	return nil
}

func TestIntegrationAllWorkerFlexFlows(t *testing.T) {
	f := fixtureUsecases(t)
	s, ctx := f.service, context.Background()
	_, err := f.pool.Exec(ctx, `UPDATE contract_tor SET project_name='TEST · งานดูแลสวนตาม Use Case' WHERE tor_id=$1`, f.tor)
	if err != nil {
		t.Fatal(err)
	}
	eq, err := s.Equipment(ctx, "", EquipmentInput{Name: "TEST · พลั่วดูแลสวน " + uuid.NewString()[:6]})
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.pool.Exec(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty,existing_qty) SELECT requisition_id,$2,2,2 FROM equipment_requisition WHERE assignment_id=$1 AND requisition_type='tor_base'`, f.assignment, eq)
	if err != nil {
		t.Fatal(err)
	}
	capture := &workerFlowCapture{line: f.worker.UserID.String(), ch: make(chan string, 20)}
	s.Notify = capture
	push := &LinePush{DashboardURL: "https://unfiled-bunkhouse-float.ngrok-free.dev/web", LogoURL: "https://unfiled-bunkhouse-float.ngrok-free.dev/brand/chaum-logo.png", WorkerURL: "https://miniapp.line.me/2011922658-kb0yFlBE"}
	var outputs []map[string]any
	receive := func(name string) {
		t.Helper()
		select {
		case msg := <-capture.ch:
			card := push.messageCard(msg)
			bubble := card["contents"].(map[string]any)
			bubble["body"].(map[string]any)["contents"].([]map[string]any)[0]["text"] = "ชะอุ่ม · ทดสอบตาม Use Case"
			card["altText"] = "ทดสอบ Worker · " + name
			outputs = append(outputs, map[string]any{"case": name, "source": "isolated PostgreSQL real business service flow", "message": card})
		case <-time.After(8 * time.Second):
			t.Fatal("missing Worker event", name)
		}
	}
	schedule, err := s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, "2026-10-01", "08:00"})
	if err != nil {
		t.Fatal(err)
	}
	receive("1W ตารางงานใหม่")
	for i, state := range []string{"approved", "rejected"} {
		date := []string{"2026-10-02", "2026-10-03"}[i]
		if _, err = s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.workerID}, date, "17:00"}); err != nil {
			t.Fatal(err)
		}
		receive("1W ตารางกะข้ามคืน " + state)
		leave, err := s.Leave(ctx, f.worker, LeaveInput{date, "TEST · ทดสอบขอลาตาม Use Case"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = s.Leave(ctx, f.worker, LeaveInput{date, "TEST ซ้ำ"}); err == nil {
			t.Fatal("duplicate leave accepted")
		}
		if _, err = s.ReviewLeaveAndReplace(ctx, leave, LeaveReviewInput{Status: state}); err != nil {
			t.Fatal(err)
		}
		receive("2W ผลใบลา " + state)
	}
	// Replacement scheduling produces the same detailed Worker schedule event.
	if _, err = s.Schedule(ctx, ScheduleInput{f.assignment, []string{f.otherID}, "2026-10-04", "08:00"}); err != nil {
		t.Fatal(err)
	}
	otherLeave, err := s.Leave(ctx, f.other, LeaveInput{"2026-10-04", "TEST คนเดิมลา"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.ReviewLeaveAndReplace(ctx, otherLeave, LeaveReviewInput{Status: "approved", WorkerID: f.workerID}); err != nil {
		t.Fatal(err)
	}
	receive("1W ได้รับมอบหมายให้ทำงานแทน")
	lat, lon, accuracy := 13.75, 100.5, 10.0
	qr, err := s.QR(ctx, f.assignment)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.CheckIn(ctx, f.worker, CheckIn{schedule[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); err != nil {
		t.Fatal(err)
	}
	if err = s.CheckIn(ctx, f.worker, CheckIn{schedule[0], &lat, &lon, &accuracy, qr["qr_token"].(string)}); err == nil {
		t.Fatal("duplicate check-in accepted")
	}
	request := func() string {
		t.Helper()
		id, err := s.Requisition(ctx, f.worker, RequestInput{f.assignment, "TEST · อุปกรณ์ไม่พอตาม Use Case", []RequestItem{{eq, 1, ""}}})
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	noBuy := request()
	reason := "TEST · อุปกรณ์เดิมเพียงพอ สามารถจัดสรรจากพื้นที่เดิมได้"
	if err = s.DecideRequest(ctx, f.assistant, noBuy, ProcurementDecision{"no_purchase", reason}); err != nil {
		t.Fatal(err)
	}
	if _, err = s.RetryNotification(ctx, f.assistant, noBuy, NotificationInput{Kind: "no_purchase", Reason: reason}); err != nil {
		t.Fatal(err)
	}
	receive("5W/6W ไม่ต้องจัดซื้อพร้อมเหตุผล")
	rejected := request()
	if err = s.DecideRequest(ctx, f.assistant, rejected, ProcurementDecision{"purchase", "TEST จำเป็นต้องจัดซื้อ"}); err != nil {
		t.Fatal(err)
	}
	if err = s.ReviewRequest(ctx, f.supervisor, rejected, ReviewInput{"reject", "TEST ขอให้ใช้รายการที่จัดสรรไว้แล้ว"}); err != nil {
		t.Fatal(err)
	}
	receive("5W/6W ไม่อนุมัติจัดซื้อพร้อมเหตุผล")
	delivery := request()
	if err = s.DecideRequest(ctx, f.assistant, delivery, ProcurementDecision{"purchase", "TEST จัดซื้อตาม TOR"}); err != nil {
		t.Fatal(err)
	}
	if err = s.ReviewRequest(ctx, f.supervisor, delivery, ReviewInput{"approve", ""}); err != nil {
		t.Fatal(err)
	}
	receive("5W/6W อนุมัติจัดซื้อ ยังไม่ส่งมอบ")
	var item string
	if err = f.pool.QueryRow(ctx, `SELECT item_id::text FROM requisition_item WHERE requisition_id=$1`, delivery).Scan(&item); err != nil {
		t.Fatal(err)
	}
	zero := 0
	if _, err = s.Purchase(ctx, f.assistant, delivery, PurchaseInput{Items: []PurchaseItem{{item, 1, "0.00", &zero}}}); err != nil {
		t.Fatal(err)
	}
	if err = s.Deliver(ctx, f.assistant, delivery, DeliveryInput{schedule[0], "TEST · ส่งมอบให้ผู้แจ้งขอ", []string{f.assistant.UserID.String() + "/" + uuid.NewString()}}); err != nil {
		t.Fatal(err)
	}
	if _, err = s.RetryNotification(ctx, f.assistant, delivery, NotificationInput{Kind: "delivery"}); err != nil {
		t.Fatal(err)
	}
	receive("6W ส่งมอบอุปกรณ์พร้อมรายการและจำนวน")
	if err = s.CheckOut(ctx, f.worker, CheckOut{schedule[0], "TEST ออกก่อนเวลา", f.worker.UserID.String() + "/" + uuid.NewString()}); err == nil {
		t.Fatal("early checkout accepted")
	}
	*f.now = time.Date(2026, 10, 1, 16, 0, 0, 0, Bangkok)
	if err = s.CheckOut(ctx, f.worker, CheckOut{schedule[0], "TEST ส่งผลงานครบกะ", f.worker.UserID.String() + "/" + uuid.NewString()}); err != nil {
		t.Fatal(err)
	}
	*f.now = time.Date(2026, 10, 16, 8, 0, 0, 0, Bangkok)
	if _, err = s.Finalize(ctx); err != nil {
		t.Fatal(err)
	}
	payroll, err := s.Payroll(ctx, f.supervisor, PayrollInput{f.worker.UserID.String(), "2026-10-01", "2026-10-15"})
	if err != nil {
		t.Fatal(err)
	}
	receive("7W ประมวลผลค่าจ้าง ยังไม่ยืนยันการจ่าย")
	if err = s.Pay(ctx, payroll); err != nil {
		t.Fatal(err)
	}
	// Explicit notification presentation test for committed paid state. Pay has
	// no automatic notification hook in the existing app; do not imply otherwise.
	s.Event("payroll", payroll)
	receive("7W สลิปที่ยืนยันการจ่ายแล้ว — เรียกแจ้งเตือนทดสอบเอง")
	summary, err := s.WorkerProjectSummaryMessage(ctx, f.worker, f.tor, "2026-10-01", "2026-10-15")
	if err != nil {
		t.Fatal(err)
	}
	if !s.notify(ctx, f.worker.UserID.String(), summary) {
		t.Fatal("personal summary capture failed")
	}
	receive("ส่วนเพิ่มตามคำขอ สรุปงานส่วนตัวในโครงการ")
	if dir := os.Getenv("CHAUM_WORKER_FLOW_ARTIFACTS"); dir != "" {
		data, err := json.MarshalIndent(outputs, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, "all-worker-flow-cards.json"), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("%d Worker Flex cards generated by actual service workflows; authentication, storage transport and GPS device are mocked; LINE transport captured locally", len(outputs))
}
