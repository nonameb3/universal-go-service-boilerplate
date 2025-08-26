package v1

import "github.com/gofiber/fiber/v2"

func (v *V1) CreateItem(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (v *V1) GetItem(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
