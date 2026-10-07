package contract

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/db"
	"chrisfoong/chaum-work-management-backend/internal/notify"
)

const (
	locationSearchLimit     = 20
	requisitionNoMaxRetries = 3
	notifyTimeout           = 10 * time.Second
)

// Store is the data access the service needs; TORRepository implements it.
type Store interface {
	LockRequisitionNumbering(context.Context, db.DBTX, int32) error
	CheckDuplicateContractNo(ctx context.Context, q db.DBTX, contractNo string) (bool, error)
	ValidateLocation(ctx context.Context, q db.DBTX, locationID uuid.UUID) (bool, error)
	CheckDuplicateLocationName(ctx context.Context, q db.DBTX, name string) (bool, error)
	CreateContract(ctx context.Context, q db.DBTX, userID uuid.UUID, c ContractInfo) (uuid.UUID, error)
	CreateLocation(ctx context.Context, q db.DBTX, name, address string) (uuid.UUID, error)
	CreateTORLocationAssignment(ctx context.Context, q db.DBTX, torID, locationID uuid.UUID, requiredWorkers int) (uuid.UUID, string, error)
	NextRequisitionSeq(ctx context.Context, q db.DBTX, prefix string) (int, error)
	CreateInitialRequisition(ctx context.Context, q db.DBTX, assignmentID, userID uuid.UUID, requisitionNo string) (uuid.UUID, error)
	CreateEquipment(ctx context.Context, q db.DBTX, name string) error
	FindEquipmentByName(ctx context.Context, q db.DBTX, name string) (uuid.UUID, error)
	CreateRequisitionItem(ctx context.Context, q db.DBTX, requisitionID, equipmentID uuid.UUID, remark string, requiredQty int) (uuid.UUID, error)
	SearchLocations(ctx context.Context, q db.DBTX, query string, limit int) ([]Location, error)
	FindActiveAssistantLineIDs(ctx context.Context, q db.DBTX) ([]string, error)
}

// Service holds the 1S controller operations (ContractFormController and
// ConfirmContractController in the class diagram).
type Service struct {
	store    Store
	conn     db.DBTX
	withTx   db.TxRunner
	notifier notify.Notifier
	now      func() time.Time
	runAsync func(func())
}

// NewService wires the 1S service. conn is used outside transactions.
func NewService(store Store, conn db.DBTX, withTx db.TxRunner, notifier notify.Notifier) *Service {
	return &Service{
		store:    store,
		conn:     conn,
		withTx:   withTx,
		notifier: notifier,
		now:      time.Now,
		runAsync: func(f func()) { go f() },
	}
}

// SubmitContractInfo is uc 1S steps 4–5: validateContractFormat, then
// checkDuplicateContractNo. Nothing is stored.
func (s *Service) SubmitContractInfo(ctx context.Context, c ContractInfo) error {
	c = c.normalized()
	if errs := ValidateContractFormat(c); len(errs) > 0 {
		return apperr.Validation(errs...)
	}
	dup, err := s.store.CheckDuplicateContractNo(ctx, s.conn, c.ContractNo)
	if err != nil {
		return apperr.Internal(err)
	}
	if dup {
		return duplicateContractNo()
	}
	return nil
}

// SubmitScopeData is uc 1S steps 7–8: validateScopeFormat, then validateLocation
// for existing locations and checkDuplicateLocationName for new ones. Nothing is stored.
func (s *Service) SubmitScopeData(ctx context.Context, sc Scope) error {
	sc = sc.normalized()
	if errs := ValidateScopeFormat(sc); len(errs) > 0 {
		return apperr.Validation(errs...)
	}
	var missing []apperr.FieldError
	for i, a := range sc.Areas {
		if a.NewLocation != nil {
			continue
		}
		ok, err := s.store.ValidateLocation(ctx, s.conn, uuid.MustParse(a.LocationID))
		if err != nil {
			return apperr.Internal(err)
		}
		if !ok {
			missing = append(missing, locationNotFound(i))
		}
	}
	if len(missing) > 0 {
		return apperr.Validation(missing...)
	}
	for i, a := range sc.Areas {
		if a.NewLocation == nil {
			continue
		}
		dup, err := s.store.CheckDuplicateLocationName(ctx, s.conn, a.NewLocation.Name)
		if err != nil {
			return apperr.Internal(err)
		}
		if dup {
			return duplicateLocationName(i)
		}
	}
	return nil
}

