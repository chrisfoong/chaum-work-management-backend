package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
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
	sup := web.Group("", auth.RequireRole(auth.RoleSupervisor))
	asst := web.Group("", auth.RequireRole(auth.RoleAssistant))

	registerReadRoutes(web, worker, s)
	registerCatalogRoutes(web, sup, s)
	registerSchedulingRoutes(web, worker, sup, asst, s)
	registerProcurementRoutes(worker, sup, asst, s)
	registerReportingRoutes(worker, sup, asst, s)
}
