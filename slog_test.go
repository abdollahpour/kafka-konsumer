package kafka

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlogLogger(t *testing.T) {
	logger := NewSlogLogger()
	assert.NotNil(t, logger, "NewSlogLogger() should not return nil")
	assert.NotNil(t, logger.logger, "logger.logger should not be nil")
}

func TestSlogLogger_With(t *testing.T) {
	logger := NewSlogLogger()

	t.Run("With even number of args", func(t *testing.T) {
		newLogger := logger.With("key1", "value1", "key2", "value2")
		assert.NotNil(t, newLogger, "With should not return nil")
	})

	t.Run("With odd number of args", func(t *testing.T) {
		newLogger := logger.With("key1", "value1", "key2")
		assert.NotNil(t, newLogger, "With should not return nil")
	})

	t.Run("With empty args", func(t *testing.T) {
		newLogger := logger.With()
		assert.NotNil(t, newLogger, "With should not return nil")
	})
}

func TestSlogLogger_LogMethods(t *testing.T) {
	// Create a logger that writes to a buffer to verify logging works
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := &SlogLogger{logger: slog.New(handler)}

	tests := []struct {
		name string
		fn   func()
	}{
		{"Debug", func() { logger.Debug("test message") }},
		{"Info", func() { logger.Info("test message") }},
		{"Warn", func() { logger.Warn("test message") }},
		{"Error", func() { logger.Error("test message") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			assert.Greater(t, buf.Len(), 0, "%s should have logged something", tt.name)
			buf.Reset()
		})
	}
}

func TestSlogLogger_LogMethodsWithArgs(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	t.Run("Debug with multiple args", func(t *testing.T) {
		logger.Debug("arg1", "arg2", "arg3")
	})

	t.Run("Info with multiple args", func(t *testing.T) {
		logger.Info("arg1", "arg2", "arg3")
	})

	t.Run("Warn with multiple args", func(t *testing.T) {
		logger.Warn("arg1", "arg2", "arg3")
	})

	t.Run("Error with multiple args", func(t *testing.T) {
		logger.Error("arg1", "arg2", "arg3")
	})
}

func TestSlogLogger_LogfMethods(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	tests := []struct {
		name string
		fn   func()
	}{
		{"Debugf", func() { logger.Debugf("format: %s", "value") }},
		{"Infof", func() { logger.Infof("format: %s", "value") }},
		{"Warnf", func() { logger.Warnf("format: %s", "value") }},
		{"Errorf", func() { logger.Errorf("format: %s", "value") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			tt.fn()
		})
	}
}

func TestSlogLogger_Infow(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	t.Run("With key-value pairs", func(t *testing.T) {
		logger.Infow("message", "key1", "value1", "key2", "value2")
	})

	t.Run("With odd number of key-value pairs", func(t *testing.T) {
		logger.Infow("message", "key1", "value1", "key2")
	})

	t.Run("Without key-value pairs", func(t *testing.T) {
		logger.Infow("message")
	})
}

func TestSlogLogger_Errorw(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	t.Run("With key-value pairs", func(t *testing.T) {
		logger.Errorw("message", "key1", "value1", "key2", "value2")
	})

	t.Run("With odd number of key-value pairs", func(t *testing.T) {
		logger.Errorw("message", "key1", "value1", "key2")
	})

	t.Run("Without key-value pairs", func(t *testing.T) {
		logger.Errorw("message")
	})
}

func TestSlogLogger_Warnw(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	t.Run("With key-value pairs", func(t *testing.T) {
		logger.Warnw("message", "key1", "value1", "key2", "value2")
	})

	t.Run("With odd number of key-value pairs", func(t *testing.T) {
		logger.Warnw("message", "key1", "value1", "key2")
	})

	t.Run("Without key-value pairs", func(t *testing.T) {
		logger.Warnw("message")
	})
}

func TestSlogLogger_buildAttrs(t *testing.T) {
	logger := &SlogLogger{logger: slog.Default()}

	t.Run("Even number of key-value pairs", func(t *testing.T) {
		attrs := logger.buildAttrs("key1", "value1", "key2", "value2")
		assert.Equal(t, 4, len(attrs), "Should have 4 attrs")
	})

	t.Run("Odd number of key-value pairs", func(t *testing.T) {
		attrs := logger.buildAttrs("key1", "value1", "key2")
		assert.Equal(t, 4, len(attrs), "Should have 4 attrs")
	})

	t.Run("Empty key-value pairs", func(t *testing.T) {
		attrs := logger.buildAttrs()
		assert.Equal(t, 0, len(attrs), "Should have 0 attrs")
	})

	t.Run("Non-string key conversion", func(t *testing.T) {
		attrs := logger.buildAttrs(123, "value1", "key2", "value2")
		assert.Equal(t, 4, len(attrs), "Should have 4 attrs")
	})
}

func TestSlogLogger_Integration(t *testing.T) {
	// Test that SlogLogger implements LoggerInterface
	var _ LoggerInterface = &SlogLogger{}

	logger := NewSlogLogger()

	// Test chaining
	logger = logger.With("context1", "value1").(*SlogLogger)
	logger = logger.With("context2", "value2").(*SlogLogger)

	// Test all log levels with chained logger
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	logger.Debugf("debug format: %s", "value")
	logger.Infof("info format: %s", "value")
	logger.Warnf("warn format: %s", "value")
	logger.Errorf("error format: %s", "value")

	logger.Infow("info with key-value", "key", "value")
	logger.Warnw("warn with key-value", "key", "value")
	logger.Errorw("error with key-value", "key", "value")
}

func TestSlogLogger_FileOutput(t *testing.T) {
	// Test that logging to a file doesn't panic
	file, err := os.CreateTemp("", "slog-test-*.log")
	assert.NoError(t, err, "Should create temp file")
	defer os.Remove(file.Name())
	defer file.Close()

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := &SlogLogger{logger: slog.New(handler)}

	logger.Debug("test debug")
	logger.Info("test info")
	logger.Warn("test warn")
	logger.Error("test error")

	// Verify file has content
	info, err := file.Stat()
	assert.NoError(t, err, "Should stat file")
	assert.Greater(t, info.Size(), int64(0), "File should have content")
}
