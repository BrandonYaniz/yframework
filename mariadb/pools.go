// Package mariadb provides MariaDB connection and configuration utilities.
package mariadb

import (
	"database/sql"
	"time"

	"github.com/BrandonYaniz/yframework/errors"
)

// Pool represents a database connection pool.
type Pool struct {
	db *sql.DB
}

// ForegroundPool represents a connection pool for foreground operations (typically web requests).
type ForegroundPool struct {
	*Pool
}

// BackgroundPool represents a connection pool for background operations (jobs, cron tasks).
type BackgroundPool struct {
	*Pool
}

// NewForegroundPool creates a new foreground connection pool.
func NewForegroundPool(db *sql.DB) *ForegroundPool {
	return &ForegroundPool{Pool: &Pool{db: db}}
}

// NewBackgroundPool creates a new background connection pool.
func NewBackgroundPool(db *sql.DB) *BackgroundPool {
	return &BackgroundPool{Pool: &Pool{db: db}}
}

// GetDB returns the underlying database connection.
func (p *Pool) GetDB() *sql.DB {
	return p.db
}

// Close closes all connections in the pool.
func (p *Pool) Close() error {
	return p.db.Close()
}

// Ping checks if the database is reachable.
func (p *Pool) Ping() error {
	return p.db.Ping()
}

// SetMaxIdleConns sets the maximum number of idle connections in the pool.
func (p *Pool) SetMaxIdleConns(maxIdleConns int) {
	p.db.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns sets the maximum number of open connections in the pool.
func (p *Pool) SetMaxOpenConns(maxOpenConns int) {
	p.db.SetMaxOpenConns(maxOpenConns)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (p *Pool) SetConnMaxLifetime(connMaxLifetime time.Duration) {
	p.db.SetConnMaxLifetime(connMaxLifetime)
}

// GetForegroundPool returns the foreground pool instance.
func GetForegroundPool() (*ForegroundPool, error) {
	// This would typically return a configured pool
	// For now, we're just returning a placeholder implementation
	return nil, errors.New("foreground pool not implemented", "MARIADB_FOREGROUND_POOL_NOT_IMPLEMENTED")
}

// GetBackgroundPool returns the background pool instance.
func GetBackgroundPool() (*BackgroundPool, error) {
	// This would typically return a configured pool
	// For now, we're just returning a placeholder implementation
	return nil, errors.New("background pool not implemented", "MARIADB_BACKGROUND_POOL_NOT_IMPLEMENTED")
}
