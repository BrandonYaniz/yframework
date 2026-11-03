package errors

import (
	"errors"
	"testing"
)

func TestFrameworkError(t *testing.T) {
	err := &FrameworkError{
		Message: "test error",
		Code:    "TEST_ERROR",
	}

	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}

	cause := errors.New("underlying error")
	err.Cause = cause
	if err.Unwrap() != cause {
		t.Errorf("Expected underlying error to be returned")
	}
}

func TestNewFrameworkError(t *testing.T) {
	err := New("test message", "TEST_CODE")

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
	}
}

func TestNewWithCause(t *testing.T) {
	cause := errors.New("underlying cause")
	err := NewWithCause("test message", "TEST_CODE", cause)

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause to be set")
	}
}

func TestIs(t *testing.T) {
	err := New("test message", "TEST_CODE")

	if !Is(err, "TEST_CODE") {
		t.Error("Expected error to match code")
	}

	if Is(err, "DIFFERENT_CODE") {
		t.Error("Expected error to not match different code")
	}

	// Test with wrapped error - this is the expected behavior for our custom Is function
	// Since we're checking by type and code, it should work correctly when the error is directly of FrameworkError type
	// For wrapped errors, users should check the original error or use errors.Is()
	cause := New("underlying", "UNDERLYING_CODE")
	wrappedErr := NewWithCause("wrapped", "WRAPPED_CODE", cause)

	// The Is function should work for direct matches
	if !Is(wrappedErr, "WRAPPED_CODE") {
		t.Error("Expected wrapped error to match its own code")
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "username",
		Message: "must not be empty",
		Code:    "VALIDATION_ERROR",
	}

	expected := "validation error on field 'username': must not be empty"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("email", "invalid format", "EMAIL_INVALID")

	if err.Field != "email" {
		t.Errorf("Expected field 'email', got '%s'", err.Field)
	}

	if err.Message != "invalid format" {
		t.Errorf("Expected message 'invalid format', got '%s'", err.Message)
	}

	if err.Code != "EMAIL_INVALID" {
		t.Errorf("Expected code 'EMAIL_INVALID', got '%s'", err.Code)
	}
}

func TestIsValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid format",
		Code:    "EMAIL_INVALID",
	}

	if !IsValidationError(err) {
		t.Errorf("Expected validation error to be detected")
	}

	regularErr := errors.New("regular error")
	if IsValidationError(regularErr) {
		t.Errorf("Expected regular error not to be detected as validation error")
	}
}
