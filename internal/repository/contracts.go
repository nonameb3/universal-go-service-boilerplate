// Package repo implements application outer layer logic. Each logic group in own file.
package repository

import (
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/internal/domain/types"
)

type (
	// ItemRepo -.
	ItemRepo interface {
		Create(item *entities.Item) (*entities.Item, error)
		Get(id string) (*entities.Item, error)
		GetByName(name string) (*entities.Item, error)
		GetWithPagination(page, limit int) (*types.PaginatedResult[*entities.Item], error)
		Update(item *entities.Item) (*entities.Item, error)
		Delete(id string) error
	}
	// other repositories will be added here
)
