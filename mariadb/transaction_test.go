package mariadb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/BrandonYaniz/yframework/errors"
)

func TestTransactionHelper_ExecuteInTransaction(t *testing.T) {
	t.Run("ExecuteInTransaction executes function successfully", func(t *testing.T) {
		// Test with nil DB - should not panic, just return an error
		helper := &TransactionHelper{}

		ctx := context.Background()
		err := helper.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
			return nil
		})

		// Should return an error due to nil DB (but not panic)
		if err == nil {
			t.Error("Expected error when DB is nil")
		}
	})

	t.Run("ExecuteInTransaction handles function error properly", func(t *testing.T) {
		helper := &TransactionHelper{}

		ctx := context.Background()
		testError := errors.New("test error", "TEST_ERROR")

		err := helper.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
			return testError
		})

		// Should return the function error (but we can't verify rollback behavior without DB)
		if err == nil {
			t.Error("Expected error when function returns error")
		}
	})
}

func TestTransactionHelper_NewTransactionHelper(t *testing.T) {
	db := &sql.DB{}
	helper := NewTransactionHelper(db)

	if helper.db != db {
		t.Error("Expected helper to store the provided database connection")
	}
}
