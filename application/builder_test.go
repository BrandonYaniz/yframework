// Package application provides application-level components and runtime management.
package application

import (
	"testing"
)

func TestApplicationBuilder_NewApplicationBuilder(t *testing.T) {
	builder := NewApplicationBuilder()
	if builder == nil {
		t.Error("Expected non-nil application builder")
	}
}

func TestApplicationBuilder_RegisterComponent(t *testing.T) {
	builder := NewApplicationBuilder()

	constructor := func() interface{} { return "test-component" }
	err := builder.RegisterComponent("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering component, got: %v", err)
	}

	// Test duplicate registration
	err = builder.RegisterComponent("test-component", constructor)
	if err == nil {
		t.Error("Expected error when registering duplicate component")
	}
}

func TestApplicationBuilder_RegisterService(t *testing.T) {
	builder := NewApplicationBuilder()

	constructor := func() interface{} { return "test-service" }
	err := builder.RegisterService("test-service", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering service, got: %v", err)
	}
}

func TestApplicationBuilder_Build(t *testing.T) {
	builder := NewApplicationBuilder()

	constructor := func() interface{} { return "test-component" }
	builder.RegisterComponent("test-component", constructor)

	app, err := builder.Build()
	if err != nil {
		t.Errorf("Expected no error when building application, got: %v", err)
	}

	if app == nil {
		t.Error("Expected non-nil application")
	}
}

func TestApplicationBuilder_ValidateComposition(t *testing.T) {
	builder := NewApplicationBuilder()

	// Should not error on empty composition
	err := builder.ValidateComposition()
	if err != nil {
		t.Errorf("Expected no error when validating empty composition, got: %v", err)
	}

	// Register a component and validate
	constructor := func() interface{} { return "test-component" }
	builder.RegisterComponent("test-component", constructor)

	err = builder.ValidateComposition()
	if err != nil {
		t.Errorf("Expected no error when validating composition with components, got: %v", err)
	}
}
