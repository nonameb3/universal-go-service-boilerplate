package errors

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/domain"
)

// HTTPError represents an HTTP error with status code and message
type HTTPError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
}

// ErrorMapper provides mapping between domain errors and HTTP errors
type ErrorMapper struct{}

// NewErrorMapper creates a new error mapper
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapDomainError maps domain errors to HTTP errors
func (em *ErrorMapper) MapDomainError(err error) HTTPError {
	switch err {
	case domain.ErrItemNotFound:
		return HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    "item not found",
		}

	case domain.ErrItemNameRequired:
		return HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "item name is required",
		}

	case domain.ErrItemNameTooLong:
		return HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "item name cannot exceed 100 characters",
		}

	case domain.ErrItemAmountTooLarge:
		return HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "item amount cannot exceed 999999",
		}

	case domain.ErrInvalidPagination:
		return HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid pagination parameters",
		}

	case domain.ErrLimitTooLarge:
		return HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "limit cannot exceed 100",
		}

	case domain.ErrItemAlreadyExists:
		return HTTPError{
			StatusCode: http.StatusConflict,
			Message:    "Item with same name already exists",
		}

	default:
		return HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
		}
	}
}

// SendError sends a standardized error response
func (em *ErrorMapper) SendError(c *fiber.Ctx, err error) error {
	httpErr := em.MapDomainError(err)
	return c.Status(httpErr.StatusCode).JSON(httpErr)
}
