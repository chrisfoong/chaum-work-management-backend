package work

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) scheduleMessage(ctx context.Context, id string) (string, error) {
	var project, place, date, start string
	err := s.Repo.Pool.QueryRow(ctx, `SELECT c.project_name,l.location_name,sc.work_date::text,to_char(sc.shift_start_time,'HH24:MI') FROM work_schedule sc JOIN tor_location_assignment a USING(assignment_id) JOIN contract_tor c USING(tor_id) JOIN location l USING(location_id) WHERE sc.schedule_id=$1`, id).Scan(&project, &place, &date, &start)
	if err != nil {
		return "", err
	}
	beg, err := Start(date, start)
	if err != nil {
		return "", err
	}
	end := beg.AddDate(0, 0, 0).Add(8 * 60 * 60 * 1000000000)
	endText := end.Format("15:04")
	if end.Format("2006-01-02") != date {
		endText += " (วันถัดไป)"
	}
	return fmt.Sprintf("Chaum: schedule details.\nโครงการ: %s\nพื้นที่: %s\nวันที่: %s\nเวลากะ: %s–%s\nระยะเวลา: 8 ชั่วโมง", project, place, cardThaiDate(date), start, endText), nil
}
func (s *Service) payrollMessage(ctx context.Context, id string) (string, error) {
	var number, start, end, base, deduct, net string
	var paid bool
	err := s.Repo.Pool.QueryRow(ctx, `SELECT payroll_slip_no,period_start::text,period_end::text,base_wage::text,total_deduction::text,net_wage::text,is_paid FROM payroll WHERE payroll_id=$1`, id).Scan(&number, &start, &end, &base, &deduct, &net, &paid)
	if err != nil {
		return "", err
	}
	state := "calculated"
	note := "รอยืนยันการจ่าย สลิปยังไม่ปรากฏในรายการสลิปที่จ่ายแล้ว"
	if paid {
		state = "paid"
		note = "ผู้ควบคุมงานบันทึกยืนยันการจ่ายแล้ว ตรวจสอบสลิปใน Mini App ได้"
	}
	return fmt.Sprintf("Chaum: payroll %s.\nรอบค่าจ้าง: %s – %s\nค่าจ้างก่อนหัก: %s บาท\nยอดหัก: %s บาท\nค่าจ้างสุทธิ: %s บาท\n%s", state, cardThaiDate(start), cardThaiDate(end), base, deduct, net, note), nil
}
func (s *Service) requestReviewMessage(ctx context.Context, id string) (string, error) {
	var number, state, reason string
	err := s.Repo.Pool.QueryRow(ctx, `SELECT requisition_no,status::text,reason FROM equipment_requisition WHERE requisition_id=$1`, id).Scan(&number, &state, &reason)
	if err != nil {
		return "", err
	}
	if state == "rejected" {
		if pos := strings.LastIndex(reason, "[Rejection] "); pos >= 0 {
			reason = reason[pos+len("[Rejection] "):]
		} else {
			reason = "เปิดระบบเพื่อตรวจสอบรายละเอียดคำขอ"
		}
		return "Chaum: equipment rejected.\nเลขคำขอ: " + number + "\nเหตุผล: " + reason, nil
	}
	return "Chaum: equipment approved.\nเลขคำขอ: " + number + "\nอนุมัติให้จัดซื้อแล้ว อยู่ระหว่างดำเนินการจัดหา ยังไม่ใช่การส่งมอบอุปกรณ์", nil
}
