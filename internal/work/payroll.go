package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"time"
)

type PayrollInput struct {
	UserID string `json:"user_id"`
	Start  string `json:"period_start"`
	End    string `json:"period_end"`
}
type Deduction struct {
	AttendanceID string `json:"attendance_id"`
	WorkDate     string `json:"work_date"`
	Penalty      string `json:"calculated_penalty"`
	Reason       string `json:"reason"`
	Cents        int64  `json:"-"`
}
type PayrollPreview struct {
	UserID             string      `json:"user_id"`
	Start              string      `json:"period_start"`
	End                string      `json:"period_end"`
	Workdays           int         `json:"workdays"`
	Base               string      `json:"base_wage"`
	Calculated         string      `json:"calculated_penalty"`
	Applied            string      `json:"total_deduction"`
	Net                string      `json:"net_wage"`
	Deductions         []Deduction `json:"deductions"`
	AttendanceComplete bool        `json:"attendance_complete"`
}

func (s *Service) preview(ctx context.Context, q Query, in PayrollInput, locking bool) (PayrollPreview, error) {
	out := PayrollPreview{UserID: in.UserID, Start: in.Start, End: in.End, Deductions: []Deduction{}, AttendanceComplete: true}
	worker, e := s.Worker(ctx, q, in.UserID)
	if e != nil {
		return out, e
	}
	var wage string
	if e = s.Repo.Row(ctx, q, query29, in.UserID).Scan(&wage); e != nil {
		return out, e
	}
	daily, e := Money(wage)
	if e != nil {
		return out, e
	}
	sql := query30
	if locking {
		sql += " FOR UPDATE OF s"
	}
	rows, e := s.Repo.Rows(ctx, q, sql, in.UserID, worker, in.Start, in.End)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	var total int64
	for rows.Next() {
		var schedule, date, clock string
		var id, status *string
		var check *time.Time
		var leave bool
		if e = rows.Scan(&schedule, &date, &clock, &id, &status, &check, &leave); e != nil {
			return out, e
		}
		start, e := Start(date, clock)
		if e != nil {
			return out, e
		}
		if !leave && (id == nil || status == nil) {
			out.AttendanceComplete = false
			continue
		}
		if locking && s.Now().Before(start.Add(10*time.Hour)) {
			return out, conflict("payroll shift cutoff has not passed")
		}
		if leave {
			continue
		}
		switch *status {
		case "on_time", "late":
			out.Workdays++
		case "absent", "leave":
		default:
			return out, conflict("unknown attendance status")
		}
		penalty := Penalty(*status, start, check)
		if penalty > 0 {
			total += penalty
			out.Deductions = append(out.Deductions, Deduction{*id, date, Decimal(penalty), *status, penalty})
		}
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	base := int64(out.Workdays) * daily
	if base > 9999999999 {
		return out, conflict("base wage exceeds schema precision")
	}
	applied := total
	if applied > base {
		applied = base
	}
	out.Base = Decimal(base)
	out.Calculated = Decimal(total)
	out.Applied = Decimal(applied)
	out.Net = Decimal(base - applied)
	return out, nil
}
func (s *Service) PreviewPayroll(ctx context.Context, in PayrollInput) (PayrollPreview, error) {
	if e := validID(in.UserID, "user_id"); e != nil {
		return PayrollPreview{}, e
	}
	if e := Period(in.Start, in.End); e != nil {
		return PayrollPreview{}, e
	}
	return s.preview(ctx, s.Repo.Pool, in, false)
}
func (s *Service) Payroll(ctx context.Context, p auth.Principal, in PayrollInput) (string, error) {
	if e := validID(in.UserID, "user_id"); e != nil {
		return "", e
	}
	if e := Period(in.Start, in.End); e != nil {
		return "", e
	}
	end, _ := Date(in.End)
	if s.Now().Before(end.AddDate(0, 0, 1)) {
		return "", conflict("payroll period has not closed")
	}
	var id string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "payroll:"+in.UserID); e != nil {
			return e
		}
		var overlapping bool
		if e := s.Repo.Row(ctx, q, query31, in.UserID, in.Start, in.End).Scan(&overlapping); e != nil {
			return e
		}
		if overlapping {
			return conflict("payroll period overlaps an existing slip")
		}
		out, e := s.preview(ctx, q, in, true)
		if e != nil {
			return e
		}
		if !out.AttendanceComplete {
			return conflict("finalize attendance before payroll")
		}
		if e = s.Repo.Row(ctx, q, query32, document("PAY"), in.UserID, p.UserID, in.Start, in.End, out.Base, out.Applied, out.Net).Scan(&id); e != nil {
			return e
		}
		worker, e := s.Worker(ctx, q, in.UserID)
		if e != nil {
			return e
		}
		for _, d := range out.Deductions {
			var count int
			if e = s.Repo.Row(ctx, q, query33, d.AttendanceID).Scan(&count); e != nil {
				return e
			}
			if count != 0 {
				return conflict("attendance already has a penalty ledger entry")
			}
			if _, e = s.Repo.Exec(ctx, q, query34, worker, d.AttendanceID, id, d.Penalty, d.Reason); e != nil {
				return e
			}
		}
		return nil
	})
	if e == nil {
		s.Event("payroll", id)
	}
	return id, e
}
func (s *Service) Payrolls(ctx context.Context, p auth.Principal, limit, offset int) (json.RawMessage, error) {
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return s.Repo.List(ctx, s.Repo.Pool, query35, owner, limit, offset)
}
func (s *Service) Pay(ctx context.Context, id string) error {
	if e := validID(id, "payroll_id"); e != nil {
		return e
	}
	tag, e := s.Repo.Exec(ctx, s.Repo.Pool, query36, id)
	if e == nil && tag.RowsAffected() != 1 {
		return conflict("payroll not found")
	}
	return e
}

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
