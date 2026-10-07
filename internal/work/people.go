package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"

	"regexp"
)

type UserInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
	Phone       string `json:"phone_number"`
	LineID      string `json:"line_id"`
	BankName    string `json:"bank_name"`
	BankAccount string `json:"bank_account_no"`
	DailyWage   string `json:"daily_wage"`
}

func (s *Service) CreateUser(ctx context.Context, p auth.Principal, in UserInput) (json.RawMessage, error) {
	for _, f := range []struct {
		v, n string
		max  int
	}{{in.FirstName, "first_name", 100}, {in.LastName, "last_name", 100}, {in.LineID, "line_id", 50}, {in.BankName, "bank_name", 100}, {in.BankAccount, "bank_account_no", 50}} {
		if e := text(f.v, f.n, f.max); e != nil {
			return nil, e
		}
	}
	if !regexp.MustCompile(`^[0-9]{10}$`).MatchString(in.Phone) {
		return nil, invalid("phone_number", "10 digits required")
	}
	if in.Role != "supervisor" && in.Role != "assistant" && in.Role != "worker" {
		return nil, invalid("role", "unknown role")
	}
	if in.DailyWage == "" {
		in.DailyWage = "400.00"
	}
	w, e := Money(in.DailyWage)
	if e != nil || w > 9999999999 {
		return nil, invalid("daily_wage", "invalid wage")
	}
	var out json.RawMessage
	e = s.Repo.Transaction(ctx, func(q Query) error {
		var id string
		e := s.Repo.Row(ctx, q, query41, in.FirstName, in.LastName, in.Role, in.Phone, in.LineID, in.BankName, in.BankAccount, Decimal(w)).Scan(&id)
		if e != nil {
			return e
		}
		if in.Role == "worker" {
			if e = lock(ctx, q, "worker:"+id); e != nil {
				return e
			}
			_, e = s.Repo.Exec(ctx, q, query42, id)
			if e != nil {
				return e
			}
		}
		out, e = one(ctx, q, query43, id)
		return e
	})
	return out, e
}

type UserUpdate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DailyWage string `json:"daily_wage"`
	Active    *bool  `json:"is_active"`
	Available *bool  `json:"is_available"`
}

func (s *Service) UpdateUser(ctx context.Context, id string, in UserUpdate) error {
	if e := validID(id, "user_id"); e != nil {
		return e
	}
	if e := text(in.FirstName, "first_name", 100); e != nil {
		return e
	}
	if e := text(in.LastName, "last_name", 100); e != nil {
		return e
	}
	w, e := Money(in.DailyWage)
	if e != nil || w > 9999999999 {
		return invalid("daily_wage", "invalid wage")
	}
	return s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "worker:"+id); e != nil {
			return e
		}
		tag, e := s.Repo.Exec(ctx, q, query44, id, in.FirstName, in.LastName, Decimal(w), in.Active)
		if e != nil {
			return e
		}
		if tag.RowsAffected() == 0 {
			return conflict("user missing")
		}
		if in.Available != nil {
			worker, e := s.Worker(ctx, q, id)
			if e != nil {
				return e
			}
			_, e = s.Repo.Exec(ctx, q, query45, worker, *in.Available)
			return e
		}
		return nil
	})
}
func (s *Service) Me(ctx context.Context, p auth.Principal) (json.RawMessage, error) {
	return one(ctx, s.Repo.Pool, query46, p.UserID)
}
func (s *Service) Users(ctx context.Context, limit, offset int) (json.RawMessage, error) {
	return s.Repo.List(ctx, s.Repo.Pool, query47, limit, offset)
}

type ScheduleInput struct {
	AssignmentID string   `json:"assignment_id"`
	WorkerIDs    []string `json:"worker_ids"`
	WorkDate     string   `json:"work_date"`
	Start        string   `json:"shift_start_time"`
}

func (s *Service) Schedule(ctx context.Context, in ScheduleInput) ([]string, error) {
	if e := validID(in.AssignmentID, "assignment_id"); e != nil {
		return nil, e
	}
	if _, e := Start(in.WorkDate, in.Start); e != nil {
		return nil, e
	}
	if len(in.WorkerIDs) == 0 || len(in.WorkerIDs) > 100 {
		return nil, invalid("worker_ids", "1-100 workers required")
	}
	seen := map[string]bool{}
	for _, id := range in.WorkerIDs {
		if e := validID(id, "worker_ids"); e != nil {
			return nil, e
		}
		if seen[id] {
			return nil, invalid("worker_ids", "duplicate worker")
		}
		seen[id] = true
	}
	var out []string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "staffing:"+in.AssignmentID+":"+in.WorkDate); e != nil {
			return e
		}
		var capacity int
		var start, end string
		var state string
		if e := s.Repo.Row(ctx, q, query48, in.AssignmentID).Scan(&capacity, &start, &end, &state); e != nil {
			return e
		}
		if state == "cancelled" || state == "complete" || in.WorkDate < start || in.WorkDate > end {
			return conflict("outside contract dates or contract closed")
		}
		var count int
		if e := s.Repo.Row(ctx, q, query49, in.AssignmentID, in.WorkDate).Scan(&count); e != nil {
			return e
		}
		if count+len(in.WorkerIDs) > capacity {
			return conflict("required_workers exceeded")
		}
		for _, worker := range in.WorkerIDs {
			var user string
			var available, active bool
			if e := s.Repo.Row(ctx, q, query50, worker).Scan(&user, &available, &active); e != nil {
				return e
			}
			id, e := s.Worker(ctx, q, user)
			if e != nil {
				return e
			}
			if id != worker || !available || !active {
				return conflict("worker unavailable")
			}
			if e = lock(ctx, q, "payroll:"+user); e != nil {
				return e
			}
			var paidPeriod bool
			if e = s.Repo.Row(ctx, q, query51, user, in.WorkDate).Scan(&paidPeriod); e != nil {
				return e
			}
			if paidPeriod {
				return conflict("cannot change staffing in a calculated payroll period")
			}
			var newID string
			if e = s.Repo.Row(ctx, q, query52, in.AssignmentID, worker, in.WorkDate, in.Start).Scan(&newID); e != nil {
				return e
			}
			out = append(out, newID)
		}
		return nil
	})
	if e == nil {
		for _, id := range out {
			s.Event("schedule", id)
		}
	}
	return out, e
}
func (s *Service) Schedules(ctx context.Context, p auth.Principal, limit, offset int) (json.RawMessage, error) {
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return s.Repo.List(ctx, s.Repo.Pool, query53, owner, limit, offset)
}

