// Package contract implements 1S (create TOR contract and work scope), the 1S
// location search and 9A (work continuation). Operation names follow the team's
// class diagrams (ContractFormController, ConfirmContractController, TORRepository).
package contract

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// bangkok is Asia/Bangkok (UTC+7 all year, no daylight saving).
var bangkok = time.FixedZone("Asia/Bangkok", 7*60*60)

// ContractInfo is contractData: the step-1 contract form (uc 1S step 3).
type ContractInfo struct {
	ContractNo    string `json:"contract_no"`
	ProjectName   string `json:"project_name"`
	PartnerAgency string `json:"partner_agency"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	// ContractValue is a decimal string (max 2 decimal places) so money never passes through float.
	ContractValue string `json:"contract_value"`
	// TODO(decision-file): the contract file (uc 1S step 3) is not accepted until its type and upload are decided.
}

// Scope is scopeData: the step-2 work scope (uc 1S step 6).
type Scope struct {
	Areas []Area `json:"areas"`
}

// Area is one work area of the contract: an existing location or a new one.
// Each area gets its own TOR_LOCATION_ASSIGNMENT and initial requisition.
type Area struct {
	LocationID      string       `json:"location_id,omitempty"`
	NewLocation     *NewLocation `json:"new_location,omitempty"`
	RequiredWorkers int          `json:"required_workers"`
	Items           []Item       `json:"items"`
}

// NewLocation is a location to create (uc 1S "เพิ่มสถานที่ใหม่").
type NewLocation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Item is one equipment line of an area's initial requisition.
type Item struct {
	EquipmentName string `json:"equipment_name"`
	RequiredQty   int    `json:"required_qty"`
	Remark        string `json:"remark,omitempty"`
}

// ConfirmRequest is confirmContract(contractData, scopeData).
type ConfirmRequest struct {
	Contract ContractInfo `json:"contract"`
	Scope    Scope        `json:"scope"`
}

// ConfirmedContract is the 201 body of confirmContract; it carries the data for
// the contract detail page (redirectToContractDetailPage(torId)).
type ConfirmedContract struct {
	TorID         uuid.UUID `json:"tor_id"`
	ContractNo    string    `json:"contract_no"`
	ProjectName   string    `json:"project_name"`
	PartnerAgency string    `json:"partner_agency"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	ContractValue string    `json:"contract_value"`
	Status        string    `json:"status"`
	// ContractFileURL is always null until TODO(decision-file).
	ContractFileURL *string         `json:"contract_file_url"`
	Areas           []ConfirmedArea `json:"areas"`
}

// ConfirmedArea is one saved TOR_LOCATION_ASSIGNMENT with its initial requisition.
type ConfirmedArea struct {
	AssignmentID    uuid.UUID            `json:"assignment_id"`
	LocationID      uuid.UUID            `json:"location_id"`
	LocationName    string               `json:"location_name"`
	RequiredWorkers int                  `json:"required_workers"`
	Requisition     ConfirmedRequisition `json:"requisition"`
}

// ConfirmedRequisition is the saved tor_base requisition of an area.
type ConfirmedRequisition struct {
	RequisitionID   uuid.UUID       `json:"requisition_id"`
	RequisitionNo   string          `json:"requisition_no"`
	RequisitionType string          `json:"requisition_type"`
	Status          string          `json:"status"`
	Items           []ConfirmedItem `json:"items"`
}

// ConfirmedItem is one saved REQUISITION_ITEM.
type ConfirmedItem struct {
	ItemID        uuid.UUID `json:"item_id"`
	EquipmentID   uuid.UUID `json:"equipment_id"`
	EquipmentName string    `json:"equipment_name"`
	RequiredQty   int       `json:"required_qty"`
	Remark        string    `json:"remark,omitempty"`
}

// Location is one row of the location search.
type Location struct {
	LocationID   uuid.UUID `json:"location_id"`
	LocationName string    `json:"location_name"`
	Address      *string   `json:"address"`
}

func (c ContractInfo) normalized() ContractInfo {
	c.ContractNo = strings.TrimSpace(c.ContractNo)
	c.ProjectName = strings.TrimSpace(c.ProjectName)
	c.PartnerAgency = strings.TrimSpace(c.PartnerAgency)
	c.StartDate = strings.TrimSpace(c.StartDate)
	c.EndDate = strings.TrimSpace(c.EndDate)
	c.ContractValue = strings.TrimSpace(c.ContractValue)
	return c
}

func (s Scope) normalized() Scope {
	areas := make([]Area, len(s.Areas))
	for i, a := range s.Areas {
		a.LocationID = strings.TrimSpace(a.LocationID)
		if a.NewLocation != nil {
			a.NewLocation = &NewLocation{
				Name:    strings.TrimSpace(a.NewLocation.Name),
				Address: strings.TrimSpace(a.NewLocation.Address),
			}
		}
		items := make([]Item, len(a.Items))
		for j, it := range a.Items {
			it.EquipmentName = strings.TrimSpace(it.EquipmentName)
			it.Remark = strings.TrimSpace(it.Remark)
			items[j] = it
		}
		a.Items = items
		areas[i] = a
	}
	return Scope{Areas: areas}
}
