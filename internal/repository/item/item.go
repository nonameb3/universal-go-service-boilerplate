package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/pkg/errors"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"gorm.io/gorm"
)

type itemRepository struct {
	db         *gorm.DB
	logger     logger.Logger
	errHandler *errors.ErrorHandler
}

func NewItemRepository(db *gorm.DB, logger logger.Logger) ItemRepository {
	return &itemRepository{
		db:         db,
		logger:     logger,
		errHandler: errors.NewErrorHandler(),
	}
}

func (r *itemRepository) Create(item *entities.Item) (*entities.Item, error) {
	return r.CreateWithTx(r.db, item)
}

func (r *itemRepository) CreateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error) {
	err := tx.Create(item).Error
	if err != nil {
		// Map database errors to domain errors
		mappedErr := r.errHandler.MapDatabaseError(err)
		r.logger.Error("failed to create item", mappedErr)
		return nil, mappedErr
	}
	return item, nil
}

func (r *itemRepository) Get(id string) (*entities.Item, error) {
	return r.GetWithTx(r.db, id)
}

func (r *itemRepository) GetWithTx(tx *gorm.DB, id string) (*entities.Item, error) {
	item := &entities.Item{}
	if err := tx.Where("id = ?", id).First(item).Error; err != nil {
		r.logger.Error("failed to get item", err)
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) GetByName(name string) (*entities.Item, error) {
	return r.GetByNameWithTx(r.db, name)
}

func (r *itemRepository) GetByNameWithTx(tx *gorm.DB, name string) (*entities.Item, error) {
	item := &entities.Item{}
	if err := tx.Where("name = ?", name).First(item).Error; err != nil {
		return nil, err // Don't log "not found" as error - it's expected business case
	}
	return item, nil
}

// GetByNameForUpdate uses SELECT FOR UPDATE for pessimistic locking
func (r *itemRepository) GetByNameForUpdate(tx *gorm.DB, name string) (*entities.Item, error) {
	item := &entities.Item{}
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("name = ?", name).First(item).Error; err != nil {
		return nil, err // Don't log "not found" as error - it's expected business case
	}
	return item, nil
}

func (r *itemRepository) GetByNames(names []string) ([]*entities.Item, error) {
	return r.GetByNamesWithTx(r.db, names)
}

func (r *itemRepository) GetByNamesWithTx(tx *gorm.DB, names []string) ([]*entities.Item, error) {
	if len(names) == 0 {
		return []*entities.Item{}, nil
	}

	var items []*entities.Item
	if err := tx.Where("name IN ?", names).Find(&items).Error; err != nil {
		r.logger.Error("failed to get items by names", err)
		return nil, err
	}
	return items, nil
}

func (r *itemRepository) Update(item *entities.Item) (*entities.Item, error) {
	return r.UpdateWithTx(r.db, item)
}

func (r *itemRepository) UpdateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error) {
	if err := tx.Save(item).Error; err != nil {
		r.logger.Error("failed to update item", err)
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) GetWithPagination(page, limit int) (*types.PaginatedResult[*entities.Item], error) {
	var items []*entities.Item
	var total int64

	// Count total records
	if err := r.db.Model(&entities.Item{}).Count(&total).Error; err != nil {
		r.logger.Error("failed to count items", err)
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated items
	if err := r.db.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		r.logger.Error("failed to get paginated items", err)
		return nil, err
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &types.PaginatedResult[*entities.Item]{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *itemRepository) Delete(id string) error {
	return r.DeleteWithTx(r.db, id)
}

func (r *itemRepository) DeleteWithTx(tx *gorm.DB, id string) error {
	return tx.Where("id = ?", id).Delete(&entities.Item{}).Error
}
