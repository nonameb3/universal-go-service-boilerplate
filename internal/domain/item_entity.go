package domain

import (
	"strings"
)

type Item struct {
	BaseEntity
	Amount uint   `json:"amount"`
	Name   string `json:"name"`
}

// Validate performs business validation on the Item entity
func (i *Item) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return ErrItemNameRequired
	}
	
	if len(i.Name) > 100 {
		return ErrItemNameTooLong
	}
	
	// Business rule: Amount cannot exceed 999999
	if i.Amount > 999999 {
		return ErrItemAmountTooLarge
	}
	
	return nil
}

// UpdateFrom applies partial updates to the item with business rules
func (i *Item) UpdateFrom(name *string, amount *uint) {
	if name != nil {
		i.Name = strings.TrimSpace(*name)
	}
	
	if amount != nil {
		i.Amount = *amount
	}
}

// IsEmpty checks if the item has meaningful data
func (i *Item) IsEmpty() bool {
	return strings.TrimSpace(i.Name) == "" && i.Amount == 0
}
