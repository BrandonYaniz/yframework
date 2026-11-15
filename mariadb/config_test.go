package mariadb

import (
	"os"
	"testing"
	"time"

	"github.com/BrandonYaniz/yframework/config"
	"github.com/BrandonYaniz/yframework/errors"
)

func TestConfig_Validate(t *testing.T) {
	// Test valid config
	validConfig := &Config{
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
	}

	if err := validConfig.Validate(); err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}

	// Test missing host
	invalidConfig := &Config{
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
	}

	if err := invalidConfig.Validate(); err == nil {
		t.Error("Expected error for missing host")
	}

	// Test invalid port
	invalidPortConfig := &Config{
		Host:     "localhost",
		Port:     0,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
	}

	if err := invalidPortConfig.Validate(); err == nil {
		t.Error("Expected error for invalid port")
	}

	// Test missing database
	missingDBConfig := &Config{
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
	}

	if err := missingDBConfig.Validate(); err == nil {
		t.Error("Expected error for missing database")
	}

	// Test missing user
	missingUserConfig := &Config{
		Host:     "localhost",
		Port:     3306,
		Password: "testpass",
		Database: "testdb",
	}

	if err := missingUserConfig.Validate(); err == nil {
		t.Error("Expected error for missing user")
	}
}

func TestConfig_DSN(t *testing.T) {
	config := &Config{
		Host:              "localhost",
		Port:              3306,
		User:              "testuser",
		Password:          "testpass",
		Database:          "testdb",
		SSLMode:           "require",
		ConnectionTimeout: 5 * time.Second,
	}

	dsn := config.DSN()
	expected := "testuser:testpass@tcp(localhost:3306)/testdb?timeout=5&readTimeout=5&writeTimeout=5&parseTime=true&loc=UTC&tls=require"

	if dsn != expected {
		t.Errorf("Expected DSN %s, got %s", expected, dsn)
	}
}

func TestConfig_DSN_DefaultSSL(t *testing.T) {
	config := &Config{
		Host:              "localhost",
		Port:              3306,
		User:              "testuser",
		Password:          "testpass",
		Database:          "testdb",
		ConnectionTimeout: 5 * time.Second,
	}

	dsn := config.DSN()
	expected := "testuser:testpass@tcp(localhost:3306)/testdb?timeout=5&readTimeout=5&writeTimeout=5&parseTime=true&loc=UTC&tls=prefer"

	if dsn != expected {
		t.Errorf("Expected DSN %s, got %s", expected, dsn)
	}
}

func TestLoadConfig(t *testing.T) {
	// Set up environment variables to simulate what would be in the config
	os.Setenv("YFRAMEWORK_DB_HOST", "localhost")
	os.Setenv("YFRAMEWORK_DB_PORT", "3306")
	os.Setenv("YFRAMEWORK_DB_USER", "testuser")
	os.Setenv("YFRAMEWORK_DB_PASSWORD", "testpass")
	os.Setenv("YFRAMEWORK_DB_NAME", "testdb")

	// Load config
	cfg, err := LoadConfig(&config.Config{})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values
	if cfg.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Port != 3306 {
		t.Errorf("Expected port 3306, got %d", cfg.Port)
	}
	if cfg.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", cfg.User)
	}
	if cfg.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", cfg.Password)
	}
	if cfg.Database != "testdb" {
		t.Errorf("Expected database 'testdb', got '%s'", cfg.Database)
	}

	// Clean up
	os.Unsetenv("YFRAMEWORK_DB_HOST")
	os.Unsetenv("YFRAMEWORK_DB_PORT")
	os.Unsetenv("YFRAMEWORK_DB_USER")
	os.Unsetenv("YFRAMEWORK_DB_PASSWORD")
	os.Unsetenv("YFRAMEWORK_DB_NAME")

	// Test with missing required fields
	os.Unsetenv("YFRAMEWORK_DB_HOST") // Remove host to test validation error
	_, err = LoadConfig(&config.Config{})
	if err == nil {
		t.Error("Expected error when loading config without required fields")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Set up environment variables for minimal config
	os.Setenv("YFRAMEWORK_DB_HOST", "localhost")
	os.Setenv("YFRAMEWORK_DB_USER", "testuser")
	os.Setenv("YFRAMEWORK_DB_PASSWORD", "testpass")
	os.Setenv("YFRAMEWORK_DB_NAME", "testdb")

	// Load config
	mariadbCfg, err := LoadConfig(&config.Config{})
	if err != nil {
		t.Fatalf("Expected no error when loading config with defaults, got: %v", err)
	}

	// Test that defaults were applied
	if mariadbCfg.Port != 3306 {
		t.Errorf("Expected default port 3306, got %d", mariadbCfg.Port)
	}

	if mariadbCfg.ConnectionTimeout != 5*time.Second {
		t.Errorf("Expected default connection timeout 5s, got %v", mariadbCfg.ConnectionTimeout)
	}

	if mariadbCfg.MaxIdleConns != 2 {
		t.Errorf("Expected default max idle conns 2, got %d", mariadbCfg.MaxIdleConns)
	}

	if mariadbCfg.MaxOpenConns != 10 {
		t.Errorf("Expected default max open conns 10, got %d", mariadbCfg.MaxOpenConns)
	}

	if mariadbCfg.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("Expected default conn max lifetime 5m, got %v", mariadbCfg.ConnMaxLifetime)
	}

	// Clean up
	os.Unsetenv("YFRAMEWORK_DB_HOST")
	os.Unsetenv("YFRAMEWORK_DB_USER")
	os.Unsetenv("YFRAMEWORK_DB_PASSWORD")
	os.Unsetenv("YFRAMEWORK_DB_NAME")
}

func TestConfig_NewDB(t *testing.T) {
	// This is a basic test - in practice we would need an actual DB connection
	// but we can at least test that it doesn't panic with invalid config

	config := &Config{
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
	}

	// This should not panic, but may return an error due to connection failure
	db, err := config.NewDB()

	// We expect an error because we don't have a real DB running
	// But we want to make sure it doesn't panic
	if err != nil {
		// This is expected - we're testing that the function doesn't panic
		t.Logf("Expected error (no real DB): %v", err)
	} else if db != nil {
		// If it somehow succeeded, close it
		db.Close()
	}

	// Test with invalid config
	invalidConfig := &Config{
		Host: "localhost",
		// Missing other required fields
	}

	_, err = invalidConfig.NewDB()
	if err == nil {
		t.Error("Expected error for invalid config")
	}

	if !errors.IsValidationError(err) {
		t.Logf("Error type: %T", err)
	}
}
