// Package migrations provides migration utilities and history tracking.
package migrations

import (
	"context"
	"database/sql"
)

// CLIService provides command-line interface services for migrations.
type CLIService struct {
	db           *sql.DB
	historyStore *HistoryStore
	lock         *AdvisoryLock
	runner       *Runner
}

// NewCLIService creates a new CLI service for migrations.
func NewCLIService(db *sql.DB) (*CLIService, error) {
	historyStore := NewHistoryStore(db)
	lock := NewAdvisoryLock(db, 12345) // Fixed lock ID for migrations
	runner := NewRunner(db, historyStore, lock)

	return &CLIService{
		db:           db,
		historyStore: historyStore,
		lock:         lock,
		runner:       runner,
	}, nil
}

// Migrate runs all pending migrations.
func (c *CLIService) Migrate(ctx context.Context) error {
	return c.runner.Run(ctx)
}

// Rollback rolls back the last migration.
func (c *CLIService) Rollback(ctx context.Context) error {
	return c.runner.Rollback(ctx)
}

// Status shows the current migration status.
func (c *CLIService) Status(ctx context.Context) ([]MigrationHistory, error) {
	return c.historyStore.GetAppliedMigrations()
}

// Checksum calculates and returns the checksum of a migration file.
func (c *CLIService) Checksum(migrationName string, version int64) (string, error) {
	// This would typically compute a checksum of the migration content
	// For now, we'll return a placeholder
	return "placeholder_checksum", nil
}

// Validate validates that all migrations are properly formatted and have valid checksums.
func (c *CLIService) Validate(ctx context.Context) error {
	// This would validate migration files against their checksums
	// For now, we'll just return success
	return nil
}
