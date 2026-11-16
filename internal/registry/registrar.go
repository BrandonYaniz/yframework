// Package registry provides component registration and dependency management.
package registry

import (
	"fmt"
)

// Registrar is the interface for registering components with the framework.
type Registrar interface {
	// RegisterComponent registers a component with its dependencies.
	RegisterComponent(id string, constructor func() interface{}, dependencies ...string) error
	
	// RegisterService registers a service with its dependencies.
	RegisterService(id string, constructor func() interface{}, dependencies ...string) error
	
	// FreezeRegistry prevents further registrations.
	FreezeRegistry()
	
	// IsRegistryFrozen returns true if the registry is frozen.
	IsRegistryFrozen() bool
}

// ComponentRegistrar is a concrete implementation of Registrar for component registration.
type ComponentRegistrar struct {
	registry *FrozenRegistry
}

// NewComponentRegistrar creates a new component registrar.
func NewComponentRegistrar() *ComponentRegistrar {
	return &ComponentRegistrar{
		registry: NewFrozenRegistry(),
	}
}

// RegisterComponent registers a component with the given ID and constructor.
func (r *ComponentRegistrar) RegisterComponent(id string, constructor func() interface{}, dependencies ...string) error {
	if r.registry.IsFrozen() {
		return fmt.Errorf("cannot register component in frozen registry")
	}
	
	return r.registry.Register(id, constructor, dependencies...)
}

// RegisterService registers a service with the given ID and constructor.
func (r *ComponentRegistrar) RegisterService(id string, constructor func() interface{}, dependencies ...string) error {
	if r.registry.IsFrozen() {
		return fmt.Errorf("cannot register service in frozen registry")
	}
	
	return r.registry.Register(id, constructor, dependencies...)
}

// FreezeRegistry prevents further registrations.
func (r *ComponentRegistrar) FreezeRegistry() {
	r.registry.Freeze()
}

// IsRegistryFrozen returns true if the registry is frozen.
func (r *ComponentRegistrar) IsRegistryFrozen() bool {
	return r.registry.IsFrozen()
}

// GetRegistry returns the underlying frozen registry.
func (r *ComponentRegistrar) GetRegistry() *FrozenRegistry {
	return r.registry
}