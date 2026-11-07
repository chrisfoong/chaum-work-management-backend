package work

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
)

type ContractDetailsInput struct {
	ProjectName       string `json:"project_name"`
	Agency            string `json:"partner_agency"`
	Start             string `json:"start_date"`
	End               string `json:"end_date"`
	Value             string `json:"contract_value"`
	File              string `json:"contract_file_path"`
	ExpectedUpdatedAt string `json:"expected_updated_at"`
}

func validateContractEdit(in ContractDetailsInput) error {
	if e := text(in.ProjectName, "project_name", 255); e != nil {
		return e
	}
	if e := text(in.Agency, "partner_agency", 255); e != nil {
		return e
	}
	if e := validRange(in.Start, in.End); e != nil {
		return e
	}
	v, e := Money(in.Value)
	if e != nil || v > 999999999999 {
		return invalid("contract_value", "จำนวนเงินไม่ถูกต้อง")
	}
	if _, e := time.Parse(time.RFC3339Nano, in.ExpectedUpdatedAt); e != nil {
		return invalid("expected_updated_at", "กรุณาโหลดข้อมูลสัญญาล่าสุดก่อนแก้ไข")
	}
	return nil
}

func (s *Service) UpdateContractDetails(ctx context.Context, p auth.Principal, id string, in ContractDetailsInput) error {
	if p.Role != auth.RoleSupervisor {
		return invalid("role", "Supervisor เท่านั้น")
	}
	if e := validID(id, "tor_id"); e != nil {
		return e
	}
	if e := validateContractEdit(in); e != nil {
		return e
	}
	expected, _ := time.Parse(time.RFC3339Nano, in.ExpectedUpdatedAt)
	return s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "contract-edit:"+id); e != nil {
			return e
		}
		var oldStart, oldEnd, oldValue, status string
		var oldFile *string
		var updated time.Time
		if e := q.QueryRow(ctx, `SELECT start_date::text,end_date::text,contract_value::text,contract_file_url,status::text,updated_at FROM contract_tor WHERE tor_id=$1 FOR UPDATE`, id).Scan(&oldStart, &oldEnd, &oldValue, &oldFile, &status, &updated); e != nil {
			return e
		}
		if !updated.Equal(expected) {
			return conflict("ข้อมูลสัญญาถูกแก้ไขแล้ว กรุณารีเฟรช")
		}
		if status == "complete" || status == "cancelled" {
			return conflict("สัญญาสิ้นสุดหรือยกเลิกแล้ว ไม่สามารถแก้รายละเอียด")
		}
		var outside bool
		if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule s JOIN tor_location_assignment a USING(assignment_id) WHERE a.tor_id=$1 AND s.work_date NOT BETWEEN $2::date AND $3::date) OR EXISTS(SELECT 1 FROM company_invoice WHERE tor_id=$1 AND (billing_month<left($2::text,7) OR billing_month>left($3::text,7)))`, id, in.Start, in.End).Scan(&outside); e != nil {
			return e
		}
		if outside {
			return conflict("ช่วงวันที่ใหม่ไม่ครอบคลุมกะงานหรือใบวางบิลเดิม")
		}
		if in.Value != oldValue {
			old, _ := Money(oldValue)
			next, _ := Money(in.Value)
			if old != next {
				var financial bool
				if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM company_invoice WHERE tor_id=$1) OR EXISTS(SELECT 1 FROM expense_claim e JOIN equipment_requisition r USING(requisition_id) JOIN tor_location_assignment a USING(assignment_id) WHERE a.tor_id=$1) OR EXISTS(SELECT 1 FROM payroll p JOIN worker w ON w.user_id=p.user_id JOIN work_schedule s ON s.worker_id=w.worker_id JOIN tor_location_assignment a USING(assignment_id) WHERE a.tor_id=$1 AND s.work_date BETWEEN p.period_start AND p.period_end)`, id).Scan(&financial); e != nil {
					return e
				}
				if financial {
					return conflict("มีรายการเงินแล้ว ไม่สามารถแก้มูลค่าสัญญาย้อนหลัง")
				}
			}
		}
		file := ""
		if oldFile != nil {
			file = *oldFile
		}
		if in.File != "" && in.File != file {
			verifier, ok := s.Files.(interface {
				VerifyPNG(context.Context, auth.Principal, string) error
			})
			if !ok {
				return conflict("storage is not configured")
			}
			if e := verifier.VerifyPNG(ctx, p, in.File); e != nil {
				return e
			}
			file = in.File
		}
		_, e := q.Exec(ctx, `UPDATE contract_tor SET project_name=$2,partner_agency=$3,start_date=$4::date,end_date=$5::date,contract_value=$6::numeric,contract_file_url=$7,updated_at=now() WHERE tor_id=$1`, id, in.ProjectName, in.Agency, in.Start, in.End, in.Value, nullableText(file))
		return e
	})
}

