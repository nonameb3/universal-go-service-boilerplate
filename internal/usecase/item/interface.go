package item

import "github.com/universal-go-service/boilerplate/internal/domain"

type ItemUseCase interface {
	Create(req *CreateItemRequest) (*domain.Item, error)
	Get(id string) (*domain.Item, error)
	GetWithPagination(req *PaginationRequest) (*domain.PaginatedResult[*domain.Item], error)
	Update(id string, req *UpdateItemRequest) (*domain.Item, error)
	Delete(id string) error
}
