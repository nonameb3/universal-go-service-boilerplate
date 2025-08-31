package item

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// SetupRoutes sets up item routes
func SetupRoutes(apiV1Group fiber.Router, itemUseCase usecase.ItemUseCase, logger logger.Logger) {
	handler := New(itemUseCase, logger)

	itemGroup := apiV1Group.Group("/items")
	{
		// Cache GET routes for better performance
		itemGroup.Get("/", cache.New(cache.Config{
			Expiration:   30 * time.Second, // List/pagination changes more frequently
			CacheControl: true,             // Send proper HTTP cache headers
		}), handler.ListItems)

		// Non-cached routes (mutations should always execute)
		itemGroup.Post("/", handler.CreateItem)
		itemGroup.Post("/bulk", handler.BulkCreateItems)

		// Cache individual item GET with longer TTL
		itemGroup.Get("/:id", cache.New(cache.Config{
			Expiration:   30 * time.Second, // Individual items change less frequently
			CacheControl: true,             // Send proper HTTP cache headers
		}), handler.GetItem)

		itemGroup.Put("/:id", handler.UpdateItem)
		itemGroup.Delete("/:id", handler.DeleteItem)
	}
}
