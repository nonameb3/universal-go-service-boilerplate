package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *entities.Item) (*entities.Item, error)
	CreateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error)
	Get(id string) (*entities.Item, error)
	GetWithTx(tx *gorm.DB, id string) (*entities.Item, error)
	GetByName(name string) (*entities.Item, error)
	GetByNameWithTx(tx *gorm.DB, name string) (*entities.Item, error)
	GetByNameForUpdate(tx *gorm.DB, name string) (*entities.Item, error)
	GetByNames(names []string) ([]*entities.Item, error)
	GetByNamesWithTx(tx *gorm.DB, names []string) ([]*entities.Item, error)
	GetWithPagination(page, limit int) (*types.PaginatedResult[*entities.Item], error)
	Update(item *entities.Item) (*entities.Item, error)
	UpdateWithTx(tx *gorm.DB, item *entities.Item) (*entities.Item, error)
	Delete(id string) error
	DeleteWithTx(tx *gorm.DB, id string) error
}
