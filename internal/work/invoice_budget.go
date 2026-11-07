package work

import "context"

// Lock the contract before invoice rows so concurrent months share one budget.
func (s *Service) checkInvoiceBudget(ctx context.Context, q Query, tor, exclude string, amount int64) error {
	var value string
	if e := q.QueryRow(ctx, `SELECT contract_value::text FROM contract_tor WHERE tor_id=$1 FOR UPDATE`, tor).Scan(&value); e != nil {
		return e
	}
	budget, e := Money(value)
	if e != nil {
		return e
	}
	var used string
	if e = q.QueryRow(ctx, `SELECT COALESCE(sum(expected_amount),0)::text FROM company_invoice WHERE tor_id=$1 AND invoice_id::text<>$2`, tor, exclude).Scan(&used); e != nil {
		return e
	}
	billed, e := Money(used)
	if e != nil {
		return e
	}
	if amount <= 0 || amount > budget-billed {
		return invalid("expected_amount", "ยอดวางบิลต้องมากกว่า 0 และยอดสะสมต้องไม่เกินมูลค่าสัญญา คงเหลือ "+Decimal(budget-billed)+" บาท")
	}
	return nil
}
