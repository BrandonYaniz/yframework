// Package migrations provides migration utilities and history tracking.
package migrations

import (
	"database/sql"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// MigrationHistory represents a record of a migration that has been applied.
type MigrationHistory struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Version   int64     `db:"version"`
	Checksum  string    `db:"checksum"`
	AppliedAt time.Time `db:"applied_at"`
}

// HistoryStore manages migration history in the database.
type HistoryStore struct {
	db *sql.DB
}

// NewHistoryStore creates a new HistoryStore with the given database connection.
func NewHistoryStore(db *sql.DB) *HistoryStore {
	return &HistoryStore{db: db}
}

// CreateSchema creates the migration history table if it doesn't exist.
func (h *HistoryStore) CreateSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS migration_history (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		version BIGINT NOT NULL,
		checksum VARCHAR(64) NOT NULL,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY unique_migration (name, version)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`

	_, err := h.db.Exec(query)
	if err != nil {
		return errors.NewWithCause("failed to create migration history schema", "MIGRATION_HISTORY_SCHEMA_CREATE_FAILED", err)
	}

	return nil
}

// GetAppliedMigrations returns all applied migrations.
func (h *HistoryStore) GetAppliedMigrations() ([]MigrationHistory, error) {
	query := `SELECT id, name, version, checksum, applied_at FROM migration_history ORDER BY version, id`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, errors.NewWithCause("failed to query migration history", "MIGRATION_HISTORY_QUERY_FAILED", err)
	}
	defer rows.Close()

	var migrations []MigrationHistory
	for rows.Next() {
		var m MigrationHistory
		if err := rows.Scan(&m.ID, &m.Name, &m.Version, &m.Checksum, &m.AppliedAt); err != nil {
			return nil, errors.NewWithCause("failed to scan migration history row", "MIGRATION_HISTORY_SCAN_FAILED", err)
		}
		migrations = append(migrations, m)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewWithCause("error iterating migration history rows", "MIGRATION_HISTORY_ITERATE_FAILED", err)
	}

	return migrations, nil
}

// MarkMigrationApplied records a migration as applied in the history.
func (h *HistoryStore) MarkMigrationApplied(name string, version int64, checksum string) error {
	query := `INSERT INTO migration_history (name, version, checksum) VALUES (?, ?, ?)`

	_, err := h.db.Exec(query, name, version, checksum)
	if err != nil {
		return errors.NewWithCause("failed to mark migration as applied", "MIGRATION_HISTORY_MARK_APPLIED_FAILED", err)
	}

	return nil
}

// MarkMigrationRolledBack removes a migration from the history.
func (h *HistoryStore) MarkMigrationRolledBack(name string, version int64) error {
	query := `DELETE FROM migration_history WHERE name = ? AND version = ?`

	_, err := h.db.Exec(query, name, version)
	if err != nil {
		return errors.NewWithCause("failed to mark migration as rolled back", "MIGRATION_HISTORY_MARK_ROLLED_BACK_FAILED", err)
	}

	return nil
}
