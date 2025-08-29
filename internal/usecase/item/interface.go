package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
)

type ItemUseCase interface {
	Create(req *dto.CreateItemRequest) (*entities.Item, error)
	BulkCreate(req *dto.BulkCreateRequest) ([]*entities.Item, error)
	Get(id string) (*entities.Item, error)
	GetWithPagination(req *dto.PaginationRequest) (*types.PaginatedResult[*entities.Item], error)
	Update(id string, req *dto.UpdateItemRequest) (*entities.Item, error)
	Delete(id string) error
}
