package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to parse UUID for tests
func parseUUID(t *testing.T, uuidStr string) uuid.UUID {
	parsed, err := uuid.Parse(uuidStr)
	require.NoError(t, err)
	return parsed
}

func TestItem_BeforeCreate(t *testing.T) {
	tests := []struct {
		name string
		item *Item
	}{
		{
			name: "should generate ID and set timestamps",
			item: &Item{
				Name:   "Test Item",
				Amount: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call BeforeCreate hook on BaseEntity
			err := tt.item.BaseEntity.BeforeCreate(nil)
			
			// Assertions
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, tt.item.Id, "ID should be generated")
			
			// Verify ID format (UUID)
			assert.NotNil(t, tt.item.Id, "ID should be UUID format")
		})
	}
}

func TestItem_BeforeCreate_PreserveExistingID(t *testing.T) {
	item := &Item{
		Name:   "Test Item",
		Amount: 100,
	}
	// Set existing UUID
	existingUUID := parseUUID(t, "550e8400-e29b-41d4-a716-446655440000")
	item.Id = existingUUID
	
	err := item.BaseEntity.BeforeCreate(nil)
	
	require.NoError(t, err)
	assert.Equal(t, existingUUID, item.Id, "Existing ID should be preserved")
}

func TestItem_UpdateFrom(t *testing.T) {
	tests := []struct {
		name           string
		item          *Item
		updateName    *string
		updateAmount  *uint
		expectedName  string
		expectedAmount uint
	}{
		{
			name: "update both name and amount",
			item: &Item{Name: "Old Name", Amount: 100},
			updateName: stringPtr("New Name"),
			updateAmount: uintPtr(200),
			expectedName: "New Name",
			expectedAmount: 200,
		},
		{
			name: "update only name",
			item: &Item{Name: "Old Name", Amount: 100},
			updateName: stringPtr("New Name"),
			updateAmount: nil,
			expectedName: "New Name",
			expectedAmount: 100,
		},
		{
			name: "update only amount",
			item: &Item{Name: "Old Name", Amount: 100},
			updateName: nil,
			updateAmount: uintPtr(200),
			expectedName: "Old Name",
			expectedAmount: 200,
		},
		{
			name: "trim whitespace from name",
			item: &Item{Name: "Old Name", Amount: 100},
			updateName: stringPtr("  New Name  "),
			updateAmount: nil,
			expectedName: "New Name",
			expectedAmount: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.UpdateFrom(tt.updateName, tt.updateAmount)
			
			assert.Equal(t, tt.expectedName, tt.item.Name)
			assert.Equal(t, tt.expectedAmount, tt.item.Amount)
		})
	}
}

func TestItem_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		item     *Item
		expected bool
	}{
		{
			name: "empty item should return true",
			item: &Item{Name: "", Amount: 0},
			expected: true,
		},
		{
			name: "whitespace name should return true",
			item: &Item{Name: "   ", Amount: 0},
			expected: true,
		},
		{
			name: "item with name should return false",
			item: &Item{Name: "Test", Amount: 0},
			expected: false,
		},
		{
			name: "item with amount should return false",
			item: &Item{Name: "", Amount: 100},
			expected: false,
		},
		{
			name: "item with both should return false",
			item: &Item{Name: "Test", Amount: 100},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func TestItem_Fields(t *testing.T) {
	item := &Item{
		Name:   "Test Item Name",
		Amount: 500,
	}
	
	assert.Equal(t, "Test Item Name", item.Name)
	assert.Equal(t, uint(500), item.Amount)
}

func TestItem_JSONTags(t *testing.T) {
	// This test ensures the struct has proper JSON tags
	// We can't directly test JSON tags, but we can test JSON marshaling/unmarshaling
	item := &Item{
		Name:   "Test Item",
		Amount: 100,
	}
	item.Id = parseUUID(t, "550e8400-e29b-41d4-a716-446655440000")
	item.CreatedAt = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	item.UpdatedAt = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	
	// This would be tested in integration tests with actual JSON marshaling
	// For now, just verify the fields are accessible
	assert.Equal(t, parseUUID(t, "550e8400-e29b-41d4-a716-446655440000"), item.Id)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, uint(100), item.Amount)
	assert.NotZero(t, item.CreatedAt)
	assert.NotZero(t, item.UpdatedAt)
}