package entities

import (
	"strings"
)

// Item represents the item business entity
type Item struct {
	BaseEntity
	Amount uint   `json:"amount"`
	Name   string `json:"name"`
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