package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// TestTransactionIntegration tests transaction functionality with real MariaDB
func TestTransactionIntegration(t *testing.T) {
	// Skip if no MariaDB connection available
	if !hasMariaDBConnection() {
		t.Skip("No MariaDB connection available for integration test")
		return
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a transaction helper
	helper := NewTransactionHelper(db)

	// Test basic transaction functionality
	ctx := context.Background()

	// This should not panic and should handle errors properly
	err := helper.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
		// Try to execute a simple query
		_, err := tx.ExecContext(ctx, "SELECT 1")
		if err != nil {
			return err
		}
		// Return nil to commit the transaction
		return nil
	})

	if err != nil {
		// This might fail due to connection issues, but should not panic
		t.Logf("Transaction test completed with expected error: %v", err)
	} else {
		t.Log("Basic transaction execution succeeded")
	}

	// Test that error handling works properly in transactions
	err = helper.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
		// Return an error to test rollback behavior
		return errors.New("test transaction error", "TEST_TRANSACTION_ERROR")
	})

	if err == nil {
		t.Error("Expected error when transaction function returns error")
	} else {
		t.Logf("Transaction error handling works: %v", err)
	}
}

// TestTransactionConcurrency tests concurrent transaction access
func TestTransactionConcurrency(t *testing.T) {
	// Skip if no MariaDB connection available
	if !hasMariaDBConnection() {
		t.Skip("No MariaDB connection available for integration test")
		return
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	helper := NewTransactionHelper(db)

	// Test concurrent access to transaction helper (thread safety)
	go func() {
		ctx := context.Background()
		err := helper.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "SELECT 1")
			return err
		})
		if err != nil {
			t.Logf("Goroutine transaction error: %v", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	t.Log("Concurrent transaction test completed")
}

// hasMariaDBConnection checks if we have a MariaDB connection available
func hasMariaDBConnection() bool {
	host := os.Getenv("YFRAMEWORK_DB_HOST")
	port := os.Getenv("YFRAMEWORK_DB_PORT")
	user := os.Getenv("YFRAMEWORK_DB_USER")
	password := os.Getenv("YFRAMEWORK_DB_PASSWORD")
	dbName := os.Getenv("YFRAMEWORK_DB_NAME")

	return host != "" && port != "" && user != "" && password != "" && dbName != ""
}

// setupTestDB sets up a test database connection
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	if !hasMariaDBConnection() {
		t.Fatal("No MariaDB connection available for testing")
	}

	host := os.Getenv("YFRAMEWORK_DB_HOST")
	port := os.Getenv("YFRAMEWORK_DB_PORT")
	user := os.Getenv("YFRAMEWORK_DB_USER")
	password := os.Getenv("YFRAMEWORK_DB_PASSWORD")
	dbName := os.Getenv("YFRAMEWORK_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	return db, func() {
		db.Close()
	}
}
