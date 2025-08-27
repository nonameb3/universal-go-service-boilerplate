package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/item"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// SetupRoutes sets up all v1 API routes
func SetupRoutes(apiV1Group fiber.Router, itemUseCase usecase.ItemUseCase, logger logger.Logger) {
	// Setup item routes
	item.SetupRoutes(apiV1Group, itemUseCase, logger)
	
	// Add more domain routes here:
	// user.SetupRoutes(apiV1Group, userUseCase, logger)
	// order.SetupRoutes(apiV1Group, orderUseCase, logger)
}
