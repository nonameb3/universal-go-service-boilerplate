package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

func NewItemRouter(apiV1Group fiber.Router, itemUseCase usecase.ItemUseCase, l logger.Logger) {
	itemHandler := &V1{
		t: itemUseCase,
		l: l,
	}

	itemGroup := apiV1Group.Group("/items")
	{
		itemGroup.Post("/", itemHandler.CreateItem)
		itemGroup.Get("/:id", itemHandler.GetItem)
	}
}
