// Package validation provides framework-wide validation utilities.
package validation

import (
	"errors"
	"fmt"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the error message.
func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", v.Field, v.Message)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
