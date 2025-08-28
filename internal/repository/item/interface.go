package item

import "github.com/universal-go-service/boilerplate/internal/domain"

type ItemRepository interface {
	Create(item *domain.Item) (*domain.Item, error)
	Get(id string) (*domain.Item, error)
	Update(item *domain.Item) (*domain.Item, error)
	Delete(id string) error
}
