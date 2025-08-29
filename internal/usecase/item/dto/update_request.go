package dto

import (
	"strings"
	
	"github.com/universal-go-service/boilerplate/internal/domain"
)

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