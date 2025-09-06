package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	service "github.com/typescript-any/llm-playground/internal/services"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{
		service: s,
	}
}

// handler/message_handler.go
func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid conversation_id")
	}
	var req struct {
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "user_id content model are required")
	}

	reply, err := h.service.SendMessage(c.Context(), convID, req.Content, req.Model)
	if err != nil {
		return err
	}

	return c.JSON(reply)
}
