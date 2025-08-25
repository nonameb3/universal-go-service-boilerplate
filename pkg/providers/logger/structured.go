package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// structuredLogger uses Go's structured logging (slog)
type structuredLogger struct {
	logger        *slog.Logger
	serviceName   string
	correlationID string
	fields        []types.Field
}

// NewStructured creates a new structured logger using slog
func NewStructured(config LoggerConfig) (Logger, error) {
	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	var handler slog.Handler
	level := slogLevel(config.Level)

	// Choose handler based on format
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if config.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)

	// Add service name to all logs if configured
	if config.ServiceName != "" {
		logger = logger.With("service", config.ServiceName)
	}

	// Add any configured fields
	for key, value := range config.Fields {
		logger = logger.With(key, value)
	}

	return &structuredLogger{
		logger:      logger,
		serviceName: config.ServiceName,
	}, nil
}

// Info logs an info message
func (l *structuredLogger) Info(msg string, fields ...types.Field) {
	l.log(context.Background(), slog.LevelInfo, msg, nil, fields...)
}

// Error logs an error message
func (l *structuredLogger) Error(msg string, err error, fields ...types.Field) {
	l.log(context.Background(), slog.LevelError, msg, err, fields...)
}

// Debug logs a debug message
func (l *structuredLogger) Debug(msg string, fields ...types.Field) {
	l.log(context.Background(), slog.LevelDebug, msg, nil, fields...)
}

// Warn logs a warning message
func (l *structuredLogger) Warn(msg string, fields ...types.Field) {
	l.log(context.Background(), slog.LevelWarn, msg, nil, fields...)
}

// WithContext returns a logger with context
func (l *structuredLogger) WithContext(ctx context.Context) Logger {
	// Extract correlation ID from context if present
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return l.WithCorrelationID(correlationID)
	}
	return l
}

// WithCorrelationID returns a logger with correlation ID
func (l *structuredLogger) WithCorrelationID(id string) Logger {
	newLogger := l.logger.With("correlation_id", id)
	
	return &structuredLogger{
		logger:        newLogger,
		serviceName:   l.serviceName,
		correlationID: id,
		fields:        l.fields,
	}
}

// WithFields returns a logger with additional fields
func (l *structuredLogger) WithFields(fields ...types.Field) Logger {
	// Convert fields to slog attributes
	attrs := make([]any, 0, len(fields)*2)
	for _, field := range fields {
		attrs = append(attrs, field.Key, field.Value)
	}

	newLogger := l.logger.With(attrs...)
	
	// Combine fields for the new logger instance
	newFields := make([]types.Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)

	return &structuredLogger{
		logger:        newLogger,
		serviceName:   l.serviceName,
		correlationID: l.correlationID,
		fields:        newFields,
	}
}

// log performs the actual logging
func (l *structuredLogger) log(ctx context.Context, level slog.Level, msg string, err error, fields ...types.Field) {
	// Convert fields to slog attributes
	attrs := make([]slog.Attr, 0, len(fields))
	
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}

	// Add error as attribute if present
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	// Log with attributes
	l.logger.LogAttrs(ctx, level, msg, attrs...)
}

// slogLevel converts our LogLevel to slog.Level
func slogLevel(level types.LogLevel) slog.Level {
	switch level {
	case types.DebugLevel:
		return slog.LevelDebug
	case types.InfoLevel:
		return slog.LevelInfo
	case types.WarnLevel:
		return slog.LevelWarn
	case types.ErrorLevel:
		return slog.LevelError
	case types.FatalLevel:
		return slog.LevelError // slog doesn't have fatal, use error
	default:
		return slog.LevelInfo
	}
}