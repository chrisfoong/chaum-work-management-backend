// Package server builds the Gin engine and the route groups workflow slices register on.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"

	"chrisfoong/chaum-work-management-backend/internal/health"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
)

// Deps are the collaborators the router needs.
type Deps struct {
	AllowedOrigins string
	DB             health.Pinger
	WebAuth        gin.HandlerFunc // Supervisor and Assistant (desktop web)
	WorkerAuth     gin.HandlerFunc // Worker (LINE Mini App only)
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
	r.Use(func(c *gin.Context) {
		id := uuid.NewString()
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 6*1024*1024)
		defer func() {
			if recover() != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": gin.H{"code": "internal", "message": "internal server error", "request_id": id}})
			}
		}()
		c.Next()
	}, httpx.RequestLogger())
	r.Use(CORS(d.AllowedOrigins))
	r.GET("/health", health.Handler(d.DB))
	r.GET("/health/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/health/ready", health.Handler(d.DB))

	return &Router{
		Engine: r,
		Web:    r.Group("/api/web", d.WebAuth),
		Worker: r.Group("/api/liff", d.WorkerAuth),
	}
}
