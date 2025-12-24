package logger

import (
	"go.uber.org/zap"
)

// Logger is the interface we will use across the system
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

type zapLogger struct {
	log *zap.Logger
}

// New initializes a new structured logger
func New(cfg Config) Logger {
	// 1. Build the Zap config (JSON vs Console, Level)
	zapCfg := buildZapConfig(cfg)

	// 2. Create the actual Zap logger
	l, err := zapCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		// Fallback to a basic logger if zap fails to initialize
		return &zapLogger{log: zap.NewExample()}
	}

	// 3. AUTOMATICALLY add the service name
	// Add the service name to every log automatically
	l = l.With(zap.String("service", cfg.ServiceName))

	// 4. Return our wrapper
	return &zapLogger{log: l}
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) { l.log.Debug(msg, fields...) }
func (l *zapLogger) Info(msg string, fields ...zap.Field)  { l.log.Info(msg, fields...) }
func (l *zapLogger) Warn(msg string, fields ...zap.Field)  { l.log.Warn(msg, fields...) }
func (l *zapLogger) Error(msg string, fields ...zap.Field) { l.log.Error(msg, fields...) }
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) { l.log.Fatal(msg, fields...) }
func (l *zapLogger) Sync() error                           { return l.log.Sync() }

// With allows us to add permanent context to a logger instance (e.g., a specific RequestID)
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{log: l.log.With(fields...)}
}
