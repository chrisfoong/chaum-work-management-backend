package contract

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"chrisfoong/chaum-work-management-backend/internal/db"
)

// Constraint names the service maps to HTTP errors (migrations/0001_init.up.sql).
const (
	constraintContractNo     = "uq_contract_tor_contract_no"
	constraintLocationName   = "uq_location_location_name"
	constraintRequisitionNo  = "uq_equipment_requisition_requisition_no"
	constraintAssignmentLoc  = "tor_location_assignment_location_id_fkey"
	initialRequisitionReason = "จัดเตรียมอุปกรณ์เริ่มต้นสำหรับสัญญาใหม่"
)

// TORRepository holds the 1S data operations (class diagram TORRepository), plus
// the location search and the notification-recipient lookup. Every method takes
// the connection or transaction to run on.
type TORRepository struct{}

// CheckDuplicateContractNo is Q1S.1.
func (TORRepository) CheckDuplicateContractNo(ctx context.Context, q db.DBTX, contractNo string) (bool, error) {
	return exists(ctx, q, `SELECT COUNT(contract_no) FROM contract_tor WHERE contract_no = $1`, contractNo)
}

// ValidateLocation is Q1S.2: reports whether the location exists.
func (TORRepository) ValidateLocation(ctx context.Context, q db.DBTX, locationID uuid.UUID) (bool, error) {
	return exists(ctx, q, `SELECT COUNT(location_id) FROM location WHERE location_id = $1`, locationID)
}

// CheckDuplicateLocationName is Q1S.2.1.
func (TORRepository) CheckDuplicateLocationName(ctx context.Context, q db.DBTX, name string) (bool, error) {
	return exists(ctx, q, `SELECT COUNT(location_id) FROM location WHERE location_name = $1`, name)
}

// CreateContract is Q1S.3. TODO(decision-file): contract_file_url stays NULL.
func (TORRepository) CreateContract(ctx context.Context, q db.DBTX, userID uuid.UUID, c ContractInfo) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `
		INSERT INTO contract_tor (contract_no, user_id, project_name, partner_agency, contract_value,
		                          start_date, end_date, contract_file_url, status)
		VALUES ($1, $2, $3, $4, $5::numeric, $6::date, $7::date, NULL, 'registered')
		RETURNING tor_id`,
		c.ContractNo, userID, c.ProjectName, c.PartnerAgency, c.ContractValue, c.StartDate, c.EndDate,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create contract: %w", err)
	}
	return id, nil
}

// CreateLocation is Q1S.3.1.
func (TORRepository) CreateLocation(ctx context.Context, q db.DBTX, name, address string) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx,
		`INSERT INTO location (location_name, address) VALUES ($1, $2) RETURNING location_id`,
		name, address,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create location: %w", err)
	}
	return id, nil
}

