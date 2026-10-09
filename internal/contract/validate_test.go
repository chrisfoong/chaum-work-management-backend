package contract

import (
	"strings"
	"testing"
)

func validInfo() ContractInfo {
	return ContractInfo{
		ContractNo: "6700001148", ProjectName: "โครงการปรับปรุงภูมิทัศน์", PartnerAgency: "การทางพิเศษแห่งประเทศไทย",
		StartDate: "2026-01-15", EndDate: "2026-12-31", ContractValue: "1500000.00", ContractFilePath: "test-contract.png",
	}
}

func TestValidateContractFormat(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ContractInfo)
		want   []string // failing fields
	}{
		{name: "valid", mutate: func(*ContractInfo) {}},
		{name: "contract no 9 digits", mutate: func(c *ContractInfo) { c.ContractNo = "670000114" }, want: []string{"contract_no"}},
		{name: "contract no letters", mutate: func(c *ContractInfo) { c.ContractNo = "67000011AB" }, want: []string{"contract_no"}},
		{name: "missing project and agency", mutate: func(c *ContractInfo) { c.ProjectName, c.PartnerAgency = "", "" }, want: []string{"project_name", "partner_agency"}},
		{name: "project too long", mutate: func(c *ContractInfo) { c.ProjectName = strings.Repeat("ก", 256) }, want: []string{"project_name"}},
		{name: "project 255 thai chars ok", mutate: func(c *ContractInfo) { c.ProjectName = strings.Repeat("ก", 255) }},
		{name: "bad date format", mutate: func(c *ContractInfo) { c.StartDate = "15/01/2026" }, want: []string{"start_date"}},
		{name: "impossible date", mutate: func(c *ContractInfo) { c.EndDate = "2026-02-30" }, want: []string{"end_date"}},
		{name: "end before start", mutate: func(c *ContractInfo) { c.EndDate = "2026-01-14" }, want: []string{"end_date"}},
		{name: "end equals start", mutate: func(c *ContractInfo) { c.EndDate = c.StartDate }},
		{name: "value 3 decimals", mutate: func(c *ContractInfo) { c.ContractValue = "100.123" }, want: []string{"contract_value"}},
		{name: "value negative", mutate: func(c *ContractInfo) { c.ContractValue = "-1" }, want: []string{"contract_value"}},
		{name: "value integer ok", mutate: func(c *ContractInfo) { c.ContractValue = "850000" }},
		{name: "value too large", mutate: func(c *ContractInfo) { c.ContractValue = "1234567890123" }, want: []string{"contract_value"}},
		{name: "value missing", mutate: func(c *ContractInfo) { c.ContractValue = "" }, want: []string{"contract_value"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validInfo()
			tt.mutate(&c)
			assertFields(t, fieldNames(ValidateContractFormat(c.normalized())), tt.want)
		})
	}
}

func validScope() Scope {
	return Scope{Areas: []Area{
		{NewLocation: &NewLocation{Name: "สวนลุมพินี", Address: "ปทุมวัน กรุงเทพฯ"}, RequiredWorkers: 3,
			Items: []Item{{EquipmentName: "กรรไกรตัดกิ่ง", RequiredQty: 5}, {EquipmentName: "รถเข็น", RequiredQty: 2, Remark: "ล้อยาง"}}},
		{LocationID: "a1b2c3d4-e5f6-4a8b-9c0d-1e2f3a4b5c6d", RequiredWorkers: 1,
			Items: []Item{{EquipmentName: "กรรไกรตัดกิ่ง", RequiredQty: 1}}},
	}}
}

func TestValidateScopeFormat(t *testing.T) {
	const id = "a1b2c3d4-e5f6-4a8b-9c0d-1e2f3a4b5c6d"
	tests := []struct {
		name   string
		mutate func(*Scope)
		want   []string
	}{
		{name: "valid", mutate: func(*Scope) {}},
		{name: "no areas", mutate: func(s *Scope) { s.Areas = nil }, want: []string{"areas"}},
		{name: "neither location", mutate: func(s *Scope) { s.Areas[1].LocationID = "" }, want: []string{"areas[1].location_id"}},
		{name: "both locations", mutate: func(s *Scope) { s.Areas[1].NewLocation = &NewLocation{Name: "x", Address: "y"} }, want: []string{"areas[1]"}},
		{name: "location not uuid", mutate: func(s *Scope) { s.Areas[1].LocationID = "42" }, want: []string{"areas[1].location_id"}},
		{name: "blank new location", mutate: func(s *Scope) { s.Areas[0].NewLocation = &NewLocation{Name: "   ", Address: ""} },
			want: []string{"areas[0].new_location.name", "areas[0].new_location.address"}},
		{name: "same existing location twice", mutate: func(s *Scope) {
			s.Areas[0] = Area{LocationID: strings.ToUpper(id), RequiredWorkers: 1, Items: []Item{{EquipmentName: "a", RequiredQty: 1}}}
		}, want: []string{"areas[1].location_id"}},
		{name: "same new location name twice", mutate: func(s *Scope) {
			s.Areas[1] = Area{NewLocation: &NewLocation{Name: " สวนลุมพินี ", Address: "z"}, RequiredWorkers: 1, Items: []Item{{EquipmentName: "a", RequiredQty: 1}}}
		}, want: []string{"areas[1].new_location.name"}},
		{name: "workers zero", mutate: func(s *Scope) { s.Areas[0].RequiredWorkers = 0 }, want: []string{"areas[0].required_workers"}},
		{name: "no items", mutate: func(s *Scope) { s.Areas[1].Items = nil }, want: []string{"areas[1].items"}},
		{name: "item blank name and zero qty", mutate: func(s *Scope) { s.Areas[0].Items[1] = Item{EquipmentName: " ", RequiredQty: 0} },
			want: []string{"areas[0].items[1].equipment_name", "areas[0].items[1].required_qty"}},
		{name: "same equipment twice in one area", mutate: func(s *Scope) { s.Areas[0].Items[1].EquipmentName = "กรรไกรตัดกิ่ง " },
			want: []string{"areas[0].items[1].equipment_name"}},
		{name: "same equipment in different areas ok", mutate: func(*Scope) {}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validScope()
			tt.mutate(&s)
			assertFields(t, fieldNames(ValidateScopeFormat(s.normalized())), tt.want)
		})
	}
}
