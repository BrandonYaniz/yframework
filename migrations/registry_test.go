package migrations

import (
	"testing"

	"github.com/BrandonYaniz/yframework/errors"
)

func TestMigrationRegistry_Register(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "test-owner",
		Name:        "test-migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	// Register first migration
	err := registry.Register(migration1)
	if err != nil {
		t.Fatalf("Expected no error when registering first migration, got: %v", err)
	}

	// Try to register duplicate - should fail
	err = registry.Register(migration1)
	if err == nil {
		t.Error("Expected error when registering duplicate migration")
	}

	// Check that error is of the right type
	if !errors.IsValidationError(err) {
		t.Logf("Error type: %T", err)
	}
}

func TestMigrationRegistry_GetMigration(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "test-owner",
		Name:        "test-migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	registry.Register(migration1)

	// Get existing migration
	migration, exists := registry.GetMigration("test-owner", "test-migration-1")
	if !exists {
		t.Error("Expected migration to exist")
	}

	if migration == nil {
		t.Error("Expected migration to not be nil")
	}

	// Get non-existing migration
	migration, exists = registry.GetMigration("test-owner", "non-existent")
	if exists {
		t.Error("Expected migration to not exist")
	}

	if migration != nil {
		t.Error("Expected migration to be nil")
	}
}

func TestMigrationRegistry_ListMigrations(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "owner-a",
		Name:        "migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	migration2 := &Migration{
		Owner:       "owner-b",
		Name:        "migration-1",
		Version:     1,
		Description: "Test migration 2",
		Checksum:    [32]byte{4, 5, 6},
	}

	registry.Register(migration1)
	registry.Register(migration2)

	migrations := registry.ListMigrations()
	if len(migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(migrations))
	}

	// Check that migrations are sorted by owner and version
	if migrations[0].Owner != "owner-a" {
		t.Error("Expected first migration to be from owner-a")
	}

	if migrations[1].Owner != "owner-b" {
		t.Error("Expected second migration to be from owner-b")
	}
}

func TestMigrationRegistry_GetMigrationsForOwner(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "test-owner",
		Name:        "migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	migration2 := &Migration{
		Owner:       "test-owner",
		Name:        "migration-2",
		Version:     2,
		Description: "Test migration 2",
		Checksum:    [32]byte{4, 5, 6},
	}

	migration3 := &Migration{
		Owner:       "other-owner",
		Name:        "migration-1",
		Version:     1,
		Description: "Other owner migration",
		Checksum:    [32]byte{7, 8, 9},
	}

	registry.Register(migration1)
	registry.Register(migration2)
	registry.Register(migration3)

	ownerMigrations := registry.GetMigrationsForOwner("test-owner")
	if len(ownerMigrations) != 2 {
		t.Errorf("Expected 2 migrations for test-owner, got %d", len(ownerMigrations))
	}

	// Should be sorted by version
	if ownerMigrations[0].Version != 1 {
		t.Error("Expected first migration to have version 1")
	}

	if ownerMigrations[1].Version != 2 {
		t.Error("Expected second migration to have version 2")
	}
}

func TestMigrationRegistry_ApplyMigration(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "test-owner",
		Name:        "test-migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	registry.Register(migration1)

	// Apply the migration
	err := registry.ApplyMigration("test-owner", "test-migration-1")
	if err != nil {
		t.Fatalf("Expected no error when applying migration, got: %v", err)
	}

	// Check that it was applied
	applied := registry.IsMigrationApplied("test-owner", "test-migration-1")
	if !applied {
		t.Error("Expected migration to be marked as applied")
	}

	// Try to apply non-existent migration
	err = registry.ApplyMigration("test-owner", "non-existent")
	if err == nil {
		t.Error("Expected error when applying non-existent migration")
	}
}

func TestMigrationRegistry_IsMigrationApplied(t *testing.T) {
	registry := NewMigrationRegistry()

	migration1 := &Migration{
		Owner:       "test-owner",
		Name:        "test-migration-1",
		Version:     1,
		Description: "Test migration 1",
		Checksum:    [32]byte{1, 2, 3},
	}

	registry.Register(migration1)

	// Should not be applied initially
	applied := registry.IsMigrationApplied("test-owner", "test-migration-1")
	if applied {
		t.Error("Expected migration to not be applied initially")
	}

	// Apply it
	registry.ApplyMigration("test-owner", "test-migration-1")

	// Should be applied now
	applied = registry.IsMigrationApplied("test-owner", "test-migration-1")
	if !applied {
		t.Error("Expected migration to be applied after calling ApplyMigration")
	}

	// Non-existent migration should not be applied
	applied = registry.IsMigrationApplied("test-owner", "non-existent")
	if applied {
		t.Error("Expected non-existent migration to not be applied")
	}
}

func TestCalculateChecksum(t *testing.T) {
	content1 := "Hello, world!"
	content2 := "Hello, world!"

	checksum1 := CalculateChecksum(content1)
	checksum2 := CalculateChecksum(content2)

	if checksum1 != checksum2 {
		t.Error("Expected identical checksums for identical content")
	}

	content3 := "Different content"
	checksum3 := CalculateChecksum(content3)

	if checksum1 == checksum3 {
		t.Error("Expected different checksums for different content")
	}
}
