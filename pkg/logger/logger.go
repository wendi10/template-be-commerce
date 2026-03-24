package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initialises the global zap logger.
func Init(level string, debug bool) {
	var cfg zap.Config

	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	lvl := zap.InfoLevel
	if err := lvl.UnmarshalText([]byte(level)); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	var err error
	log, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		// Fallback to stderr and exit
		_, _ = os.Stderr.WriteString("failed to initialise logger: " + err.Error())
		os.Exit(1)
	}
}

// Get returns the global logger instance.
func Get() *zap.Logger {
	if log == nil {
		Init("info", true)
	}
	return log
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
