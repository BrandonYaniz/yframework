// Package testapp provides a test application for demonstrating component usage.
package testapp

import (
	"testing"
)

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	if app == nil {
		t.Error("Expected non-nil TestApp")
	}

	registry := app.GetRegistry()
	if registry == nil {
		t.Error("Expected non-nil registry")
	}
}

func TestTestApp_Run(t *testing.T) {
	app := NewTestApp()
	err := app.Run()
	if err != nil {
		t.Errorf("Expected no error when running test app, got: %v", err)
	}
}
