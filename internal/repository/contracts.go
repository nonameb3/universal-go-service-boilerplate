// Package repo implements application outer layer logic. Each logic group in own file.
package repository

import "github.com/universal-go-service/boilerplate/internal/domain"

type (
	// ItemRepo -.
	ItemRepo interface {
		Create(item *domain.Item) (*domain.Item, error)
		Get(id string) (*domain.Item, error)
		Update(item *domain.Item) (*domain.Item, error)
		Delete(id string) error
	}
	// other repositories will be added here
)
