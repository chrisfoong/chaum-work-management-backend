package server

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowlist(t *testing.T) {
	r := gin.New()
	r.Use(CORS("https://trusted.example"))
	r.GET("/x", func(c *gin.Context) { c.Status(200) })
	r.OPTIONS("/x", func(c *gin.Context) { c.Status(401) })
	for _, tt := range []struct {
		method, origin string
		status         int
	}{{"GET", "https://trusted.example", 200}, {"GET", "https://evil.example", 403}, {"OPTIONS", "https://trusted.example", 204}, {"GET", "", 200}} {
		q := httptest.NewRequest(tt.method, "/x", nil)
		q.Header.Set("Origin", tt.origin)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		if w.Code != tt.status {
			t.Errorf("%s %s: %d", tt.method, tt.origin, w.Code)
		}
	}
}
