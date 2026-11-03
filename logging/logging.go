// Package logging provides framework-wide logging utilities.
package logging

import (
	"log"
	"os"
)

// Logger interface defines the logging methods used by the framework.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// StdLogger implements the Logger interface using Go's standard log package.
type StdLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// NewStdLogger creates a new StdLogger instance.
func NewStdLogger() *StdLogger {
	return &StdLogger{
		debugLogger: log.New(os.Stderr, "DEBUG: ", log.LstdFlags),
		infoLogger:  log.New(os.Stderr, "INFO: ", log.LstdFlags),
		warnLogger:  log.New(os.Stderr, "WARN: ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// Debug logs a debug message.
func (l *StdLogger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Printf(msg, args...)
}

// Info logs an info message.
func (l *StdLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

// Warn logs a warning message.
func (l *StdLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}

// Error logs an error message.
func (l *StdLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}
