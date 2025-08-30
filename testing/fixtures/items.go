package fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
)

// ValidItem returns a valid test item
func ValidItem() *entities.Item {
	return &entities.Item{
		BaseEntity: entities.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "Test Item",
		Amount: 100,
	}
}

// ValidItemWithName returns a valid test item with custom name
func ValidItemWithName(name string) *entities.Item {
	item := ValidItem()
	item.Name = name
	return item
}

// ValidItemWithAmount returns a valid test item with custom amount
func ValidItemWithAmount(amount uint) *entities.Item {
	item := ValidItem()
	item.Amount = amount
	return item
}

// ValidItems returns multiple valid test items
func ValidItems(count int) []*entities.Item {
	items := make([]*entities.Item, count)
	for i := 0; i < count; i++ {
		items[i] = ValidItemWithName(fmt.Sprintf("Test Item %d", i+1))
		items[i].Amount = uint((i + 1) * 10)
	}
	return items
}

// EmptyItem returns an item with empty required fields
func EmptyItem() *entities.Item {
	return &entities.Item{
		Name:   "",
		Amount: 0,
	}
}

// InvalidNameItem returns an item with invalid name (too long)
func InvalidNameItem() *entities.Item {
	longName := ""
	for i := 0; i < 101; i++ { // 101 characters (exceeds 100 limit)
		longName += "a"
	}
	
	return &entities.Item{
		Name:   longName,
		Amount: 100,
	}
}

// InvalidAmountItem returns an item with invalid amount (too large)
func InvalidAmountItem() *entities.Item {
	return &entities.Item{
		Name:   "Test Item",
		Amount: 1000000, // Exceeds 999999 limit
	}
}

// ItemWithID returns an item with specific ID
func ItemWithID(id uuid.UUID) *entities.Item {
	item := ValidItem()
	item.Id = id
	return item
}