package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// Logger interface - defined locally to avoid import cycle
type Logger interface {
	Info(msg string, fields ...types.Field)
	Error(msg string, err error, fields ...types.Field)
	Debug(msg string, fields ...types.Field)
	Warn(msg string, fields ...types.Field)
	WithContext(ctx context.Context) Logger
	WithCorrelationID(id string) Logger
	WithFields(fields ...types.Field) Logger
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Type        string            `yaml:"type"`
	Level       types.LogLevel    `yaml:"level"`
	ServiceName string            `yaml:"service_name"`
	Format      string            `yaml:"format"` // json, text
	Output      io.Writer         `yaml:"-"`
	Fields      map[string]string `yaml:"fields"`
}

// simpleLogger is a basic logger implementation using Go's standard log package
type simpleLogger struct {
	logger        *log.Logger
	level         types.LogLevel
	serviceName   string
	correlationID string
	fields        []types.Field
}

// NewSimple creates a new simple logger
func NewSimple(config LoggerConfig) (Logger, error) {
	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	logger := log.New(output, "", log.LstdFlags)

	return &simpleLogger{
		logger:      logger,
		level:       config.Level,
		serviceName: config.ServiceName,
	}, nil
}

// Info logs an info message
func (l *simpleLogger) Info(msg string, fields ...types.Field) {
	if l.shouldLog(types.InfoLevel) {
		l.log("INFO", msg, nil, fields...)
	}
}

// Error logs an error message
func (l *simpleLogger) Error(msg string, err error, fields ...types.Field) {
	if l.shouldLog(types.ErrorLevel) {
		l.log("ERROR", msg, err, fields...)
	}
}

// Debug logs a debug message
func (l *simpleLogger) Debug(msg string, fields ...types.Field) {
	if l.shouldLog(types.DebugLevel) {
		l.log("DEBUG", msg, nil, fields...)
	}
}

// Warn logs a warning message
func (l *simpleLogger) Warn(msg string, fields ...types.Field) {
	if l.shouldLog(types.WarnLevel) {
		l.log("WARN", msg, nil, fields...)
	}
}

// WithContext returns a logger with context (simple implementation ignores context)
func (l *simpleLogger) WithContext(ctx context.Context) Logger {
	// Extract correlation ID from context if present
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return l.WithCorrelationID(correlationID)
	}
	return l
}

// WithCorrelationID returns a logger with correlation ID
func (l *simpleLogger) WithCorrelationID(id string) Logger {
	return &simpleLogger{
		logger:        l.logger,
		level:         l.level,
		serviceName:   l.serviceName,
		correlationID: id,
		fields:        l.fields,
	}
}

// WithFields returns a logger with additional fields
func (l *simpleLogger) WithFields(fields ...types.Field) Logger {
	newFields := make([]types.Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)

	return &simpleLogger{
		logger:        l.logger,
		level:         l.level,
		serviceName:   l.serviceName,
		correlationID: l.correlationID,
		fields:        newFields,
	}
}

// log performs the actual logging
func (l *simpleLogger) log(level, msg string, err error, fields ...types.Field) {
	// Build log message
	var parts []string

	// Add service name if present
	if l.serviceName != "" {
		parts = append(parts, fmt.Sprintf("service=%s", l.serviceName))
	}

	// Add correlation ID if present
	if l.correlationID != "" {
		parts = append(parts, fmt.Sprintf("correlation_id=%s", l.correlationID))
	}

	// Add persistent fields
	for _, field := range l.fields {
		parts = append(parts, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	// Add current fields
	for _, field := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	// Add error if present
	if err != nil {
		parts = append(parts, fmt.Sprintf("error=%q", err.Error()))
	}

	// Combine all parts
	fieldsStr := ""
	if len(parts) > 0 {
		fieldsStr = " " + strings.Join(parts, " ")
	}

	// Log the message
	l.logger.Printf("[%s] %s%s", level, msg, fieldsStr)
}

// shouldLog determines if a message should be logged based on level
func (l *simpleLogger) shouldLog(msgLevel types.LogLevel) bool {
	levelPriority := map[types.LogLevel]int{
		types.DebugLevel: 0,
		types.InfoLevel:  1,
		types.WarnLevel:  2,
		types.ErrorLevel: 3,
		types.FatalLevel: 4,
	}

	configPriority, ok := levelPriority[l.level]
	if !ok {
		configPriority = levelPriority[types.InfoLevel]
	}

	msgPriority, ok := levelPriority[msgLevel]
	if !ok {
		msgPriority = levelPriority[types.InfoLevel]
	}

	return msgPriority >= configPriority
}

// Field helper functions for convenience

// StringField creates a string field
func StringField(key, value string) types.Field {
	return types.Field{Key: key, Value: value}
}

// IntField creates an int field
func IntField(key string, value int) types.Field {
	return types.Field{Key: key, Value: value}
}

// Int64Field creates an int64 field
func Int64Field(key string, value int64) types.Field {
	return types.Field{Key: key, Value: value}
}

// Float64Field creates a float64 field
func Float64Field(key string, value float64) types.Field {
	return types.Field{Key: key, Value: value}
}

// BoolField creates a bool field
func BoolField(key string, value bool) types.Field {
	return types.Field{Key: key, Value: value}
}

// DurationField creates a duration field
func DurationField(key string, value time.Duration) types.Field {
	return types.Field{Key: key, Value: value.String()}
}

// TimeField creates a time field
func TimeField(key string, value time.Time) types.Field {
	return types.Field{Key: key, Value: value.Format(time.RFC3339)}
}

// ErrorField creates an error field
func ErrorField(err error) types.Field {
	if err == nil {
		return types.Field{Key: "error", Value: nil}
	}
	return types.Field{Key: "error", Value: err.Error()}
}