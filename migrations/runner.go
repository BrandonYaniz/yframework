// Package migrations provides migration utilities and history tracking.
package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// Runner executes migrations against the database.
type Runner struct {
	db           *sql.DB
	historyStore *HistoryStore
	lock         *AdvisoryLock
	migrations   []*Migration // This references Migration from registry.go
	timeout      time.Duration
	maxRetries   int
	retryDelay   time.Duration
}

// NewRunner creates a new migration runner.
func NewRunner(db *sql.DB, historyStore *HistoryStore, lock *AdvisoryLock) *Runner {
	return &Runner{
		db:           db,
		historyStore: historyStore,
		lock:         lock,
		migrations:   make([]*Migration, 0),
		timeout:      30 * time.Second,
		maxRetries:   3,
		retryDelay:   1 * time.Second,
	}
}

// AddMigration adds a migration to the runner.
func (r *Runner) AddMigration(migration *Migration) {
	r.migrations = append(r.migrations, migration)
}

// Run applies all pending migrations.
func (r *Runner) Run(ctx context.Context) error {
	// Acquire lock
	acquired, err := r.lock.Acquire(ctx, r.timeout)
	if err != nil {
		return errors.NewWithCause("failed to acquire migration lock", "MIGRATION_LOCK_ACQUIRE_FAILED", err)
	}
	if !acquired {
		return errors.New("migration lock already held by another process", "MIGRATION_LOCK_HELD")
	}
	defer func() {
		// Release lock on exit
		lockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r.lock.Release(lockCtx)
	}()

	// Get applied migrations
	appliedMigrations, err := r.historyStore.GetAppliedMigrations()
	if err != nil {
		return errors.NewWithCause("failed to get applied migrations", "MIGRATION_HISTORY_GET_FAILED", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMap[fmt.Sprintf("%s:%d", m.Name, m.Version)] = true
	}

	// Sort pending migrations by version
	sort.Slice(r.migrations, func(i, j int) bool {
		return r.migrations[i].Version < r.migrations[j].Version
	})

	// Run pending migrations
	for _, migration := range r.migrations {
		key := fmt.Sprintf("%s:%d", migration.Name, migration.Version)
		if _, exists := appliedMap[key]; exists {
			continue // Migration already applied
		}

		// Execute the actual DDL operations first
		err = r.executeMigration(ctx, migration)
		if err != nil {
			return errors.NewWithCause("failed to execute migration", "MIGRATION_EXECUTE_FAILED", err)
		}

		// Now mark as applied in history store AFTER executing DDL
		if err := r.historyStore.MarkMigrationApplied(migration.Name, int64(migration.Version), string(migration.Checksum[:])); err != nil {
			return errors.NewWithCause(fmt.Sprintf("failed to mark migration %s as applied", migration.Name), "MIGRATION_HISTORY_MARK_FAILED", err)
		}
	}

	return nil
}

// Rollback rolls back the last applied migration.
func (r *Runner) Rollback(ctx context.Context) error {
	// Acquire lock
	acquired, err := r.lock.Acquire(ctx, r.timeout)
	if err != nil {
		return errors.NewWithCause("failed to acquire migration lock", "MIGRATION_LOCK_ACQUIRE_FAILED", err)
	}
	if !acquired {
		return errors.New("migration lock already held by another process", "MIGRATION_LOCK_HELD")
	}
	defer func() {
		// Release lock on exit
		lockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r.lock.Release(lockCtx)
	}()

	// Get applied migrations ordered by version descending (most recent first)
	appliedMigrations, err := r.historyStore.GetAppliedMigrations()
	if err != nil {
		return errors.NewWithCause("failed to get applied migrations", "MIGRATION_HISTORY_GET_FAILED", err)
	}

	if len(appliedMigrations) == 0 {
		return errors.New("no migrations to rollback", "MIGRATION_NO_MIGRATIONS")
	}

	// Get the most recent migration
	lastMigration := appliedMigrations[len(appliedMigrations)-1]

	// Find the matching migration definition (this would normally be in a registry)
	// For now, we just mark it as rolled back in history
	if err := r.historyStore.MarkMigrationRolledBack(lastMigration.Name, int64(lastMigration.Version)); err != nil {
		return errors.NewWithCause("failed to mark migration as rolled back", "MIGRATION_HISTORY_MARK_ROLLED_BACK_FAILED", err)
	}

	return nil
}

// executeMigration executes the actual DDL operations for a migration.
// This is a stub implementation - in a real system, this would contain
// the actual SQL statements for each migration.
func (r *Runner) executeMigration(ctx context.Context, migration *Migration) error {
	// In a real implementation, this would:
	// 1. Execute the actual DDL operations for the migration
	// 2. Handle errors appropriately
	// 3. Potentially support transaction rollback

	// For now, we'll just return nil as a placeholder
	// In practice, this would contain logic like:
	//
	// tx, err := r.db.BeginTx(ctx, nil)
	// if err != nil {
	//     return err
	// }
	// defer tx.Rollback()
	//
	// // Execute migration SQL statements
	// _, err = tx.ExecContext(ctx, migration.SQLStatement)
	// if err != nil {
	//     return err
	// }
	//
	// return tx.Commit()

	return nil
}
