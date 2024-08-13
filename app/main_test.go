package main

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/superboomer/maptile/app/options"
)

func TestCreateLogger(t *testing.T) {

	tests := []struct {
		name       string
		opts       options.Log
		expectFile bool
		expectPath string
	}{
		{
			name:       "No log file",
			opts:       options.Log{Save: false},
			expectFile: false,
		},
		{
			name: "Log to file",
			opts: options.Log{
				Save:       true,
				Path:       "test.log",
				MaxBackups: 10,
				MaxSize:    5,
				MaxAge:     30,
			},
			expectFile: true,
			expectPath: "test.log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := createLogger(&tt.opts)

			// Check if logger is not nil
			if logger == nil {
				t.Fatalf("Expected logger to be non-nil")
			}

			// If logging to a file is expected, check if the file exists
			if tt.expectFile {
				if _, err := os.Stat(tt.expectPath); os.IsNotExist(err) {
					t.Fatalf("Expected log file %s to exist", tt.expectPath)
				}
				// Clean up after test
				defer os.Remove(tt.expectPath)
			}
		})
	}
}

func TestRun(t *testing.T) {
	opts := &options.Opts{} // Populate with valid options for your tests
	logger := zap.NewNop()  // Use a no-op logger for testing

	// Simulate running the app
	err := run(opts, logger)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}
