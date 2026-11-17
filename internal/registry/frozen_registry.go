// Package registry provides component registration and dependency management.
package registry

import (
	"fmt"
	"sync"
)

// Component represents a registered component with its dependencies.
type Component struct {
	ID           string
	Constructor  func() interface{}
	Dependencies []string
}

// FrozenRegistry is a generic private frozen registry that prevents further mutations after freezing.
type FrozenRegistry struct {
	mu         sync.RWMutex
	components map[string]*Component
	frozen     bool
}

// NewFrozenRegistry creates a new empty frozen registry.
func NewFrozenRegistry() *FrozenRegistry {
	return &FrozenRegistry{
		components: make(map[string]*Component),
		frozen:     false,
	}
}

// Register registers a component with the given ID and constructor.
// Returns an error if the registry is frozen or if a component with the same ID already exists.
func (r *FrozenRegistry) Register(id string, constructor func() interface{}, dependencies ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("cannot register component in frozen registry")
	}

	if _, exists := r.components[id]; exists {
		return fmt.Errorf("component %s already registered", id)
	}

	r.components[id] = &Component{
		ID:           id,
		Constructor:  constructor,
		Dependencies: dependencies,
	}

	return nil
}

// Get retrieves a component by its ID.
// Returns false if the component does not exist.
func (r *FrozenRegistry) Get(id string) (*Component, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	component, exists := r.components[id]
	return component, exists
}

// List returns all registered components.
func (r *FrozenRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.components))
	for id := range r.components {
		ids = append(ids, id)
	}
	return ids
}

// Freeze prevents further component registrations and mutations.
// This method can only be called once.
func (r *FrozenRegistry) Freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.frozen {
		r.frozen = true
	}
}

// IsFrozen returns true if the registry is frozen.
func (r *FrozenRegistry) IsFrozen() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.frozen
}
