package work

import (
	"context"
	"encoding/json"
)

func (s *Service) Continuation(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "tor_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT c.tor_id,c.contract_no,c.project_name,c.start_date,c.end_date,c.status,
 COALESCE((SELECT jsonb_agg(jsonb_build_object('assignment_id',a.assignment_id,'location_id',a.location_id,'location_name',l.location_name,'required_workers',a.required_workers)) FROM tor_location_assignment a JOIN location l USING(location_id) WHERE a.tor_id=c.tor_id),'[]'::jsonb) AS assignments
 FROM contract_tor c WHERE c.tor_id=$1 AND c.status='active' AND c.end_date >= $2::date`, id, s.Now().In(Bangkok).Format("2006-01-02"))
}

// AssignmentContinuation is a read model for the area selected in 9A.
// Notification history cannot be reconstructed from the unchanged schema.
func (s *Service) AssignmentContinuation(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "assignment_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, `SELECT c.tor_id,c.contract_no,c.project_name,c.start_date,c.end_date,c.status,a.assignment_id,a.required_workers,l.location_id,l.location_name,l.address,
 (c.status='active' AND $2::date BETWEEN c.start_date AND c.end_date) AS can_continue,
 CASE WHEN c.status='active' AND $2::date BETWEEN c.start_date AND c.end_date THEN '' ELSE 'กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ' END AS notice,
 false AS summary_notification_history_available
 FROM tor_location_assignment a JOIN contract_tor c USING(tor_id) JOIN location l USING(location_id) WHERE a.assignment_id=$1`, id, s.Now().In(Bangkok).Format("2006-01-02"))
}
