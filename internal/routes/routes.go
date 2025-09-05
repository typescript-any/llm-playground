package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typescript-any/llm-playground/internal/handler"
)

func RegisterConversationRoutes(app *fiber.App, h *handler.ConversationHandler) {
	app.Post("/conversations", h.CreateConversation)
	app.Get("/conversations/:user_id", h.ListConversations)
}
