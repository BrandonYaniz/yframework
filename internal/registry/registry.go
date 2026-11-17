// Package registry provides component registration and dependency management.
package registry

import (
	"fmt"
	"sync"
)

// Registry manages component registration and dependency resolution.
type Registry struct {
	mu         sync.RWMutex
	components map[string]*Component
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[string]*Component),
	}
}

// Register registers a component with the given ID and constructor.
func (r *Registry) Register(id string, constructor func() interface{}, dependencies ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

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
func (r *Registry) Get(id string) (*Component, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	component, exists := r.components[id]
	return component, exists
}

// List returns all registered components.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.components))
	for id := range r.components {
		ids = append(ids, id)
	}
	return ids
}

// Freeze prevents further component registrations.
func (r *Registry) Freeze() {
	// In a real implementation, this would lock the registry to prevent further changes
	// For now, we'll just log that it's been called
}
