// Package universalgo_integration demonstrates how Universal-Go Finance would integrate
// their centralized logging library with the universal Go service boilerplate.

package main

import (
	"context"
	"fmt"

	// Universal-Go's centralized logging library (example import)
	// ugoLogger "github.com/universal-go-finance/universal-go-log-library/ugoLogger"
	// ugoLoggerConstants "github.com/universal-go-finance/universal-go-log-library/constants"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// UniversalGoLogger wraps Universal-Go's centralized logging library to implement
// the universal Logger interface. This allows Universal-Go services to use
// their existing logging infrastructure seamlessly.
type UniversalGoLogger struct {
	// ugoLogger    *logrus.Logger  // Universal-Go's actual logger
	serviceName   string
	correlationID string
	fields        []types.Field
}

// NewUniversalGoLogger creates a new Universal-Go logger that implements the universal Logger interface
func NewUniversalGoLogger(config LoggerConfig) (Logger, error) {
	// In real implementation, this would initialize the Universal-Go logger:
	// universalGoLogger := universalGoLogger.NewUniversalGoLogger(universalGoLogger.LoggerConfig{
	//     ServiceName: config.ServiceName,
	//     LogLevel:    universalGoLoggerConstants.InfoLevel,
	// })

	return &UniversalGoLogger{
		serviceName: config.ServiceName,
	}, nil
}

// Universal Logger interface implementation
// These methods are called by all business logic - they stay the same
// regardless of which company logging library is used underneath.

// Info logs an info message through Universal-Go's centralized system
func (u *UniversalGoLogger) Info(msg string, fields ...types.Field) {
	// Universal-Go implementation would forward to their OpenSearch cluster:
	// u.universalGoLogger.WithFields(convertFieldsToLogrus(fields)).Info(msg)

	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] INFO: %s | service=%s", msg, u.serviceName)
	if u.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", u.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Error logs an error message through Universal-Go's centralized system
func (u *UniversalGoLogger) Error(msg string, err error, fields ...types.Field) {
	// Universal-Go implementation:
	// u.universalGoLogger.WithFields(convertFieldsToLogrus(fields)).WithError(err).Error(msg)

	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] ERROR: %s | service=%s", msg, u.serviceName)
	if u.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", u.correlationID)
	}
	if err != nil {
		fmt.Printf(" | error=%q", err.Error())
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Debug logs a debug message (Universal-Go might filter based on environment)
func (u *UniversalGoLogger) Debug(msg string, fields ...types.Field) {
	// Universal-Go implementation:
	// u.universalGoLogger.WithFields(convertFieldsToLogrus(fields)).Debug(msg)

	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] DEBUG: %s | service=%s", msg, u.serviceName)
	if u.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", u.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Warn logs a warning message
func (u *UniversalGoLogger) Warn(msg string, fields ...types.Field) {
	// Universal-Go implementation:
	// u.universalGoLogger.WithFields(convertFieldsToLogrus(fields)).Warn(msg)

	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] WARN: %s | service=%s", msg, u.serviceName)
	if u.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", u.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// WithContext extracts Universal-Go-specific context values (like trace IDs)
func (u *UniversalGoLogger) WithContext(ctx context.Context) Logger {
	// Universal-Go might extract specific correlation IDs or trace information
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return u.WithCorrelationID(correlationID)
	}
	return u
}

// WithCorrelationID adds Universal-Go's correlation ID for request tracing
func (u *UniversalGoLogger) WithCorrelationID(id string) Logger {
	return &UniversalGoLogger{
		serviceName:   u.serviceName,
		correlationID: id,
		fields:        u.fields,
	}
}

// WithFields adds persistent fields (Universal-Go might add environment, region, etc.)
func (u *UniversalGoLogger) WithFields(fields ...types.Field) Logger {
	newFields := make([]types.Field, len(u.fields)+len(fields))
	copy(newFields, u.fields)
	copy(newFields[len(u.fields):], fields)

	return &UniversalGoLogger{
		serviceName:   u.serviceName,
		correlationID: u.correlationID,
		fields:        newFields,
	}
}

// Helper interfaces to avoid import cycles (normally defined in providers package)
type Logger interface {
	Info(msg string, fields ...types.Field)
	Error(msg string, err error, fields ...types.Field)
	Debug(msg string, fields ...types.Field)
	Warn(msg string, fields ...types.Field)
	WithContext(ctx context.Context) Logger
	WithCorrelationID(id string) Logger
	WithFields(fields ...types.Field) Logger
}

type LoggerConfig struct {
	Type        string
	Level       string
	ServiceName string
	Format      string
}

// Example: How Universal-Go would register and use their logger

func main() {
	fmt.Println("🏢 Universal-Go Logger Integration Example")
	fmt.Println("=====================================\n")

	// 1. Universal-Go registers their logger implementation
	fmt.Println("1️⃣ Register Universal-Go logger implementation:")
	fmt.Println("   providers.RegisterCustomLogger(\"universal-go\", NewUniversalGoLogger)")

	// 2. Update configuration to use Universal-Go logger
	fmt.Println("\n2️⃣ Update config/environments/production.yaml:")
	fmt.Println(`   providers:
     logger:
       type: "universal-go"
       service_name: "user-service"`)

	// 3. Business logic uses the universal interface - no changes needed!
	fmt.Println("\n3️⃣ Business logic remains unchanged:")
	config := LoggerConfig{
		Type:        "universal-go",
		ServiceName: "user-service",
	}

	logger, _ := NewUniversalGoLogger(config)

	// This is how your business logic would use logging - same interface,
	// but now it goes through Universal-Go's centralized system to OpenSearch!
	logger = logger.WithCorrelationID("req-abc-123")

	logger.Info("User login successful",
		types.Field{Key: "user_id", Value: "user-456"},
		types.Field{Key: "method", Value: "oauth2"},
		types.Field{Key: "duration_ms", Value: 234})

	logger.Error("Database connection failed", fmt.Errorf("connection timeout"),
		types.Field{Key: "database", Value: "users-db"},
		types.Field{Key: "retry_count", Value: 3})

	fmt.Println("\n✅ Benefits for Universal-Go:")
	fmt.Printf("   • All logs automatically go to OpenSearch\n")
	fmt.Printf("   • Consistent format across all Go services\n")
	fmt.Printf("   • No code changes in business logic\n")
	fmt.Printf("   • Easy to switch back to simple logger for testing\n")
	fmt.Printf("   • Correlation IDs work across microservices\n")

	fmt.Println("\n🔄 Easy Migration:")
	fmt.Printf("   • Deploy with simple logger first\n")
	fmt.Printf("   • Change config to use Universal-Go logger\n")
	fmt.Printf("   • Restart service - no code deployment needed!\n")
}
