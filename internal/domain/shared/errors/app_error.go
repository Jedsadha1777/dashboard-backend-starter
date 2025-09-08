package errors

import (
	"fmt"
	"time"
)

// AppError represents standardized application error
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	StatusCode int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Path       string      `json:"path,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
}

// Error codes
const (
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeRateLimit         = "RATE_LIMIT_EXCEEDED"
	ErrCodeDatabaseError     = "DATABASE_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeTokenExpired      = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid      = "TOKEN_INVALID"
	ErrCodeDuplicateEntry    = "DUPLICATE_ENTRY"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
)

// Error implements error interface
func (e AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code string, statusCode int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

// WithDetails adds details to error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// WithPath adds request path to error
func (e *AppError) WithPath(path string) *AppError {
	e.Path = path
	return e
}

// WithRequestID adds request ID to error
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// Common errors
var (
	ErrInternalServer = NewAppError(ErrCodeInternal, 500, "Internal server error")
	ErrUnauthorized   = NewAppError(ErrCodeUnauthorized, 401, "Unauthorized")
	ErrForbidden      = NewAppError(ErrCodeForbidden, 403, "Forbidden")
	ErrNotFound       = NewAppError(ErrCodeNotFound, 404, "Resource not found")
	ErrBadRequest     = NewAppError(ErrCodeBadRequest, 400, "Bad request")
	ErrConflict       = NewAppError(ErrCodeConflict, 409, "Resource already exists")
	ErrRateLimit      = NewAppError(ErrCodeRateLimit, 429, "Rate limit exceeded")
)

// ValidationError creates validation error with field details
func ValidationError(fields map[string]string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		StatusCode: 400,
		Details:    fields,
		Timestamp:  time.Now(),
	}
}

// DatabaseError wraps database errors
func DatabaseError(err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    "Database operation failed",
		StatusCode: 500,
		Details:    err.Error(),
		Timestamp:  time.Now(),
	}
}

// WrapError wraps standard error into AppError
func WrapError(err error, code string, statusCode int) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return &AppError{
		Code:       code,
		Message:    err.Error(),
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

// IsAppError checks if error is AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetStatusCode returns HTTP status code from error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return 500
}

// ToResponse converts AppError to response format
func (e *AppError) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"success":   false,
		"error":     e.Message,
		"code":      e.Code,
		"timestamp": e.Timestamp,
	}

	if e.Details != nil {
		response["details"] = e.Details
	}

	if e.Path != "" {
		response["path"] = e.Path
	}

	if e.RequestID != "" {
		response["request_id"] = e.RequestID
	}

	return response
}

// Specific business errors
func EmailAlreadyExists(email string) *AppError {
	return &AppError{
		Code:       ErrCodeDuplicateEntry,
		Message:    fmt.Sprintf("Email %s already exists", email),
		StatusCode: 409,
		Timestamp:  time.Now(),
	}
}

func InvalidCredentials() *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Invalid email or password",
		StatusCode: 401,
		Timestamp:  time.Now(),
	}
}

func TokenExpired() *AppError {
	return &AppError{
		Code:       ErrCodeTokenExpired,
		Message:    "Token has expired",
		StatusCode: 401,
		Timestamp:  time.Now(),
	}
}

func InsufficientStock(productName string, available int) *AppError {
	return &AppError{
		Code:       ErrCodeInsufficientStock,
		Message:    fmt.Sprintf("Insufficient stock for %s", productName),
		StatusCode: 400,
		Details: map[string]interface{}{
			"product":   productName,
			"available": available,
		},
		Timestamp: time.Now(),
	}
}
