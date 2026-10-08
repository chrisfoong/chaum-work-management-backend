package work

import (
	"context"

	"log/slog"
	"time"
)

// Events are best effort after commit. There is no durable outbox in the schema.
func (s *Service) Event(kind, id string) {
	if s.Notify == nil {
		slog.Info("notification skipped", "event", kind)
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var sql string
		switch kind {
		case "schedule":
			sql = query23
		case "leave":
			sql = query24
		case "leave_review":
			sql = query25
		case "delivery", "request_review":
			sql = query26
		case "approval_needed", "purchase_funding", "replacement_needed":
			sql = `SELECT user_id::text FROM public."USER" WHERE role='supervisor' AND is_active AND $1::uuid IS NOT NULL`
		case "funded", "continuation":
			sql = query27
		case "request":
			sql = query27
		case "payroll":
			sql = query28
		default:
			return
		}
		rows, e := s.Repo.Rows(ctx, s.Repo.Pool, sql, id)
		if e != nil {
			slog.Error("notification recipients unavailable", "event", kind)
			return
		}
		var ids []string
		for rows.Next() {
			var user string
			if rows.Scan(&user) == nil {
				ids = append(ids, user)
			}
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			slog.Error("notification recipient query failed", "event", kind)
			return
		}
		for _, user := range ids {
			message := "Chaum: " + kind + " updated. Please open the app to view details."
			if kind == "request_review" {
				var state, reason string
				if s.Repo.Pool.QueryRow(ctx, `SELECT status::text,reason FROM equipment_requisition WHERE requisition_id=$1`, id).Scan(&state, &reason) == nil && state == "rejected" {
					message = "Chaum: equipment request rejected. " + reason
				}
			}
			s.notify(ctx, user, message)
		}
	}()
}
func (s *Service) Send(ctx context.Context, line string, msg string) error {
	if s.Notify == nil {
		return conflict("LINE notifications disabled")
	}
	return s.Notify.Send(ctx, line, msg)
}
