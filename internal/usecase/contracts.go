package usecase

import (
	"github.com/universal-go-service/boilerplate/internal/domain"
	itemUseCase "github.com/universal-go-service/boilerplate/internal/usecase/item"
)

type (
	// ItemUseCase -.
	ItemUseCase interface {
		Create(req *itemUseCase.CreateItemRequest) (*domain.Item, error)
		Get(id string) (*domain.Item, error)
		GetWithPagination(req *itemUseCase.PaginationRequest) (*domain.PaginatedResult[*domain.Item], error)
		Update(id string, req *itemUseCase.UpdateItemRequest) (*domain.Item, error)
		Delete(id string) error
	}
	// other UseCases will be added here
)
