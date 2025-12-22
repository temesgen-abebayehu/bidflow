package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// 1. Set temporary environment variables
	os.Setenv("DB_HOST", "test-db-host")
	os.Setenv("KAFKA_BROKERS", "kafka1:9092,kafka2:9092")
	
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("KAFKA_BROKERS")

	// 2. Load the config
	cfg := LoadConfig("test-service")

	// 3. Assert values
	assert.Equal(t, "test-service", cfg.ServiceName)
	assert.Equal(t, "test-db-host", cfg.DBHost)
	assert.Equal(t, []string{"kafka1:9092", "kafka2:9092"}, cfg.KafkaBrokers)
}

func TestGetEnvFallback(t *testing.T) {
	// Should return default because KEY_NOT_EXIST is not set
	val := getEnv("KEY_NOT_EXIST", "fallback_val")
	assert.Equal(t, "fallback_val", val)
}