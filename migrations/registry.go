// Package migrations provides migration registration and management.
package migrations

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// Migration represents a single migration step.
type Migration struct {
	// Owner is the name of the owner (component, service, etc.) that owns this migration
	Owner string

	// Name is the unique name of this migration within the owner
	Name string

	// Version is the version number of this migration
	Version int

	// Checksum is the SHA256 checksum of the migration content
	Checksum [32]byte

	// AppliedAt is when this migration was applied (set after successful application)
	AppliedAt *time.Time

	// Description provides a human-readable description of what this migration does
	Description string
}

// MigrationRegistry manages registered migrations.
type MigrationRegistry struct {
	mu       sync.RWMutex
	registry map[string]*Migration
}

// NewMigrationRegistry creates a new migration registry.
func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		registry: make(map[string]*Migration),
	}
}

// Register registers a migration with the registry.
// Returns an error if a migration with the same owner and name is already registered.
func (mr *MigrationRegistry) Register(migration *Migration) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	key := fmt.Sprintf("%s:%s", migration.Owner, migration.Name)

	if _, exists := mr.registry[key]; exists {
		return errors.New("migration already registered", "MIGRATION_ALREADY_REGISTERED")
	}

	mr.registry[key] = migration
	return nil
}

// GetMigration retrieves a migration by owner and name.
func (mr *MigrationRegistry) GetMigration(owner, name string) (*Migration, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", owner, name)
	migration, exists := mr.registry[key]
	return migration, exists
}

// ListMigrations returns all registered migrations sorted by owner and version.
func (mr *MigrationRegistry) ListMigrations() []*Migration {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var migrations []*Migration
	for _, migration := range mr.registry {
		migrations = append(migrations, migration)
	}

	// Sort by owner, then by version
	sort.Slice(migrations, func(i, j int) bool {
		if migrations[i].Owner != migrations[j].Owner {
			return migrations[i].Owner < migrations[j].Owner
		}
		return migrations[i].Version < migrations[j].Version
	})

	return migrations
}

// GetMigrationsForOwner returns all migrations for a specific owner.
func (mr *MigrationRegistry) GetMigrationsForOwner(owner string) []*Migration {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var ownerMigrations []*Migration
	for _, migration := range mr.registry {
		if migration.Owner == owner {
			ownerMigrations = append(ownerMigrations, migration)
		}
	}

	// Sort by version
	sort.Slice(ownerMigrations, func(i, j int) bool {
		return ownerMigrations[i].Version < ownerMigrations[j].Version
	})

	return ownerMigrations
}

// ApplyMigration marks a migration as applied.
func (mr *MigrationRegistry) ApplyMigration(owner, name string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	key := fmt.Sprintf("%s:%s", owner, name)
	migration, exists := mr.registry[key]
	if !exists {
		return errors.New("migration not found", "MIGRATION_NOT_FOUND")
	}

	now := time.Now()
	migration.AppliedAt = &now
	return nil
}

// IsMigrationApplied checks if a migration has been applied.
func (mr *MigrationRegistry) IsMigrationApplied(owner, name string) bool {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", owner, name)
	migration, exists := mr.registry[key]
	if !exists {
		return false
	}

	return migration.AppliedAt != nil
}

// CalculateChecksum calculates the SHA256 checksum of migration content.
func CalculateChecksum(content string) [32]byte {
	return sha256.Sum256([]byte(content))
}
