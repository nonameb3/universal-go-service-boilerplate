package item

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
	"github.com/universal-go-service/boilerplate/internal/usecase"
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
	reqBody := new(request.AddItem)

	if err := c.BodyParser(reqBody); err != nil {
		h.logger.Error("Request Error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	createItem, err := h.itemUseCase.Create(&domain.Item{
		Amount: reqBody.Amount,
		Name:   reqBody.Name,
	})

	if err != nil {
		h.logger.Error("Create item error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(http.StatusCreated).JSON(createItem)
}

// GetItem retrieves an item by ID
func (h *Handler) GetItem(c *fiber.Ctx) error {
	params := new(request.GetItem)

	if err := c.ParamsParser(params); err != nil {
		h.logger.Error("Request Error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	item, err := h.itemUseCase.Get(params.Id)
	if err != nil {
		h.logger.Error("Get item error", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "item not found",
		})
	}

	return c.Status(http.StatusOK).JSON(item)
}

func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	params := new(request.GetItem)
	body := new(request.UpdateItem)

	if err := c.ParamsParser(params); err != nil {
		h.logger.Error("Request Error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request: params id is invalid",
		})
	}

	// get item
	existingItem, err := h.itemUseCase.Get(params.Id)
	if err != nil {
		h.logger.Error("Get item error", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "item not found",
		})
	}

	// Update only provided fields
	if body.Name != nil {
		existingItem.Name = *body.Name
	}
	if body.Amount != nil {
		existingItem.Amount = *body.Amount
	}

	// Update the existing item
	updatedItem, err := h.itemUseCase.Update(existingItem)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}

	return c.Status(http.StatusOK).JSON(updatedItem)
}

func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	params := new(request.GetItem)

	if err := c.ParamsParser(params); err != nil {
		h.logger.Error("Request Error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request: params id is invalid",
		})
	}

	// get item
	_, err := h.itemUseCase.Get(params.Id)
	if err != nil {
		h.logger.Error("Get item error", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "item not found",
		})
	}

	// delete item
	if err := h.itemUseCase.Delete(params.Id); err != nil {
		h.logger.Error("Delete item error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "delete item error",
		})
	}

	return c.Status(http.StatusOK).JSON(params)
}
