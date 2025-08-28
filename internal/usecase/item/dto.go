package item

import (
	"strings"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
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
func (r *CreateItemRequest) ToEntity() *domain.Item {
	return &domain.Item{
		Name:   strings.TrimSpace(r.Name),
		Amount: r.Amount,
	}
}

// UpdateItemRequest represents the business request to update an item
type UpdateItemRequest struct {
	Name   *string `json:"name,omitempty"`
	Amount *uint   `json:"amount,omitempty"`
}

// Validate performs business validation on the update request
func (r *UpdateItemRequest) Validate() error {
	if r.Name != nil {
		if strings.TrimSpace(*r.Name) == "" {
			return domain.ErrItemNameRequired
		}
		
		if len(*r.Name) > 100 {
			return domain.ErrItemNameTooLong
		}
	}
	
	if r.Amount != nil && *r.Amount > 999999 {
		return domain.ErrItemAmountTooLarge
	}
	
	return nil
}

// HasUpdates checks if the request contains any updates
func (r *UpdateItemRequest) HasUpdates() bool {
	return r.Name != nil || r.Amount != nil
}

// PaginationRequest represents the business request for pagination
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Validate performs business validation and applies business rules for pagination
func (r *PaginationRequest) Validate() error {
	if r.Page < 0 {
		return domain.ErrInvalidPagination
	}
	
	if r.Limit < 0 {
		return domain.ErrInvalidPagination
	}
	
	// Business rule: Maximum limit is 100
	if r.Limit > 100 {
		return domain.ErrLimitTooLarge
	}
	
	return nil
}

// ApplyDefaults applies business default values
func (r *PaginationRequest) ApplyDefaults() {
	// Business rule: Default page is 1
	if r.Page <= 0 {
		r.Page = 1
	}
	
	// Business rule: Default limit is 10
	if r.Limit <= 0 {
		r.Limit = 10
	}
	
	// Business rule: Maximum limit is 100
	if r.Limit > 100 {
		r.Limit = 100
	}
}