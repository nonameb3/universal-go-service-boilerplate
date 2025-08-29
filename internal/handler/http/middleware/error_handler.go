package middleware

import (
	"net/http"
	
	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/domain"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// ErrorHandler provides centralized error handling for HTTP handlers
type ErrorHandler struct {
	logger logger.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError maps business errors to appropriate HTTP responses
func (eh *ErrorHandler) HandleError(c *fiber.Ctx, err error) error {
	eh.logger.Error("Handler error occurred", err)
	
	// Map domain errors to HTTP status codes
	switch err {
	case domain.ErrItemNotFound:
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "item not found",
		})
		
	case domain.ErrItemNameRequired:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "item name is required",
		})
		
	case domain.ErrItemNameTooLong:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "item name cannot exceed 100 characters",
		})
		
	case domain.ErrItemAmountTooLarge:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "item amount cannot exceed 999999",
		})
		
	case domain.ErrInvalidPagination:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid pagination parameters",
		})
		
	case domain.ErrLimitTooLarge:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "limit cannot exceed 100",
		})
		
	default:
		// Generic server error for unknown errors
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
}

// ErrorResponse creates standardized error responses
func (eh *ErrorHandler) ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}

// SuccessResponse creates standardized success responses
func (eh *ErrorHandler) SuccessResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}