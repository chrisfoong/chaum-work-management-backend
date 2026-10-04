// Package httpx holds Gin helpers shared by all handlers: JSON binding with
// field-level validation errors, the error response writer and request logging.
package httpx

import (
	"errors"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
)

// ErrorBody is the JSON shape of every error response.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail describes one error response.
type ErrorDetail struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Fields  []apperr.FieldError `json:"fields,omitempty"`
}

// WriteError aborts the request with err as an ErrorBody. Errors that are not
// *apperr.Error become 500s; their detail is logged, never sent.
func WriteError(c *gin.Context, err error) {
	var ae *apperr.Error
	if !errors.As(err, &ae) {
		ae = apperr.Internal(err)
	}
	if ae.Kind == apperr.KindInternal {
		slog.ErrorContext(c.Request.Context(), "request failed", "route", c.FullPath(), "error", err)
	}
	c.AbortWithStatusJSON(apperr.HTTPStatus(ae.Kind), ErrorBody{Error: ErrorDetail{
		Code:    ae.Code,
		Message: ae.Message,
		Fields:  ae.Fields,
	}})
}

// UseJSONFieldNames makes validation errors name fields by their json tag.
// Call once at startup.
func UseJSONFieldNames() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		switch name {
		case "-":
			return ""
		case "":
			return f.Name
		}
		return name
	})
}

// BindJSON decodes and validates the request body into dst. Failures are
// returned as an apperr validation error that names each failing field.
func BindJSON(c *gin.Context, dst any) error {
	err := c.ShouldBindJSON(dst)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return apperr.Validation(apperr.FieldError{Field: "body", Message: "must be valid JSON of the expected shape"})
	}
	fields := make([]apperr.FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, apperr.FieldError{Field: fe.Field(), Message: fieldMessage(fe)})
	}
	return apperr.Validation(fields...)
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min", "gte":
		return "must be at least " + fe.Param()
	case "max", "lte":
		return "must be at most " + fe.Param()
	case "oneof":
		return "must be one of: " + fe.Param()
	case "uuid":
		return "must be a UUID"
	default:
		return "is invalid (" + fe.Tag() + ")"
	}
}

// RequestLogger logs one line per request. It logs the route template, not the
// raw path or query string, so tokens or personal data in URLs are not recorded.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.InfoContext(c.Request.Context(), "http request",
			"method", c.Request.Method,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}
