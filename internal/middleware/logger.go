package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func RequestResponseLogger(c *fiber.Ctx) error {
	// --- Log request ---
	reqBody := ""
	if c.Request().Body() != nil {
		reqBody = string(c.Request().Body())
	}
	log.Printf("[REQUEST] %s %s\nHeaders: %v\nBody: %s\n", c.Method(), c.Path(), c.GetReqHeaders(), reqBody)

	// --- Process request ---
	err := c.Next() // call next handler/middleware
	if err != nil {
		return err // let error handler handle it
	}

	// --- Log response ---
	resBody := string(c.Response().Body())
	log.Printf("[RESPONSE] %s %s\nStatus: %d\nBody: %s\n", c.Method(), c.Path(), c.Response().StatusCode(), resBody)

	return nil
}
