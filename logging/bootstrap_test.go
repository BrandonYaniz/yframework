package logging

import (
	"testing"
)

func TestBootstrapLoggerInterface(t *testing.T) {
	// Verify BootstrapLogger implements the Logger interface
	var _ Logger = BootstrapLogger

	// This test will fail at compile time if BootstrapLogger doesn't implement Logger
	// We're just ensuring it compiles correctly
	_ = BootstrapLogger
}

func TestBootstrapLoggerBasic(t *testing.T) {
	// Just verify that we can call the methods without errors
	// The actual output capture is difficult in this context
	BootstrapLogger.Debug("test debug")
	BootstrapLogger.Info("test info")
	BootstrapLogger.Warn("test warn")
	BootstrapLogger.Error("test error")

	// If we reach here without panic, the test passes
	t.Log("Bootstrap logger methods executed successfully")
}
