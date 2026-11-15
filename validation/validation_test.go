package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "username",
		Message: "must not be empty",
		Code:    "USERNAME_REQUIRED",
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

func TestValidationErrors(t *testing.T) {
	errors := &ValidationErrors{}

	// Test empty validation errors
	if errors.Error() != "" {
		t.Errorf("Expected empty error string, got '%s'", errors.Error())
	}

	if errors.HasErrors() {
		t.Error("Expected no errors")
	}

	if errors.Len() != 0 {
		t.Errorf("Expected 0 errors, got %d", errors.Len())
	}

	// Add some errors
	errors.Add("username", "must not be empty", "USERNAME_REQUIRED")
	errors.Add("email", "invalid format", "EMAIL_INVALID")

	if !errors.HasErrors() {
		t.Error("Expected errors to be present")
	}

	if errors.Len() != 2 {
		t.Errorf("Expected 2 errors, got %d", errors.Len())
	}

	// Test error string
	errorString := errors.Error()
	if errorString == "" {
		t.Error("Expected non-empty error string")
	}

	if !strings.Contains(errorString, "username") || !strings.Contains(errorString, "email") {
		t.Errorf("Error string should contain field names: %s", errorString)
	}
}

func TestValidationErrorsIs(t *testing.T) {
	errors := &ValidationErrors{}
	errors.Add("username", "must not be empty", "USERNAME_REQUIRED")
	errors.Add("email", "invalid format", "EMAIL_INVALID")

	if !errors.HasCode("USERNAME_REQUIRED") {
		t.Error("Expected to find USERNAME_REQUIRED code")
	}

	if errors.HasCode("NON_EXISTENT_CODE") {
		t.Error("Expected to not find NON_EXISTENT_CODE code")
	}
}

func TestValidate(t *testing.T) {
	// Test with no validators
	result := Validate()
	if result.HasErrors() {
		t.Error("Expected no errors with no validators")
	}

	// Test with one validator that passes
	validator1 := func() *ValidationErrors {
		return &ValidationErrors{}
	}

	result = Validate(validator1)
	if result.HasErrors() {
		t.Error("Expected no errors with passing validator")
	}

	// Test with one validator that fails
	validator2 := func() *ValidationErrors {
		errors := &ValidationErrors{}
		errors.Add("username", "must not be empty", "USERNAME_REQUIRED")
		return errors
	}

	result = Validate(validator2)
	if !result.HasErrors() {
		t.Error("Expected errors with failing validator")
	}

	if result.Len() != 1 {
		t.Errorf("Expected 1 error, got %d", result.Len())
	}
}

func TestRequired(t *testing.T) {
	// Test valid value
	errors := Required("valid", "username", "USERNAME_REQUIRED")
	if errors.HasErrors() {
		t.Error("Expected no errors for valid value")
	}

	// Test empty value
	errors = Required("", "username", "USERNAME_REQUIRED")
	if !errors.HasErrors() {
		t.Error("Expected error for empty value")
	}

	if errors.Len() != 1 {
		t.Errorf("Expected 1 error, got %d", errors.Len())
	}

	// Test whitespace-only value
	errors = Required("   ", "username", "USERNAME_REQUIRED")
	if !errors.HasErrors() {
		t.Error("Expected error for whitespace-only value")
	}
}

func TestEmail(t *testing.T) {
	// Test valid email
	errors := Email("test@example.com", "email", "EMAIL_INVALID")
	if errors.HasErrors() {
		t.Error("Expected no errors for valid email")
	}

	// Test invalid email
	errors = Email("invalid-email", "email", "EMAIL_INVALID")
	if !errors.HasErrors() {
		t.Error("Expected error for invalid email")
	}

	// Test empty email (should not fail validation)
	errors = Email("", "email", "EMAIL_INVALID")
	if errors.HasErrors() {
		t.Error("Expected no errors for empty email")
	}
}

func TestMinLength(t *testing.T) {
	// Test valid length
	errors := MinLength("valid", "username", "USERNAME_TOO_SHORT", 3)
	if errors.HasErrors() {
		t.Error("Expected no errors for valid length")
	}

	// Test too short
	errors = MinLength("ab", "username", "USERNAME_TOO_SHORT", 3)
	if !errors.HasErrors() {
		t.Error("Expected error for too short value")
	}
}

func TestMaxLength(t *testing.T) {
	// Test valid length
	errors := MaxLength("valid", "username", "USERNAME_TOO_LONG", 10)
	if errors.HasErrors() {
		t.Error("Expected no errors for valid length")
	}

	// Test too long
	errors = MaxLength("this-is-too-long", "username", "USERNAME_TOO_LONG", 5)
	if !errors.HasErrors() {
		t.Error("Expected error for too long value")
	}
}

func TestBetweenLength(t *testing.T) {
	// Test valid length
	errors := BetweenLength("valid", "username", "USERNAME_INVALID_LENGTH", 3, 10)
	if errors.HasErrors() {
		t.Error("Expected no errors for valid length")
	}

	// Test too short
	errors = BetweenLength("ab", "username", "USERNAME_INVALID_LENGTH", 3, 10)
	if !errors.HasErrors() {
		t.Error("Expected error for too short value")
	}

	// Test too long
	errors = BetweenLength("this-is-too-long", "username", "USERNAME_INVALID_LENGTH", 3, 10)
	if !errors.HasErrors() {
		t.Error("Expected error for too long value")
	}
}

func TestInteger(t *testing.T) {
	// Test valid integer
	errors := Integer("123", "age", "AGE_NOT_INTEGER")
	if errors.HasErrors() {
		t.Error("Expected no errors for valid integer")
	}

	// Test invalid integer
	errors = Integer("not-an-integer", "age", "AGE_NOT_INTEGER")
	if !errors.HasErrors() {
		t.Error("Expected error for invalid integer")
	}

	// Test empty value (should not fail validation)
	errors = Integer("", "age", "AGE_NOT_INTEGER")
	if errors.HasErrors() {
		t.Error("Expected no errors for empty value")
	}
}

func TestFloat(t *testing.T) {
	// Test valid float
	errors := Float("123.45", "price", "PRICE_NOT_FLOAT")
	if errors.HasErrors() {
		t.Error("Expected no errors for valid float")
	}

	// Test invalid float
	errors = Float("not-a-float", "price", "PRICE_NOT_FLOAT")
	if !errors.HasErrors() {
		t.Error("Expected error for invalid float")
	}

	// Test empty value (should not fail validation)
	errors = Float("", "price", "PRICE_NOT_FLOAT")
	if errors.HasErrors() {
		t.Error("Expected no errors for empty value")
	}
}

func TestValidationAggregation(t *testing.T) {
	// Test combining multiple validators
	validateUsername := func() *ValidationErrors {
		return Validate(
			func() *ValidationErrors { return Required("test", "username", "USERNAME_REQUIRED") },
			func() *ValidationErrors { return MinLength("test", "username", "USERNAME_TOO_SHORT", 3) },
		)
	}

	result := validateUsername()
	if result.HasErrors() {
		t.Error("Expected no errors for valid username")
	}

	validateInvalidUsername := func() *ValidationErrors {
		return Validate(
			func() *ValidationErrors { return Required("", "username", "USERNAME_REQUIRED") },
			func() *ValidationErrors { return MinLength("ab", "username", "USERNAME_TOO_SHORT", 3) },
		)
	}

	result = validateInvalidUsername()
	if !result.HasErrors() {
		t.Error("Expected errors for invalid username")
	}

	if result.Len() != 2 {
		t.Errorf("Expected 2 errors, got %d", result.Len())
	}
}
