package errors

import (
	"net/http"
	
	"github.com/gofiber/fiber/v2"
)

// StandardResponses provides standard HTTP response helpers
type StandardResponses struct{}

// NewStandardResponses creates standard response helpers
func NewStandardResponses() *StandardResponses {
	return &StandardResponses{}
}

// BadRequest returns a 400 Bad Request response
func (sr *StandardResponses) BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"error": message,
	})
}

// NotFound returns a 404 Not Found response
func (sr *StandardResponses) NotFound(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}

// InternalServerError returns a 500 Internal Server Error response
func (sr *StandardResponses) InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": message,
	})
}

// Created returns a 201 Created response
func (sr *StandardResponses) Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(data)
}

// OK returns a 200 OK response
func (sr *StandardResponses) OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(data)
}

// SuccessMessage returns a standardized success message
func (sr *StandardResponses) SuccessMessage(c *fiber.Ctx, message string, additionalData ...fiber.Map) error {
	response := fiber.Map{
		"message": message,
	}
	
	// Merge additional data if provided
	if len(additionalData) > 0 {
		for key, value := range additionalData[0] {
			response[key] = value
		}
	}
	
	return c.Status(http.StatusOK).JSON(response)
}