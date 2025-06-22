package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hinphansa/7-solutions-challenge/pkg/utils"
)

// validate request
func MustValid(c *fiber.Ctx, req any) error {
	// parse request body
	if err := c.BodyParser(req); err != nil {
		if errors := utils.ValidateStruct(req); len(errors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}

		// unknown error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return nil
}
