// Package health serves the liveness/readiness endpoint.
package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Pinger checks the database connection. *pgxpool.Pool satisfies it.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Handler reports 200 when the database answers and 503 when it does not.
func Handler(db Pinger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			slog.WarnContext(ctx, "health check: database unreachable", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "database": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "up"})
	}
}
