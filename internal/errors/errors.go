package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeUnauthorized represents unauthorized errors
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	// ErrorTypeForbidden represents forbidden errors
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeConflict represents conflict errors
	ErrorTypeConflict ErrorType = "conflict"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeExternal represents external service errors
	ErrorTypeExternal ErrorType = "external"
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout ErrorType = "timeout"
	// ErrorTypeRateLimit represents rate limit errors
	ErrorTypeRateLimit ErrorType = "rate_limit"
	// ErrorTypeBadRequest represents bad request errors
	ErrorTypeBadRequest ErrorType = "bad_request"
)

// AppError represents an application error with context
type AppError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Cause     error                  `json:"-"`
	Operation string                 `json:"operation,omitempty"`
	Component string                 `json:"component,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// HTTPStatusCode returns the appropriate HTTP status code for the error
func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeExternal:
		return http.StatusBadGateway
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WithField adds a field to the error
func (e *AppError) WithField(key string, value interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[key] = value
	return e
}

// WithFields adds multiple fields to the error
func (e *AppError) WithFields(fields map[string]interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// WithComponent sets the component where the error occurred
func (e *AppError) WithComponent(component string) *AppError {
	e.Component = component
	return e
}

// WithOperation sets the operation where the error occurred
func (e *AppError) WithOperation(operation string) *AppError {
	e.Operation = operation
	return e
}

// WithStack captures the current stack trace
func (e *AppError) WithStack() *AppError {
	e.Stack = captureStack()
	return e
}

// New creates a new AppError
func New(errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Stack:   captureStack(),
	}
}

// Wrap wraps an existing error as an AppError
func Wrap(err error, errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Cause:   err,
		Stack:   captureStack(),
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, errorType ErrorType, code, format string, args ...interface{}) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
		Stack:   captureStack(),
	}
}

// ValidationError creates a validation error
func ValidationError(code, message string) *AppError {
	return New(ErrorTypeValidation, code, message)
}

// NotFoundError creates a not found error
func NotFoundError(code, message string) *AppError {
	return New(ErrorTypeNotFound, code, message)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(code, message string) *AppError {
	return New(ErrorTypeUnauthorized, code, message)
}

// ForbiddenError creates a forbidden error
func ForbiddenError(code, message string) *AppError {
	return New(ErrorTypeForbidden, code, message)
}

// ConflictError creates a conflict error
func ConflictError(code, message string) *AppError {
	return New(ErrorTypeConflict, code, message)
}

// InternalError creates an internal server error
func InternalError(code, message string) *AppError {
	return New(ErrorTypeInternal, code, message)
}

// ExternalError creates an external service error
func ExternalError(code, message string) *AppError {
	return New(ErrorTypeExternal, code, message)
}

// TimeoutError creates a timeout error
func TimeoutError(code, message string) *AppError {
	return New(ErrorTypeTimeout, code, message)
}

// RateLimitError creates a rate limit error
func RateLimitError(code, message string) *AppError {
	return New(ErrorTypeRateLimit, code, message)
}

// BadRequestError creates a bad request error
func BadRequestError(code, message string) *AppError {
	return New(ErrorTypeBadRequest, code, message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts an error to AppError if possible
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// IsType checks if an error is of a specific type
func IsType(err error, errorType ErrorType) bool {
	if appErr, ok := AsAppError(err); ok {
		return appErr.Type == errorType
	}
	return false
}

// IsCode checks if an error has a specific code
func IsCode(err error, code string) bool {
	if appErr, ok := AsAppError(err); ok {
		return appErr.Code == code
	}
	return false
}

// captureStack captures the current stack trace
func captureStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	
	// Parse stack trace to remove internal error handling frames
	stack := string(buf)
	lines := strings.Split(stack, "\n")
	
	// Find the first frame that's not in this package
	var filtered []string
	skip := true
	for _, line := range lines {
		if strings.Contains(line, "go-transport-prac/internal/errors") && skip {
			continue
		}
		skip = false
		filtered = append(filtered, line)
	}
	
	return strings.Join(filtered, "\n")
}

// Common error codes
const (
	// Validation error codes
	CodeValidationFailed    = "VALIDATION_FAILED"
	CodeInvalidInput        = "INVALID_INPUT"
	CodeMissingField        = "MISSING_FIELD"
	CodeInvalidFormat       = "INVALID_FORMAT"
	CodeInvalidValue        = "INVALID_VALUE"
	
	// Authentication/Authorization codes
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeInvalidCredentials  = "INVALID_CREDENTIALS"
	CodeTokenExpired        = "TOKEN_EXPIRED"
	CodeInvalidToken        = "INVALID_TOKEN"
	CodeForbidden           = "FORBIDDEN"
	CodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
	
	// Resource error codes
	CodeNotFound            = "NOT_FOUND"
	CodeAlreadyExists       = "ALREADY_EXISTS"
	CodeConflict            = "CONFLICT"
	CodeResourceLocked      = "RESOURCE_LOCKED"
	
	// System error codes
	CodeInternalError       = "INTERNAL_ERROR"
	CodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	CodeTimeout             = "TIMEOUT"
	CodeRateLimit           = "RATE_LIMIT_EXCEEDED"
	CodeDatabaseError       = "DATABASE_ERROR"
	CodeExternalService     = "EXTERNAL_SERVICE_ERROR"
	
	// Network error codes
	CodeConnectionError     = "CONNECTION_ERROR"
	CodeNetworkTimeout      = "NETWORK_TIMEOUT"
	CodeDNSError           = "DNS_ERROR"
	
	// Data processing codes
	CodeSerializationError  = "SERIALIZATION_ERROR"
	CodeDeserializationError = "DESERIALIZATION_ERROR"
	CodeEncodingError       = "ENCODING_ERROR"
	CodeDecodingError       = "DECODING_ERROR"
)

// Predefined common errors
var (
	ErrValidationFailed    = ValidationError(CodeValidationFailed, "Validation failed")
	ErrInvalidInput        = ValidationError(CodeInvalidInput, "Invalid input provided")
	ErrUnauthorized        = UnauthorizedError(CodeUnauthorized, "Unauthorized access")
	ErrForbidden           = ForbiddenError(CodeForbidden, "Access forbidden")
	ErrNotFound            = NotFoundError(CodeNotFound, "Resource not found")
	ErrConflict            = ConflictError(CodeConflict, "Resource conflict")
	ErrInternalError       = InternalError(CodeInternalError, "Internal server error")
	ErrServiceUnavailable  = InternalError(CodeServiceUnavailable, "Service unavailable")
	ErrTimeout             = TimeoutError(CodeTimeout, "Operation timed out")
	ErrRateLimit           = RateLimitError(CodeRateLimit, "Rate limit exceeded")
)