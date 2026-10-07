package contract

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
)

// Formats from uc 1S steps 4 and 7.
var (
	contractNoPattern    = regexp.MustCompile(`^[0-9]{10}$`)
	datePattern          = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	contractValuePattern = regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`)
)

const (
	maxTextLen = 255 // varchar(255) in the data dictionary
	// maxContractValueDigits is the integer part allowed by numeric(14, 2).
	maxContractValueDigits = 12
)

type fieldErrors []apperr.FieldError

func (f *fieldErrors) add(field, msg string) {
	*f = append(*f, apperr.FieldError{Field: field, Message: msg})
}

// ValidateContractFormat checks contractData (uc 1S step 4). Input must be normalized.
func ValidateContractFormat(c ContractInfo) []apperr.FieldError {
	var errs fieldErrors
	if !contractNoPattern.MatchString(c.ContractNo) {
		errs.add("contract_no", "must be exactly 10 digits")
	}
	errs.checkText("project_name", c.ProjectName)
	errs.checkText("partner_agency", c.PartnerAgency)
	start, startOK := errs.checkDate("start_date", c.StartDate)
	end, endOK := errs.checkDate("end_date", c.EndDate)
	if startOK && endOK && end.Before(start) {
		errs.add("end_date", "must be on or after start_date")
	}
	switch {
	case !contractValuePattern.MatchString(c.ContractValue):
		errs.add("contract_value", "must be a number with at most 2 decimal places")
	case len(strings.SplitN(c.ContractValue, ".", 2)[0]) > maxContractValueDigits:
		errs.add("contract_value", "is too large")
	}
	return errs
}

// ValidateScopeFormat checks scopeData (uc 1S step 7) plus the user's duplicate
// rule: the same location twice in one contract, or the same equipment name
// twice in one area, is rejected. Input must be normalized.
func ValidateScopeFormat(s Scope) []apperr.FieldError {
	var errs fieldErrors
	if len(s.Areas) == 0 {
		errs.add("areas", "at least one area is required")
		return errs
	}
	seenLocation := map[uuid.UUID]int{}
	seenNewName := map[string]int{}
	for i, a := range s.Areas {
		p := fmt.Sprintf("areas[%d]", i)
		switch {
		case a.LocationID != "" && a.NewLocation != nil:
			errs.add(p, "choose an existing location or enter a new one, not both")
		case a.LocationID == "" && a.NewLocation == nil:
			errs.add(p+".location_id", "is required unless new_location is given")
		case a.LocationID != "":
			id, err := uuid.Parse(a.LocationID)
			if err != nil {
				errs.add(p+".location_id", "must be a UUID")
				break
			}
			if j, dup := seenLocation[id]; dup {
				errs.add(p+".location_id", fmt.Sprintf("is already used by areas[%d]", j))
				break
			}
			seenLocation[id] = i
		default:
			errs.checkText(p+".new_location.name", a.NewLocation.Name)
			errs.checkText(p+".new_location.address", a.NewLocation.Address)
			if a.NewLocation.Name == "" {
				break
			}
			if j, dup := seenNewName[a.NewLocation.Name]; dup {
				errs.add(p+".new_location.name", fmt.Sprintf("is already used by areas[%d]", j))
				break
			}
			seenNewName[a.NewLocation.Name] = i
		}
		if a.RequiredWorkers <= 0 {
			errs.add(p+".required_workers", "must be a whole number greater than 0")
		}
		if len(a.Items) == 0 {
			errs.add(p+".items", "at least one item is required")
			continue
		}
		seenEquipment := map[string]int{}
		for j, it := range a.Items {
			q := fmt.Sprintf("%s.items[%d]", p, j)
			errs.checkText(q+".equipment_name", it.EquipmentName)
			if it.EquipmentName != "" {
				if k, dup := seenEquipment[it.EquipmentName]; dup {
					errs.add(q+".equipment_name", fmt.Sprintf("is already listed as items[%d] in this area", k))
				} else {
					seenEquipment[it.EquipmentName] = j
				}
			}
			if it.RequiredQty <= 0 {
				errs.add(q+".required_qty", "must be a whole number greater than 0")
			}
		}
	}
	return errs
}

// checkText: required, not only spaces (input is trimmed), at most 255 characters.
func (f *fieldErrors) checkText(field, s string) {
	switch {
	case s == "":
		f.add(field, "is required")
	case utf8.RuneCountInString(s) > maxTextLen:
		f.add(field, fmt.Sprintf("must be at most %d characters", maxTextLen))
	}
}

func (f *fieldErrors) checkDate(field, s string) (time.Time, bool) {
	if s == "" {
		f.add(field, "is required")
		return time.Time{}, false
	}
	d, err := time.Parse("2006-01-02", s)
	if !datePattern.MatchString(s) || err != nil {
		f.add(field, "must be a date in YYYY-MM-DD format")
		return time.Time{}, false
	}
	return d, true
}
