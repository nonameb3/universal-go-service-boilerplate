package domain

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
		Id uint `json:"id"`
	}
)
