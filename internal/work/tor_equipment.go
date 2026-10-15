package work

import (
	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
)

// AssignmentEquipment lists active initial TOR equipment without changing additional-request rules.
func (s *Service) AssignmentEquipment(ctx context.Context, p auth.Principal, assignment string, limit, offset int) (json.RawMessage, error) {
	if assignment == "" {
		return s.Equipments(ctx, limit, offset)
	}
	if e := validID(assignment, "assignment_id"); e != nil {
		return nil, e
	}
	if p.Role == auth.RoleWorker {
		worker, e := s.Worker(ctx, s.Repo.Pool, p.UserID.String())
		if e != nil {
			return nil, e
		}
		var own bool
		if e = s.Repo.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM work_schedule WHERE worker_id=$1 AND assignment_id=$2 AND shift_status<>'cancelled')`, worker, assignment).Scan(&own); e != nil {
			return nil, e
		}
		if !own {
			return nil, apperr.Forbidden("assignment is not assigned to worker")
		}
	}
	return s.Repo.List(ctx, s.Repo.Pool, `SELECT e.equipment_id,e.equipment_name,e.is_active FROM equipment e WHERE e.is_active AND EXISTS(SELECT 1 FROM requisition_item i JOIN equipment_requisition r ON r.requisition_id=i.requisition_id JOIN tor_location_assignment a ON a.assignment_id=r.assignment_id JOIN tor_location_assignment selected ON selected.tor_id=a.tor_id WHERE selected.assignment_id=$1 AND r.requisition_type='tor_base' AND i.equipment_id=e.equipment_id) ORDER BY e.equipment_name,e.equipment_id LIMIT $2 OFFSET $3`, assignment, limit, offset)
}