func nullableText(v string) any {
	if v == "" {
		return nil
	}
	return v
}

type ScopeItemEdit struct {
	EquipmentID string `json:"equipment_id"`
	Required    int    `json:"required_qty"`
}
type AssignmentEditInput struct {
	LocationID      string           `json:"location_id"`
	RequiredWorkers int              `json:"required_workers"`
	ExpectedItems   *[]ScopeItemEdit `json:"expected_items"`
	ExpectedWorkers *int             `json:"expected_required_workers"`
	Items           []ScopeItemEdit  `json:"items"`
}

func validateAssignmentEdit(in AssignmentEditInput) error {
	if e := validID(in.LocationID, "location_id"); e != nil {
		return e
	}
	if in.RequiredWorkers < 1 || in.RequiredWorkers > 10000 {
		return invalid("required_workers", "จำนวนคนต้องอยู่ระหว่าง 1–10000")
	}
	if len(in.Items) == 0 || len(in.Items) > 100 {
		return invalid("items", "ระบุอุปกรณ์ตั้งต้น 1–100 รายการ")
	}
	seen := map[string]bool{}
	for _, item := range in.Items {
		if e := validID(item.EquipmentID, "equipment_id"); e != nil {
			return e
		}
		if seen[strings.ToLower(item.EquipmentID)] || item.Required < 1 || item.Required > 1000000 {
			return invalid("items", "อุปกรณ์ห้ามซ้ำและจำนวนต้องเป็นบวกไม่เกิน 1000000")
		}
		seen[strings.ToLower(item.EquipmentID)] = true
	}
	return nil
}

func (s *Service) AssignmentScope(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "assignment_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT a.assignment_id,a.tor_id,a.location_id,a.required_workers,r.requisition_id,r.status,COALESCE((SELECT jsonb_agg(jsonb_build_object('equipment_id',i.equipment_id,'equipment_name',e.equipment_name,'required_qty',i.required_qty) ORDER BY e.equipment_name) FROM requisition_item i JOIN equipment e USING(equipment_id) WHERE i.requisition_id=r.requisition_id),'[]'::jsonb) AS items FROM tor_location_assignment a LEFT JOIN equipment_requisition r ON r.assignment_id=a.assignment_id AND r.requisition_type='tor_base' WHERE a.assignment_id=$1`, id)
}

