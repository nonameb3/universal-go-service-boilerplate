package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/pkg/types"
)

// MockLogger is a mock implementation of Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string, fields ...types.Field) {
	args := []interface{}{message}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Error(message string, err error, fields ...types.Field) {
	args := []interface{}{message, err}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Debug(message string, fields ...types.Field) {
	args := []interface{}{message}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) Warn(message string, fields ...types.Field) {
	args := []interface{}{message}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithCorrelationID(correlationID string) logger.Logger {
	args := m.Called(correlationID)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithFields(fields ...types.Field) logger.Logger {
	args := make([]interface{}, len(fields))
	for i, field := range fields {
		args[i] = field
	}
	result := m.Called(args...)
	return result.Get(0).(logger.Logger)
}