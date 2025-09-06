package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error with request info
	log.Printf("[ERROR] %s %s - %d - %s\n", c.Method(), c.Path(), code, err.Error())
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
