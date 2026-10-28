package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func registerSchedulingRoutes(web, worker, sup, asst *gin.RouterGroup, s *Service) {
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
}
