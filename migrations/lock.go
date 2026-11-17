// Package migrations provides migration utilities and history tracking.
package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// AdvisoryLock represents an advisory lock for migration operations.
type AdvisoryLock struct {
	db     *sql.DB
	lockID int64
}

// NewAdvisoryLock creates a new advisory lock with the given database connection and lock ID.
func NewAdvisoryLock(db *sql.DB, lockID int64) *AdvisoryLock {
	return &AdvisoryLock{db: db, lockID: lockID}
}

// Acquire attempts to acquire the advisory lock.
// Returns true if the lock was acquired, false if it was already held by another process.
func (l *AdvisoryLock) Acquire(ctx context.Context, timeout time.Duration) (bool, error) {
	// Create a dedicated connection for the lock operations
	conn, err := l.db.Conn(ctx)
	if err != nil {
		return false, errors.NewWithCause("failed to create dedicated connection for advisory lock", "MIGRATION_LOCK_CONN_CREATE_FAILED", err)
	}
	defer conn.Close()

	// Use MySQL's GET_LOCK function with dedicated connection
	query := `SELECT GET_LOCK(?, ?)`

	var result int64
	err = conn.QueryRowContext(ctx, query, l.lockID, int(timeout.Seconds())).Scan(&result)
	if err != nil {
		return false, errors.NewWithCause("failed to acquire advisory lock", "MIGRATION_LOCK_ACQUIRE_FAILED", err)
	}

	return result == 1, nil
}

// Release releases the advisory lock.
func (l *AdvisoryLock) Release(ctx context.Context) error {
	// Create a dedicated connection for the lock operations
	conn, err := l.db.Conn(ctx)
	if err != nil {
		return errors.NewWithCause("failed to create dedicated connection for advisory lock release", "MIGRATION_LOCK_CONN_CREATE_FAILED", err)
	}
	defer conn.Close()

	// Use MySQL's RELEASE_LOCK function with dedicated connection
	query := `SELECT RELEASE_LOCK(?)`

	_, err = conn.ExecContext(ctx, query, l.lockID)
	if err != nil {
		return errors.NewWithCause("failed to release advisory lock", "MIGRATION_LOCK_RELEASE_FAILED", err)
	}

	return nil
}

// IsLocked checks if the advisory lock is currently held by any process.
func (l *AdvisoryLock) IsLocked(ctx context.Context) (bool, error) {
	// Create a dedicated connection for the lock operations
	conn, err := l.db.Conn(ctx)
	if err != nil {
		return false, errors.NewWithCause("failed to create dedicated connection for advisory lock status check", "MIGRATION_LOCK_CONN_CREATE_FAILED", err)
	}
	defer conn.Close()

	// Use MySQL's IS_USED_LOCK function to check if lock is held with dedicated connection
	query := `SELECT IS_USED_LOCK(?)`

	var result int64
	err = conn.QueryRowContext(ctx, query, l.lockID).Scan(&result)
	if err != nil {
		return false, errors.NewWithCause("failed to check advisory lock status", "MIGRATION_LOCK_STATUS_CHECK_FAILED", err)
	}

	return result != 0, nil
}
