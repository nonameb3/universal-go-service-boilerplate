package dto

import (
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/domain/entities"
)

type BulkCreateRequest struct {
	Items []CreateItemRequest `json:"items"`
}

func (req *BulkCreateRequest) Validate() error {
	if len(req.Items) == 0 {
		return domain.ErrItemNameRequired
	}

	for i, item := range req.Items {
		if err := item.Validate(); err != nil {
			return err
		}
		if i >= 1000 {
			return domain.ErrInvalidInput
		}
	}

	return nil
}

func (req *BulkCreateRequest) ToEntities() []*entities.Item {
	items := make([]*entities.Item, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = itemReq.ToEntity()
	}
	return items
}