// CreateTORLocationAssignment is Q1S.4. It also returns the location name for the
// 201 response, in the same statement.
func (TORRepository) CreateTORLocationAssignment(ctx context.Context, q db.DBTX, torID, locationID uuid.UUID, requiredWorkers int) (uuid.UUID, string, error) {
	var (
		id   uuid.UUID
		name string
	)
	err := q.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO tor_location_assignment (tor_id, location_id, required_workers)
			VALUES ($1, $2, $3)
			RETURNING assignment_id, location_id
		)
		SELECT ins.assignment_id, l.location_name
		FROM ins JOIN location l ON l.location_id = ins.location_id`,
		torID, locationID, requiredWorkers,
	).Scan(&id, &name)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("create assignment: %w", err)
	}
	return id, name, nil
}

// NextRequisitionSeq returns the next NNN for requisition numbers starting with
// prefix ("REQ-YYYYMMDD"). Two concurrent transactions can get the same number;
// the unique key then rejects one and the service retries.
func (TORRepository) NextRequisitionSeq(ctx context.Context, q db.DBTX, prefix string) (int, error) {
	var next int
	err := q.QueryRow(ctx, `
		SELECT COALESCE(MAX(split_part(requisition_no, '-', 3)::int), 0) + 1
		FROM equipment_requisition
		WHERE requisition_no LIKE $1 || '-%'`,
		prefix,
	).Scan(&next)
	if err != nil {
		return 0, fmt.Errorf("next requisition number: %w", err)
	}
	return next, nil
}

// CreateInitialRequisition is Q1S.5 (requested_by per the data dictionary).
func (TORRepository) CreateInitialRequisition(ctx context.Context, q db.DBTX, assignmentID, userID uuid.UUID, requisitionNo string) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `
		INSERT INTO equipment_requisition (assignment_id, requested_by, reason, status, requisition_type, requisition_no)
		VALUES ($1, $2, $3, 'pending_survey', 'tor_base', $4)
		RETURNING requisition_id`,
		assignmentID, userID, initialRequisitionReason, requisitionNo,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create initial requisition: %w", err)
	}
	return id, nil
}

// CreateEquipment is Q1S.6.1, made safe for concurrent requests: an existing name
// is left untouched (ON CONFLICT DO NOTHING) and FindEquipmentByName reads the id.
func (TORRepository) CreateEquipment(ctx context.Context, q db.DBTX, name string) error {
	_, err := q.Exec(ctx,
		`INSERT INTO equipment (equipment_name, is_active) VALUES ($1, true) ON CONFLICT (equipment_name) DO NOTHING`,
		name)
	if err != nil {
		return fmt.Errorf("create equipment: %w", err)
	}
	return nil
}

// FindEquipmentByName is Q1S.6.
func (TORRepository) FindEquipmentByName(ctx context.Context, q db.DBTX, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `SELECT equipment_id FROM equipment WHERE equipment_name = $1`, name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("find equipment by name: %w", err)
	}
	return id, nil
}

// CreateRequisitionItem is Q1S.7.
func (TORRepository) CreateRequisitionItem(ctx context.Context, q db.DBTX, requisitionID, equipmentID uuid.UUID, remark string, requiredQty int) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `
		INSERT INTO requisition_item (requisition_id, equipment_id, remark, required_qty)
		VALUES ($1, $2, NULLIF($3, ''), $4)
		RETURNING item_id`,
		requisitionID, equipmentID, remark, requiredQty,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create requisition item: %w", err)
	}
	return id, nil
}

// SearchLocations backs GET /api/web/locations?q= (an addition by user decision,
// not in the use case or class diagram): up to limit locations whose name contains query.
func (TORRepository) SearchLocations(ctx context.Context, q db.DBTX, query string, limit int) ([]Location, error) {
	rows, err := q.Query(ctx, `
		SELECT location_id, location_name, address
		FROM location
		WHERE location_name ILIKE '%' || $1 || '%' ESCAPE '\'
		ORDER BY location_name
		LIMIT $2`,
		escapeLike(query), limit)
	if err != nil {
		return nil, fmt.Errorf("search locations: %w", err)
	}
	locations, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Location, error) {
		var l Location
		err := r.Scan(&l.LocationID, &l.LocationName, &l.Address)
		return l, err
	})
	if err != nil {
		return nil, fmt.Errorf("search locations: %w", err)
	}
	return locations, nil
}

// FindActiveAssistantLineIDs returns the LINE userIds of all active assistants.
// TODO(decision-12): until assistants are linked to contracts, every active
// assistant receives the 1S notification.
func (TORRepository) FindActiveAssistantLineIDs(ctx context.Context, q db.DBTX) ([]string, error) {
	rows, err := q.Query(ctx, `SELECT line_id FROM users WHERE role = 'assistant' AND is_active ORDER BY user_id`)
	if err != nil {
		return nil, fmt.Errorf("find assistants: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("find assistants: %w", err)
	}
	return ids, nil
}

func exists(ctx context.Context, q db.DBTX, sql string, arg any) (bool, error) {
	var n int
	if err := q.QueryRow(ctx, sql, arg).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

// errRequisitionNoTaken signals that another transaction took the generated number.
var errRequisitionNoTaken = errors.New("requisition number already taken")
