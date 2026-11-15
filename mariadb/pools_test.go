// Package mariadb provides MariaDB connection and configuration utilities.
package mariadb

import (
	"testing"
)

func TestPools(t *testing.T) {
	t.Skip("Skipping pool tests - requires database connection")
}

func TestForegroundPool(t *testing.T) {
	t.Skip("Skipping foreground pool test - requires database connection")
}

func TestBackgroundPool(t *testing.T) {
	t.Skip("Skipping background pool test - requires database connection")
}
