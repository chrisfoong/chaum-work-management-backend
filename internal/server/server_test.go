package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type okPinger struct{}

func (okPinger) Ping(context.Context) error { return nil }

func TestRouterGroupsUseTheirAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deny := func(c *gin.Context) { c.AbortWithStatus(http.StatusUnauthorized) }
	allow := func(c *gin.Context) { c.Next() }

	r := New(Deps{DB: okPinger{}, WebAuth: deny, WorkerAuth: allow})
	r.Web.GET("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.Worker.GET("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })

	tests := []struct {
		path string
		want int
	}{
		{"/health", http.StatusOK},
		{"/api/web/probe", http.StatusUnauthorized},
		{"/api/liff/probe", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}