// A scope edit is deliberately conservative: no rewriting schedules, purchases,
// survey results or attendance. Adding a new area leaves existing areas intact.
func (s *Service) SaveAssignment(ctx context.Context, p auth.Principal, tor, id string, in AssignmentEditInput) (string, error) {
	if p.Role != auth.RoleSupervisor {
		return "", invalid("role", "Supervisor เท่านั้น")
	}
	if e := validID(tor, "tor_id"); e != nil {
		return "", e
	}
	if id != "" {
		if e := validID(id, "assignment_id"); e != nil {
			return "", e
		}
		if in.ExpectedWorkers == nil {
			return "", invalid("expected_required_workers", "กรุณาโหลดข้อมูลพื้นที่ล่าสุด")
		}
	}
	if e := validateAssignmentEdit(in); e != nil {
		return "", e
	}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "contract-edit:"+tor); e != nil {
			return e
		}
		var state string
		if e := q.QueryRow(ctx, `SELECT status::text FROM contract_tor WHERE tor_id=$1 FOR UPDATE`, tor).Scan(&state); e != nil {
			return e
		}
		if state != "registered" && state != "active" {
			return conflict("สัญญาสิ้นสุดหรือยกเลิกแล้ว")
		}
		var req string
		if id != "" {
			// Schedule locks the same assignment row before creating its shifts.
			var current int
			if e := q.QueryRow(ctx, `SELECT required_workers FROM tor_location_assignment WHERE assignment_id=$1 AND tor_id=$2 FOR UPDATE`, id, tor).Scan(&current); e != nil {
				return e
			}
			if current != *in.ExpectedWorkers {
				return conflict("จำนวนคนถูกแก้แล้ว กรุณารีเฟรช")
			}
			var used bool
			if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule WHERE assignment_id=$1)`, id).Scan(&used); e != nil {
				return e
			}
			if used {
				return conflict("พื้นที่มีกะงานแล้ว ไม่แก้ขอบเขตย้อนหลัง")
			}
			var baseCount int
			if e := q.QueryRow(ctx, `SELECT count(*) FROM equipment_requisition WHERE assignment_id=$1 AND requisition_type='tor_base'`, id).Scan(&baseCount); e != nil {
				return e
			}
			if baseCount != 1 {
				return conflict("พื้นที่ต้องมีใบอุปกรณ์ตั้งต้นเพียงหนึ่งใบ")
			}
			var rs string
			if e := q.QueryRow(ctx, `SELECT requisition_id::text,status::text FROM equipment_requisition WHERE assignment_id=$1 AND requisition_type='tor_base' FOR UPDATE`, id).Scan(&req, &rs); e != nil {
				return e
			}
			if rs != "pending_survey" {
				return conflict("เริ่มสำรวจหรือจัดซื้อแล้ว ไม่แก้อุปกรณ์ตั้งต้น")
			}
			var moved bool
			if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1 AND (existing_qty<>0 OR actual_qty<>0)) OR EXISTS(SELECT 1 FROM expense_claim WHERE requisition_id=$1)`, req).Scan(&moved); e != nil {
				return e
			}
			if moved {
				return conflict("มีการสำรวจหรือรายการเงินแล้ว")
			}
		}
		var exists bool
		if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM location WHERE location_id=$1)`, in.LocationID).Scan(&exists); e != nil {
			return e
		}
		if !exists {
			return invalid("location_id", "ไม่พบสถานที่")
		}
		if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tor_location_assignment WHERE tor_id=$1 AND location_id=$2 AND ($3::uuid IS NULL OR assignment_id<>$3))`, tor, in.LocationID, nullableID(id)).Scan(&exists); e != nil {
			return e
		}
		if exists {
			return conflict("สถานที่นี้ผูกกับ TOR แล้ว")
		}
		if id == "" {
			if e := q.QueryRow(ctx, `INSERT INTO tor_location_assignment(tor_id,location_id,required_workers) VALUES($1,$2,$3) RETURNING assignment_id::text`, tor, in.LocationID, in.RequiredWorkers).Scan(&id); e != nil {
				return e
			}
			if e := q.QueryRow(ctx, `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type,status) VALUES($1,$2,$3,'อุปกรณ์ตั้งต้นของพื้นที่ TOR','tor_base','pending_survey') RETURNING requisition_id::text`, document("REQ"), p.UserID, id).Scan(&req); e != nil {
				return e
			}
		} else {
			if _, e := q.Exec(ctx, `UPDATE tor_location_assignment SET location_id=$2,required_workers=$3 WHERE assignment_id=$1`, id, in.LocationID, in.RequiredWorkers); e != nil {
				return e
			}
		}
		// Omission is not deletion: include every original item explicitly. There is
		// no DELETE endpoint for historical TOR requirements.
		originals := map[string]int{}
		rows, e := q.Query(ctx, `SELECT equipment_id::text,required_qty FROM requisition_item WHERE requisition_id=$1 FOR UPDATE`, req)
		if e != nil {
			return e
		}
		for rows.Next() {
			var equipment string
			var qty int
			if e = rows.Scan(&equipment, &qty); e != nil {
				rows.Close()
				return e
			}
			if _, duplicate := originals[equipment]; duplicate {
				rows.Close()
				return conflict("อุปกรณ์ตั้งต้นซ้ำ กรุณาตรวจข้อมูลก่อนแก้ไข")
			}
			originals[equipment] = qty
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			return e
		}
		provided := map[string]bool{}
		for _, it := range in.Items {
			provided[strings.ToLower(it.EquipmentID)] = true
		}
		if id != "" && in.ExpectedWorkers != nil {
			if in.ExpectedItems == nil {
				return invalid("expected_items", "กรุณาโหลดรายการอุปกรณ์ล่าสุด")
			}
			expected := map[string]int{}
			for _, it := range *in.ExpectedItems {
				expected[strings.ToLower(it.EquipmentID)] = it.Required
			}
			if len(expected) != len(originals) {
				return conflict("รายการอุปกรณ์เปลี่ยนแล้ว กรุณารีเฟรช")
			}
			for equipment, qty := range originals {
				if expected[equipment] != qty {
					return conflict("จำนวนอุปกรณ์ถูกแก้แล้ว กรุณารีเฟรช")
				}
			}
		}
		for equipment := range originals {
			if !provided[equipment] {
				return invalid("items", "ห้ามลบรายการอุปกรณ์เดิม")
			}
		}
		for _, item := range in.Items {
			if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM equipment WHERE equipment_id=$1 AND is_active)`, item.EquipmentID).Scan(&exists); e != nil {
				return e
			}
			if !exists {
				return invalid("equipment_id", "ไม่พบอุปกรณ์ที่ใช้งานได้")
			}
			tag, e := q.Exec(ctx, `UPDATE requisition_item SET required_qty=$3 WHERE requisition_id=$1 AND equipment_id=$2`, req, item.EquipmentID, item.Required)
			if e != nil {
				return e
			}
			if tag.RowsAffected() == 0 {
				if _, e = q.Exec(ctx, `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty) VALUES($1,$2,$3)`, req, item.EquipmentID, item.Required); e != nil {
					return e
				}
			}
		}
		return nil
	})
	return id, err
}

type InvoiceEditInput struct {
	Expected         string `json:"expected_amount"`
	Deduction        string `json:"exat_deduction_amount"`
	Reason           string `json:"deduction_reason"`
	PreviousExpected string `json:"previous_expected_amount"`
}

func (s *Service) UpdateInvoice(ctx context.Context, id string, in InvoiceEditInput) error {
	return s.updateInvoice(ctx, id, in, false)
}
func (s *Service) updateInvoice(ctx context.Context, id string, in InvoiceEditInput, automatic bool) error {
	if e := validID(id, "invoice_id"); e != nil {
		return e
	}
	expected, e := Money(in.Expected)
	if e != nil || expected > 999999999999 {
		return invalid("expected_amount", "จำนวนเงินไม่ถูกต้อง")
	}
	deduction, e := Money(in.Deduction)
	if e != nil || deduction > expected {
		return invalid("exat_deduction_amount", "ยอดหักเกินยอดวางบิล")
	}
	previous, e := Money(in.PreviousExpected)
	if e != nil {
		return invalid("previous_expected_amount", "กรุณาโหลดข้อมูลล่าสุด")
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		var budgetTor string
		if e := q.QueryRow(ctx, `SELECT tor_id::text FROM company_invoice WHERE invoice_id=$1`, id).Scan(&budgetTor); e != nil {
			return e
		}
		if e := s.checkInvoiceBudget(ctx, q, budgetTor, id, expected); e != nil {
			return e
		}
		var status, amount, tor, month string
		if e := q.QueryRow(ctx, `SELECT status::text,expected_amount::text,tor_id::text,billing_month FROM company_invoice WHERE invoice_id=$1 FOR UPDATE`, id).Scan(&status, &amount, &tor, &month); e != nil {
			return e
		}
		if status != "pending" {
			return conflict("รับเงินจริงแล้ว ไม่สามารถแก้ใบวางบิล")
		}
		current, _ := Money(amount)
		if current != previous {
			return conflict("ยอดถูกแก้แล้ว กรุณารีเฟรช")
		}
		if automatic {
			if e := s.lockInvoicePayroll(ctx, q, tor, month); e != nil {
				return e
			}
			preview, e := s.invoiceDeductions(ctx, q, tor, month)
			if e != nil {
				return e
			}
			if preview.Pending > 0 {
				return conflict("ยังประมวลผลค่าจ้างของโครงการ/เดือนนี้ไม่ครบ")
			}
			deduction, e = Money(preview.Total)
			if e != nil {
				return e
			}
			if deduction > expected {
				return invalid("expected_amount", "ยอดวางบิลน้อยกว่าเงินหักอัตโนมัติ")
			}
		}
		_, e := q.Exec(ctx, `UPDATE company_invoice SET expected_amount=$2::numeric,exat_deduction_amount=$3::numeric,deduction_reason=$4 WHERE invoice_id=$1`, id, Decimal(expected), Decimal(deduction), in.Reason)
		return e
	})
}

// DeleteAssignment removes only a pristine TOR area. Location master data and
// every area with staffing, surveys or financial history are preserved.
func (s *Service) DeleteAssignment(ctx context.Context, p auth.Principal, id string) error {
	if p.Role != auth.RoleSupervisor {
		return invalid("role", "Supervisor เท่านั้น")
	}
	if e := validID(id, "assignment_id"); e != nil {
		return e
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		var tor, state string
		if e := q.QueryRow(ctx, `SELECT tor_id::text FROM tor_location_assignment WHERE assignment_id=$1`, id).Scan(&tor); e != nil {
			return e
		}
		if e := lock(ctx, q, "contract-edit:"+tor); e != nil {
			return e
		}
		if e := q.QueryRow(ctx, `SELECT status::text FROM contract_tor WHERE tor_id=$1 FOR UPDATE`, tor).Scan(&state); e != nil {
			return e
		}
		if state != "registered" && state != "active" {
			return conflict("สัญญาสิ้นสุดหรือยกเลิกแล้ว")
		}
		var locked string
		if e := q.QueryRow(ctx, `SELECT assignment_id::text FROM tor_location_assignment WHERE assignment_id=$1 FOR UPDATE`, id).Scan(&locked); e != nil {
			return e
		}
		// Serialize with survey/funding/purchase writers before deciding whether
		// cascading deletion is safe. New child inserts also wait on the area lock.
		requests, e := q.Query(ctx, `SELECT requisition_id FROM equipment_requisition WHERE assignment_id=$1 ORDER BY requisition_id FOR UPDATE`, id)
		if e != nil {
			return e
		}
		for requests.Next() {
		}
		e = requests.Err()
		requests.Close()
		if e != nil {
			return e
		}
		var used bool
		if e := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule WHERE assignment_id=$1) OR EXISTS(SELECT 1 FROM equipment_requisition WHERE assignment_id=$1 AND (requisition_type<>'tor_base' OR status<>'pending_survey')) OR EXISTS(SELECT 1 FROM requisition_item i JOIN equipment_requisition r USING(requisition_id) WHERE r.assignment_id=$1 AND (i.existing_qty<>0 OR i.actual_qty<>0)) OR EXISTS(SELECT 1 FROM expense_claim x JOIN equipment_requisition r USING(requisition_id) WHERE r.assignment_id=$1)`, id).Scan(&used); e != nil {
			return e
		}
		if used {
			return conflict("ลบไม่ได้: พื้นที่มีกะงาน การสำรวจ หรือรายการจัดซื้อแล้ว")
		}
		_, e = q.Exec(ctx, `DELETE FROM tor_location_assignment WHERE assignment_id=$1`, id)
		return e
	})
}
