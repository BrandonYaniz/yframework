package validation

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "username",
		Message: "must not be empty",
	}

	expected := "validation error on field 'username': must not be empty"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestIsValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid format",
	}

	if !IsValidationError(err) {
		t.Errorf("Expected validation error to be detected")
	}

	regularErr := errors.New("regular error")
	if IsValidationError(regularErr) {
		t.Errorf("Expected regular error not to be detected as validation error")
	}
}
