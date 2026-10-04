// Package server builds the Gin engine and the route groups workflow slices register on.
package server

import (
	"github.com/gin-gonic/gin"

	"chrisfoong/chaum-work-management-backend/internal/health"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
)

// Deps are the collaborators the router needs.
type Deps struct {
	DB         health.Pinger
	WebAuth    gin.HandlerFunc // Supervisor and Assistant (desktop web)
	WorkerAuth gin.HandlerFunc // Worker (LINE Mini App only)
}

// Router holds the engine and its authenticated groups.
type Router struct {
	Engine *gin.Engine
	Web    *gin.RouterGroup // /api/web: Supervisor and Assistant endpoints
	Worker *gin.RouterGroup // /api/liff: Worker endpoints
}

// New builds the router with /health and the two authenticated groups.
func New(d Deps) *Router {
	r := gin.New()
	r.Use(gin.Recovery(), httpx.RequestLogger())
	r.GET("/health", health.Handler(d.DB))

	return &Router{
		Engine: r,
		Web:    r.Group("/api/web", d.WebAuth),
		Worker: r.Group("/api/liff", d.WorkerAuth),
	}
}
