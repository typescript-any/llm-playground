package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typescript-any/llm-playground/internal/handler"
	"github.com/typescript-any/llm-playground/internal/middleware"
)

func RegisterConversationRoutes(app *fiber.App, convHandler *handler.ConversationHandler, messageHandler *handler.MessageHandler) {
	convGroup := app.Group("/conversations", middleware.AuthMiddleware)

	// Conversations
	convGroup.Post("/", convHandler.CreateConversation)
	convGroup.Get("/:user_id", convHandler.ListConversations)

	// Messages inside conversation
	convGroup.Post("/:id/messages", messageHandler.SendMessage)
	convGroup.Post("/:id/messages/stream", messageHandler.StreamMessage)
}
