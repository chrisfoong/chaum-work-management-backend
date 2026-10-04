package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type pingerFunc func(context.Context) error

func (f pingerFunc) Ping(ctx context.Context) error { return f(ctx) }

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		pingErr error
		want    int
	}{
		{name: "database up", want: http.StatusOK},
		{name: "database down", pingErr: errors.New("refused"), want: http.StatusServiceUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/health", Handler(pingerFunc(func(context.Context) error { return tt.pingErr })))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}
