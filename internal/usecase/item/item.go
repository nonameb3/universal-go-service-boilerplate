package item

import (
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/repository"
)

type itemUseCase struct {
	itemRepo repository.ItemRepo
}

func NewItemUseCase(itemRepo repository.ItemRepo) ItemUseCase {
	return &itemUseCase{itemRepo: itemRepo}
}

func (uc *itemUseCase) Create(item *domain.Item) (*domain.Item, error) {
	return uc.itemRepo.Create(item)
}

func (uc *itemUseCase) Get(id string) (*domain.Item, error) {
	return uc.itemRepo.Get(id)
}

func (uc *itemUseCase) Update(item *domain.Item) (*domain.Item, error) {
	return uc.itemRepo.Update(item)
}

func (uc *itemUseCase) Delete(id string) error {
	return uc.itemRepo.Delete(id)
}
