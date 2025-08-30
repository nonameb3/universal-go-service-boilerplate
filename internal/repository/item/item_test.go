package item

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/universal-go-service/boilerplate/internal/domain"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/testing/fixtures"
	"github.com/universal-go-service/boilerplate/testing/helpers"
	"gorm.io/gorm"
)

func TestItemRepository_Create(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return // Skip test if no test database
	}
	defer testDB.CleanupTestDB(t)

	// Create repository instance with noop logger
	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should create item successfully", func(t *testing.T) {
		item := fixtures.ValidItem()
		item.Name = "Create Test Item"

		createdItem, err := repo.Create(item)

		require.NoError(t, err)
		require.NotNil(t, createdItem)
		helpers.AssertItemNotEmpty(t, createdItem)
		assert.Equal(t, "Create Test Item", createdItem.Name)
		assert.Equal(t, uint(100), createdItem.Amount)
	})

	t.Run("should fail on duplicate name", func(t *testing.T) {
		// Create first item
		item1 := fixtures.ValidItemWithName("Duplicate Test")
		_, err := repo.Create(item1)
		require.NoError(t, err)

		// Try to create second item with same name
		item2 := fixtures.ValidItemWithName("Duplicate Test")
		_, err = repo.Create(item2)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrItemAlreadyExists, err)
	})
}

func TestItemRepository_CreateWithTx(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should create item with transaction", func(t *testing.T) {
		item := fixtures.ValidItemWithName("Transaction Test")

		err := testDB.DB.Transaction(func(tx *gorm.DB) error {
			_, err := repo.CreateWithTx(tx, item)
			return err
		})

		require.NoError(t, err)

		// Verify item was created
		foundItem, err := repo.GetByName("Transaction Test")
		require.NoError(t, err)
		assert.Equal(t, "Transaction Test", foundItem.Name)
	})
}

func TestItemRepository_Get(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should get existing item", func(t *testing.T) {
		// Create test item
		createdItem := testDB.CreateTestItem("Get Test Item", 150)

		// Get the item
		foundItem, err := repo.Get(createdItem.Id.String())

		require.NoError(t, err)
		helpers.AssertItemEqual(t, createdItem, foundItem)
	})

	t.Run("should return error for non-existent item", func(t *testing.T) {
		_, err := repo.Get("non-existent-id")
		assert.Error(t, err)
	})
}

func TestItemRepository_GetByName(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should get item by name", func(t *testing.T) {
		createdItem := testDB.CreateTestItem("GetByName Test", 200)

		foundItem, err := repo.GetByName("GetByName Test")

		require.NoError(t, err)
		helpers.AssertItemEqual(t, createdItem, foundItem)
	})

	t.Run("should return error for non-existent name", func(t *testing.T) {
		_, err := repo.GetByName("non-existent-name")
		assert.Error(t, err)
	})
}

func TestItemRepository_GetByNameForUpdate(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should get item with FOR UPDATE lock", func(t *testing.T) {
		createdItem := testDB.CreateTestItem("Lock Test Item", 300)

		err := testDB.DB.Transaction(func(tx *gorm.DB) error {
			foundItem, err := repo.GetByNameForUpdate(tx, "Lock Test Item")
			if err != nil {
				return err
			}
			
			helpers.AssertItemEqual(t, createdItem, foundItem)
			return nil
		})

		require.NoError(t, err)
	})
}

func TestItemRepository_GetByNames(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should get multiple items by names", func(t *testing.T) {
		// Create test items
		testDB.CreateTestItem("Batch Item 1", 100)
		testDB.CreateTestItem("Batch Item 2", 200)
		testDB.CreateTestItem("Other Item", 300) // Should not be returned

		names := []string{"Batch Item 1", "Batch Item 2"}
		items, err := repo.GetByNames(names)

		require.NoError(t, err)
		assert.Len(t, items, 2)
		
		// Verify items are correct
		foundNames := make(map[string]bool)
		for _, item := range items {
			foundNames[item.Name] = true
		}
		assert.True(t, foundNames["Batch Item 1"])
		assert.True(t, foundNames["Batch Item 2"])
	})

	t.Run("should return empty slice for empty names", func(t *testing.T) {
		items, err := repo.GetByNames([]string{})

		require.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestItemRepository_Update(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should update item successfully", func(t *testing.T) {
		// Create item
		createdItem := testDB.CreateTestItem("Update Test", 100)

		// Update item
		createdItem.Name = "Updated Name"
		createdItem.Amount = 500

		updatedItem, err := repo.Update(createdItem)

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedItem.Name)
		assert.Equal(t, uint(500), updatedItem.Amount)
	})
}

func TestItemRepository_GetWithPagination(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should return paginated results", func(t *testing.T) {
		// Create test items
		testDB.CreateTestItems(15, "Page Test")

		// Get first page
		result, err := repo.GetWithPagination(1, 10)

		require.NoError(t, err)
		assert.Len(t, result.Items, 10)
		assert.Equal(t, int64(15), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, 2, result.TotalPages)
	})

	t.Run("should return second page", func(t *testing.T) {
		// Clear and create fresh items
		testDB.CleanData(t)
		testDB.CreateTestItems(25, "Page2 Test")

		result, err := repo.GetWithPagination(2, 10)

		require.NoError(t, err)
		assert.Len(t, result.Items, 10)
		assert.Equal(t, int64(25), result.Total)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 3, result.TotalPages)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	if testDB == nil {
		return
	}
	defer testDB.CleanupTestDB(t)

	noopLogger, _ := logger.NewNoop(logger.LoggerConfig{})
	repo := NewItemRepository(testDB.DB, noopLogger)

	t.Run("should delete item successfully", func(t *testing.T) {
		// Create item
		createdItem := testDB.CreateTestItem("Delete Test", 100)

		// Delete item
		err := repo.Delete(createdItem.Id.String())
		require.NoError(t, err)

		// Verify item is deleted
		_, err = repo.Get(createdItem.Id.String())
		assert.Error(t, err)
	})
}