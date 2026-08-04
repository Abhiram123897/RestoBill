package errs

import (
	"net/http"
)

// Error type constants
const (
	BadRequestErrType     = "Bad request"
	AuthenticationErrType = "Authentication error"
	AuthorizationErrType  = "Authorization error"
	NotFoundErrType       = "Not found"
	ConflictErrType       = "Conflict"
	ValidationErrType     = "Validation error"
	InternalErrType       = "Internal server error"
)

// AppError represents a structured API error
type AppError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewBadRequestError returns a 400 error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    BadRequestErrType,
		Code:    http.StatusBadRequest,
	}
}

// NewUnauthorizedError returns a 401 error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    AuthenticationErrType,
		Code:    http.StatusUnauthorized,
	}
}

// NewForbiddenError returns a 403 error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    AuthorizationErrType,
		Code:    http.StatusForbidden,
	}
}

// NewNotFoundError returns a 404 error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    NotFoundErrType,
		Code:    http.StatusNotFound,
	}
}

// NewConflictError returns a 409 error
func NewConflictError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    ConflictErrType,
		Code:    http.StatusConflict,
	}
}

// NewValidationError returns a 422 error
func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    ValidationErrType,
		Code:    http.StatusUnprocessableEntity,
	}
}

// NewInternalServerError returns a 500 error
func NewInternalServerError(message string) *AppError {
	return &AppError{
		Message: message,
		Type:    InternalErrType,
		Code:    http.StatusInternalServerError,
	}
}