// ConfirmContract is uc 1S step 10 (confirmContract): all inserts in one
// transaction (technical atomicity; the use case does not state a boundary).
// It repeats no checks as queries: duplicates and missing locations are caught
// by the database constraints. A taken requisition number retries the whole
// transaction. The notification is sent afterwards in the background.
func (s *Service) ConfirmContract(ctx context.Context, userID uuid.UUID, req ConfirmRequest) (ConfirmedContract, error) {
	c := req.Contract.normalized()
	sc := req.Scope.normalized()
	errs := append(ValidateContractFormat(c), ValidateScopeFormat(sc)...)
	if len(errs) > 0 {
		return ConfirmedContract{}, withConfirmPrefixes(apperr.Validation(errs...))
	}

	prefix := "REQ-" + s.now().In(bangkok).Format("20060102")
	var (
		result ConfirmedContract
		err    error
	)
	for attempt := 1; attempt <= requisitionNoMaxRetries; attempt++ {
		err = s.withTx(ctx, func(q db.DBTX) error {
			var txErr error
			result, txErr = s.confirmInTx(ctx, q, userID, c, sc, prefix)
			return txErr
		})
		if ctx.Err() != nil {
			err = errRequisitionNoTaken
			break
		}
		if !errors.Is(err, errRequisitionNoTaken) {
			break
		}
		slog.WarnContext(ctx, "requisition number taken, retrying", "attempt", attempt)
	}
	if err != nil {
		if errors.Is(err, errRequisitionNoTaken) {
			return ConfirmedContract{}, apperr.Conflict("requisition_no_busy", "could not allocate a requisition number, please retry")
		}
		var ae *apperr.Error
		if errors.As(err, &ae) {
			return ConfirmedContract{}, withConfirmPrefixes(ae)
		}
		return ConfirmedContract{}, apperr.Internal(err)
	}

	s.SendNewContractNotification(result.TorID, result.ProjectName)
	return result, nil
}

func (s *Service) confirmInTx(ctx context.Context, q db.DBTX, userID uuid.UUID, c ContractInfo, sc Scope, prefix string) (ConfirmedContract, error) {
	day, _ := strconv.Atoi(strings.TrimPrefix(prefix, "REQ-"))
	if err := s.store.LockRequisitionNumbering(ctx, q, int32(day)); err != nil {
		return ConfirmedContract{}, err
	}
	torID, err := s.store.CreateContract(ctx, q, userID, c)
	if name, ok := db.UniqueViolation(err); ok && name == constraintContractNo {
		return ConfirmedContract{}, duplicateContractNo()
	}
	if err != nil {
		return ConfirmedContract{}, err
	}
	seq, err := s.store.NextRequisitionSeq(ctx, q, prefix)
	if err != nil {
		return ConfirmedContract{}, err
	}

	out := ConfirmedContract{
		TorID: torID, ContractNo: c.ContractNo, ProjectName: c.ProjectName, PartnerAgency: c.PartnerAgency,
		StartDate: c.StartDate, EndDate: c.EndDate, ContractValue: c.ContractValue, Status: "registered",
		Areas: make([]ConfirmedArea, 0, len(sc.Areas)),
	}
	for i, a := range sc.Areas {
		area, err := s.createArea(ctx, q, userID, torID, i, a, fmt.Sprintf("%s-%03d", prefix, seq+i))
		if err != nil {
			return ConfirmedContract{}, err
		}
		out.Areas = append(out.Areas, area)
	}
	return out, nil
}

// createArea: [not found] createLocation → createTORLocationAssignment →
// createInitialRequisition → per item: equipment find-or-create → createRequisitionItem.
func (s *Service) createArea(ctx context.Context, q db.DBTX, userID, torID uuid.UUID, i int, a Area, requisitionNo string) (ConfirmedArea, error) {
	var locationID uuid.UUID
	if a.NewLocation != nil {
		id, err := s.store.CreateLocation(ctx, q, a.NewLocation.Name, a.NewLocation.Address)
		if errors.Is(err, errDuplicateLocation) {
			return ConfirmedArea{}, duplicateLocationName(i)
		}
		if name, ok := db.UniqueViolation(err); ok && name == constraintLocationName {
			return ConfirmedArea{}, duplicateLocationName(i)
		}
		if err != nil {
			return ConfirmedArea{}, err
		}
		locationID = id
	} else {
		locationID = uuid.MustParse(a.LocationID)
	}

	assignmentID, locationName, err := s.store.CreateTORLocationAssignment(ctx, q, torID, locationID, a.RequiredWorkers)
	if name, ok := db.ForeignKeyViolation(err); ok && name == constraintAssignmentLoc {
		return ConfirmedArea{}, apperr.Validation(locationNotFound(i))
	}
	if err != nil {
		return ConfirmedArea{}, err
	}

	requisitionID, err := s.store.CreateInitialRequisition(ctx, q, assignmentID, userID, requisitionNo)
	if name, ok := db.UniqueViolation(err); ok && name == constraintRequisitionNo {
		return ConfirmedArea{}, errRequisitionNoTaken
	}
	if err != nil {
		return ConfirmedArea{}, err
	}

	items := make([]ConfirmedItem, 0, len(a.Items))
	for _, it := range a.Items {
		if err := s.store.CreateEquipment(ctx, q, it.EquipmentName); err != nil {
			return ConfirmedArea{}, err
		}
		equipmentID, err := s.store.FindEquipmentByName(ctx, q, it.EquipmentName)
		if err != nil {
			return ConfirmedArea{}, err
		}
		itemID, err := s.store.CreateRequisitionItem(ctx, q, requisitionID, equipmentID, it.Remark, it.RequiredQty)
		if err != nil {
			return ConfirmedArea{}, err
		}
		items = append(items, ConfirmedItem{
			ItemID: itemID, EquipmentID: equipmentID, EquipmentName: it.EquipmentName,
			RequiredQty: it.RequiredQty, Remark: it.Remark,
		})
	}

	return ConfirmedArea{
		AssignmentID: assignmentID, LocationID: locationID, LocationName: locationName,
		RequiredWorkers: a.RequiredWorkers,
		Requisition: ConfirmedRequisition{
			RequisitionID: requisitionID, RequisitionNo: requisitionNo,
			RequisitionType: "tor_base", Status: "pending_survey", Items: items,
		},
	}, nil
}

