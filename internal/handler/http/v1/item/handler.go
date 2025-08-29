package item

import (
	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/handler/http/errors"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
	"github.com/universal-go-service/boilerplate/internal/usecase"
	"github.com/universal-go-service/boilerplate/internal/usecase/item/dto"
	logger "github.com/universal-go-service/boilerplate/pkg/providers/logger"
)

// Handler represents item handler
type Handler struct {
	itemUseCase     usecase.ItemUseCase
	logger          logger.Logger
	errorMapper     *errors.ErrorMapper
	stdResponses    *errors.StandardResponses
}

// New creates a new item handler
func New(itemUseCase usecase.ItemUseCase, logger logger.Logger) *Handler {
	return &Handler{
		itemUseCase:  itemUseCase,
		logger:       logger,
		errorMapper:  errors.NewErrorMapper(),
		stdResponses: errors.NewStandardResponses(),
	}
}

// CreateItem creates a new item
func (h *Handler) CreateItem(c *fiber.Ctx) error {
	// HTTP request parsing
	var httpReq request.AddItem
	if err := c.BodyParser(&httpReq); err != nil {
		h.logger.Error("Request parsing error", err)
		return h.stdResponses.BadRequest(c, "invalid request format")
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &dto.CreateItemRequest{
		Name:   httpReq.Name,
		Amount: httpReq.Amount,
	}

	// Delegate ALL business logic to UseCase
	item, err := h.itemUseCase.Create(useCaseReq)
	if err != nil {
		return h.errorMapper.SendError(c, err)
	}

	// HTTP response formatting
	return h.stdResponses.Created(c, item)
}

// GetItem retrieves an item by ID
func (h *Handler) GetItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return h.stdResponses.BadRequest(c, "id parameter is required")
	}

	// Delegate ALL business logic to UseCase
	item, err := h.itemUseCase.Get(id)
	if err != nil {
		return h.errorMapper.SendError(c, err)
	}

	// HTTP response formatting
	return h.stdResponses.OK(c, item)
}

// ListItems retrieves items with pagination
func (h *Handler) ListItems(c *fiber.Ctx) error {
	// HTTP query parameter parsing
	var httpReq request.ListItems
	if err := c.QueryParser(&httpReq); err != nil {
		h.logger.Error("Query parsing error", err)
		return h.stdResponses.BadRequest(c, "invalid query parameters")
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &dto.PaginationRequest{
		Page:  httpReq.Page,
		Limit: httpReq.Limit,
	}

	// Delegate ALL business logic (including defaults) to UseCase
	items, err := h.itemUseCase.GetWithPagination(useCaseReq)
	if err != nil {
		return h.errorMapper.SendError(c, err)
	}

	// HTTP response formatting
	return h.stdResponses.OK(c, items)
}

// UpdateItem updates an existing item
func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return h.stdResponses.BadRequest(c, "id parameter is required")
	}

	// HTTP body parsing
	var httpReq request.UpdateItem
	if err := c.BodyParser(&httpReq); err != nil {
		h.logger.Error("Request parsing error", err)
		return h.stdResponses.BadRequest(c, "invalid request format")
	}

	// Convert HTTP request to UseCase request
	useCaseReq := &dto.UpdateItemRequest{
		Name:   httpReq.Name,
		Amount: httpReq.Amount,
	}

	// Delegate ALL business logic to UseCase
	updatedItem, err := h.itemUseCase.Update(id, useCaseReq)
	if err != nil {
		return h.errorMapper.SendError(c, err)
	}

	// HTTP response formatting
	return h.stdResponses.OK(c, updatedItem)
}

// DeleteItem deletes an existing item
func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	// HTTP parameter parsing
	id := c.Params("id")
	if id == "" {
		return h.stdResponses.BadRequest(c, "id parameter is required")
	}

	// Delegate ALL business logic to UseCase
	if err := h.itemUseCase.Delete(id); err != nil {
		return h.errorMapper.SendError(c, err)
	}

	// HTTP response formatting
	return h.stdResponses.SuccessMessage(c, "item deleted successfully", fiber.Map{
		"id": id,
	})
}
