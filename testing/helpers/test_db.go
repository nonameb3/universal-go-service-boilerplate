package helpers

import (
	"fmt"
	"testing"

	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/pkg/providers/database"
	"gorm.io/gorm"
)

// TestDatabase represents a test database instance
type TestDatabase struct {
	DB       *gorm.DB
	Provider database.DatabaseProvider
}

// SetupTestDB creates a test database connection for testing
func SetupTestDB(t *testing.T) *TestDatabase {
	// Use test configuration matching the actual container setup
	config := database.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres", // Actual password from container environment
		Database: "universal_service_test",
		SSLMode:  "disable",
		Timezone: "UTC",
	}

	provider, err := database.NewPostgres(config)
	if err != nil {
		t.Skipf("Test database not available: %v", err)
		return nil
	}

	// Test connection
	if err := provider.Health(); err != nil {
		t.Skipf("Test database health check failed: %v", err)
		return nil
	}

	db := provider.GetDB()

	// Auto-migrate test tables
	err = db.AutoMigrate(&entities.Item{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return &TestDatabase{
		DB:       db,
		Provider: provider,
	}
}

// CleanupTestDB cleans up test database and closes connection
func (td *TestDatabase) CleanupTestDB(t *testing.T) {
	// Clean up all test data
	td.CleanData(t)
	
	// Close connection
	if td.Provider != nil {
		td.Provider.Close()
	}
}

// CleanData cleans up test data without closing the connection
func (td *TestDatabase) CleanData(t *testing.T) {
	td.DB.Exec("TRUNCATE TABLE items RESTART IDENTITY CASCADE")
}

// CreateTestItem creates a test item in the database
func (td *TestDatabase) CreateTestItem(name string, amount uint) *entities.Item {
	item := &entities.Item{
		Name:   name,
		Amount: amount,
	}
	
	td.DB.Create(item)
	return item
}

// CreateTestItems creates multiple test items
func (td *TestDatabase) CreateTestItems(count int, namePrefix string) []*entities.Item {
	items := make([]*entities.Item, count)
	
	for i := 0; i < count; i++ {
		items[i] = td.CreateTestItem(
			fmt.Sprintf("%s_%d", namePrefix, i+1),
			uint((i+1)*10),
		)
	}
	
	return items
}