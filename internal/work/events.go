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
		case "funded":
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
			if kind == "leave_review" {
				var err error
				message, err = s.leaveReviewMessage(ctx, id)
				if err != nil {
					slog.Error("leave notification details unavailable")
					continue
				}
			}
			if kind == "delivery" {
				var summary string
				if e := s.Repo.Pool.QueryRow(ctx, `SELECT 'Chaum: ส่งมอบอุปกรณ์คำขอ ' || r.requisition_no || E'\n' || COALESCE(string_agg(e.equipment_name || ' × ' || COALESCE(i.actual_qty,0)::text,E'\n' ORDER BY i.item_id),'') FROM equipment_requisition r JOIN requisition_item i USING(requisition_id) JOIN equipment e USING(equipment_id) WHERE r.requisition_id=$1 GROUP BY r.requisition_no`, id).Scan(&summary); e != nil {
					slog.Error("delivery summary unavailable")
					continue
				}
				message = summary
			}
			if kind == "schedule" || kind == "payroll" || kind == "request_review" {
				var err error
				switch kind {
				case "schedule":
					message, err = s.scheduleMessage(ctx, id)
				case "payroll":
					message, err = s.payrollMessage(ctx, id)
				case "request_review":
					message, err = s.requestReviewMessage(ctx, id)
				}
				if err != nil {
					slog.Error("worker notification details unavailable", "event", kind)
					continue
				}
			}

			s.notifyKey(ctx, user, message, kind+":"+id)
		}
	}()
}
func (s *Service) Send(ctx context.Context, line string, msg string) error {
	if s.Notify == nil {
		return conflict("LINE notifications disabled")
	}
	return s.Notify.Send(ctx, line, msg)
}
