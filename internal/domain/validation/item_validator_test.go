package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/testing/fixtures"
)

func TestItemValidator_ValidateItem(t *testing.T) {
	validator := NewItemValidator()

	tests := []struct {
		name        string
		item        *entities.Item
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid item should pass",
			item:        fixtures.ValidItem(),
			expectError: false,
		},
		{
			name:        "empty name should fail",
			item:        fixtures.EmptyItem(),
			expectError: true,
			expectedErr: domain.ErrItemNameRequired,
		},
		{
			name: "whitespace-only name should fail",
			item: &entities.Item{
				Name:   "   ",
				Amount: 100,
			},
			expectError: true,
			expectedErr: domain.ErrItemNameRequired,
		},
		{
			name:        "name too long should fail",
			item:        fixtures.InvalidNameItem(),
			expectError: true,
			expectedErr: domain.ErrItemNameTooLong,
		},
		{
			name:        "amount too large should fail",
			item:        fixtures.InvalidAmountItem(),
			expectError: true,
			expectedErr: domain.ErrItemAmountTooLarge,
		},
		{
			name: "zero amount should pass",
			item: fixtures.ValidItemWithAmount(0),
			expectError: false,
		},
		{
			name: "minimum amount should pass",
			item: fixtures.ValidItemWithAmount(1),
			expectError: false,
		},
		{
			name: "maximum valid amount should pass",
			item: fixtures.ValidItemWithAmount(999999),
			expectError: false,
		},
		{
			name: "exact name length limit should pass",
			item: func() *entities.Item {
				name := ""
				for i := 0; i < 100; i++ { // Exactly 100 characters
					name += "a"
				}
				return fixtures.ValidItemWithName(name)
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateItem(tt.item)

			if tt.expectError {
				require.Error(t, err, "Expected validation error")
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err, "Error should match expected error")
				}
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}

func TestItemValidator_ValidateItemUpdate(t *testing.T) {
	validator := NewItemValidator()

	tests := []struct {
		name        string
		item        *entities.Item
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid item should pass",
			item:        fixtures.ValidItem(),
			expectError: false,
		},
		{
			name:        "empty name should fail",
			item:        fixtures.EmptyItem(),
			expectError: true,
			expectedErr: domain.ErrItemNameRequired,
		},
		{
			name:        "name too long should fail",
			item:        fixtures.InvalidNameItem(),
			expectError: true,
			expectedErr: domain.ErrItemNameTooLong,
		},
		{
			name:        "amount too large should fail",
			item:        fixtures.InvalidAmountItem(),
			expectError: true,
			expectedErr: domain.ErrItemAmountTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateItem(tt.item)

			if tt.expectError {
				require.Error(t, err, "Expected validation error")
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err, "Error should match expected error")
				}
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}

func TestItemValidator_ValidatePagination(t *testing.T) {
	validator := NewItemValidator()

	tests := []struct {
		name        string
		page        int
		limit       int
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid pagination should pass",
			page:        1,
			limit:       10,
			expectError: false,
		},
		{
			name:        "zero page should fail",
			page:        0,
			limit:       10,
			expectError: true,
			expectedErr: domain.ErrInvalidPagination,
		},
		{
			name:        "negative page should fail",
			page:        -1,
			limit:       10,
			expectError: true,
			expectedErr: domain.ErrInvalidPagination,
		},
		{
			name:        "zero limit should fail",
			page:        1,
			limit:       0,
			expectError: true,
			expectedErr: domain.ErrInvalidPagination,
		},
		{
			name:        "negative limit should fail",
			page:        1,
			limit:       -5,
			expectError: true,
			expectedErr: domain.ErrInvalidPagination,
		},
		{
			name:        "limit too large should fail",
			page:        1,
			limit:       101,
			expectError: true,
			expectedErr: domain.ErrLimitTooLarge,
		},
		{
			name:        "maximum valid limit should pass",
			page:        1,
			limit:       100,
			expectError: false,
		},
		{
			name:        "large page number should pass",
			page:        1000,
			limit:       50,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePagination(tt.page, tt.limit)

			if tt.expectError {
				require.Error(t, err, "Expected validation error")
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err, "Error should match expected error")
				}
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}

func TestItemValidator_NilItem(t *testing.T) {
	validator := NewItemValidator()

	t.Run("ValidateItem with nil item should fail", func(t *testing.T) {
		err := validator.ValidateItem(nil)
		assert.Error(t, err, "Should return error for nil item")
	})
}