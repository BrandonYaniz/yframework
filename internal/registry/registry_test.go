// Package registry provides component registration and dependency management.
package registry

import (
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	// Test registering a component
	constructor := func() interface{} { return "test" }
	err := r.Register("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering component, got: %v", err)
	}

	// Test registering duplicate component
	err = r.Register("test-component", constructor)
	if err == nil {
		t.Error("Expected error when registering duplicate component")
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	constructor := func() interface{} { return "test" }
	r.Register("test-component", constructor)

	component, exists := r.Get("test-component")
	if !exists {
		t.Error("Expected component to exist")
	}

	if component.ID != "test-component" {
		t.Errorf("Expected component ID 'test-component', got: %s", component.ID)
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()

	constructor := func() interface{} { return "test" }
	r.Register("component-1", constructor)
	r.Register("component-2", constructor)

	components := r.List()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got: %d", len(components))
	}
}

func TestRegistry_Freeze(t *testing.T) {
	r := NewRegistry()

	// Just test that it doesn't panic
	r.Freeze()

	// Should still be able to register after freezing (for this simple implementation)
	constructor := func() interface{} { return "test" }
	err := r.Register("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error even after freeze, got: %v", err)
	}
}
