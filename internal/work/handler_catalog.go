package work

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func registerCatalogRoutes(web, sup *gin.RouterGroup, s *Service) {
	web.GET("/contracts", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
		return s.Contracts(ctx, l, o)
	}))
	web.GET("/contracts/:id", func(c *gin.Context) { out, e := s.Contract(c.Request.Context(), c.Param("id")); respond(c, out, e) })
	web.GET("/users", list(func(ctx context.Context, _ auth.Principal, l, o int) (json.RawMessage, error) {
		return s.Users(ctx, l, o)
	}))
	sup.POST("/invitations", input(func(c *gin.Context, p auth.Principal, in InvitationUserInput) (any, error) {
		return s.CreateInvitation(c.Request.Context(), p, in)
	}))
	sup.POST("/users/:id/invitation", input(func(c *gin.Context, p auth.Principal, _ struct{}) (any, error) {
		return s.RenewInvitation(c.Request.Context(), p, c.Param("id"))
	}))
	sup.POST("/users/:id/invitation/cancel", input(func(c *gin.Context, p auth.Principal, _ struct{}) (any, error) {
		return s.CancelInvitation(c.Request.Context(), p, c.Param("id"))
	}))
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
}
