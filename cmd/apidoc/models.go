package main

import (
	"chrisfoong/chaum-work-management-backend/internal/contract"
	"chrisfoong/chaum-work-management-backend/internal/work"
	"reflect"
	"strings"
)

func model(t reflect.Type) map[string]any {
	if t.Kind() == reflect.Pointer {
		s := model(t.Elem())
		s["nullable"] = true
		return s
	}
	switch t.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int64:
		return map[string]any{"type": "integer"}
	case reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Slice:
		return map[string]any{"type": "array", "items": model(t.Elem())}
	case reflect.Struct:
		props := map[string]any{}
		required := []string{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			parts := strings.Split(tag, ",")
			name := parts[0]
			if name == "-" || name == "" {
				continue
			}
			props[name] = model(f.Type)
			if !strings.Contains(tag, "omitempty") && (f.Type.Kind() != reflect.Pointer || t.Name() == "CheckIn") && name != "daily_wage" && name != "remark" && !(name == "reason" && (t.Name() == "ReviewInput" || t.Name() == "InvoiceInput")) {
				required = append(required, name)
			}
		}
		out := map[string]any{"type": "object", "properties": props, "additionalProperties": false}
		if len(required) > 0 {
			out["required"] = required
		}
		return out
	}
	return map[string]any{"type": "object"}
}
func requestModels() map[string]any {
	state := struct {
		Status string `json:"status"`
	}{}
	review := work.ReviewInput{}
	replacement := struct {
		Worker string `json:"worker_id"`
	}{}
	survey := struct {
		Items []work.SurveyItem `json:"items"`
	}{}
	paid := struct {
		Amount string `json:"net_received"`
	}{}
	return map[string]any{
		"POST /api/web/contracts/info": contract.ContractInfo{}, "POST /api/web/contracts/scope": contract.Scope{}, "POST /api/web/contracts/confirm": contract.ConfirmRequest{},
		"POST /api/web/users": work.UserInput{}, "PATCH /api/web/users/{id}": work.UserUpdate{}, "PATCH /api/web/contracts/{id}": state,
		"POST /api/web/locations": work.LocationInput{}, "PATCH /api/web/locations/{id}": work.LocationInput{}, "POST /api/web/equipment": work.EquipmentInput{}, "PATCH /api/web/equipment/{id}": work.EquipmentInput{},
		"POST /api/web/schedules": work.ScheduleInput{}, "POST /api/liff/leave-requests": work.LeaveInput{}, "POST /api/web/leave-requests/{id}/review": state, "POST /api/web/leave-requests/{id}/replacement": replacement,
		"POST /api/liff/attendance/check-in": work.CheckIn{}, "POST /api/liff/attendance/check-out": work.CheckOut{}, "POST /api/liff/requisitions": work.RequestInput{}, "POST /api/web/requisitions/{id}/survey": survey, "POST /api/web/requisitions/{id}/review": review, "POST /api/web/requisitions/{id}/fund-transfers": work.FundInput{}, "POST /api/web/requisitions/{id}/purchase": work.PurchaseInput{},
		"POST /api/web/payroll/preview": work.PayrollInput{}, "POST /api/web/payroll": work.PayrollInput{}, "POST /api/web/invoices": work.InvoiceInput{}, "POST /api/web/invoices/{id}/mark-paid": paid,
	}
}
func describeRequest(method, path string, op map[string]any) {
	if v, ok := requestModels()[method+" "+path]; ok {
		op["requestBody"] = map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": model(reflect.TypeOf(v))}}}
	}
	if method == "POST" && strings.HasSuffix(path, "/files") {
		op["requestBody"] = map[string]any{"required": true, "content": map[string]any{"application/octet-stream": map[string]any{"schema": map[string]string{"type": "string", "format": "binary"}}}}
	}
	if method == "GET" {
		if strings.HasSuffix(path, "/reports/profit") || strings.HasSuffix(path, "/reports/profit.csv") {
			op["parameters"] = []map[string]any{{"name": "month", "in": "query", "required": true, "schema": map[string]string{"type": "string", "pattern": `^\d{4}-\d{2}$`}}}
		}
		if strings.HasSuffix(path, "/files") {
			op["parameters"] = []map[string]any{{"name": "path", "in": "query", "required": true, "schema": map[string]string{"type": "string"}}}
		}
	}
}
