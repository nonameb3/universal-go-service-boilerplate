// Package tokenx_integration demonstrates how TokenX Finance would integrate
// their centralized logging library with the universal Go service boilerplate.

package main

import (
	"context"
	"fmt"

	// TokenX's centralized logging library (example import)
	// tkxLogger "github.com/tokenx-finance/tkx-golang-log-library/tkxLogger"
	// tkxLoggerConstants "github.com/tokenx-finance/tkx-golang-log-library/constants"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// TokenXLogger wraps TokenX's centralized logging library to implement
// the universal Logger interface. This allows TokenX services to use
// their existing logging infrastructure seamlessly.
type TokenXLogger struct {
	// tkxLogger    *logrus.Logger  // TokenX's actual logger
	serviceName   string
	correlationID string
	fields        []types.Field
}

// NewTokenXLogger creates a new TokenX logger that implements the universal Logger interface
func NewTokenXLogger(config LoggerConfig) (Logger, error) {
	// In real implementation, this would initialize the TokenX logger:
	// tkxLogger := tkxLogger.NewTkxLogger(tkxLogger.LoggerConfig{
	//     ServiceName: config.ServiceName,
	//     LogLevel:    tkxLoggerConstants.InfoLevel,
	// })

	return &TokenXLogger{
		serviceName: config.ServiceName,
	}, nil
}

// Universal Logger interface implementation
// These methods are called by all business logic - they stay the same
// regardless of which company logging library is used underneath.

// Info logs an info message through TokenX's centralized system
func (t *TokenXLogger) Info(msg string, fields ...types.Field) {
	// TokenX implementation would forward to their OpenSearch cluster:
	// t.tkxLogger.WithFields(convertFieldsToLogrus(fields)).Info(msg)
	
	// Demo implementation:
	fmt.Printf("[TokenX Logger] INFO: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Error logs an error message through TokenX's centralized system
func (t *TokenXLogger) Error(msg string, err error, fields ...types.Field) {
	// TokenX implementation:
	// t.tkxLogger.WithFields(convertFieldsToLogrus(fields)).WithError(err).Error(msg)
	
	// Demo implementation:
	fmt.Printf("[TokenX Logger] ERROR: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	if err != nil {
		fmt.Printf(" | error=%q", err.Error())
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Debug logs a debug message (TokenX might filter based on environment)
func (t *TokenXLogger) Debug(msg string, fields ...types.Field) {
	// TokenX implementation:
	// t.tkxLogger.WithFields(convertFieldsToLogrus(fields)).Debug(msg)
	
	// Demo implementation:
	fmt.Printf("[TokenX Logger] DEBUG: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Warn logs a warning message
func (t *TokenXLogger) Warn(msg string, fields ...types.Field) {
	// TokenX implementation:
	// t.tkxLogger.WithFields(convertFieldsToLogrus(fields)).Warn(msg)
	
	// Demo implementation:
	fmt.Printf("[TokenX Logger] WARN: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// WithContext extracts TokenX-specific context values (like trace IDs)
func (t *TokenXLogger) WithContext(ctx context.Context) Logger {
	// TokenX might extract specific correlation IDs or trace information
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return t.WithCorrelationID(correlationID)
	}
	return t
}

// WithCorrelationID adds TokenX's correlation ID for request tracing
func (t *TokenXLogger) WithCorrelationID(id string) Logger {
	return &TokenXLogger{
		serviceName:   t.serviceName,
		correlationID: id,
		fields:        t.fields,
	}
}

// WithFields adds persistent fields (TokenX might add environment, region, etc.)
func (t *TokenXLogger) WithFields(fields ...types.Field) Logger {
	newFields := make([]types.Field, len(t.fields)+len(fields))
	copy(newFields, t.fields)
	copy(newFields[len(t.fields):], fields)

	return &TokenXLogger{
		serviceName:   t.serviceName,
		correlationID: t.correlationID,
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

// Example: How TokenX would register and use their logger

func main() {
	fmt.Println("🏢 TokenX Logger Integration Example")
	fmt.Println("=====================================\n")

	// 1. TokenX registers their logger implementation
	fmt.Println("1️⃣ Register TokenX logger implementation:")
	fmt.Println("   providers.RegisterCustomLogger(\"tokenx\", NewTokenXLogger)")

	// 2. Update configuration to use TokenX logger
	fmt.Println("\n2️⃣ Update config/environments/production.yaml:")
	fmt.Println(`   providers:
     logger:
       type: "tokenx"
       service_name: "user-service"`)

	// 3. Business logic uses the universal interface - no changes needed!
	fmt.Println("\n3️⃣ Business logic remains unchanged:")
	config := LoggerConfig{
		Type:        "tokenx",
		ServiceName: "user-service",
	}
	
	logger, _ := NewTokenXLogger(config)
	
	// This is how your business logic would use logging - same interface,
	// but now it goes through TokenX's centralized system to OpenSearch!
	logger = logger.WithCorrelationID("req-abc-123")
	
	logger.Info("User login successful", 
		types.Field{Key: "user_id", Value: "user-456"},
		types.Field{Key: "method", Value: "oauth2"},
		types.Field{Key: "duration_ms", Value: 234})

	logger.Error("Database connection failed", fmt.Errorf("connection timeout"),
		types.Field{Key: "database", Value: "users-db"},
		types.Field{Key: "retry_count", Value: 3})

	fmt.Println("\n✅ Benefits for TokenX:")
	fmt.Printf("   • All logs automatically go to OpenSearch\n")
	fmt.Printf("   • Consistent format across all Go services\n")
	fmt.Printf("   • No code changes in business logic\n") 
	fmt.Printf("   • Easy to switch back to simple logger for testing\n")
	fmt.Printf("   • Correlation IDs work across microservices\n")

	fmt.Println("\n🔄 Easy Migration:")
	fmt.Printf("   • Deploy with simple logger first\n")
	fmt.Printf("   • Change config to use TokenX logger\n")
	fmt.Printf("   • Restart service - no code deployment needed!\n")
}