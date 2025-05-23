package errors

import (
	"fmt"
	"net/http"

	stdErrors "errors" // Aliasing standard errors package
)

// Export standard errors functions
var (
	Is = stdErrors.Is
)

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(message string) *Error {
	return NewError(http.StatusUnauthorized, "authorization_error", message)
}

// Error represents a custom error with additional context
type Error struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	err        error  // Internal error for wrapping
}

// NewError creates a new Error instance
// Supports both:
// - NewError(statusCode int, message string) to create error with auto-generated code
// - NewError(statusCode int, code, message string) to create error with custom code
func NewError(statusCode int, messageOptional ...string) *Error {
	if len(messageOptional) == 0 {
		code := fmt.Sprintf("%d", statusCode)
		return &Error{
			StatusCode: statusCode,
			Code:       code,
			Message:    code,
		}
	}

	codeOrMessage := messageOptional[0]
	err := &Error{
		StatusCode: statusCode,
		Code:       codeOrMessage,
		Message:    codeOrMessage,
	}
	if len(messageOptional) > 1 {
		err.Message = messageOptional[1]
	}

	return err
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// WithError wraps an existing error
func (e *Error) WithError(err error) *Error {
	if err == nil {
		return e
	}
	return &Error{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    e.Message,
		err:        err,
	}
}

// WithMessage creates a new error with a custom message
func (e *Error) WithMessage(message string) *Error {
	return &Error{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    message,
		err:        e.err,
	}
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.err
}

// Is implements error Is interface for error comparison
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Helper functions for creating common errors

func NewValidationError(message string) *Error {
	return NewError(http.StatusBadRequest, message)
}

func NewBadRequestError(message string) *Error {
	return NewError(http.StatusBadRequest, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(http.StatusNotFound, message)
}

func NewConflictError(message string) *Error {
	return NewError(http.StatusConflict, message)
}

func NewInternalServerError(message string) *Error {
	return NewError(http.StatusInternalServerError, message)
}

func NewUnauthorizedError(message string) *Error {
	return NewError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *Error {
	return NewError(http.StatusForbidden, message)
}

// WrapError wraps an error with additional context while preserving the original error
func WrapError(err error, code int, message string) error {
	if err == nil {
		return nil
	}

	// If it's already our custom error type, wrap the message
	if customErr, ok := err.(*Error); ok {
		return customErr.WithMessage(fmt.Sprintf("%s: %s", message, customErr.Message))
	}

	// Otherwise create a new error with the given code
	return NewError(code, "wrapped_error", fmt.Sprintf("%s: %v", message, err))
}

// As returns true if the target implements the Error interface
func As(err error) (*Error, bool) {
	var e *Error
	if stdErrors.As(err, &e) {
		return e, true
	}
	return nil, false
}
