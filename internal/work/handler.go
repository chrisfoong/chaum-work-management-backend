package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type stateInput struct {
	Status string `json:"status"`
}
type replacementInput struct {
	WorkerID string `json:"worker_id"`
}
type surveyInput struct {
	Items []SurveyItem `json:"items"`
}
type paidInput struct {
	Amount string `json:"net_received"`
}

func Register(web, worker *gin.RouterGroup, s *Service) {
	worker.Use(func(c *gin.Context) {
		if _, e := s.Worker(c.Request.Context(), s.Repo.Pool, actor(c).UserID.String()); e != nil {
			respond(c, nil, e)
			return
		}
		c.Next()
	})
	for _, g := range []*gin.RouterGroup{web, worker} {
		g.GET("/dashboard", func(c *gin.Context) { out, e := s.Dashboard(c.Request.Context(), actor(c)); respond(c, out, e) })
		g.GET("/me", func(c *gin.Context) { out, e := s.Me(c.Request.Context(), actor(c)); respond(c, out, e) })
		g.GET("/schedules", list(s.Schedules))
		g.GET("/leave-requests", list(s.Leaves))
		g.GET("/attendance", list(s.Attendance))
		g.GET("/requisitions", list(s.Requisitions))
		g.GET("/requisitions/:id", func(c *gin.Context) {
			out, e := s.RequisitionDetail(c.Request.Context(), actor(c), c.Param("id"))
			respond(c, out, e)
		})
		g.GET("/equipment", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
			return s.Equipments(ctx, l, o)
		}))
	}
	web.GET("/contracts", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
		return s.Contracts(ctx, l, o)
	}))
	web.GET("/contracts/:id", func(c *gin.Context) { out, e := s.Contract(c.Request.Context(), c.Param("id")); respond(c, out, e) })
	web.GET("/users", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
		return s.Users(ctx, l, o)
	}))
	sup := web.Group("", auth.RequireRole(auth.RoleSupervisor))
	asst := web.Group("", auth.RequireRole(auth.RoleAssistant))
	sup.POST("/users", input(func(c *gin.Context, p auth.Principal, in UserInput) (any, error) {
		return s.CreateUser(c.Request.Context(), p, in)
	}))
	sup.PATCH("/users/:id", input(func(c *gin.Context, _ auth.Principal, in UserUpdate) (any, error) {
		return nil, s.UpdateUser(c.Request.Context(), c.Param("id"), in)
	}))
	sup.PATCH("/contracts/:id", input(func(c *gin.Context, _ auth.Principal, in stateInput) (any, error) {
		return nil, s.ContractState(c.Request.Context(), c.Param("id"), in.Status)
	}))
	sup.POST("/locations", input(func(c *gin.Context, _ auth.Principal, in LocationInput) (any, error) {
		id, e := s.Location(c.Request.Context(), "", in)
		return gin.H{"location_id": id}, e
	}))
	sup.PATCH("/locations/:id", input(func(c *gin.Context, _ auth.Principal, in LocationInput) (any, error) {
		id, e := s.Location(c.Request.Context(), c.Param("id"), in)
		return gin.H{"location_id": id}, e
	}))
	sup.POST("/equipment", input(func(c *gin.Context, _ auth.Principal, in EquipmentInput) (any, error) {
		id, e := s.Equipment(c.Request.Context(), "", in)
		return gin.H{"equipment_id": id}, e
	}))
	sup.PATCH("/equipment/:id", input(func(c *gin.Context, _ auth.Principal, in EquipmentInput) (any, error) {
		id, e := s.Equipment(c.Request.Context(), c.Param("id"), in)
		return gin.H{"equipment_id": id}, e
	}))
	asst.POST("/schedules", input(func(c *gin.Context, _ auth.Principal, in ScheduleInput) (any, error) {
		ids, e := s.Schedule(c.Request.Context(), in)
		return gin.H{"schedule_ids": ids}, e
	}))
	sup.POST("/leave-requests/:id/review", input(func(c *gin.Context, _ auth.Principal, in stateInput) (any, error) {
		return nil, s.ReviewLeave(c.Request.Context(), c.Param("id"), in.Status)
	}))
	asst.POST("/leave-requests/:id/replacement", input(func(c *gin.Context, _ auth.Principal, in replacementInput) (any, error) {
		ids, e := s.Replacement(c.Request.Context(), c.Param("id"), in.WorkerID)
		return gin.H{"schedule_ids": ids}, e
	}))
	worker.POST("/leave-requests", input(func(c *gin.Context, p auth.Principal, in LeaveInput) (any, error) {
		id, e := s.Leave(c.Request.Context(), p, in)
		return gin.H{"request_id": id}, e
	}))
	worker.POST("/attendance/check-in", input(func(c *gin.Context, p auth.Principal, in CheckIn) (any, error) {
		return nil, s.CheckIn(c.Request.Context(), p, in)
	}))
	worker.POST("/attendance/check-out", input(func(c *gin.Context, p auth.Principal, in CheckOut) (any, error) {
		return nil, s.CheckOut(c.Request.Context(), p, in)
	}))
	web.POST("/assignments/:id/qr", func(c *gin.Context) { out, e := s.QR(c.Request.Context(), c.Param("id")); respond(c, out, e) })
	sup.POST("/attendance/finalize", func(c *gin.Context) { n, e := s.Finalize(c.Request.Context()); respond(c, gin.H{"changed": n}, e) })
	worker.POST("/requisitions", input(func(c *gin.Context, p auth.Principal, in RequestInput) (any, error) {
		id, e := s.Requisition(c.Request.Context(), p, in)
		return gin.H{"requisition_id": id}, e
	}))
	asst.POST("/requisitions/:id/survey", input(func(c *gin.Context, _ auth.Principal, in surveyInput) (any, error) {
		return nil, s.Survey(c.Request.Context(), c.Param("id"), in.Items)
	}))
	sup.POST("/requisitions/:id/review", input(func(c *gin.Context, p auth.Principal, in ReviewInput) (any, error) {
		return nil, s.ReviewRequest(c.Request.Context(), p, c.Param("id"), in)
	}))
	sup.POST("/requisitions/:id/fund-transfers", input(func(c *gin.Context, p auth.Principal, in FundInput) (any, error) {
		id, e := s.Fund(c.Request.Context(), p, c.Param("id"), in)
		return gin.H{"expense_id": id}, e
	}))
	asst.POST("/requisitions/:id/purchase", input(func(c *gin.Context, p auth.Principal, in PurchaseInput) (any, error) {
		id, e := s.Purchase(c.Request.Context(), p, c.Param("id"), in)
		return gin.H{"expense_id": id}, e
	}))
	sup.POST("/payroll/preview", input(func(c *gin.Context, _ auth.Principal, in PayrollInput) (any, error) {
		return s.PreviewPayroll(c.Request.Context(), in)
	}))
	sup.POST("/payroll", input(func(c *gin.Context, p auth.Principal, in PayrollInput) (any, error) {
		id, e := s.Payroll(c.Request.Context(), p, in)
		return gin.H{"payroll_id": id}, e
	}))
	sup.GET("/payroll", list(s.Payrolls))
	worker.GET("/payroll", list(s.Payrolls))
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
		out, e := s.Profit(c.Request.Context(), c.Query("month"))
		respond(c, gin.H{"data": out}, e)
	})
	sup.GET("/reports/profit.csv", s.ExportProfit)
}
