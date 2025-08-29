package item

import (
	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// SetupRoutes sets up item routes
func SetupRoutes(apiV1Group fiber.Router, itemUseCase usecase.ItemUseCase, logger logger.Logger) {
	handler := New(itemUseCase, logger)

	itemGroup := apiV1Group.Group("/items")
	{
		itemGroup.Get("/", handler.ListItems)
		itemGroup.Post("/", handler.CreateItem)
		itemGroup.Post("/bulk", handler.BulkCreateItems)
		itemGroup.Get("/:id", handler.GetItem)
		itemGroup.Put("/:id", handler.UpdateItem)
		itemGroup.Delete("/:id", handler.DeleteItem)
	}
}
