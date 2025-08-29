package usecase

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
)

type (
	// ItemUseCase -.
	ItemUseCase interface {
		Create(req *dto.CreateItemRequest) (*entities.Item, error)
		Get(id string) (*entities.Item, error)
		GetWithPagination(req *dto.PaginationRequest) (*types.PaginatedResult[*entities.Item], error)
		Update(id string, req *dto.UpdateItemRequest) (*entities.Item, error)
		Delete(id string) error
	}
	// other UseCases will be added here
)
