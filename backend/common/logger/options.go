package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config defines the configuration for the logger
type Config struct {
	Level       string // debug, info, warn, error
	Development bool   // true for console-friendly logs, false for JSON
	ServiceName string // name of the microservice (e.g., "auction-service")
}

// buildZapConfig creates a zap.Config based on our custom Config
func buildZapConfig(cfg Config) zap.Config {
	var zapCfg zap.Config

	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colors for terminal
	} else {
		zapCfg = zap.NewProductionConfig() // JSON for Azure/AKS
	}

	// Set the log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zap.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	return zapCfg
}