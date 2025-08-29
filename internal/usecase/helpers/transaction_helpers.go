package helpers

import (
	"gorm.io/gorm"
	
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/pkg/providers"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// TransactionHelper provides reusable transaction utilities for boilerplate pattern
type TransactionHelper struct {
	db     providers.DatabaseProvider
	logger logger.Logger
}

// NewTransactionHelper creates a new transaction helper
func NewTransactionHelper(db providers.DatabaseProvider, logger logger.Logger) *TransactionHelper {
	return &TransactionHelper{
		db:     db,
		logger: logger,
	}
}

// WithTransaction executes a function within a database transaction
// Similar to NestJS: await this.dataSource.manager.transaction(async (manager) => {...})
func (h *TransactionHelper) WithTransaction(fn func(tx *gorm.DB) error) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		h.logger.Debug("Starting database transaction")
		
		err := fn(tx)
		if err != nil {
			h.logger.Error("Transaction failed, rolling back", err)
			return err
		}
		
		h.logger.Debug("Transaction completed successfully")
		return nil
	})
}

// AtomicCreateItem performs atomic create with duplicate checking for items
// Enterprise pattern for race-condition-safe creation
func (h *TransactionHelper) AtomicCreateItem(
	checkFn func(tx *gorm.DB) error,
	createFn func(tx *gorm.DB) (*entities.Item, error),
) (*entities.Item, error) {
	var result *entities.Item
	err := h.WithTransaction(func(tx *gorm.DB) error {
		// Check for duplicates/constraints first
		if err := checkFn(tx); err != nil {
			return err
		}
		
		// Create within same transaction
		var err error
		result, err = createFn(tx)
		return err
	})
	return result, err
}

// AtomicBulkCreateItems performs atomic bulk create for items
// All-or-nothing pattern for bulk operations
func (h *TransactionHelper) AtomicBulkCreateItems(
	items []*entities.Item,
	validateFn func(items []*entities.Item) error,
	createFn func(tx *gorm.DB, items []*entities.Item) ([]*entities.Item, error),
) ([]*entities.Item, error) {
	// Pre-validate outside transaction for fast failure
	if err := validateFn(items); err != nil {
		return nil, err
	}
	
	var results []*entities.Item
	err := h.WithTransaction(func(tx *gorm.DB) error {
		var err error
		results, err = createFn(tx, items)
		return err
	})
	return results, err
}