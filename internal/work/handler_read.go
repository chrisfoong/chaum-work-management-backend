package work

import "github.com/gin-gonic/gin"

func registerReadRoutes(web, worker *gin.RouterGroup, s *Service) {
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
		g.GET("/equipment", func(c *gin.Context) {
			l, o, e := paging(c)
			if e != nil {
				respond(c, nil, e)
				return
			}
			out, e := s.AssignmentEquipment(c.Request.Context(), actor(c), c.Query("assignment_id"), l, o)
			respond(c, out, e)
		})
	}
}
