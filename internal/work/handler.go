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
			c.Abort()
			return
		}
		c.Next()
	})
	for _, g := range []*gin.RouterGroup{web, worker} {
		g.GET("/dashboard", func(c *gin.Context) { out, e := s.Dashboard(c.Request.Context(), actor(c)); respond(c, out, e) })
		g.GET("/me", func(c *gin.Context) { out, e := s.Me(c.Request.Context(), actor(c)); respond(c, out, e) })
		g.GET("/schedules", filteredList(s, "schedules"))
		g.GET("/leave-requests", filteredList(s, "leave"))
		g.GET("/attendance", filteredList(s, "attendance"))
		g.GET("/requisitions", filteredList(s, "requisitions"))
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
	web.GET("/assignments", func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		out, e := s.Assignments(c.Request.Context(), c.Query("tor_id"), l, o)
		respond(c, gin.H{"data": out, "limit": l, "offset": o}, e)
	})
	asst.GET("/workers/available", func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		out, e := s.AvailableWorkers(c.Request.Context(), c.Query("work_date"), "", l, o)
		respond(c, gin.H{"data": out, "limit": l, "offset": o}, e)
	})
	asst.GET("/leave-requests/:id/candidates", func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		out, e := s.LeaveCandidates(c.Request.Context(), c.Param("id"), l, o)
		respond(c, gin.H{"data": out, "limit": l, "offset": o}, e)
	})
	asst.GET("/leave-requests/:id", func(c *gin.Context) { out, e := s.LeaveDetail(c.Request.Context(), c.Param("id")); respond(c, out, e) })
	asst.GET("/requisitions/:id/delivery-schedules", func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		out, e := s.DeliverySchedules(c.Request.Context(), c.Param("id"), l, o)
		respond(c, gin.H{"data": out, "limit": l, "offset": o}, e)
	})
	asst.POST("/schedules", input(func(c *gin.Context, _ auth.Principal, in ScheduleInput) (any, error) {
		ids, e := s.Schedule(c.Request.Context(), in)
		return gin.H{"schedule_ids": ids}, e
	}))
	asst.POST("/leave-requests/:id/review", input(func(c *gin.Context, _ auth.Principal, in LeaveReviewInput) (any, error) {
		ids, e := s.ReviewLeaveAndReplace(c.Request.Context(), c.Param("id"), in)
		return gin.H{"schedule_ids": ids}, e
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
	asst.GET("/requisitions/:id/inspection", func(c *gin.Context) {
		out, e := s.InspectRequest(c.Request.Context(), c.Param("id"))
		respond(c, out, e)
	})
	asst.POST("/requisitions/:id/decision", input(func(c *gin.Context, p auth.Principal, in ProcurementDecision) (any, error) {
		return nil, s.DecideRequest(c.Request.Context(), p, c.Param("id"), in)
	}))
	asst.POST("/requisitions/:id/delivery", input(func(c *gin.Context, p auth.Principal, in DeliveryInput) (any, error) {
		return nil, s.Deliver(c.Request.Context(), p, c.Param("id"), in)
	}))
	asst.GET("/requisitions/:id/deliveries", func(c *gin.Context) { out, e := s.Deliveries(c.Request.Context(), c.Param("id")); respond(c, out, e) })
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
