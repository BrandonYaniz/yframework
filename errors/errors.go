// Package errors provides framework-wide error handling utilities.
package errors

import (
	"errors"
	"fmt"
)

// FrameworkError represents a generic framework error.
type FrameworkError struct {
	// Message describes the error.
	Message string

	// Cause holds the underlying error, if any.
	Cause error

	// Code is an optional error code for categorization.
	Code string
}

// Error returns the error message.
func (e *FrameworkError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error.
func (e *FrameworkError) Unwrap() error {
	return e.Cause
}

// New creates a new FrameworkError with the given message and code.
func New(message, code string) *FrameworkError {
	return &FrameworkError{
		Message: message,
		Code:    code,
	}
}

// NewWithCause creates a new FrameworkError with the given message, code, and cause.
func NewWithCause(message, code string, cause error) *FrameworkError {
	return &FrameworkError{
		Message: message,
		Cause:   cause,
		Code:    code,
	}
}

// Is checks if an error is of type FrameworkError and has the specified code.
func Is(err error, code string) bool {
	var fe *FrameworkError
	if errors.As(err, &fe) {
		return fe.Code == code
	}
	return false
}

// ValidationError represents a validation error with field information.
type ValidationError struct {
	// Field is the name of the field that failed validation.
	Field string

	// Message describes the validation failure.
	Message string

	// Code is an optional error code for categorization.
	Code string
}

// Error returns the error message.
func (v *ValidationError) Error() string {
	if v.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", v.Field, v.Message)
	}
	return fmt.Sprintf("validation error: %s", v.Message)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message, code string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	}
}
