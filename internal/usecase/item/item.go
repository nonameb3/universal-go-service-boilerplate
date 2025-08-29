package item

import (
	"errors"
	"sync"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/domain/validation"
	"github.com/universal-go-service/boilerplate/internal/repository"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

type itemUseCase struct {
	itemRepo  repository.ItemRepo
	logger    logger.Logger
	validator *validation.ItemValidator
}

func NewItemUseCase(itemRepo repository.ItemRepo, logger logger.Logger) ItemUseCase {
	return &itemUseCase{
		itemRepo:  itemRepo,
		logger:    logger,
		validator: validation.NewItemValidator(),
	}
}

// Create implements business logic for creating an item
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
	
	// Business rule: Check for duplicates (example business logic)
	// Note: This would require a repository method like GetByName
	// For now, we'll just proceed with creation
	
	createdItem, err := uc.itemRepo.Create(item)
	if err != nil {
		uc.logger.Error("Failed to create item in repository", err)
		return nil, err
	}
	
	uc.logger.Info("Item created successfully")
	return createdItem, nil
}

// BulkCreate implements business logic for creating multiple items concurrently
func (uc *itemUseCase) BulkCreate(req *dto.BulkCreateRequest) ([]*entities.Item, error) {
	// Business validation
	if err := req.Validate(); err != nil {
		uc.logger.Error("Bulk create validation failed", err)
		return nil, err
	}

	// Convert to entities for processing
	itemsToCreate := req.ToEntities()

	// Process items concurrently using goroutines
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	results := make([]*entities.Item, 0, len(itemsToCreate))
	errs := make([]error, 0)

	// Create worker goroutines
	for _, item := range itemsToCreate {
		wg.Add(1)
		go func(itemToCreate *entities.Item) {
			defer wg.Done()

			// Domain validation in goroutine
			if err := uc.validator.ValidateItem(itemToCreate); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			// Create item in repository
			createdItem, err := uc.itemRepo.Create(itemToCreate)
			
			// Thread-safe result collection
			mu.Lock()
			if err != nil {
				errs = append(errs, err)
			} else {
				results = append(results, createdItem)
			}
			mu.Unlock()
		}(item)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for errors
	if len(errs) > 0 {
		uc.logger.Error("Bulk create had errors", errs[0])
		return results, errs[0] // Return first error
	}

	uc.logger.Info("Bulk create completed successfully")
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
