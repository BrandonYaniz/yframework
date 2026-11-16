// Package registry provides component registration and dependency management.
package registry

import (
	"fmt"
	"sort"
)

// DependencyGraph represents a directed graph of component dependencies.
type DependencyGraph struct {
	components map[string]*Component
}

// NewDependencyGraph creates a new empty dependency graph.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		components: make(map[string]*Component),
	}
}

// AddComponent adds a component to the dependency graph.
func (dg *DependencyGraph) AddComponent(component *Component) error {
	if _, exists := dg.components[component.ID]; exists {
		return fmt.Errorf("component %s already exists in graph", component.ID)
	}
	
	dg.components[component.ID] = component
	return nil
}

// TopologicalSort returns a topologically sorted list of component IDs.
// This provides a deterministic order for component initialization.
func (dg *DependencyGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	result := []string{}
	
	// Helper function for DFS
	var dfs func(string) error
	dfs = func(id string) error {
		if recStack[id] {
			return fmt.Errorf("cycle detected in dependency graph involving component %s", id)
		}
		
		if visited[id] {
			return nil // Already processed
		}
		
		recStack[id] = true
		visited[id] = true
		
		component, exists := dg.components[id]
		if !exists {
			return fmt.Errorf("component %s not found in graph", id)
		}
		
		// Process dependencies first
		for _, depID := range component.Dependencies {
			if err := dfs(depID); err != nil {
				return err
			}
		}
		
		recStack[id] = false
		result = append(result, id)
		return nil
	}
	
	// Process all components
	for id := range dg.components {
		if !visited[id] {
			if err := dfs(id); err != nil {
				return nil, err
			}
		}
	}
	
	// Reverse to get topological order (dependencies first)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	
	return result, nil
}

// ValidateDependencies checks that all component dependencies exist.
func (dg *DependencyGraph) ValidateDependencies() error {
	for id, component := range dg.components {
		for _, depID := range component.Dependencies {
			if _, exists := dg.components[depID]; !exists {
				return fmt.Errorf("component %s depends on non-existent component %s", id, depID)
			}
		}
	}
	return nil
}

// GetComponent returns a component by ID.
func (dg *DependencyGraph) GetComponent(id string) (*Component, bool) {
	component, exists := dg.components[id]
	return component, exists
}

// ListComponents returns all component IDs in the graph.
func (dg *DependencyGraph) ListComponents() []string {
	ids := make([]string, 0, len(dg.components))
	for id := range dg.components {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}