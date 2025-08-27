package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/universal-go-service/boilerplate/internal/domain"
	"github.com/universal-go-service/boilerplate/internal/handler/http/v1/request"
)

func (v *V1) CreateItem(c *fiber.Ctx) error {
	reqBody := new(request.AddItem)

	if err := c.BodyParser(reqBody); err != nil {
		v.l.Error("Request Error", err)
		return errorResponse(c, http.StatusBadRequest, "bad request")
	}

	v.iUC.Create(&domain.Item{
		Amount: reqBody.Amount,
		Name:   reqBody.Name,
	})

	return c.Status(http.StatusOK).JSON(reqBody)
}

func (v *V1) GetItem(c *fiber.Ctx) error {
	request := new(request.GetItem)

	if err := c.ParamsParser(request); err != nil {
		v.l.Error("Request Error", err)
		return errorResponse(c, http.StatusBadRequest, "bad request")
	}

	item, _ := v.iUC.Get(request.Id)

	return c.Status(http.StatusOK).JSON(item)
}
