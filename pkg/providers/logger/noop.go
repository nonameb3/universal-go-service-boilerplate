package logger

import (
	"context"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// noopLogger is a logger that doesn't log anything - useful for testing
type noopLogger struct{}

// NewNoop creates a new no-op logger
func NewNoop(config LoggerConfig) (Logger, error) {
	return &noopLogger{}, nil
}

// Info does nothing
func (l *noopLogger) Info(msg string, fields ...types.Field) {}

// Error does nothing
func (l *noopLogger) Error(msg string, err error, fields ...types.Field) {}

// Debug does nothing
func (l *noopLogger) Debug(msg string, fields ...types.Field) {}

// Warn does nothing
func (l *noopLogger) Warn(msg string, fields ...types.Field) {}

// WithContext returns the same no-op logger
func (l *noopLogger) WithContext(ctx context.Context) Logger {
	return l
}

// WithCorrelationID returns the same no-op logger
func (l *noopLogger) WithCorrelationID(id string) Logger {
	return l
}

// WithFields returns the same no-op logger
func (l *noopLogger) WithFields(fields ...types.Field) Logger {
	return l
}