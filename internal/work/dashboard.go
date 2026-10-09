package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
)

func (s *Service) Dashboard(ctx context.Context, p auth.Principal) (json.RawMessage, error) {
	if p.Role == auth.RoleWorker {
		return one(ctx, s.Repo.Pool, query21, p.UserID, s.Now())
	}
	return one(ctx, s.Repo.Pool, frontendDashboard, s.Now())
}
