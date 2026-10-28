package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func registerReportingRoutes(worker, sup, asst *gin.RouterGroup, s *Service) {
	asst.GET("/assignments/:id/operations-summary", func(c *gin.Context) {
		out, e := s.OperationsSummary(c.Request.Context(), actor(c), c.Param("id"), c.Query("period_start"), c.Query("period_end"))
		respond(c, out, e)
	})
	sup.POST("/contracts/:id/operations-summary/notify", input(func(c *gin.Context, p auth.Principal, in OperationsNotificationInput) (any, error) {
		return s.NotifyOperationsSummary(c.Request.Context(), p, c.Param("id"), in)
	}))
	asst.GET("/assignments/:id/continuation", func(c *gin.Context) {
		out, e := s.AssignmentContinuation(c.Request.Context(), c.Param("id"))
		respond(c, out, e)
	})
	asst.GET("/contracts/:id/continuation", func(c *gin.Context) { out, e := s.Continuation(c.Request.Context(), c.Param("id")); respond(c, out, e) })
	sup.POST("/payroll/batch", input(func(c *gin.Context, p auth.Principal, in PayrollBatchInput) (any, error) {
		return s.PayrollBatch(c.Request.Context(), p, in)
	}))
	sup.POST("/reports/profit/confirm", input(func(c *gin.Context, _ auth.Principal, in ProfitConfirmation) (any, error) {
		return s.CloseSummary(c.Request.Context(), in.TorID, in.Start, in.End)
	}))
	sup.GET("/reports/profit.pdf", s.ExportProfitPDF)
	sup.POST("/payroll/preview", input(func(c *gin.Context, _ auth.Principal, in PayrollInput) (any, error) {
		return s.PreviewPayroll(c.Request.Context(), in)
	}))
	sup.POST("/payroll", input(func(c *gin.Context, p auth.Principal, in PayrollInput) (any, error) {
		id, e := s.Payroll(c.Request.Context(), p, in)
		return gin.H{"payroll_id": id}, e
	}))
	sup.GET("/payroll", payrollList(s))
	worker.GET("/payroll", payrollList(s))
	sup.POST("/payroll/:id/mark-paid", func(c *gin.Context) { respond(c, nil, s.Pay(c.Request.Context(), c.Param("id"))) })
	sup.POST("/invoices", input(func(c *gin.Context, _ auth.Principal, in InvoiceInput) (any, error) {
		id, e := s.Invoice(c.Request.Context(), in)
		return gin.H{"invoice_id": id}, e
	}))
	sup.GET("/invoices", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
		return s.Invoices(ctx, l, o)
	}))
	sup.POST("/invoices/:id/mark-paid", input(func(c *gin.Context, _ auth.Principal, in paidInput) (any, error) {
		return nil, s.Receive(c.Request.Context(), c.Param("id"), in.Amount)
	}))
	sup.GET("/reports/profit", func(c *gin.Context) {
		out, e := s.profitRequest(c)
		respond(c, gin.H{"data": out}, e)
	})
	sup.GET("/reports/profit.csv", s.ExportProfit)
}
