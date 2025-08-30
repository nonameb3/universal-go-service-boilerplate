package helpers

import (
	"context"
	
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/pkg/types"
)

// MockLogger is a simple mock logger for testing
type MockLogger struct{}

func (m *MockLogger) Info(message string, fields ...types.Field) {
	// No-op for tests
}

func (m *MockLogger) Error(message string, err error, fields ...types.Field) {
	// No-op for tests
}

func (m *MockLogger) Debug(message string, fields ...types.Field) {
	// No-op for tests
}

func (m *MockLogger) Warn(message string, fields ...types.Field) {
	// No-op for tests
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	return m
}

func (m *MockLogger) WithCorrelationID(correlationID string) logger.Logger {
	return m
}

func (m *MockLogger) WithFields(fields ...types.Field) logger.Logger {
	return m
}