package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
)

func init() {
	gin.SetMode(gin.TestMode)
	UseJSONFieldNames()
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantBody ErrorDetail
	}{
		{
			name:     "validation names fields",
			err:      apperr.Validation(apperr.FieldError{Field: "contract_no", Message: "is required"}),
			wantCode: http.StatusBadRequest,
			wantBody: ErrorDetail{Code: "validation_failed", Message: "request validation failed",
				Fields: []apperr.FieldError{{Field: "contract_no", Message: "is required"}}},
		},
		{
			name:     "conflict",
			err:      apperr.Conflict("duplicate_contract_no", "contract number already exists"),
			wantCode: http.StatusConflict,
			wantBody: ErrorDetail{Code: "duplicate_contract_no", Message: "contract number already exists"},
		},
		{
			name:     "plain error hides detail",
			err:      errors.New("pq: password=secret"),
			wantCode: http.StatusInternalServerError,
			wantBody: ErrorDetail{Code: "internal", Message: "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			WriteError(c, tt.err)

			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantCode)
			}
			var got ErrorBody
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.Error.Code != tt.wantBody.Code || got.Error.Message != tt.wantBody.Message ||
				len(got.Error.Fields) != len(tt.wantBody.Fields) {
				t.Fatalf("body = %+v, want %+v", got.Error, tt.wantBody)
			}
			if strings.Contains(w.Body.String(), "secret") {
				t.Fatal("internal detail leaked to client")
			}
		})
	}
}

func TestBindJSON(t *testing.T) {
	type input struct {
		ContractNo string `json:"contract_no" binding:"required"`
		Workers    int    `json:"required_workers" binding:"min=1"`
	}

	tests := []struct {
		name       string
		body       string
		wantFields []string
	}{
		{name: "valid", body: `{"contract_no":"C-1","required_workers":2}`},
		{name: "missing and too small", body: `{"required_workers":0}`, wantFields: []string{"contract_no", "required_workers"}},
		{name: "malformed json", body: `{`, wantFields: []string{"body"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			var in input
			err := BindJSON(c, &in)
			if len(tt.wantFields) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			var ae *apperr.Error
			if !errors.As(err, &ae) || ae.Kind != apperr.KindValidation {
				t.Fatalf("err = %v, want validation error", err)
			}
			if len(ae.Fields) != len(tt.wantFields) {
				t.Fatalf("fields = %+v, want %v", ae.Fields, tt.wantFields)
			}
			for i, f := range tt.wantFields {
				if ae.Fields[i].Field != f {
					t.Fatalf("field[%d] = %q, want %q", i, ae.Fields[i].Field, f)
				}
			}
		})
	}
}
