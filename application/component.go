// Package application provides application-level components and runtime management.
package application

import (
	"fmt"
)

// ComponentID uniquely identifies a component within the framework.
type ComponentID struct {
	// Name is the unique name of the component.
	Name string

	// Version is the semantic version of the component.
	Version string
}

// NewComponentID creates a new ComponentID with the given name and version.
func NewComponentID(name, version string) ComponentID {
	return ComponentID{
		Name:    name,
		Version: version,
	}
}

// String returns a string representation of the ComponentID.
func (id ComponentID) String() string {
	return fmt.Sprintf("%s:%s", id.Name, id.Version)
}

// IsZero returns true if the ComponentID is zero-valued.
func (id ComponentID) IsZero() bool {
	return id.Name == "" && id.Version == ""
}