// SendNewContractNotification is sendNewContractNotification(torId, projectName).
// It runs in the background with its own context and timeout, so a slow or failed
// send never delays or undoes the saved contract. The outcome is logged.
func (s *Service) SendNewContractNotification(torID uuid.UUID, projectName string) {
	s.runAsync(func() {
		ctx, cancel := context.WithTimeout(context.Background(), notifyTimeout)
		defer cancel()

		// TODO(decision-12): recipients are every active assistant until area permission is decided.
		lineIDs, err := s.store.FindActiveAssistantLineIDs(ctx, s.conn)
		if err != nil {
			slog.Error("new contract notification failed", "tor_id", torID, "error", err)
			return
		}
		if len(lineIDs) == 0 {
			slog.Warn("new contract notification skipped: no active assistant", "tor_id", torID)
			return
		}
		text := fmt.Sprintf("มีโครงการใหม่ถูกบันทึกเข้าระบบ (%s) โปรดลงพื้นที่เพื่อตรวจสอบรายการอุปกรณ์", projectName)
		counts := map[notify.Status]int{}
		for _, id := range lineIDs {
			status, err := s.notifier.Send(ctx, notify.Message{LineUserID: id, Text: text})
			if err != nil {
				status = notify.StatusFailed
				slog.Error("new contract notification send failed", "tor_id", torID, "error", err)
			}
			counts[status]++
		}
		slog.Info("new contract notification",
			"tor_id", torID,
			"sent", counts[notify.StatusSent],
			"skipped", counts[notify.StatusSkipped],
			"failed", counts[notify.StatusFailed])
	})
}

// SearchLocations backs the 1S location picker (user addition; not in the use case).
func (s *Service) SearchLocations(ctx context.Context, query string) ([]Location, error) {
	if utf8.RuneCountInString(query) > maxTextLen {
		return nil, apperr.Validation(apperr.FieldError{Field: "q", Message: fmt.Sprintf("must be at most %d characters", maxTextLen)})
	}
	locations, err := s.store.SearchLocations(ctx, s.conn, query, locationSearchLimit)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if locations == nil {
		locations = []Location{}
	}
	return locations, nil
}

// withConfirmPrefixes names fields as they appear in the confirm body
// ({"contract": {...}, "scope": {"areas": [...]}}).
func withConfirmPrefixes(e *apperr.Error) *apperr.Error {
	for i, f := range e.Fields {
		if strings.HasPrefix(f.Field, "areas") {
			e.Fields[i].Field = "scope." + f.Field
		} else {
			e.Fields[i].Field = "contract." + f.Field
		}
	}
	return e
}

func duplicateContractNo() *apperr.Error {
	return apperr.Conflict("duplicate_contract_no", "contract number already exists").
		WithFields(apperr.FieldError{Field: "contract_no", Message: "already exists"})
}

func duplicateLocationName(i int) *apperr.Error {
	return apperr.Conflict("duplicate_location_name", "location name already exists").
		WithFields(apperr.FieldError{Field: fmt.Sprintf("areas[%d].new_location.name", i), Message: "already exists"})
}

func locationNotFound(i int) apperr.FieldError {
	return apperr.FieldError{Field: fmt.Sprintf("areas[%d].location_id", i), Message: "location does not exist"}
}
