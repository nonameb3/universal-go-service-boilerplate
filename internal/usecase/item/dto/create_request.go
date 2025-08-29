package dto

import (
	"strings"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
)

// CreateItemRequest represents the business request to create an item
type CreateItemRequest struct {
	Name   string `json:"name"`
	Amount uint   `json:"amount"`
}

// Validate performs business validation on the create request
func (r *CreateItemRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return domain.ErrItemNameRequired
	}
	
	if len(r.Name) > 100 {
		return domain.ErrItemNameTooLong
	}
	
	if r.Amount > 999999 {
		return domain.ErrItemAmountTooLarge
	}
	
	return nil
}

// ToEntity converts the request to a domain entity
func (r *CreateItemRequest) ToEntity() *entities.Item {
	return &entities.Item{
		Name:   strings.TrimSpace(r.Name),
		Amount: r.Amount,
	}
}