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
