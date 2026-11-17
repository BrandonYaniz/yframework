// Package application provides application-level components and runtime management.
package application

import (
	"github.com/BrandonYaniz/yframework/internal/registry"
)

// ApplicationBuilder is responsible for building applications with registered components.
type ApplicationBuilder struct {
	registrar *registry.ComponentRegistrar
}

// NewApplicationBuilder creates a new application builder.
func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{
		registrar: registry.NewComponentRegistrar(),
	}
}

// RegisterComponent registers a component with the application builder.
func (ab *ApplicationBuilder) RegisterComponent(id string, constructor func() interface{}, dependencies ...string) error {
	return ab.registrar.RegisterComponent(id, constructor, dependencies...)
}

// RegisterService registers a service with the application builder.
func (ab *ApplicationBuilder) RegisterService(id string, constructor func() interface{}, dependencies ...string) error {
	return ab.registrar.RegisterService(id, constructor, dependencies...)
}

// Build validates and builds the application.
// This method validates the composition before any database mutations.
func (ab *ApplicationBuilder) Build() (*Application, error) {
	// Validate that all dependencies can be resolved
	// In a real implementation, this would check for cycles and missing dependencies
	
	// Freeze the registry to prevent further modifications
	ab.registrar.FreezeRegistry()
	
	return &Application{
		registrar: ab.registrar,
	}, nil
}

// Application represents a built application with its registered components.
type Application struct {
	registrar *registry.ComponentRegistrar
}

// GetComponent retrieves a component by ID.
func (app *Application) GetComponent(id string) (interface{}, bool) {
	// In a real implementation, this would return the actual component instance
	// For now, we'll just return nil to indicate it's not implemented
	return nil, false
}

// Run executes the application.
func (app *Application) Run() error {
	// In a real implementation, this would initialize and run the application
	// This is intentionally left as a stub for now
	return nil
}

// ValidateComposition performs validation of the application composition.
func (ab *ApplicationBuilder) ValidateComposition() error {
	// In a real implementation, this would validate:
	// - No cycles in dependency graph
	// - All dependencies exist
	// - No duplicate registrations
	// - Valid component configurations
	
	return nil
}