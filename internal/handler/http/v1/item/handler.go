package item

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	itemUseCase "github.com/universal-go-service/boilerplate/internal/usecase/item"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// Handler represents item handler
type Handler struct {
	itemUseCase usecase.ItemUseCase
	logger      logger.Logger
}

// New creates a new item handler
func New(itemUseCase usecase.ItemUseCase, logger logger.Logger) *Handler {
	return &Handler{
		itemUseCase: itemUseCase,
		logger:      logger,
	}
}

// CreateItem creates a new item
func (h *Handler) CreateItem(c *fiber.Ctx) error {
	// HTTP request parsing
	var httpReq request.AddItem
	if err := c.BodyParser(&httpReq); err != nil {
		h.logger.Error("Request parsing error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &itemUseCase.CreateItemRequest{
		Name:   httpReq.Name,
		Amount: httpReq.Amount,
	}

	// Delegate ALL business logic to UseCase
	item, err := h.itemUseCase.Create(useCaseReq)
	if err != nil {
		return h.handleError(c, err)
	}

	// HTTP response formatting
	return c.Status(http.StatusCreated).JSON(item)
}

// GetItem retrieves an item by ID
func (h *Handler) GetItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	// Delegate ALL business logic to UseCase
	item, err := h.itemUseCase.Get(id)
	if err != nil {
		return h.handleError(c, err)
	}

	// HTTP response formatting
	return c.Status(http.StatusOK).JSON(item)
}

// ListItems retrieves items with pagination
func (h *Handler) ListItems(c *fiber.Ctx) error {
	// HTTP query parameter parsing
	var httpReq request.ListItems
	if err := c.QueryParser(&httpReq); err != nil {
		h.logger.Error("Query parsing error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid query parameters",
		})
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &itemUseCase.PaginationRequest{
		Page:  httpReq.Page,
		Limit: httpReq.Limit,
	}

	// Delegate ALL business logic (including defaults) to UseCase
	items, err := h.itemUseCase.GetWithPagination(useCaseReq)
	if err != nil {
		return h.handleError(c, err)
	}

	// HTTP response formatting
	return c.Status(http.StatusOK).JSON(items)
}

// UpdateItem updates an existing item
func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	// HTTP body parsing
	var httpReq request.UpdateItem
	if err := c.BodyParser(&httpReq); err != nil {
		h.logger.Error("Request parsing error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &itemUseCase.UpdateItemRequest{
		Name:   httpReq.Name,
		Amount: httpReq.Amount,
	}

	// Delegate ALL business logic to UseCase
	updatedItem, err := h.itemUseCase.Update(id, useCaseReq)
	if err != nil {
		return h.handleError(c, err)
	}

	// HTTP response formatting
	return c.Status(http.StatusOK).JSON(updatedItem)
}

// DeleteItem deletes an existing item
func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	// Delegate ALL business logic to UseCase
	if err := h.itemUseCase.Delete(id); err != nil {
		return h.handleError(c, err)
	}

	// HTTP response formatting
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "item deleted successfully",
		"id":      id,
	})
}

// handleError maps business errors to appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	h.logger.Error("Handler error occurred", err)

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
