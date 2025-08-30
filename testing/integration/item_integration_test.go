package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/handler/http"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
	"github.com/universal-go-service/boilerplate/internal/repository/item"
	itemUC "github.com/universal-go-service/boilerplate/internal/usecase/item"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/testing/helpers"
)

type ItemIntegrationTestSuite struct {
	suite.Suite
	app    *fiber.App
	testDB *helpers.TestDatabase
	logger logger.Logger
}

func (s *ItemIntegrationTestSuite) SetupSuite() {
	// Setup test database - FAIL if not available for integration tests
	s.testDB = helpers.SetupTestDB(s.T())
	if s.testDB == nil {
		s.T().Fatalf("Integration tests require test database - ensure PostgreSQL is running with correct credentials")
	}
	
	// Setup logger
	noopLogger, err := logger.NewNoop(logger.LoggerConfig{})
	require.NoError(s.T(), err)
	s.logger = noopLogger
	
	// Setup the FULL application stack (mimicking app.Run)
	s.app = fiber.New()
	
	// Setup health check middleware
	s.app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/health",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return s.testDB.Provider.Health() == nil
		},
		ReadinessEndpoint: "/readiness",
	}))
	
	// Setup full dependency injection chain
	itemRepository := item.NewItemRepository(s.testDB.DB, s.logger)
	itemUseCase := itemUC.NewItemUseCase(itemRepository, s.testDB.Provider, s.logger)
	
	// Setup actual HTTP routes
	http.NewRouter(s.app, itemUseCase, s.logger)
}

func (s *ItemIntegrationTestSuite) TearDownSuite() {
	if s.testDB != nil {
		s.testDB.CleanupTestDB(s.T())
	}
}

func (s *ItemIntegrationTestSuite) SetupTest() {
	if s.testDB != nil {
		// Clean database before each test using the safer method
		s.testDB.CleanData(s.T())
	}
}

