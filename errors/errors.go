// Package errors provides framework-wide error handling utilities.
package errors

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
