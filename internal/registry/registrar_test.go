// Package registry provides component registration and dependency management.
package registry

import (
	"testing"
)

func TestComponentRegistrar_RegisterComponent(t *testing.T) {
	registrar := NewComponentRegistrar()
	
	constructor := func() interface{} { return "test-component" }
	err := registrar.RegisterComponent("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering component, got: %v", err)
	}
	
	// Test duplicate registration
	err = registrar.RegisterComponent("test-component", constructor)
	if err == nil {
		t.Error("Expected error when registering duplicate component")
	}
}

func TestComponentRegistrar_RegisterService(t *testing.T) {
	registrar := NewComponentRegistrar()
	
	constructor := func() interface{} { return "test-service" }
	err := registrar.RegisterService("test-service", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering service, got: %v", err)
	}
}

func TestComponentRegistrar_FreezeRegistry(t *testing.T) {
	registrar := NewComponentRegistrar()
	
	constructor := func() interface{} { return "test" }
	err := registrar.RegisterComponent("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering before freeze, got: %v", err)
	}
	
	// Freeze the registry
	registrar.FreezeRegistry()
	
	// Verify it's frozen
	if !registrar.IsRegistryFrozen() {
		t.Error("Expected registry to be frozen")
	}
	
	// Test that registration fails after freezing
	err = registrar.RegisterComponent("another-component", constructor)
	if err == nil {
		t.Error("Expected error when registering after freeze")
	}
}

func TestComponentRegistrar_GetRegistry(t *testing.T) {
	registrar := NewComponentRegistrar()
	
	registry := registrar.GetRegistry()
	if registry == nil {
		t.Error("Expected non-nil registry")
	}
	
	// Verify it's the same registry
	constructor := func() interface{} { return "test" }
	err := registrar.RegisterComponent("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering, got: %v", err)
	}
	
	components := registry.List()
	if len(components) != 1 {
		t.Errorf("Expected 1 component in registry, got: %d", len(components))
	}
}