// Test complete CRUD workflow end-to-end
func (s *ItemIntegrationTestSuite) TestItemCRUDFlow() {
	// === CREATE ===
	createRequest := request.AddItem{
		Name:   "Integration Test Item",
		Amount: 150,
	}
	
	bodyBytes, err := json.Marshal(createRequest)
	s.Require().NoError(err)
	
	req := httptest.NewRequest("POST", "/api/v1/items", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.app.Test(req, 10000) // 10 second timeout
	s.Require().NoError(err)
	s.Assert().Equal(201, resp.StatusCode)
	
	// Parse create response
	var createResp entities.Item
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	s.Require().NoError(err)
	s.Assert().Equal("Integration Test Item", createResp.Name)
	s.Assert().Equal(uint(150), createResp.Amount)
	s.Assert().NotEmpty(createResp.Id)
	itemID := createResp.Id.String()
	
	// === READ ===
	req = httptest.NewRequest("GET", "/api/v1/items/"+itemID, nil)
	resp, err = s.app.Test(req, 10000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	var getResp entities.Item
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	s.Require().NoError(err)
	s.Assert().Equal(itemID, getResp.Id.String())
	s.Assert().Equal("Integration Test Item", getResp.Name)
	
	// === UPDATE ===
	updateRequest := request.UpdateItem{
		Name:   stringPtr("Updated Integration Item"),
		Amount: uintPtr(300),
	}
	
	bodyBytes, err = json.Marshal(updateRequest)
	s.Require().NoError(err)
	
	req = httptest.NewRequest("PUT", "/api/v1/items/"+itemID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = s.app.Test(req, 10000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	var updateResp entities.Item
	err = json.NewDecoder(resp.Body).Decode(&updateResp)
	s.Require().NoError(err)
	s.Assert().Equal("Updated Integration Item", updateResp.Name)
	s.Assert().Equal(uint(300), updateResp.Amount)
	
	// === DELETE ===
	req = httptest.NewRequest("DELETE", "/api/v1/items/"+itemID, nil)
	resp, err = s.app.Test(req, 10000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	// Verify deletion - should return 404
	req = httptest.NewRequest("GET", "/api/v1/items/"+itemID, nil)
	resp, err = s.app.Test(req, 10000)
	s.Require().NoError(err)
	s.Assert().Equal(404, resp.StatusCode)
}

// Test concurrent item creation (race condition prevention)
func (s *ItemIntegrationTestSuite) TestConcurrentItemCreation() {
	const numGoroutines = 10
	const itemName = "Concurrent Test Item"
	
	var wg sync.WaitGroup
	results := make(chan int, numGoroutines)
	
	// Launch multiple goroutines trying to create the same item
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			createRequest := request.AddItem{
				Name:   itemName,
				Amount: 100,
			}
			
			bodyBytes, _ := json.Marshal(createRequest)
			req := httptest.NewRequest("POST", "/api/v1/items", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := s.app.Test(req, 5000)
			if err == nil {
				results <- resp.StatusCode
			}
		}()
	}
	
	wg.Wait()
	close(results)
	
	// Collect results
	successCount := 0
	conflictCount := 0
	
	for statusCode := range results {
		switch statusCode {
		case 201:
			successCount++
		case 409:
			conflictCount++
		}
	}
	
	// Only ONE should succeed, others should get 409 Conflict
	s.Assert().Equal(1, successCount, "Only one concurrent create should succeed")
	s.Assert().Equal(numGoroutines-1, conflictCount, "Other creates should get 409 Conflict")
}

// Test pagination with real data
func (s *ItemIntegrationTestSuite) TestPaginationIntegration() {
	// Create 25 items for pagination testing
	for i := 1; i <= 25; i++ {
		createRequest := request.AddItem{
			Name:   fmt.Sprintf("Pagination Item %02d", i),
			Amount: uint(i * 10),
		}
		
		bodyBytes, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/api/v1/items", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := s.app.Test(req, 5000)
		s.Require().NoError(err)
		s.Assert().Equal(201, resp.StatusCode)
	}
	
	// Test first page
	req := httptest.NewRequest("GET", "/api/v1/items?page=1&limit=10", nil)
	resp, err := s.app.Test(req, 5000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	var page1Response types.PaginatedResult[*entities.Item]
	err = json.NewDecoder(resp.Body).Decode(&page1Response)
	s.Require().NoError(err)
	s.Assert().Equal(10, len(page1Response.Items))
	s.Assert().Equal(int64(25), page1Response.Total)
	s.Assert().Equal(1, page1Response.Page)
	s.Assert().Equal(3, page1Response.TotalPages) // 25/10 = 3 pages
	
	// Test second page
	req = httptest.NewRequest("GET", "/api/v1/items?page=2&limit=10", nil)
	resp, err = s.app.Test(req, 5000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	var page2Response types.PaginatedResult[*entities.Item]
	err = json.NewDecoder(resp.Body).Decode(&page2Response)
	s.Require().NoError(err)
	s.Assert().Equal(10, len(page2Response.Items))
	s.Assert().Equal(2, page2Response.Page)
	
	// Test last page (should have 5 items)
	req = httptest.NewRequest("GET", "/api/v1/items?page=3&limit=10", nil)
	resp, err = s.app.Test(req, 5000)
	s.Require().NoError(err)
	s.Assert().Equal(200, resp.StatusCode)
	
	var page3Response types.PaginatedResult[*entities.Item]
	err = json.NewDecoder(resp.Body).Decode(&page3Response)
	s.Require().NoError(err)
	s.Assert().Equal(5, len(page3Response.Items)) // Remaining items
	s.Assert().Equal(3, page3Response.Page)
}

// Test bulk creation with transaction rollback
func (s *ItemIntegrationTestSuite) TestBulkCreationIntegration() {
	bulkRequest := request.BulkCreateItems{
		Items: []request.AddItem{
			{Name: "Bulk Item 1", Amount: 100},
			{Name: "Bulk Item 2", Amount: 200},
			{Name: "Bulk Item 3", Amount: 300},
		},
	}
	
	bodyBytes, err := json.Marshal(bulkRequest)
	s.Require().NoError(err)
	
	req := httptest.NewRequest("POST", "/api/v1/items/bulk", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.app.Test(req, 10000)
	s.Require().NoError(err)
	s.Assert().Equal(201, resp.StatusCode)
	
	var bulkResponse []*entities.Item
	err = json.NewDecoder(resp.Body).Decode(&bulkResponse)
	s.Require().NoError(err)
	s.Assert().Equal(3, len(bulkResponse))
	
	// Verify all items were created in database
	for i, item := range bulkResponse {
		req = httptest.NewRequest("GET", "/api/v1/items/"+item.Id.String(), nil)
		resp, err = s.app.Test(req, 5000)
		s.Require().NoError(err)
		s.Assert().Equal(200, resp.StatusCode)
		
		var verifyResp entities.Item
		err = json.NewDecoder(resp.Body).Decode(&verifyResp)
		s.Require().NoError(err)
		s.Assert().Equal(fmt.Sprintf("Bulk Item %d", i+1), verifyResp.Name)
	}
}

// Run the test suite
func TestItemIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ItemIntegrationTestSuite))
}

// Simple health check test (separate from the main integration suite)
func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	
	// Add a simple health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})
	
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "healthy", response["status"])
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}