type LeaveInput struct {
	Date   string `json:"leave_date"`
	Reason string `json:"reason"`
}

func (s *Service) Leave(ctx context.Context, p auth.Principal, in LeaveInput) (string, error) {
	if _, e := Date(in.Date); e != nil {
		return "", e
	}
	if e := text(in.Reason, "reason", 500); e != nil {
		return "", e
	}
	var id string
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "leave:"+p.UserID.String()+":"+in.Date); e != nil {
			return e
		}
		worker, e := s.Worker(ctx, q, p.UserID.String())
		if e != nil {
			return e
		}
		var clock string
		var count int
		if e = s.Repo.Row(ctx, q, query54, p.UserID, in.Date).Scan(&count); e != nil {
			return e
		}
		if count > 0 {
			return conflict("leave already exists for date")
		}
		if e = s.Repo.Row(ctx, q, query55, worker, in.Date).Scan(&clock); e != nil {
			return e
		}
		start, e := Start(in.Date, clock)
		if e != nil {
			return e
		}
		if start.Sub(s.Now()) < 24*60*60*1000000000 {
			return invalid("leave_date", "at least 24 hours notice required")
		}
		return s.Repo.Row(ctx, q, query56, document("LEAVE"), p.UserID, in.Date, in.Reason).Scan(&id)
	})
	if e == nil {
		s.Event("leave", id)
	}
	return id, e
}
func (s *Service) Leaves(ctx context.Context, p auth.Principal, limit, offset int) (json.RawMessage, error) {
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	return s.Repo.List(ctx, s.Repo.Pool, query57, owner, limit, offset)
}
func (s *Service) ReviewLeave(ctx context.Context, id, state string) error {
	if e := validID(id, "request_id"); e != nil {
		return e
	}
	if state != "approved" && state != "rejected" {
		return invalid("status", "approved or rejected required")
	}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		var user, date string
		if e := s.Repo.Row(ctx, q, query58, id).Scan(&user, &date); e != nil {
			return e
		}
		worker, e := s.Worker(ctx, q, user)
		if e != nil {
			return e
		}
		if e = lock(ctx, q, "payroll:"+user); e != nil {
			return e
		}
		var calculated bool
		if e = s.Repo.Row(ctx, q, query59, user, date).Scan(&calculated); e != nil {
			return e
		}
		if calculated {
			return conflict("leave period already calculated")
		}
		var schedule string
		if e = s.Repo.Row(ctx, q, query60, worker, date).Scan(&schedule); e != nil {
			return e
		}
		if state == "approved" {
			var checked bool
			if e = s.Repo.Row(ctx, q, query61, schedule).Scan(&checked); e != nil {
				return e
			}
			if checked {
				return conflict("cannot approve leave after check-in")
			}
		}
		tag, e := s.Repo.Exec(ctx, q, query62, id, state)
		if e != nil {
			return e
		}
		if tag.RowsAffected() != 1 {
			return conflict("leave already reviewed")
		}
		if state == "approved" {
			_, e = s.Repo.Exec(ctx, q, query63, schedule, worker, date)
			return e
		}
		return nil
	})
	if err == nil {
		s.Event("leave_review", id)
	}
	return err
}
func (s *Service) Replacement(ctx context.Context, leaveID, workerID string) ([]string, error) {
	if e := validID(leaveID, "request_id"); e != nil {
		return nil, e
	}
	if e := validID(workerID, "worker_id"); e != nil {
		return nil, e
	}
	var in ScheduleInput
	var original string
	var state string
	var user string
	if e := s.Repo.Row(ctx, s.Repo.Pool, `SELECT user_id::text FROM leave_request WHERE request_id=$1`, leaveID).Scan(&user); e != nil {
		return nil, e
	}
	if _, e := s.Worker(ctx, s.Repo.Pool, user); e != nil {
		return nil, e
	}
	e := s.Repo.Row(ctx, s.Repo.Pool, query64, leaveID).Scan(&state, &in.AssignmentID, &in.WorkDate, &in.Start, &original)
	if e != nil {
		return nil, e
	}
	if state != "approved" || workerID == original {
		return nil, conflict("approved leave and a different worker required")
	}
	in.WorkerIDs = []string{workerID}
	return s.Schedule(ctx, in)
}
