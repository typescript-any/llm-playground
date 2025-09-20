package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typescript-any/llm-playground/internal/handler"
	"github.com/typescript-any/llm-playground/internal/middleware"
)

func RegisterConversationRoutes(router fiber.Router, convHandler *handler.ConversationHandler, messageHandler *handler.MessageHandler) {
	convGroup := router.Group("/conversations", middleware.AuthMiddleware)

	// Conversations
	convGroup.Post("/", convHandler.CreateConversation)
	convGroup.Get("/:user_id", convHandler.ListConversations)
	convGroup.Post("/new", convHandler.CreateNewConversation)

	// Messages inside conversation
	convGroup.Post("/:id/messages", messageHandler.SendMessage)
	convGroup.Post("/:id/messages/stream", messageHandler.StreamMessage)
}
