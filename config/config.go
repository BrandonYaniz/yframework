// Package config provides framework-wide configuration utilities.
package config

import (
	"fmt"
	"os"
)

// Config represents the framework's configuration structure.
type Config struct {
	// Environment specifies the environment (e.g., "development", "production").
	Environment string

	// LogLevel specifies the logging level.
	LogLevel string

	// DatabaseURL specifies the database connection URL.
	DatabaseURL string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	return &Config{
		Environment: getEnv("YFRAMEWORK_ENV", "development"),
		LogLevel:    getEnv("YFRAMEWORK_LOG_LEVEL", "info"),
		DatabaseURL: getEnv("YFRAMEWORK_DATABASE_URL", ""),
	}, nil
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Environment == "" {
		return fmt.Errorf("environment cannot be empty")
	}

	if c.LogLevel == "" {
		return fmt.Errorf("log level cannot be empty")
	}

	return nil
}
