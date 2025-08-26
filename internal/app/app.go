package app

import (
	"github.com/universal-go-service/boilerplate/config"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

func Run(cfg *config.Config) {
	config := logger.LoggerConfig{
		Type:        "boilerplate",
		ServiceName: "go-service",
	}

	l := logger.NewCentralizedLogger(config)

	// Initial Repository
}
