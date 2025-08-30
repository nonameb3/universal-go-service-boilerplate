package mocks

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabaseProvider is a mock implementation of DatabaseProvider
type MockDatabaseProvider struct {
	mock.Mock
}

func (m *MockDatabaseProvider) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return &gorm.DB{}
	}
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabaseProvider) GetSQLDB() *sql.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sql.DB)
}

func (m *MockDatabaseProvider) Health() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabaseProvider) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabaseProvider) Migrate(models ...interface{}) error {
	args := m.Called(models)
	return args.Error(0)
}

func (m *MockDatabaseProvider) Transaction(fn func(*gorm.DB) error) error {
	args := m.Called(fn)
	return args.Error(0)
}