package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/typescript-any/llm-playground/internal/repository"
	service "github.com/typescript-any/llm-playground/internal/services"
	"github.com/valyala/fasthttp"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
	messageService      *service.MessageService
}

func NewConversationHandler(conversationService *service.ConversationService, messageService *service.MessageService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		messageService:      messageService,
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

	conv, err := h.conversationService.CreateConversation(ctx, service.ConversationCreateParams{
		UserID: userID,
		Title:  body.Title,
	})
	if err == repository.ErrInternal {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not create conversation")
	}
	return c.Status(http.StatusCreated).JSON(conv)
}

// GET /conversations/:user_id
func (h *ConversationHandler) ListConversations(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("user_id"))
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conversations, err := h.conversationService.ListConversations(ctx, service.ConversationListParams{
		UserID: userID,
		Offset: skip,
		Limit:  limit,
	})
	if err == repository.ErrNotFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no conversations found"})
	}
	if err == repository.ErrInternal {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "could not get conversation")
	}

	return c.JSON(conversations)
}

// POST /conversations/new
func (h *ConversationHandler) CreateNewConversation(c *fiber.Ctx) error {
	type reqBody struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	var req reqBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conv, err := h.conversationService.CreateNewConversation(ctx, service.ConversationNewParams{
		UserID:  userID,
		Content: req.Content,
		Model:   req.Model,
	})
	if err == repository.ErrInternal {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not create conversation")
	}

	convID := conv.ID

	stream, acc, err := h.messageService.StreamMessage(c.Context(), service.MessageStreamParams{
		ConversationID: convID,
		Content:        req.Content,
		Model:          req.Model,
	})
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*") // Add CORS if needed

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		// 1. Send start event
		messageStart := MessageStart{
			Type:      "message_start",
			Model:     req.Model,
			ConvID:    convID.String(),
			CreatedAt: time.Now().Unix(),
		}
		h.sendEvent(w, "message_start", messageStart)

		// 2. Stream content deltas
		// Use the accumulator returned from service instead of creating new one
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta.Content
				if delta != "" {
					// Send as JSON structured event
					contentDelta := ContentDelta{
						Type:  "content_block_delta",
						Index: 0, // Assuming single message for now
					}
					contentDelta.Delta.Type = "text_delta"
					contentDelta.Delta.Value = delta
					h.sendEvent(w, "content_block_delta", contentDelta)
				}
			}
		}

		if stream.Err() != nil {
			h.sendErrorEvent(w, stream.Err().Error())
			return
		}

		// Save AI full message after completion
		if len(acc.Choices) > 0 {
			aiContent := acc.Choices[0].Message.Content
			if _, err := h.messageService.SaveAssistantMessage(context.Background(), service.MessageSaveParams{
				ConversationID: convID,
				Content:        aiContent,
			}); err != nil {
				fmt.Fprintf(w, "event: error\ndata: %v\n\n", err)
				w.Flush()
				return
			}
		}

		// 3. Send completion event
		messageComplete := MessageComplete{
			Type:       "message_complete",
			StopReason: "end_turn",
		}
		h.sendEvent(w, "message_complete", messageComplete)
	}))

	return nil
}

// Helper function to send structured events
func (h *ConversationHandler) sendEvent(w *bufio.Writer, eventType string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, string(jsonData))
	w.Flush()
}

// Helper function to send error events
func (h *ConversationHandler) sendErrorEvent(w *bufio.Writer, errorMsg string) {
	errorData := map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	}
	jsonData, _ := json.Marshal(errorData)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(jsonData))
	w.Flush()
}
