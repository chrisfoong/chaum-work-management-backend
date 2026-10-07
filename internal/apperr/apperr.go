// Package apperr defines typed application errors that map to HTTP responses.
package apperr

import "net/http"

// Kind classifies an error for HTTP mapping.
type Kind uint8

const (
	KindInternal Kind = iota
	KindValidation
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
)

// FieldError names one failing input field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error is the error type services return. Message is safe to show to clients;
// Err holds internal detail that is logged but never sent.
type Error struct {
	Kind    Kind
	Code    string
	Message string
	Fields  []FieldError
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

// WithFields attaches failing fields to e and returns it.
func (e *Error) WithFields(fields ...FieldError) *Error {
	e.Fields = append(e.Fields, fields...)
	return e
}

// Validation reports invalid input, naming each failing field.
func Validation(fields ...FieldError) *Error {
	return &Error{Kind: KindValidation, Code: "validation_failed", Message: "request validation failed", Fields: fields}
}

// Unauthorized reports a missing or invalid identity.
func Unauthorized(msg string) *Error {
	return &Error{Kind: KindUnauthorized, Code: "unauthorized", Message: msg}
}

// Forbidden reports an identity that may not perform the action.
func Forbidden(msg string) *Error {
	return &Error{Kind: KindForbidden, Code: "forbidden", Message: msg}
}

// NotFound reports a missing resource.
func NotFound(msg string) *Error {
	return &Error{Kind: KindNotFound, Code: "not_found", Message: msg}
}

// Conflict reports a rule violation against current data, e.g. code "duplicate_contract_no".
func Conflict(code, msg string) *Error {
	return &Error{Kind: KindConflict, Code: code, Message: msg}
}

// Internal wraps an unexpected error. Clients only see a generic message.
func Internal(err error) *Error {
	return &Error{Kind: KindInternal, Code: "internal", Message: "internal server error", Err: err}
}

// HTTPStatus returns the status code for a kind.
func HTTPStatus(k Kind) int {
	switch k {
	case KindValidation:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
