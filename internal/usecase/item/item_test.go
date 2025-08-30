package item

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/testing/fixtures"
	"github.com/universal-go-service/boilerplate/testing/mocks"
	"gorm.io/gorm"
)

func TestItemUseCase_Create(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	tests := []struct {
		name          string
		request       *dto.CreateItemRequest
		mockSetup     func()
		expectedError error
		expectedName  string
	}{
		{
			name: "should create item successfully",
			request: &dto.CreateItemRequest{
				Name:   "Test Item",
				Amount: 100,
			},
			mockSetup: func() {
				// Mock validation (no existing item)
				mockRepo.On("GetByNameForUpdate", mock.Anything, "Test Item").Return(nil, gorm.ErrRecordNotFound)
				
				// Mock successful creation
				expectedItem := fixtures.ValidItemWithName("Test Item")
				expectedItem.Amount = 100
				mockRepo.On("CreateWithTx", mock.Anything, mock.MatchedBy(func(item *entities.Item) bool {
					return item.Name == "Test Item" && item.Amount == 100
				})).Return(expectedItem, nil)

				// Mock transaction
				mockDB.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(0).(func(*gorm.DB) error)
					fn(&gorm.DB{}) // Call with mock DB
				})
			},
			expectedError: nil,
			expectedName:  "Test Item",
		},
		{
			name: "should fail on invalid name",
			request: &dto.CreateItemRequest{
				Name:   "", // Empty name
				Amount: 100,
			},
			mockSetup:     func() {}, // No mocks needed for validation error
			expectedError: domain.ErrItemNameRequired,
		},
		{
			name: "should fail on duplicate name",
			request: &dto.CreateItemRequest{
				Name:   "Duplicate Item",
				Amount: 100,
			},
			mockSetup: func() {
				// Mock existing item found
				existingItem := fixtures.ValidItemWithName("Duplicate Item")
				mockRepo.On("GetByNameForUpdate", mock.Anything, "Duplicate Item").Return(existingItem, nil)

				// Mock transaction
				mockDB.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(domain.ErrItemAlreadyExists).Run(func(args mock.Arguments) {
					fn := args.Get(0).(func(*gorm.DB) error)
					fn(&gorm.DB{}) // This will return the duplicate error
				})
			},
			expectedError: domain.ErrItemAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockDB.ExpectedCalls = nil

			// Setup mocks
			tt.mockSetup()

			// Execute
			result, err := useCase.Create(tt.request)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedName, result.Name)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestItemUseCase_Get(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	tests := []struct {
		name          string
		itemID        string
		mockSetup     func()
		expectedError error
		expectedName  string
	}{
		{
			name:   "should get existing item",
			itemID: "existing-id",
			mockSetup: func() {
				expectedItem := fixtures.ValidItemWithName("Found Item")
				mockRepo.On("Get", "existing-id").Return(expectedItem, nil)
			},
			expectedError: nil,
			expectedName:  "Found Item",
		},
		{
			name:   "should fail for non-existent item",
			itemID: "non-existent-id",
			mockSetup: func() {
				mockRepo.On("Get", "non-existent-id").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: domain.ErrItemNotFound, // UseCase converts to domain error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.mockSetup()

			result, err := useCase.Get(tt.itemID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedName, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUseCase_Update(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	t.Run("should update item successfully", func(t *testing.T) {
		request := &dto.UpdateItemRequest{
			Name:   strPtr("Updated Item"),
			Amount: uintPtr(200),
		}

		existingItem := fixtures.ValidItemWithName("Original Item")
		existingItem.Amount = 100

		updatedItem := fixtures.ValidItemWithName("Updated Item")
		updatedItem.Amount = 200

		// Mock get existing item
		mockRepo.On("Get", "item-id").Return(existingItem, nil)
		
		// Mock duplicate check (should not find any duplicates)
		mockRepo.On("GetByName", "Updated Item").Return(nil, gorm.ErrRecordNotFound)
		
		// Mock update
		mockRepo.On("Update", mock.MatchedBy(func(item *entities.Item) bool {
			return item.Name == "Updated Item" && item.Amount == 200
		})).Return(updatedItem, nil)

		result, err := useCase.Update("item-id", request)

		require.NoError(t, err)
		assert.Equal(t, "Updated Item", result.Name)
		assert.Equal(t, uint(200), result.Amount)

		mockRepo.AssertExpectations(t)
	})
}

func TestItemUseCase_Delete(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	t.Run("should delete item successfully", func(t *testing.T) {
		existingItem := fixtures.ValidItemWithName("Item to Delete")
		
		// Mock get to verify item exists
		mockRepo.On("Get", "item-id").Return(existingItem, nil)
		
		// Mock delete
		mockRepo.On("Delete", "item-id").Return(nil)

		err := useCase.Delete("item-id")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail for non-existent item", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		// Mock get returns error
		mockRepo.On("Get", "non-existent-id").Return(nil, gorm.ErrRecordNotFound)

		err := useCase.Delete("non-existent-id")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrItemNotFound, err) // UseCase converts to domain error
		mockRepo.AssertExpectations(t)
	})
}

func TestItemUseCase_GetWithPagination(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	tests := []struct {
		name          string
		request       *dto.PaginationRequest
		mockSetup     func()
		expectedError error
		expectedPage  int
	}{
		{
			name: "should get paginated items with defaults",
			request: &dto.PaginationRequest{
				Page:  0, // Should default to 1
				Limit: 0, // Should default to 10
			},
			mockSetup: func() {
				items := fixtures.ValidItems(5)
				result := &types.PaginatedResult[*entities.Item]{
					Items:      items,
					Total:      5,
					Page:       1,
					Limit:      10,
					TotalPages: 1,
				}
				mockRepo.On("GetWithPagination", 1, 10).Return(result, nil)
			},
			expectedError: nil,
			expectedPage:  1,
		},
		{
			name: "should clamp large limits and succeed",
			request: &dto.PaginationRequest{
				Page:  1,
				Limit: 200, // Will be clamped to 100 by ApplyDefaults
			},
			mockSetup: func() {
				items := fixtures.ValidItems(10)
				result := &types.PaginatedResult[*entities.Item]{
					Items:      items,
					Total:      50,
					Page:       1,
					Limit:      100, // Clamped value
					TotalPages: 1,
				}
				// Expect call with clamped limit of 100
				mockRepo.On("GetWithPagination", 1, 100).Return(result, nil)
			},
			expectedError: nil,
			expectedPage:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.mockSetup()

			result, err := useCase.GetWithPagination(tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedPage, result.Page)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestItemUseCase_BulkCreate(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{}
	mockDB := &mocks.MockDatabaseProvider{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})

	useCase := NewItemUseCase(mockRepo, mockDB, noopLogger)

	t.Run("should create multiple items successfully", func(t *testing.T) {
		request := &dto.BulkCreateRequest{
			Items: []dto.CreateItemRequest{
				{Name: "Bulk Item 1", Amount: 100},
				{Name: "Bulk Item 2", Amount: 200},
			},
		}

		// Mock batch duplicate check (no duplicates)
		mockRepo.On("GetByNamesWithTx", mock.Anything, []string{"Bulk Item 1", "Bulk Item 2"}).Return([]*entities.Item{}, nil)

		// Mock individual creates
		item1 := fixtures.ValidItemWithName("Bulk Item 1")
		item1.Amount = 100
		item2 := fixtures.ValidItemWithName("Bulk Item 2")
		item2.Amount = 200

		mockRepo.On("CreateWithTx", mock.Anything, mock.MatchedBy(func(item *entities.Item) bool {
			return item.Name == "Bulk Item 1"
		})).Return(item1, nil)

		mockRepo.On("CreateWithTx", mock.Anything, mock.MatchedBy(func(item *entities.Item) bool {
			return item.Name == "Bulk Item 2"
		})).Return(item2, nil)

		// Mock transaction
		mockDB.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(*gorm.DB) error)
			fn(&gorm.DB{})
		})

		result, err := useCase.BulkCreate(request)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Bulk Item 1", result[0].Name)
		assert.Equal(t, "Bulk Item 2", result[1].Name)

		mockRepo.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}