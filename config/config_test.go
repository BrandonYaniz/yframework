package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set some test environment variables
	os.Setenv("YFRAMEWORK_ENV", "test")
	os.Setenv("YFRAMEWORK_LOG_LEVEL", "debug")
	os.Setenv("YFRAMEWORK_DATABASE_URL", "mysql://test:test@localhost/test")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", config.Environment)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.LogLevel)
	}

	if config.DatabaseURL != "mysql://test:test@localhost/test" {
		t.Errorf("Expected database URL 'mysql://test:test@localhost/test', got '%s'", config.DatabaseURL)
	}
}

func TestConfigValidate(t *testing.T) {
	// Test valid config
	validConfig := &Config{
		Environment: "development",
		LogLevel:    "info",
		DatabaseURL: "mysql://user:pass@localhost/db",
	}

	if err := validConfig.Validate(); err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test invalid config (empty environment)
	invalidConfig := &Config{
		Environment: "",
		LogLevel:    "info",
	}

	if err := invalidConfig.Validate(); err == nil {
		t.Error("Expected invalid config to fail validation")
	}
}
