package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(env string) (*ZapLogger, error) {
	var zapConfig zap.Config

	if env == "production" {
		// Production configuration
		zapConfig = zap.NewProductionConfig()

		// JSON format for log aggregation
		zapConfig.Encoding = "json"

		// Only log info and above in production
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

		// Disable stack traces for warnings (only for errors)
		zapConfig.DisableStacktrace = false

		// Remove caller information in production
		zapConfig.DisableCaller = true

		// Output to stdout (for container logs)
		zapConfig.OutputPaths = []string{"stdout"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}

		// Sampling to reduce log volume in production
		zapConfig.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	} else {
		// Development configuration
		zapConfig = zap.NewDevelopmentConfig()

		// Human-readable format for development
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		// Show all logs in development
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

		// Enable stack traces
		zapConfig.DisableStacktrace = false

		// Show caller information
		zapConfig.DisableCaller = false
	}

	// Add common fields
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.InitialFields = map[string]interface{}{
		"app":         "dashboard-backend",
		"environment": env,
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	// Replace global logger
	zap.ReplaceGlobals(logger)

	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *ZapLogger) With(fields ...zap.Field) *ZapLogger {
	return &ZapLogger{
		logger: l.logger.With(fields...),
	}
}
