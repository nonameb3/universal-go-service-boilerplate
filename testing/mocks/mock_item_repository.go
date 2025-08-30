package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"gorm.io/gorm"
)

// MockItemRepository is a mock implementation of ItemRepository
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(item *entities.Item) (*entities.Item, error) {
	args := m.Called(item)
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) CreateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error) {
	args := m.Called(tx, item)
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) Get(id string) (*entities.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetWithTx(tx *gorm.DB, id string) (*entities.Item, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByName(name string) (*entities.Item, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByNameWithTx(tx *gorm.DB, name string) (*entities.Item, error) {
	args := m.Called(tx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByNameForUpdate(tx *gorm.DB, name string) (*entities.Item, error) {
	args := m.Called(tx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByNames(names []string) ([]*entities.Item, error) {
	args := m.Called(names)
	return args.Get(0).([]*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByNamesWithTx(tx *gorm.DB, names []string) ([]*entities.Item, error) {
	args := m.Called(tx, names)
	return args.Get(0).([]*entities.Item), args.Error(1)
}

func (m *MockItemRepository) Update(item *entities.Item) (*entities.Item, error) {
	args := m.Called(item)
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) UpdateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error) {
	args := m.Called(tx, item)
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetWithPagination(page, limit int) (*types.PaginatedResult[*entities.Item], error) {
	args := m.Called(page, limit)
	return args.Get(0).(*types.PaginatedResult[*entities.Item]), args.Error(1)
}

func (m *MockItemRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockItemRepository) DeleteWithTx(tx *gorm.DB, id string) error {
	args := m.Called(tx, id)
	return args.Error(0)
}