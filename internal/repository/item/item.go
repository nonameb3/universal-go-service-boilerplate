package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"gorm.io/gorm"
)

type itemRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewItemRepository(db *gorm.DB, logger logger.Logger) ItemRepository {
	return &itemRepository{
		db:     db,
		logger: logger,
	}
}

func (r *itemRepository) Create(item *entities.Item) (*entities.Item, error) {
	if err := r.db.Create(item).Error; err != nil {
		r.logger.Error("failed to create item", err)
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) Get(id string) (*entities.Item, error) {
	item := &entities.Item{}
	if err := r.db.Where("id = ?", id).First(item).Error; err != nil {
		r.logger.Error("failed to get item", err)
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) GetByName(name string) (*entities.Item, error) {
	item := &entities.Item{}
	if err := r.db.Where("name = ?", name).First(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) Update(item *entities.Item) (*entities.Item, error) {
	if err := r.db.Save(item).Error; err != nil {
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
	return r.db.Where("id = ?", id).Delete(&entities.Item{}).Error
}
