package work

// Exported-schema SQL for catalog.
const (
	query11 = `SELECT tor_id,contract_no,project_name,partner_agency,start_date,end_date,contract_value::text,status,CASE WHEN status IN ('complete','cancelled') OR end_date < $3::date THEN 'ended' WHEN status='active' AND start_date<=$3::date THEN 'active' ELSE 'registered' END AS workflow_status,(status='active' AND $3::date BETWEEN start_date AND end_date) AS can_operate FROM contract_tor ORDER BY created_at DESC,tor_id LIMIT $1 OFFSET $2`
	query12 = `SELECT c.tor_id,c.contract_no,c.project_name,c.partner_agency,c.start_date,c.end_date,c.contract_value::text,c.contract_file_url,c.status,COALESCE((SELECT jsonb_agg(jsonb_build_object('assignment_id',a.assignment_id,'required_workers',a.required_workers,'location_id',l.location_id,'location_name',l.location_name,'latitude',l.latitude::text,'longitude',l.longitude::text)) FROM tor_location_assignment a JOIN location l ON l.location_id=a.location_id WHERE a.tor_id=c.tor_id),'[]'::jsonb) AS assignments FROM contract_tor c WHERE tor_id=$1`
	query13 = `UPDATE contract_tor SET status=$2,updated_at=now() WHERE tor_id=$1`
	query14 = `SELECT EXISTS(SELECT 1 FROM location WHERE location_name=$1 AND ($2::uuid IS NULL OR location_id<>$2))`
	query15 = `INSERT INTO location(location_name,address,latitude,longitude) VALUES($1,$2,$3,$4) RETURNING location_id::text`
	query16 = `UPDATE location SET location_name=$2,address=$3,latitude=$4,longitude=$5 WHERE location_id=$1`
	query17 = `SELECT EXISTS(SELECT 1 FROM equipment WHERE equipment_name=$1 AND ($2::uuid IS NULL OR equipment_id<>$2))`
	query18 = `INSERT INTO equipment(equipment_name,is_active) VALUES($1,COALESCE($2,true)) RETURNING equipment_id::text`
	query19 = `UPDATE equipment SET equipment_name=$2,is_active=COALESCE($3,is_active) WHERE equipment_id=$1`
	query20 = `SELECT equipment_id,equipment_name,is_active FROM equipment ORDER BY equipment_name,equipment_id LIMIT $1 OFFSET $2`
)
