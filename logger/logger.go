package logger

import (
	"go.uber.org/zap"
)

// ZapLogger wraps the Zap sugared logger.
type ZapLogger struct {
	SugaredLogger *zap.SugaredLogger
}

// NewZapLogger initializes a new ZapLogger with a sugared logger.
func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{
		SugaredLogger: logger.Sugar(),
	}, nil
}

// Define logging methods
func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.SugaredLogger.Infof(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.SugaredLogger.Errorf(msg, args...)
}

// Debug logs a debug-level message with additional fields.
func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.SugaredLogger.Debugf(msg, args...)
}

// Warn logs a warning-level message with additional fields.
func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.SugaredLogger.Warnf(msg, args...)
}

// Fatal logs a fatal-level message with additional fields and then exits the application.
func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalf(msg, args...)
}
