package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestProfitPDFPaginationAndFinancialFields(t *testing.T) {
	records := make([]map[string]string, 25)
	for i := range records {
		records[i] = map[string]string{"tor_id": "00000000-0000-0000-0000-000000000001", "period_start": "2026-10-01", "period_end": "2026-10-31", "revenue": "900.00", "labor": "400.00", "material": "50.00", "profit": "450.00"}
	}
	data, _ := json.Marshal(records)
	pdf, e := profitPDF(data)
	if e != nil {
		t.Fatal(e)
	}
	text := string(pdf)
	if !strings.HasPrefix(text, "%PDF-1.4") || !strings.HasSuffix(text, "%%EOF\n") || strings.Count(text, "/Type /Page /Parent") != 3 || strings.Count(text, "Profit THB: 450.00") != 25 {
		t.Fatal("PDF lost pages or amounts")
	}
	if path := os.Getenv("PROFIT_PDF_QA_PATH"); path != "" {
		if e = os.WriteFile(path, pdf, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestNewUsecaseInputsRejectInvalidPayloadBeforeDatabase(t *testing.T) {
	s := &Service{}
	if _, e := s.ProfitRange(t.Context(), "", "2026-11-01", "2026-10-01"); e == nil {
		t.Fatal("reversed range accepted")
	}
	if _, e := s.PayrollBatch(t.Context(), auth.Principal{}, PayrollBatchInput{"2026-10-01", "2026-10-31"}); e == nil {
		t.Fatal("non-half-month batch accepted")
	}
}
