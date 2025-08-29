package item

import (
	"errors"
	
	"gorm.io/gorm"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/domain/validation"
	"github.com/universal-go-service/boilerplate/internal/repository"
	"github.com/universal-go-service/boilerplate/internal/usecase/helpers"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
	"github.com/universal-go-service/boilerplate/pkg/providers"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

type itemUseCase struct {
	itemRepo  repository.ItemRepo
	db        providers.DatabaseProvider
	txHelper  *helpers.TransactionHelper
	logger    logger.Logger
	validator *validation.ItemValidator
}

func NewItemUseCase(itemRepo repository.ItemRepo, db providers.DatabaseProvider, logger logger.Logger) ItemUseCase {
	return &itemUseCase{
		itemRepo:  itemRepo,
		db:        db,
		txHelper:  helpers.NewTransactionHelper(db, logger),
		logger:    logger,
		validator: validation.NewItemValidator(),
	}
}

// Create implements business logic for creating an item with enterprise transaction safety
func (uc *itemUseCase) Create(req *dto.CreateItemRequest) (*entities.Item, error) {
	// Business validation
	if err := req.Validate(); err != nil {
		uc.logger.Error("Create item validation failed", err)
		return nil, err
	}
	
	// Convert to domain entity
	item := req.ToEntity()
	
	// Domain validation using validator
	if err := uc.validator.ValidateItem(item); err != nil {
		uc.logger.Error("Create item domain validation failed", err)
		return nil, err
	}
	
	// Use enterprise transaction helper for atomic create
	createdItem, err := uc.txHelper.AtomicCreateItem(
		// Check function: pessimistic locking to prevent race conditions
		func(tx *gorm.DB) error {
			existingItem, err := uc.itemRepo.GetByNameForUpdate(tx, item.Name)
			if err == nil && existingItem != nil {
				uc.logger.Error("Item with same name already exists", nil)
				return domain.ErrItemAlreadyExists
			}
			return nil
		},
		// Create function: create item within transaction
		func(tx *gorm.DB) (*entities.Item, error) {
			return uc.itemRepo.CreateWithTx(tx, item)
		},
	)
	
	if err != nil {
		return nil, err
	}
	
	uc.logger.Info("Item created successfully with enterprise transaction safety")
	return createdItem, nil
}

// BulkCreate implements business logic for creating multiple items with transaction safety
func (uc *itemUseCase) BulkCreate(req *dto.BulkCreateRequest) ([]*entities.Item, error) {
	// Business validation
	if err := req.Validate(); err != nil {
		uc.logger.Error("Bulk create validation failed", err)
		return nil, err
	}

	// Convert to entities for processing
	itemsToCreate := req.ToEntities()

	// Business rule: Check for internal duplicate names within the same request
	namesSeen := make(map[string]bool)
	for _, item := range itemsToCreate {
		if namesSeen[item.Name] {
			uc.logger.Error("Duplicate names found in bulk create request", nil)
			return nil, domain.ErrItemAlreadyExists
		}
		namesSeen[item.Name] = true
	}

	// Use single transaction for entire bulk operation (NestJS-style)
	var results []*entities.Item
	err := uc.db.Transaction(func(tx *gorm.DB) error {
		// Pre-validate all items
		for _, item := range itemsToCreate {
			if err := uc.validator.ValidateItem(item); err != nil {
				uc.logger.Error("Bulk create item validation failed", err)
				return err
			}
		}

		// Check for external duplicates in batches within transaction
		if err := uc.checkExternalDuplicatesInBatchesWithTx(tx, itemsToCreate); err != nil {
			return err
		}

		// Create all items within single transaction
		results = make([]*entities.Item, 0, len(itemsToCreate))
		for _, item := range itemsToCreate {
			createdItem, err := uc.itemRepo.CreateWithTx(tx, item)
			if err != nil {
				uc.logger.Error("Failed to create item in bulk operation", err)
				return err // This will rollback entire transaction
			}
			results = append(results, createdItem)
		}

		return nil // Success - commit transaction
	})

	if err != nil {
		return nil, err
	}

	uc.logger.Info("Bulk create completed successfully with transaction safety")
	return results, nil
}

// Get implements business logic for retrieving an item
func (uc *itemUseCase) Get(id string) (*entities.Item, error) {
	if id == "" {
		return nil, domain.ErrInvalidPagination // Using available error for now
	}
	
	item, err := uc.itemRepo.Get(id)
	if err != nil {
		uc.logger.Error("Failed to get item", err)
		return nil, domain.ErrItemNotFound
	}
	
	return item, nil
}

// GetWithPagination implements business logic for paginated retrieval
func (uc *itemUseCase) GetWithPagination(req *dto.PaginationRequest) (*types.PaginatedResult[*entities.Item], error) {
	// Apply business defaults
	req.ApplyDefaults()
	
	// Business validation
	if err := req.Validate(); err != nil {
		uc.logger.Error("Pagination validation failed", err)
		return nil, err
	}
	
	result, err := uc.itemRepo.GetWithPagination(req.Page, req.Limit)
	if err != nil {
		uc.logger.Error("Failed to get paginated items", err)
		return nil, err
	}
	
	return result, nil
}

// Update implements business logic for updating an item
func (uc *itemUseCase) Update(id string, req *dto.UpdateItemRequest) (*entities.Item, error) {
	if id == "" {
		return nil, domain.ErrInvalidPagination // Using available error for now
	}
	
	// Business validation
	if err := req.Validate(); err != nil {
		uc.logger.Error("Update item validation failed", err)
		return nil, err
	}
	
	// Business rule: Check if there are any updates
	if !req.HasUpdates() {
		return nil, errors.New("no updates provided")
	}
	
	// Get existing item (business rule: must exist)
	existingItem, err := uc.itemRepo.Get(id)
	if err != nil {
		uc.logger.Error("Failed to get existing item for update", err)
		return nil, domain.ErrItemNotFound
	}
	
	// Apply updates using business logic
	existingItem.UpdateFrom(req.Name, req.Amount)
	
	// Business rule: Check for duplicate names if name is being updated
	if req.Name != nil && *req.Name != "" {
		duplicateItem, err := uc.itemRepo.GetByName(*req.Name)
		if err == nil && duplicateItem != nil && duplicateItem.Id != existingItem.Id {
			uc.logger.Error("Item with same name already exists", nil)
			return nil, domain.ErrItemAlreadyExists
		}
	}
	
	// Domain validation after update using validator
	if err := uc.validator.ValidateItem(existingItem); err != nil {
		uc.logger.Error("Update item domain validation failed", err)
		return nil, err
	}
	
	updatedItem, err := uc.itemRepo.Update(existingItem)
	if err != nil {
		uc.logger.Error("Failed to update item in repository", err)
		return nil, err
	}
	
	uc.logger.Info("Item updated successfully")
	return updatedItem, nil
}

// Delete implements business logic for deleting an item
func (uc *itemUseCase) Delete(id string) error {
	if id == "" {
		return domain.ErrInvalidPagination // Using available error for now
	}
	
	// Business rule: Check if item exists before deletion
	_, err := uc.itemRepo.Get(id)
	if err != nil {
		uc.logger.Error("Item not found for deletion", err)
		return domain.ErrItemNotFound
	}
	
	// Business rule: Add any deletion constraints here
	// For example: Check if item is referenced by other entities
	
	if err := uc.itemRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete item", err)
		return err
	}
	
	uc.logger.Info("Item deleted successfully")
	return nil
}

