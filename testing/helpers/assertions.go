package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
)

// AssertItemEqual asserts that two items are equal
func AssertItemEqual(t *testing.T, expected, actual *entities.Item) {
	require.NotNil(t, expected, "Expected item should not be nil")
	require.NotNil(t, actual, "Actual item should not be nil")
	
	assert.Equal(t, expected.Id, actual.Id, "Item IDs should match")
	assert.Equal(t, expected.Name, actual.Name, "Item names should match")
	assert.Equal(t, expected.Amount, actual.Amount, "Item amounts should match")
}

// AssertItemNotEmpty asserts that item has required fields
func AssertItemNotEmpty(t *testing.T, item *entities.Item) {
	require.NotNil(t, item, "Item should not be nil")
	assert.NotEmpty(t, item.Id, "Item ID should not be empty")
	assert.NotEmpty(t, item.Name, "Item name should not be empty")
	assert.NotZero(t, item.CreatedAt, "CreatedAt should not be zero")
	assert.NotZero(t, item.UpdatedAt, "UpdatedAt should not be zero")
}

// AssertTimestampsValid asserts that timestamps are properly set
func AssertTimestampsValid(t *testing.T, item *entities.Item) {
	now := time.Now()
	
	assert.True(t, item.CreatedAt.Before(now) || item.CreatedAt.Equal(now), 
		"CreatedAt should be before or equal to now")
	assert.True(t, item.UpdatedAt.Before(now) || item.UpdatedAt.Equal(now), 
		"UpdatedAt should be before or equal to now")
	assert.True(t, item.UpdatedAt.After(item.CreatedAt) || item.UpdatedAt.Equal(item.CreatedAt), 
		"UpdatedAt should be after or equal to CreatedAt")
}

// AssertItemsEqual asserts that two slices of items are equal
func AssertItemsEqual(t *testing.T, expected, actual []*entities.Item) {
	require.Equal(t, len(expected), len(actual), "Item slices should have same length")
	
	for i, expectedItem := range expected {
		AssertItemEqual(t, expectedItem, actual[i])
	}
}

// AssertErrorContains asserts that error contains specific text
func AssertErrorContains(t *testing.T, err error, contains string) {
	require.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), contains, "Error should contain expected text")
}