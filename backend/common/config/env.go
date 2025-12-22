package config

import (
	"os"

	"github.com/joho/godotenv"
)

// InitEnv loads the .env file if it exists. 
// In Production (K8s/Azure), the .env file won't exist, and that's okay.
func InitEnv() {
	_ = godotenv.Load() // Ignore error if .env is missing (expected in K8s)
}

// getEnv is a helper to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}