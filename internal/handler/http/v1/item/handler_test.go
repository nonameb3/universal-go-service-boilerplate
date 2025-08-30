package item

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/testing/fixtures"
)

// MockItemUseCase for testing handlers
type MockItemUseCase struct {
	mock.Mock
}

func (m *MockItemUseCase) Create(req *dto.CreateItemRequest) (*entities.Item, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemUseCase) Get(id string) (*entities.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemUseCase) Update(id string, req *dto.UpdateItemRequest) (*entities.Item, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemUseCase) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockItemUseCase) GetWithPagination(req *dto.PaginationRequest) (*types.PaginatedResult[*entities.Item], error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.PaginatedResult[*entities.Item]), args.Error(1)
}

func (m *MockItemUseCase) BulkCreate(req *dto.BulkCreateRequest) ([]*entities.Item, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Item), args.Error(1)
}

func TestHandler_CreateItem(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Post("/items", handler.CreateItem)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "should create item successfully",
			requestBody: request.AddItem{
				Name:   "Test Item",
				Amount: 100,
			},
			mockSetup: func() {
				mockUseCase.On("Create", mock.MatchedBy(func(req *dto.CreateItemRequest) bool {
					return req.Name == "Test Item" && req.Amount == 100
				})).Return(fixtures.ValidItemWithName("Test Item"), nil)
			},
			expectedStatus: 201,
		},
		{
			name: "should return 400 for invalid request",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup: func() {
				// The JSON is valid but creates empty values, which should trigger validation error
				mockUseCase.On("Create", mock.MatchedBy(func(req *dto.CreateItemRequest) bool {
					return req.Name == "" && req.Amount == 0
				})).Return(nil, domain.ErrItemNameRequired)
			},
			expectedStatus: 400,
		},
		{
			name: "should return 409 for duplicate item",
			requestBody: request.AddItem{
				Name:   "Duplicate Item",
				Amount: 100,
			},
			mockSetup: func() {
				mockUseCase.On("Create", mock.AnythingOfType("*dto.CreateItemRequest")).Return(nil, domain.ErrItemAlreadyExists)
			},
			expectedStatus: 409,
		},
		{
			name: "should return 400 for validation error",
			requestBody: request.AddItem{
				Name:   "",
				Amount: 100,
			},
			mockSetup: func() {
				mockUseCase.On("Create", mock.AnythingOfType("*dto.CreateItemRequest")).Return(nil, domain.ErrItemNameRequired)
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			tt.mockSetup()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/items", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestHandler_GetItem(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Get("/items/:id", handler.GetItem)

	tests := []struct {
		name           string
		itemID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "should get existing item",
			itemID: "existing-id",
			mockSetup: func() {
				item := fixtures.ValidItemWithName("Found Item")
				mockUseCase.On("Get", "existing-id").Return(item, nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "should return 404 for non-existent item",
			itemID: "non-existent-id",
			mockSetup: func() {
				mockUseCase.On("Get", "non-existent-id").Return(nil, domain.ErrItemNotFound)
			},
			expectedStatus: 404,
		},
		{
			name:           "should return 404 for missing ID route",
			itemID:         "",
			mockSetup:      func() {},
			expectedStatus: 404, // Fiber returns 404 for unmatched routes like /items/
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			tt.mockSetup()

			url := "/items/" + tt.itemID
			if tt.itemID == "" {
				url = "/items/"
			}
			
			req := httptest.NewRequest("GET", url, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			
			if len(mockUseCase.ExpectedCalls) > 0 {
				mockUseCase.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_ListItems(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Get("/items", handler.ListItems)

	t.Run("should list items with pagination", func(t *testing.T) {
		items := fixtures.ValidItems(3)
		paginatedResult := &types.PaginatedResult[*entities.Item]{
			Items:      items,
			Total:      3,
			Page:       1,
			Limit:      10,
			TotalPages: 1,
		}

		mockUseCase.On("GetWithPagination", mock.MatchedBy(func(req *dto.PaginationRequest) bool {
			return req.Page == 1 && req.Limit == 10
		})).Return(paginatedResult, nil)

		req := httptest.NewRequest("GET", "/items?page=1&limit=10", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid pagination", func(t *testing.T) {
		mockUseCase.ExpectedCalls = nil
		mockUseCase.On("GetWithPagination", mock.AnythingOfType("*dto.PaginationRequest")).Return(nil, domain.ErrInvalidPagination)

		req := httptest.NewRequest("GET", "/items?page=-1&limit=10", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestHandler_UpdateItem(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Put("/items/:id", handler.UpdateItem)

	tests := []struct {
		name           string
		itemID         string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "should update item successfully",
			itemID: "item-id",
			requestBody: request.UpdateItem{
				Name:   strPtr("Updated Item"),
				Amount: uintPtr(200),
			},
			mockSetup: func() {
				updatedItem := fixtures.ValidItemWithName("Updated Item")
				updatedItem.Amount = 200
				mockUseCase.On("Update", "item-id", mock.MatchedBy(func(req *dto.UpdateItemRequest) bool {
					return *req.Name == "Updated Item" && *req.Amount == 200
				})).Return(updatedItem, nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "should return 404 for non-existent item",
			itemID: "non-existent-id",
			requestBody: request.UpdateItem{
				Name: strPtr("Updated Item"),
			},
			mockSetup: func() {
				mockUseCase.On("Update", "non-existent-id", mock.AnythingOfType("*dto.UpdateItemRequest")).Return(nil, domain.ErrItemNotFound)
			},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			tt.mockSetup()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/items/"+tt.itemID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteItem(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Delete("/items/:id", handler.DeleteItem)

	tests := []struct {
		name           string
		itemID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "should delete item successfully",
			itemID: "item-id",
			mockSetup: func() {
				mockUseCase.On("Delete", "item-id").Return(nil)
			},
			expectedStatus: 200,
		},
		{
			name:   "should return 404 for non-existent item",
			itemID: "non-existent-id",
			mockSetup: func() {
				mockUseCase.On("Delete", "non-existent-id").Return(domain.ErrItemNotFound)
			},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			tt.mockSetup()

			req := httptest.NewRequest("DELETE", "/items/"+tt.itemID, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestHandler_BulkCreateItems(t *testing.T) {
	app := fiber.New()
	mockUseCase := &MockItemUseCase{}
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	
	handler := New(mockUseCase, noopLogger)
	app.Post("/items/bulk", handler.BulkCreateItems)

	t.Run("should create multiple items successfully", func(t *testing.T) {
		requestBody := request.BulkCreateItems{
			Items: []request.AddItem{
				{Name: "Bulk Item 1", Amount: 100},
				{Name: "Bulk Item 2", Amount: 200},
			},
		}

		createdItems := []*entities.Item{
			fixtures.ValidItemWithName("Bulk Item 1"),
			fixtures.ValidItemWithName("Bulk Item 2"),
		}

		mockUseCase.On("BulkCreate", mock.MatchedBy(func(req *dto.BulkCreateRequest) bool {
			return len(req.Items) == 2 && req.Items[0].Name == "Bulk Item 1"
		})).Return(createdItems, nil)

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/items/bulk", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 201, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid request", func(t *testing.T) {
		mockUseCase.ExpectedCalls = nil
		
		req := httptest.NewRequest("POST", "/items/bulk", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}