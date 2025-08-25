package domain

type (
	InsertItemDto struct {
		BaseEntity
		Amount uint32 `json:"amount"`
		Name   string `json:"name"`
	}

	Item struct {
		BaseEntity
		Amount uint32 `json:"amount"`
		Name   string `json:"name"`
	}

	GetItemDto struct {
		Id uint32 `json:"id"`
	}
)
