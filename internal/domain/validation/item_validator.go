package validation

import (
	"strings"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
)

// ItemValidator provides validation methods for Item entities
type ItemValidator struct{}

// NewItemValidator creates a new item validator
func NewItemValidator() *ItemValidator {
	return &ItemValidator{}
}

// ValidateItem performs business validation on the Item entity
func (v *ItemValidator) ValidateItem(item *entities.Item) error {
	if strings.TrimSpace(item.Name) == "" {
		return domain.ErrItemNameRequired
	}
	
	if len(item.Name) > 100 {
		return domain.ErrItemNameTooLong
	}
	
	// Business rule: Amount cannot exceed 999999
	if item.Amount > 999999 {
		return domain.ErrItemAmountTooLarge
	}
	
	return nil
}

// ValidateName validates item name specifically
func (v *ItemValidator) ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return domain.ErrItemNameRequired
	}
	
	if len(name) > 100 {
		return domain.ErrItemNameTooLong
	}
	
	return nil
}

// ValidateAmount validates item amount specifically
func (v *ItemValidator) ValidateAmount(amount uint) error {
	if amount > 999999 {
		return domain.ErrItemAmountTooLarge
	}
	
	return nil
}