package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typescript-any/llm-playground/internal/handler"
	"github.com/typescript-any/llm-playground/internal/middleware"
)

func RegisterConversationRoutes(app *fiber.App, h *handler.ConversationHandler) {
	convGroup := app.Group("/conversations", middleware.AuthMiddleware)
	convGroup.Post("/", h.CreateConversation)
	convGroup.Get("/:user_id", h.ListConversations)
}
