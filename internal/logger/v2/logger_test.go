package v2

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogger(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "syspulse_test_logs")
	defer os.RemoveAll(tmpDir)

	t.Run("New Logger", func(t *testing.T) {
		logger, err := New(tmpDir, INFO)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		if logger.level != INFO {
			t.Errorf("Expected log level INFO, got %v", logger.level)
		}
	})

	t.Run("Log Levels", func(t *testing.T) {
		logger, _ := New(tmpDir, DEBUG)
		defer logger.Close()

		tests := []struct {
			level   LogLevel
			message string
		}{
			{DEBUG, "Debug message"},
			{INFO, "Info message"},
			{WARN, "Warning message"},
			{ERROR, "Error message"},
		}

		for _, tt := range tests {
			switch tt.level {
			case DEBUG:
				logger.Debug(tt.message)
			case INFO:
				logger.Info(tt.message)
			case WARN:
				logger.Warn(tt.message)
			case ERROR:
				logger.Error(tt.message)
			}
		}

		files, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to read log directory: %v", err)
		}
		if len(files) == 0 {
			t.Error("No log files created")
		}
	})

	t.Run("Log Rotation", func(t *testing.T) {
		logger, _ := New(tmpDir, INFO)
		defer logger.Close()

		logger.Info("Test message")
		err := logger.RotateLog()
		if err != nil {
			t.Fatalf("Failed to rotate log: %v", err)
		}

		logger.Info("Test message after rotation")
	})
}
