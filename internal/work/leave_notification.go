package work

import (
	"context"
	"fmt"
)

// Build the notification from committed leave state, shared by initial send
// and guarded resend. Never infer the decision from a client-provided label.
func (s *Service) leaveReviewMessage(ctx context.Context, id string) (string, error) {
	var number, date, reason, state string
	if err := s.Repo.Pool.QueryRow(ctx, `SELECT leave_no,leave_date::text,reason,status::text FROM leave_request WHERE request_id=$1`, id).Scan(&number, &date, &reason, &state); err != nil {
		return "", err
	}
	return formatLeaveReviewMessage(number, date, reason, state)
}

func formatLeaveReviewMessage(number, date, reason, state string) (string, error) {
	if state != "approved" && state != "rejected" {
		return "", conflict("leave review is not committed")
	}
	return fmt.Sprintf("Chaum: leave %s.\nเลขใบลา: %s\nวันที่ลา: %s\nเหตุผลการลา: %s", state, number, date, reason), nil
}
