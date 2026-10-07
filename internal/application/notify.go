package application

import (
	"chrisfoong/chaum-work-management-backend/internal/notify"
	"chrisfoong/chaum-work-management-backend/internal/work"
	"context"
)

type contractNotify struct{ s *work.Service }

func (n contractNotify) Send(ctx context.Context, m notify.Message) (notify.Status, error) {
	if n.s.Notify == nil {
		return notify.StatusSkipped, nil
	}
	if e := n.s.Send(ctx, m.LineUserID, m.Text); e != nil {
		return notify.StatusFailed, e
	}
	return notify.StatusSent, nil
}
