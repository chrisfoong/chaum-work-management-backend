package work

// Exported-schema SQL for procurement.
const (
	query65 = `SELECT EXISTS(SELECT 1 FROM work_schedule s JOIN attendance a ON a.schedule_id=s.schedule_id AND a.worker_id=s.worker_id WHERE s.worker_id=$1 AND s.assignment_id=$2 AND s.shift_status='scheduled' AND a.check_in IS NOT NULL AND a.check_out IS NULL AND a.check_in<=$3 AND $3 < ((s.work_date+s.shift_start_time) AT TIME ZONE 'Asia/Bangkok')+interval '8 hours')`
	query66 = `INSERT INTO equipment_requisition(requisition_no,requested_by,assignment_id,reason,requisition_type) VALUES($1,$2,$3,$4,'additional') RETURNING requisition_id::text`
	query67 = `INSERT INTO requisition_item(requisition_id,equipment_id,required_qty,remark) SELECT $1,equipment_id,$3,$4 FROM equipment WHERE equipment_id=$2 AND is_active`
	query68 = `SELECT r.* FROM equipment_requisition r WHERE ($1::uuid IS NULL OR r.requested_by=$1) ORDER BY r.created_at DESC,r.requisition_id LIMIT $2 OFFSET $3`
	query69 = `SELECT r.*,COALESCE((SELECT jsonb_agg(jsonb_build_object('item_id',i.item_id,'equipment_id',i.equipment_id,'equipment_name',e.equipment_name,'required_qty',i.required_qty,'existing_qty',i.existing_qty,'to_buy_qty',i.to_buy_qty,'actual_qty',i.actual_qty,'actual_price',i.actual_price::text,'remark',i.remark)) FROM requisition_item i JOIN equipment e ON e.equipment_id=i.equipment_id WHERE i.requisition_id=r.requisition_id),'[]'::jsonb) AS items FROM equipment_requisition r WHERE r.requisition_id=$1 AND ($2::uuid IS NULL OR requested_by=$2)`
	query70 = `SELECT status::text,requisition_type::text FROM equipment_requisition WHERE requisition_id=$1 FOR UPDATE`
	query71 = `SELECT count(*) FROM requisition_item WHERE requisition_id=$1`
	query72 = `UPDATE requisition_item SET existing_qty=$3 WHERE requisition_id=$1 AND item_id=$2`
	query73 = `SELECT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1 AND to_buy_qty>0)`
	query74 = `UPDATE equipment_requisition SET status=$2 WHERE requisition_id=$1`
	query75 = `SELECT status::text FROM equipment_requisition WHERE requisition_id=$1 AND requisition_type='additional' FOR UPDATE`
	query76 = `UPDATE equipment_requisition SET status=$2::text::requisition_status_enum,reviewed_by=$3,reviewed_at=$4,reason=CASE WHEN $2::text='rejected' THEN reason || E'\n[Rejection] ' || $5 ELSE reason END WHERE requisition_id=$1`
	query77 = `SELECT expense_id::text,requisition_id::text,total_amount::text,receipt_photo_url,user_id::text FROM expense_claim WHERE transfer_ref_no=$1 AND expense_type='fund_transfer'`
	query78 = `SELECT status::text FROM equipment_requisition WHERE requisition_id=$1 FOR UPDATE`
	query79 = `INSERT INTO expense_claim(expense_no,expense_type,user_id,requisition_id,total_amount,transfer_ref_no,receipt_photo_url) VALUES($1,'fund_transfer',$2,$3,$4::numeric,$5,$6) RETURNING expense_id::text`
	query80 = `UPDATE equipment_requisition SET status='pending_procurement' WHERE requisition_id=$1`
	query81 = `SELECT status::text,requisition_type::text,reviewed_by IS NOT NULL AND reviewed_at IS NOT NULL FROM equipment_requisition WHERE requisition_id=$1 FOR UPDATE`
	query82 = `SELECT to_buy_qty,COALESCE(actual_qty,0),COALESCE(actual_price,0)::text FROM requisition_item WHERE item_id=$1 AND requisition_id=$2 FOR UPDATE`
	query83 = `UPDATE requisition_item SET actual_qty=COALESCE(actual_qty,0)+$2,actual_price=CASE WHEN $4::boolean THEN round((COALESCE(actual_qty,0)*COALESCE(actual_price,0)+$2*$3::numeric)/(COALESCE(actual_qty,0)+$2),2) ELSE $3::numeric END WHERE item_id=$1`
	query84 = `INSERT INTO expense_claim(expense_no,expense_type,user_id,requisition_id,total_amount,receipt_photo_url) VALUES($1,'actual_expense',$2,$3,$4::numeric,$5) RETURNING expense_id::text`
	query85 = `SELECT EXISTS(SELECT 1 FROM requisition_item WHERE requisition_id=$1 AND to_buy_qty>COALESCE(actual_qty,0))`
	query86 = `UPDATE equipment_requisition SET status=$2 WHERE requisition_id=$1`
)
