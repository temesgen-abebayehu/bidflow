package config

import (
	"strings"
)

// Config holds all configuration for a microservice
type Config struct {
	Env         string // development, production, test
	ServiceName string

	// Service URLs (for Gateway)
	AuthServiceURL         string
	AuctionServiceURL      string
	BiddingServiceURL      string
	NotificationServiceURL string

	// Server configurations
	HTTPPort string
	GRPCPort string

	// Database configurations
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string

	// Kafka configurations
	KafkaBrokers []string

	// Auth configurations
	JWTSecret string
}

// LoadConfig merges environment variables into the Config struct
func LoadConfig(serviceName string) *Config {
	return &Config{
		Env:         getEnv("APP_ENV", "development"),
		ServiceName: serviceName,

		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
		AuctionServiceURL:      getEnv("AUCTION_SERVICE_URL", "http://auction-service:8081"),
		BiddingServiceURL:      getEnv("BIDDING_SERVICE_URL", "http://bidding-service:8082"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8084"),

		HTTPPort: getEnv("HTTP_PORT", "8080"),
		GRPCPort: getEnv("GRPC_PORT", "50051"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", serviceName+"_db"), // Default: auth_db, auction_db, etc.
		DBSslMode:  getEnv("DB_SSLMODE", "disable"),

		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),

		JWTSecret: getEnv("JWT_SECRET", "bidflow_default_secret_key_change_me"),
	}
}
