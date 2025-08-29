package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
)

type ItemRepository interface {
	Create(item *entities.Item) (*entities.Item, error)
	Get(id string) (*entities.Item, error)
	GetByName(name string) (*entities.Item, error)
	GetWithPagination(page, limit int) (*types.PaginatedResult[*entities.Item], error)
	Update(item *entities.Item) (*entities.Item, error)
	Delete(id string) error
}
