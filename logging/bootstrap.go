// Package logging provides framework-wide logging utilities.
package logging

import (
	"log"
	"os"
)

// BootstrapLogger is a minimal logger that can be used before full configuration.
// It always outputs to stderr and is available during bootstrap phase.
var BootstrapLogger = &bootstrapLogger{
	logger: log.New(os.Stderr, "", 0),
}

// bootstrapLogger implements the Logger interface for bootstrap purposes only.
type bootstrapLogger struct {
	logger *log.Logger
}

// Debug logs a debug message to stderr.
func (l *bootstrapLogger) Debug(msg string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

// Info logs an info message to stderr.
func (l *bootstrapLogger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Warn logs a warning message to stderr.
func (l *bootstrapLogger) Warn(msg string, args ...interface{}) {
	l.logger.Printf("[WARN] "+msg, args...)
}

// Error logs an error message to stderr.
func (l *bootstrapLogger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}
