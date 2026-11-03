package logging

import (
	"testing"
)

func TestStdLogger(t *testing.T) {
	logger := NewStdLogger()

	if logger == nil {
		t.Error("Expected logger to be created")
	}
}
