// Package mariadb provides MariaDB integration and transaction helpers.
package mariadb

import (
	"context"
	"database/sql"

	"github.com/BrandonYaniz/yframework/errors"
)

// TransactionHelper provides helper methods for database transactions.
type TransactionHelper struct {
	db *sql.DB
}

// NewTransactionHelper creates a new transaction helper with the given database connection.
func NewTransactionHelper(db *sql.DB) *TransactionHelper {
	return &TransactionHelper{db: db}
}

// ExecuteInTransaction executes the provided function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (th *TransactionHelper) ExecuteInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := th.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewWithCause("failed to begin transaction", "TRANSACTION_BEGIN_FAILED", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.NewWithCause("failed to rollback transaction", "TRANSACTION_ROLLBACK_FAILED", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.NewWithCause("failed to commit transaction", "TRANSACTION_COMMIT_FAILED", err)
	}

	return nil
}
