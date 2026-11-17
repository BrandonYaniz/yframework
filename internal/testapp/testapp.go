// Package testapp provides a test application for demonstrating component usage.
package testapp

import (
	"github.com/BrandonYaniz/yframework/internal/registry"
)

// TestApp represents a simple test application structure.
type TestApp struct {
	registry *registry.Registry
}

// NewTestApp creates a new test application with an empty registry.
func NewTestApp() *TestApp {
	return &TestApp{
		registry: registry.NewRegistry(),
	}
}

// GetRegistry returns the application's component registry.
func (app *TestApp) GetRegistry() *registry.Registry {
	return app.registry
}

// Run executes the test application logic.
func (app *TestApp) Run() error {
	// In a real implementation, this would initialize components and run the application
	return nil
}
