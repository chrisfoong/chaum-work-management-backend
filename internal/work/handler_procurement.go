package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func registerProcurementRoutes(worker, sup, asst *gin.RouterGroup, s *Service) {
	asst.GET("/requisitions/:id/delivery-schedules", func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		out, e := s.DeliverySchedules(c.Request.Context(), c.Param("id"), l, o)
		respond(c, gin.H{"data": out, "limit": l, "offset": o}, e)
	})
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
	asst.POST("/requisitions/:id/purchase/preview", input(func(c *gin.Context, p auth.Principal, in PurchaseInput) (any, error) {
		return s.PreviewPurchase(c.Request.Context(), p, c.Param("id"), in)
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
		if e := s.DecideRequest(c.Request.Context(), p, c.Param("id"), in); e != nil {
			return nil, e
		}
		kind := "approval_needed"
		if in.Decision == "no_purchase" {
			kind = "no_purchase"
		}
		out, _ := s.RetryNotification(c.Request.Context(), p, c.Param("id"), NotificationInput{Kind: kind, Reason: in.Reason})
		return gin.H{"business_saved": true, "notification": out}, nil
	}))
	asst.POST("/requisitions/:id/delivery", input(func(c *gin.Context, p auth.Principal, in DeliveryInput) (any, error) {
		if e := s.Deliver(c.Request.Context(), p, c.Param("id"), in); e != nil {
			return nil, e
		}
		out, _ := s.RetryNotification(c.Request.Context(), p, c.Param("id"), NotificationInput{Kind: "delivery"})
		return gin.H{"business_saved": true, "notification": out}, nil
	}))
	asst.POST("/notifications/:id/retry", input(func(c *gin.Context, p auth.Principal, in NotificationInput) (any, error) {
		return s.RetryNotification(c.Request.Context(), p, c.Param("id"), in)
	}))
	asst.GET("/requisitions/:id/deliveries", func(c *gin.Context) { out, e := s.Deliveries(c.Request.Context(), c.Param("id")); respond(c, out, e) })
}
