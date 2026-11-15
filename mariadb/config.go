// Package mariadb provides MariaDB connection and configuration utilities.
package mariadb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/BrandonYaniz/yframework/config"
	"github.com/BrandonYaniz/yframework/errors"
)

// Config represents the MariaDB configuration.
type Config struct {
	// Host is the database host
	Host string `toml:"host"`

	// Port is the database port
	Port int `toml:"port"`

	// User is the database user
	User string `toml:"user"`

	// Password is the database password
	Password string `toml:"password"`

	// Database is the database name
	Database string `toml:"database"`

	// SSLMode specifies the SSL mode (disable, prefer, require)
	SSLMode string `toml:"ssl_mode"`

	// ConnectionTimeout is the timeout for establishing connections
	ConnectionTimeout time.Duration `toml:"connection_timeout"`

	// MaxIdleConns is the maximum number of idle connections in the pool
	MaxIdleConns int `toml:"max_idle_conns"`

	// MaxOpenConns is the maximum number of open connections in the pool
	MaxOpenConns int `toml:"max_open_conns"`

	// ConnMaxLifetime is the maximum amount of time a connection may be reused
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
}

// Validate validates the MariaDB configuration.
func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is required", "MARIADB_HOST_REQUIRED")
	}

	if c.Port <= 0 {
		return errors.New("port must be positive", "MARIADB_PORT_INVALID")
	}

	if c.Database == "" {
		return errors.New("database is required", "MARIADB_DATABASE_REQUIRED")
	}

	if c.User == "" {
		return errors.New("user is required", "MARIADB_USER_REQUIRED")
	}

	return nil
}

// DSN returns the Data Source Name for the MariaDB connection.
func (c *Config) DSN() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "prefer"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%d&readTimeout=%d&writeTimeout=%d&parseTime=true&loc=UTC&tls=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		int(c.ConnectionTimeout.Seconds()),
		int(c.ConnectionTimeout.Seconds()),
		int(c.ConnectionTimeout.Seconds()),
		sslMode)
}

// NewDB creates a new database connection using the configuration.
func (c *Config) NewDB() (*sql.DB, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.NewWithCause("invalid MariaDB configuration", "MARIADB_CONFIG_INVALID", err)
	}

	db, err := sql.Open("mysql", c.DSN())
	if err != nil {
		return nil, errors.NewWithCause("failed to open database connection", "MARIADB_CONNECTION_FAILED", err)
	}

	// Set connection pool settings
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	return db, nil
}

// LoadConfig loads MariaDB configuration from the provided config source.
func LoadConfig(cfg *config.Config) (*Config, error) {
	var mariadbCfg Config

	// Extract connection details from environment variables
	mariadbCfg.Host = getEnv("YFRAMEWORK_DB_HOST", "127.0.0.1")
	mariadbCfg.Port = getIntEnv("YFRAMEWORK_DB_PORT", 3306)
	mariadbCfg.User = getEnv("YFRAMEWORK_DB_USER", "devframework")
	mariadbCfg.Password = getEnv("YFRAMEWORK_DB_PASSWORD", "")
	mariadbCfg.Database = getEnv("YFRAMEWORK_DB_NAME", "devframework")

	// Set defaults if not provided
	if mariadbCfg.ConnectionTimeout == 0 {
		mariadbCfg.ConnectionTimeout = 5 * time.Second
	}

	if mariadbCfg.MaxIdleConns == 0 {
		mariadbCfg.MaxIdleConns = 2
	}

	if mariadbCfg.MaxOpenConns == 0 {
		mariadbCfg.MaxOpenConns = 10
	}

	if mariadbCfg.ConnMaxLifetime == 0 {
		mariadbCfg.ConnMaxLifetime = 5 * time.Minute
	}

	// Validate the loaded config
	if err := mariadbCfg.Validate(); err != nil {
		return nil, errors.NewWithCause("failed to validate MariaDB configuration", "MARIADB_CONFIG_VALIDATION_FAILED", err)
	}

	return &mariadbCfg, nil
}

// Helper function to get integer environment variables
func getIntEnv(key string, defaultValue int) int {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}

	// Parse the integer value
	if val, err := strconv.Atoi(env); err == nil {
		return val
	}

	return defaultValue
}

// Helper function to get environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
