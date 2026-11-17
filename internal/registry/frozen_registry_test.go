// Package registry provides component registration and dependency management.
package registry

import (
	"sync"
	"testing"
)

func TestFrozenRegistry_Register(t *testing.T) {
	r := NewFrozenRegistry()

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

func TestFrozenRegistry_Get(t *testing.T) {
	r := NewFrozenRegistry()

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

func TestFrozenRegistry_List(t *testing.T) {
	r := NewFrozenRegistry()

	constructor := func() interface{} { return "test" }
	r.Register("component-1", constructor)
	r.Register("component-2", constructor)

	components := r.List()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got: %d", len(components))
	}
}

func TestFrozenRegistry_Freeze(t *testing.T) {
	r := NewFrozenRegistry()

	// Should be able to register before freezing
	constructor := func() interface{} { return "test" }
	err := r.Register("test-component", constructor)
	if err != nil {
		t.Errorf("Expected no error when registering before freeze, got: %v", err)
	}

	// Test that we can freeze
	r.Freeze()

	// Verify registry is frozen
	if !r.IsFrozen() {
		t.Error("Expected registry to be frozen")
	}

	// Should not be able to register after freezing
	err = r.Register("another-component", constructor)
	if err == nil {
		t.Error("Expected error when registering after freeze")
	}
}

func TestFrozenRegistry_ConcurrentAccess(t *testing.T) {
	r := NewFrozenRegistry()

	constructor := func() interface{} { return "test" }
	r.Register("test-component", constructor)

	// Freeze the registry
	r.Freeze()

	// Test concurrent reads after freezing
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, exists := r.Get("test-component")
			if !exists {
				t.Error("Component should exist")
			}
		}()
	}

	wg.Wait()
}
