package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"gorm.io/gorm"
)

var config = logger.LoggerConfig{
	Type:        "tokenx",
	ServiceName: "item-service",
}
var log = logger.NewCentralizedLogger(config)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *domain.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) Get(id string) (*domain.Item, error) {
	item := &domain.Item{}
	if err := r.db.Where("id = ?", id).First(item).Error; err != nil {
		log.Error("failed to get item", err)
		return nil, err
	}
	return item, nil
}

func (r *itemRepository) Update(item *domain.Item) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) Delete(id string) error {
	return r.db.Delete(&domain.Item{}, id).Error
}