// checkExternalDuplicatesInBatches checks for existing items with same names in batches
func (uc *itemUseCase) checkExternalDuplicatesInBatches(items []*entities.Item) error {
	const MAX_BATCH_SIZE = 1000
	
	for i := 0; i < len(items); i += MAX_BATCH_SIZE {
		end := min(i+MAX_BATCH_SIZE, len(items))
		
		// Extract names from current batch
		batch := items[i:end]
		names := make([]string, len(batch))
		for j, item := range batch {
			names[j] = item.Name
		}
		
		// Check batch for existing items
		existingItems, err := uc.itemRepo.GetByNames(names)
		if err != nil {
			uc.logger.Error("Failed to check for duplicate names", err)
			return err
		}
		
		// If any existing items found, return error
		if len(existingItems) > 0 {
			uc.logger.Error("Item with same name already exists in database", nil)
			return domain.ErrItemAlreadyExists
		}
	}
	
	return nil
}

// checkExternalDuplicatesInBatchesWithTx checks for existing items within a transaction
func (uc *itemUseCase) checkExternalDuplicatesInBatchesWithTx(tx *gorm.DB, items []*entities.Item) error {
	const MAX_BATCH_SIZE = 1000
	
	for i := 0; i < len(items); i += MAX_BATCH_SIZE {
		end := min(i+MAX_BATCH_SIZE, len(items))
		
		// Extract names from current batch
		batch := items[i:end]
		names := make([]string, len(batch))
		for j, item := range batch {
			names[j] = item.Name
		}
		
		// Check batch for existing items within transaction
		existingItems, err := uc.itemRepo.GetByNamesWithTx(tx, names)
		if err != nil {
			uc.logger.Error("Failed to check for duplicate names in transaction", err)
			return err
		}
		
		// If any existing items found, return error
		if len(existingItems) > 0 {
			uc.logger.Error("Item with same name already exists in database", nil)
			return domain.ErrItemAlreadyExists
		}
	}
	
	return nil
}
