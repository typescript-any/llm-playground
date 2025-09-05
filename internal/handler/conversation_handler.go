package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/typescript-any/llm-playground/internal/repository"
	service "github.com/typescript-any/llm-playground/internal/services"
)

type ConversationHandler struct {
	service *service.ConversationService
}

func NewConversationHandler(s *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: s,
	}
}

// POST /conversations
func (h *ConversationHandler) CreateConversation(c *fiber.Ctx) error {
	type reqBody struct {
		UserID string `json:"user_id"`
		Title  string `json:"title"`
	}

	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conv, err := h.service.CreateConversation(ctx, userID, body.Title)
	if err == repository.ErrInternal {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not create conversation"})
	}
	return c.Status(http.StatusCreated).JSON(conv)
}

// GET /conversations/:user_id
func (h *ConversationHandler) ListConversations(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conversations, err := h.service.ListConversations(ctx, userID)
	if err == repository.ErrNotFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no conversations found"})
	}
	if err == repository.ErrInternal {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch conversations"})
	}

	return c.JSON(conversations)
}
