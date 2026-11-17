// Package registry provides component registration and dependency management.
package registry

import (
	"testing"
)

func TestDependencyGraph_AddComponent(t *testing.T) {
	dg := NewDependencyGraph()

	component := &Component{
		ID: "test-component",
	}

	err := dg.AddComponent(component)
	if err != nil {
		t.Errorf("Expected no error when adding component, got: %v", err)
	}

	// Test adding duplicate
	err = dg.AddComponent(component)
	if err == nil {
		t.Error("Expected error when adding duplicate component")
	}
}

func TestDependencyGraph_TopologicalSort(t *testing.T) {
	dg := NewDependencyGraph()

	// Add components with dependencies
	componentA := &Component{
		ID: "component-a",
	}
	componentB := &Component{
		ID:           "component-b",
		Dependencies: []string{"component-a"},
	}
	componentC := &Component{
		ID:           "component-c",
		Dependencies: []string{"component-b"},
	}

	dg.AddComponent(componentA)
	dg.AddComponent(componentB)
	dg.AddComponent(componentC)

	sorted, err := dg.TopologicalSort()
	if err != nil {
		t.Errorf("Expected no error when sorting, got: %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("Expected 3 components in sorted list, got: %d", len(sorted))
	}

	// Check that dependencies come before dependents
	aIndex := -1
	bIndex := -1
	cIndex := -1

	for i, id := range sorted {
		switch id {
		case "component-a":
			aIndex = i
		case "component-b":
			bIndex = i
		case "component-c":
			cIndex = i
		}
	}

	if aIndex > bIndex {
		t.Error("Component A should come before component B")
	}

	if bIndex > cIndex {
		t.Error("Component B should come before component C")
	}
}

func TestDependencyGraph_ValidateDependencies(t *testing.T) {
	dg := NewDependencyGraph()

	// Add component with missing dependency
	componentA := &Component{
		ID:           "component-a",
		Dependencies: []string{"non-existent-component"},
	}

	dg.AddComponent(componentA)

	err := dg.ValidateDependencies()
	if err == nil {
		t.Error("Expected error when validating dependencies with missing component")
	}
}

func TestDependencyGraph_CycleDetection(t *testing.T) {
	dg := NewDependencyGraph()

	// Create a cycle: A -> B -> A
	componentA := &Component{
		ID:           "component-a",
		Dependencies: []string{"component-b"},
	}
	componentB := &Component{
		ID:           "component-b",
		Dependencies: []string{"component-a"},
	}

	dg.AddComponent(componentA)
	dg.AddComponent(componentB)

	_, err := dg.TopologicalSort()
	if err == nil {
		t.Error("Expected error when sorting with cycle")
	}
}

func TestDependencyGraph_GetComponent(t *testing.T) {
	dg := NewDependencyGraph()

	component := &Component{
		ID: "test-component",
	}
	dg.AddComponent(component)

	found, exists := dg.GetComponent("test-component")
	if !exists {
		t.Error("Expected component to exist")
	}

	if found.ID != "test-component" {
		t.Errorf("Expected component ID 'test-component', got: %s", found.ID)
	}
}

func TestDependencyGraph_ListComponents(t *testing.T) {
	dg := NewDependencyGraph()

	componentA := &Component{
		ID: "component-a",
	}
	componentB := &Component{
		ID: "component-b",
	}

	dg.AddComponent(componentA)
	dg.AddComponent(componentB)

	components := dg.ListComponents()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got: %d", len(components))
	}

	// Should be sorted
	if components[0] != "component-a" || components[1] != "component-b" {
		t.Error("Components should be returned in sorted order")
	}
}
