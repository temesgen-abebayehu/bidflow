package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoggerInitialization(t *testing.T) {
	cfg := Config{
		Level:       "debug",
		Development: true,
		ServiceName: "test-service",
	}

	l := New(cfg)
	assert.NotNil(t, l)

	// Verify it doesn't panic when logging
	l.Info("testing logger info message", zap.String("key", "value"))
}

func TestLoggerWithFields(t *testing.T) {
	cfg := Config{
		Level:       "info",
		Development: false,
		ServiceName: "test-service",
	}

	l := New(cfg)
	// Create a sub-logger with a fixed ID
	subLogger := l.With(zap.String("request_id", "123-abc"))
	
	assert.NotNil(t, subLogger)
	subLogger.Info("this log will contain the request_id")
}