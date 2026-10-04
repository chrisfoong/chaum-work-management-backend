package apperr

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		err  *Error
		want int
	}{
		{Validation(FieldError{Field: "x", Message: "is required"}), http.StatusBadRequest},
		{Unauthorized("no token"), http.StatusUnauthorized},
		{Forbidden("no"), http.StatusForbidden},
		{NotFound("gone"), http.StatusNotFound},
		{Conflict("duplicate_contract_no", "dup"), http.StatusConflict},
		{Internal(errors.New("db down")), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.err.Code, func(t *testing.T) {
			if got := HTTPStatus(tt.err.Kind); got != tt.want {
				t.Fatalf("HTTPStatus = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestInternalWrapsCause(t *testing.T) {
	cause := errors.New("db down")
	err := Internal(cause)
	if !errors.Is(err, cause) {
		t.Fatal("Internal should unwrap to its cause")
	}
	if err.Message != "internal server error" {
		t.Fatalf("client message leaks detail: %q", err.Message)
	}
}
