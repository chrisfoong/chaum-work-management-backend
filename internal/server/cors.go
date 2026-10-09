package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CORS allows only the configured frontend origins. It never reflects an
// arbitrary origin and does not authorize a request on behalf of its caller.
func CORS(origins string) gin.HandlerFunc {
	allow := map[string]bool{}
	for _, v := range strings.Split(origins, ",") {
		v = strings.TrimSpace(v)
		if v != "" && v != "*" {
			allow[v] = true
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if !allow[origin] {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
			c.Header("Access-Control-Expose-Headers", "X-Request-ID")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		c.Next()
	}
}
