// Package validation provides framework-wide validation utilities.
package validation

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
	Code    string
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

// ValidationErrors represents a collection of validation errors.
type ValidationErrors struct {
	Errors []ValidationError
}

// Error returns a formatted string of all validation errors.
func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}

	var errorMessages []string
	for _, err := range v.Errors {
		errorMessages = append(errorMessages, err.Error())
	}

	return fmt.Sprintf("validation failed: %s", strings.Join(errorMessages, "; "))
}

// Add adds a validation error to the collection.
func (v *ValidationErrors) Add(field, message, code string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// Addf adds a formatted validation error to the collection.
func (v *ValidationErrors) Addf(field, format string, args ...interface{}) {
	v.Add(field, fmt.Sprintf(format, args...), "")
}

// HasErrors returns true if there are validation errors.
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// Len returns the number of validation errors.
func (v *ValidationErrors) Len() int {
	return len(v.Errors)
}

// HasCode checks if any error in the collection matches the specified code.
func (v *ValidationErrors) HasCode(code string) bool {
	for _, err := range v.Errors {
		if err.Code == code {
			return true
		}
	}
	return false
}

// ValidateFunc is a function that performs validation and returns validation errors.
type ValidateFunc func() *ValidationErrors

// Validate runs all validation functions and returns aggregated errors.
func Validate(validators ...ValidateFunc) *ValidationErrors {
	errors := &ValidationErrors{}

	for _, validator := range validators {
		if err := validator(); err != nil {
			errors.Errors = append(errors.Errors, err.Errors...)
		}
	}

	return errors
}

// Required validates that a string value is not empty.
func Required(value, field, code string) *ValidationErrors {
	errors := &ValidationErrors{}

	if strings.TrimSpace(value) == "" {
		errors.Add(field, "field is required", code)
	}

	return errors
}

// Email validates that a string value is a valid email format.
func Email(value, field, code string) *ValidationErrors {
	errors := &ValidationErrors{}

	if value != "" { // Only validate if not empty
		// Simple email validation
		if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
			errors.Add(field, "must be a valid email address", code)
		}
	}

	return errors
}

// MinLength validates that a string has minimum length.
func MinLength(value, field, code string, min int) *ValidationErrors {
	errors := &ValidationErrors{}

	if len(strings.TrimSpace(value)) < min {
		errors.Addf(field, "must be at least %d characters long", min)
	}

	return errors
}

// MaxLength validates that a string has maximum length.
func MaxLength(value, field, code string, max int) *ValidationErrors {
	errors := &ValidationErrors{}

	if len(strings.TrimSpace(value)) > max {
		errors.Addf(field, "must be at most %d characters long", max)
	}

	return errors
}

// BetweenLength validates that a string length is between min and max.
func BetweenLength(value, field, code string, min, max int) *ValidationErrors {
	errors := &ValidationErrors{}

	length := len(strings.TrimSpace(value))
	if length < min || length > max {
		errors.Addf(field, "must be between %d and %d characters long", min, max)
	}

	return errors
}

// Integer validates that a string value represents a valid integer.
func Integer(value, field, code string) *ValidationErrors {
	errors := &ValidationErrors{}

	if value != "" { // Only validate if not empty
		// Simple integer validation (could be enhanced)
		if _, err := fmt.Sscanf(value, "%d", new(int)); err != nil {
			errors.Add(field, "must be a valid integer", code)
		}
	}

	return errors
}

// Float validates that a string value represents a valid float.
func Float(value, field, code string) *ValidationErrors {
	errors := &ValidationErrors{}

	if value != "" { // Only validate if not empty
		// Simple float validation (could be enhanced)
		if _, err := fmt.Sscanf(value, "%f", new(float64)); err != nil {
			errors.Add(field, "must be a valid number", code)
		}
	}

	return errors
}
