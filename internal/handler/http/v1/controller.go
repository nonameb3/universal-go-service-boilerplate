package v1

import (
	"github.com/universal-go-service/boilerplate/internal/usecase"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// V1 -.
type V1 struct {
	t usecase.ItemUseCase
	l logger.Logger
}
