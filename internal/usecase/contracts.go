package usecase

import "github.com/universal-go-service/boilerplate/internal/domain"

type (
	// ItemUseCase -.
	ItemUseCase interface {
		Create(item *domain.Item) error
		Get(id string) (*domain.Item, error)
		Update(item *domain.Item) error
		Delete(id string) error
	}
	// other usecases will be added here
)
