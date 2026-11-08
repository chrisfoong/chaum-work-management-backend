package work

import (
	"context"
	"time"
)

// 6S: project invoices and receipt confirmation, separate from workforce payroll.
type InvoiceInput struct {
	TorID     string `json:"tor_id"`
	Month     string `json:"billing_month"`
	Expected  string `json:"expected_amount"`
	Deduction string `json:"exat_deduction_amount"`
	Reason    string `json:"deduction_reason"`
}

func Month(v string) error {
	if _, e := time.Parse("2006-01", v); e != nil {
		return invalid("billing_month", "YYYY-MM required")
	}
	return nil
}
func (s *Service) Invoice(ctx context.Context, in InvoiceInput) (string, error) {
	return s.invoice(ctx, in, false)
}
func (s *Service) invoice(ctx context.Context, in InvoiceInput, automatic bool) (string, error) {
	if e := validID(in.TorID, "tor_id"); e != nil {
		return "", e
	}
	if e := Month(in.Month); e != nil {
		return "", e
	}
	a, e := Money(in.Expected)
	if e != nil {
		return "", e
	}
	if in.Deduction == "" {
		in.Deduction = "0"
	}
	d, e := Money(in.Deduction)
	if e != nil || d > a || a > 999999999999 {
		return "", invalid("amount", "invalid deduction or schema precision")
	}
	var id string
	e = s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "invoice:"+in.TorID+":"+in.Month); e != nil {
			return e
		}
		if e := s.checkInvoiceBudget(ctx, q, in.TorID, "", a); e != nil {
			return e
		}
		if automatic {
			if e := s.lockInvoicePayroll(ctx, q, in.TorID, in.Month); e != nil {
				return e
			}
			preview, e := s.invoiceDeductions(ctx, q, in.TorID, in.Month)
			if e != nil {
				return e
			}
			if preview.Pending > 0 {
				return conflict("ยังประมวลผลค่าจ้างของโครงการ/เดือนนี้ไม่ครบ กรุณาปิดงวดและประมวลผลก่อนสร้างใบวางบิล")
			}
			d, e = Money(preview.Total)
			if e != nil {
				return e
			}
			if d > a {
				return invalid("expected_amount", "ยอดวางบิลน้อยกว่าเงินหักอัตโนมัติ")
			}
		}
		var exists bool
		if e := s.Repo.Row(ctx, q, query37, in.TorID, in.Month).Scan(&exists); e != nil {
			return e
		}
		if exists {
			return conflict("invoice already exists for TOR/month")
		}
		return s.Repo.Row(ctx, q, query38, document("INV"), in.TorID, in.Month, Decimal(a), Decimal(d), in.Reason).Scan(&id)
	})
	return id, e
}
func (s *Service) Receive(ctx context.Context, id, amount string) error {
	if e := validID(id, "invoice_id"); e != nil {
		return e
	}
	c, e := Money(amount)
	if e != nil || c > 999999999999 {
		return invalid("net_received", "invalid amount")
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		var state string
		var current *string
		if e := s.Repo.Row(ctx, q, query39, id).Scan(&state, &current); e != nil {
			return e
		}
		if state == "paid" {
			if current != nil && *current == Decimal(c) {
				return nil
			}
			return conflict("paid invoice cannot be overwritten")
		}
		_, e := s.Repo.Exec(ctx, q, query40, id, Decimal(c))
		return e
	})
}
