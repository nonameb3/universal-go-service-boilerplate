// Package integration demonstrates how Centralized Logging would integrate
// their centralized logging library with the universal Go service boilerplate.

package logger

import (
	"context"
	"fmt"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// CentralizedLogger wraps Centralized Logging's centralized logging library to implement
// the universal Logger interface. This allows Centralized Logging services to use
// their existing logging infrastructure seamlessly.
type CentralizedLogger struct {
	serviceName   string
	correlationID string
	fields        []types.Field
}

// NewCentralizedLogger creates a new Centralized Logging logger that implements the universal Logger interface
func NewCentralizedLogger(config LoggerConfig) Logger {

	return &CentralizedLogger{
		serviceName: config.ServiceName,
	}
}

// Universal Logger interface implementation
// These methods are called by all business logic - they stay the same
// regardless of which company logging library is used underneath.

// Info logs an info message through Centralized Logging's centralized system
func (t *CentralizedLogger) Info(msg string, fields ...types.Field) {
	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] INFO: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Error logs an error message through Centralized Logging's centralized system
func (t *CentralizedLogger) Error(msg string, err error, fields ...types.Field) {
	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] ERROR: %s | service=%s", msg, t.serviceName)
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

// Debug logs a debug message (Centralized Logging might filter based on environment)
func (t *CentralizedLogger) Debug(msg string, fields ...types.Field) {
	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] DEBUG: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// Warn logs a warning message
func (t *CentralizedLogger) Warn(msg string, fields ...types.Field) {
	// Demo implementation:
	fmt.Printf("[Universal-Go Logger] WARN: %s | service=%s", msg, t.serviceName)
	if t.correlationID != "" {
		fmt.Printf(" | correlation_id=%s", t.correlationID)
	}
	for _, field := range fields {
		fmt.Printf(" | %s=%v", field.Key, field.Value)
	}
	fmt.Println()
}

// WithContext extracts Centralized Logging-specific context values (like trace IDs)
func (t *CentralizedLogger) WithContext(ctx context.Context) Logger {
	// Centralized Logging might extract specific correlation IDs or trace information
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return t.WithCorrelationID(correlationID)
	}
	return t
}

// WithCorrelationID adds Centralized Logging's correlation ID for request tracing
func (t *CentralizedLogger) WithCorrelationID(id string) Logger {
	return &CentralizedLogger{
		serviceName:   t.serviceName,
		correlationID: id,
		fields:        t.fields,
	}
}

// WithFields adds persistent fields (Centralized Logging might add environment, region, etc.)
func (t *CentralizedLogger) WithFields(fields ...types.Field) Logger {
	newFields := make([]types.Field, len(t.fields)+len(fields))
	copy(newFields, t.fields)
	copy(newFields[len(t.fields):], fields)

	return &CentralizedLogger{
		serviceName:   t.serviceName,
		correlationID: t.correlationID,
		fields:        newFields,
	}
}

// Example: How Universal-Go would register and use their logger

// func main() {
// 	fmt.Println("🏢 Universal-Go Logger Integration Example")
// 	fmt.Println("=====================================\n")

// 	// 1. Universal-Go registers their logger implementation
// 	fmt.Println("1️⃣ Register Universal-Go logger implementation:")
// 	fmt.Println("   providers.RegisterCustomLogger(\"universal-go\", NewUniversal-GoLogger)")

// 	// 2. Update configuration to use Universal-Go logger
// 	fmt.Println("\n2️⃣ Update config/environments/production.yaml:")
// 	fmt.Println(`   providers:
//      logger:
//        type: "universal-go"
//        service_name: "user-service"`)

// 	// 3. Business logic uses the universal interface - no changes needed!
// 	fmt.Println("\n3️⃣ Business logic remains unchanged:")
// 	config := LoggerConfig{
// 		Type:        "universal-go",
// 		ServiceName: "user-service",
// 	}

// 	logger, _ := NewUniversal-GoLogger(config)

// 	// This is how your business logic would use logging - same interface,
// 	// but now it goes through Universal-Go's centralized system to OpenSearch!
// 	logger = logger.WithCorrelationID("req-abc-123")

// 	logger.Info("User login successful",
// 		types.Field{Key: "user_id", Value: "user-456"},
// 		types.Field{Key: "method", Value: "oauth2"},
// 		types.Field{Key: "duration_ms", Value: 234})

// 	logger.Error("Database connection failed", fmt.Errorf("connection timeout"),
// 		types.Field{Key: "database", Value: "users-db"},
// 		types.Field{Key: "retry_count", Value: 3})

// 	fmt.Println("\n✅ Benefits for Universal-Go:")
// 	fmt.Printf("   • All logs automatically go to OpenSearch\n")
// 	fmt.Printf("   • Consistent format across all Go services\n")
// 	fmt.Printf("   • No code changes in business logic\n")
// 	fmt.Printf("   • Easy to switch back to simple logger for testing\n")
// 	fmt.Printf("   • Correlation IDs work across microservices\n")

// 	fmt.Println("\n🔄 Easy Migration:")
// 	fmt.Printf("   • Deploy with simple logger first\n")
// 	fmt.Printf("   • Change config to use Universal-Go logger\n")
// 	fmt.Printf("   • Restart service - no code deployment needed!\n")
// }
