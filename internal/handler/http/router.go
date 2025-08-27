package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/universal-go-service/boilerplate/internal/handler/http/middleware"
	v1 "github.com/universal-go-service/boilerplate/internal/handler/http/v1"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	appLog "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

func NewRouter(app *fiber.App, itemUseCase usecase.ItemUseCase, l appLog.Logger) {
	// Middleware
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(middleware.Recovery(l))

	// Initialize V1 Router
	apiV1Group := app.Group("/api/v1")
	{
		v1.NewItemRouter(apiV1Group, itemUseCase, l)
	}
}
