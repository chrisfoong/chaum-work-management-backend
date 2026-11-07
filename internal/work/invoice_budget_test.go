package work

import (
	"context"
	"sync"
	"testing"
)

func TestIntegrationInvoiceContractBudget(t *testing.T) {
	f := fixtureUsecases(t)
	ctx := context.Background()
	s := f.service
	if _, e := s.Invoice(ctx, InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "1000.01"}); e == nil {
		t.Fatal("over-contract invoice accepted")
	}
	if _, e := s.Invoice(ctx, InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "0"}); e == nil {
		t.Fatal("zero invoice accepted")
	}
	id, e := s.Invoice(ctx, InvoiceInput{TorID: f.tor, Month: "2026-10", Expected: "600"})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.UpdateInvoice(ctx, id, InvoiceEditInput{Expected: "1001", Deduction: "0", PreviousExpected: "600"}); e == nil {
		t.Fatal("over-budget edit accepted")
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, m := range []string{"2026-11", "2026-12"} {
		wg.Add(1)
		go func(month string) {
			defer wg.Done()
			_, err := s.Invoice(ctx, InvoiceInput{TorID: f.tor, Month: month, Expected: "300"})
			results <- err
		}(m)
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("concurrent invoices: successes=%d want 1", success)
	}
	if e = s.UpdateInvoice(ctx, id, InvoiceEditInput{Expected: "701", Deduction: "0", PreviousExpected: "600"}); e == nil {
		t.Fatal("edit exceeded cumulative budget")
	}
	if e = s.UpdateInvoice(ctx, id, InvoiceEditInput{Expected: "700", Deduction: "0", PreviousExpected: "600"}); e != nil {
		t.Fatal(e)
	}
}
