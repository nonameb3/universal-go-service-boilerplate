package domain

import "github.com/google/uuid"

type (
	InsertItemDto struct {
		BaseEntity
		Amount uint   `json:"amount"`
		Name   string `json:"name"`
	}

	Item struct {
		BaseEntity
		Amount uint   `json:"amount"`
		Name   string `json:"name"`
	}

	GetItemDto struct {
		Id uuid.UUID `json:"id"`
	}
)
