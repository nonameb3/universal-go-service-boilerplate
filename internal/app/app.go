package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/universal-go-service/boilerplate/config"
	"github.com/universal-go-service/boilerplate/internal/handler/http"
	"github.com/universal-go-service/boilerplate/internal/repository/item"
	itemUC "github.com/universal-go-service/boilerplate/internal/usecase/item"
	"github.com/universal-go-service/boilerplate/pkg/httpserver"
	"github.com/universal-go-service/boilerplate/pkg/providers/database"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
	"github.com/universal-go-service/boilerplate/pkg/types"
)

func Run(cfg *config.Config, db database.DatabaseProvider) {
	config := logger.LoggerConfig{
		Type:        "boilerplate",
		ServiceName: "go-service",
	}

	// Initial Logger
	l := logger.NewCentralizedLogger(config)

	// Use the database instance passed from main.go
	pg := db

	// Initial UseCase
	itemUseCase := itemUC.NewItemUseCase(item.NewItemRepository(pg.GetDB(), l), pg, l)

	// Initial Server
	httpServer := httpserver.New(cfg.Server.Port)

	// Initial HealthCheck Middleware
	httpServer.App.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/health",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			// Check database connection
			err := pg.Health()
			return err == nil
		},
		ReadinessEndpoint: "/health",
	}))

	// Initial Router
	http.NewRouter(httpServer.App, itemUseCase, l)

	// Start Server
	l.Info("🚀 Server starting",
		types.Field{Key: "host", Value: cfg.Server.Host},
		types.Field{Key: "port", Value: cfg.Server.Port})
	httpServer.Start()
}
