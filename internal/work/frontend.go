package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// ListFilter filters in SQL before pagination, never just the current client page.
type ListFilter struct{ Status, AssignmentID, TorID, From, To, Kind string }

func (f ListFilter) validate(kind string) error {
	for _, v := range []struct{ value, name string }{{f.AssignmentID, "assignment_id"}, {f.TorID, "tor_id"}} {
		if v.value != "" {
			if e := validID(v.value, v.name); e != nil {
				return e
			}
		}
	}
	for _, v := range []string{f.From, f.To} {
		if v != "" {
			if _, e := Date(v); e != nil {
				return e
			}
		}
	}
	if f.From != "" && f.To != "" && f.From > f.To {
		return invalid("period_end", "must not precede period_start")
	}
	allowed := map[string][]string{"schedules": {"scheduled", "completed", "cancelled"}, "leave": {"pending", "approved", "rejected"}, "requisitions": {"pending_survey", "pending_procurement", "pending_approval", "pending_fund", "completed", "rejected"}, "attendance": {"on_time", "late", "absent", "leave"}}
	if f.Status != "" {
		found := false
		for _, v := range allowed[kind] {
			found = found || v == f.Status
		}
		if !found {
			return invalid("status", "unknown status")
		}
	}
	if f.Kind != "" && (kind != "requisitions" || (f.Kind != "tor_base" && f.Kind != "additional")) {
		return invalid("requisition_type", "tor_base or additional required")
	}
	return nil
}
func filteredList(s *Service, kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		f := ListFilter{c.Query("status"), c.Query("assignment_id"), c.Query("tor_id"), c.Query("period_start"), c.Query("period_end"), c.Query("requisition_type")}
		out, e := s.FilteredList(c.Request.Context(), actor(c), kind, f, limit, offset)
		respond(c, gin.H{"data": out, "limit": limit, "offset": offset}, e)
	}
}
func (s *Service) FilteredList(ctx context.Context, p auth.Principal, kind string, f ListFilter, limit, offset int) (json.RawMessage, error) {
	if e := f.validate(kind); e != nil {
		return nil, e
	}
	var owner any
	if p.Role == auth.RoleWorker {
		owner = p.UserID
	}
	// UUID parameters are NULL when absent; textual filters are optional empty strings.
	args := []any{owner, limit, offset, f.Status, nullableID(f.AssignmentID), nullableID(f.TorID), f.From, f.To}
	sql := ""
	switch kind {
	case "schedules":
		sql = frontendSchedules
		args = append(args, s.Now())
	case "leave":
		sql = frontendLeaves
	case "attendance":
		sql = frontendAttendance
	case "requisitions":
		sql = frontendRequisitions
		args = append(args, f.Kind)
	default:
		return nil, invalid("list", "unknown list")
	}
	return s.Repo.List(ctx, s.Repo.Pool, sql, args...)
}

func (s *Service) AvailableWorkers(ctx context.Context, date, excludeUser string, limit, offset int) (json.RawMessage, error) {
	if _, e := Date(date); e != nil {
		return nil, e
	}
	return s.availableWorkers(ctx, s.Repo.Pool, date, excludeUser, limit, offset)
}
func (s *Service) availableWorkers(ctx context.Context, q Query, date, excludeUser string, limit, offset int) (json.RawMessage, error) {
	// Validate every eligible mapping before paging so duplicates cannot be hidden.
	rows, e := q.Query(ctx, frontendCandidateUsers, date, nullableID(excludeUser))
	if e != nil {
		return nil, e
	}
	var users []string
	for rows.Next() {
		var user string
		if e = rows.Scan(&user); e != nil {
			rows.Close()
			return nil, e
		}
		users = append(users, user)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return nil, e
	}
	for _, user := range users {
		if _, e = s.Worker(ctx, q, user); e != nil {
			return nil, e
		}
	}
	return s.Repo.List(ctx, q, frontendCandidates, date, nullableID(excludeUser), limit, offset)
}
func (s *Service) LeaveCandidates(ctx context.Context, id string, limit, offset int) (json.RawMessage, error) {
	if e := validID(id, "request_id"); e != nil {
		return nil, e
	}
	var out json.RawMessage
	e := s.Repo.Transaction(ctx, func(q Query) error {
		var user, date, state string
		if e := q.QueryRow(ctx, `SELECT user_id::text,leave_date::text,status::text FROM leave_request WHERE request_id=$1 FOR SHARE`, id).Scan(&user, &date, &state); e != nil {
			return e
		}
		if state == "rejected" {
			return conflict("rejected leave cannot have a replacement")
		}
		if _, e := s.Worker(ctx, q, user); e != nil {
			return e
		}
		var e error
		out, e = s.availableWorkers(ctx, q, date, user, limit, offset)
		return e
	})
	return out, e
}
func (s *Service) LeaveDetail(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "request_id"); e != nil {
		return nil, e
	}
	var user string
	if e := s.Repo.Pool.QueryRow(ctx, `SELECT user_id::text FROM leave_request WHERE request_id=$1`, id).Scan(&user); e != nil {
		return nil, e
	}
	if _, e := s.Worker(ctx, s.Repo.Pool, user); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, frontendLeaveDetail, id)
}
func (s *Service) Assignments(ctx context.Context, tor string, limit, offset int) (json.RawMessage, error) {
	if tor != "" {
		if e := validID(tor, "tor_id"); e != nil {
			return nil, e
		}
	}
	return s.Repo.List(ctx, s.Repo.Pool, frontendAssignments, nullableID(tor), s.Now(), limit, offset)
}
func (s *Service) DeliverySchedules(ctx context.Context, id string, limit, offset int) (json.RawMessage, error) {
	if e := validID(id, "requisition_id"); e != nil {
		return nil, e
	}
	var owner, assignment, state, kind string
	if e := s.Repo.Pool.QueryRow(ctx, `SELECT requested_by::text,assignment_id::text,status::text,requisition_type::text FROM equipment_requisition WHERE requisition_id=$1`, id).Scan(&owner, &assignment, &state, &kind); e != nil {
		return nil, e
	}
	if state != "completed" || kind != "additional" {
		return nil, conflict("only completed additional requests may be delivered")
	}
	worker, e := s.Worker(ctx, s.Repo.Pool, owner)
	if e != nil {
		return nil, e
	}
	return s.Repo.List(ctx, s.Repo.Pool, `SELECT `+frontendScheduleFields+frontendScheduleJoins+`WHERE s.worker_id=$1 AND s.assignment_id=$2 AND s.work_date<=$3::date AND s.shift_status<>'cancelled' ORDER BY s.work_date DESC,s.schedule_id LIMIT $4 OFFSET $5`, worker, assignment, s.Now().In(Bangkok).Format("2006-01-02"), limit, offset)
}

// LINE text messages accept at most 5000 UTF-16 units, including surrogate pairs.
func notificationParts(message string) []string {
	parts := []string{}
	var runes []rune
	units := 0
	for _, r := range message {
		size := 1
		if r > 0xffff {
			size = 2
		}
		if units+size > 5000 {
			parts = append(parts, string(runes))
			runes = nil
			units = 0
		}
		runes = append(runes, r)
		units += size
	}
	if len(runes) > 0 {
		parts = append(parts, string(runes))
	}
	return parts
}
