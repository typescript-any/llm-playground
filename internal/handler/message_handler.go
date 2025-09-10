package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	service "github.com/typescript-any/llm-playground/internal/services"
	"github.com/valyala/fasthttp"
)

type MessageHandler struct {
	service *service.MessageService
}

type MessageStart struct {
	Type      string `json:"type"`
	Model     string `json:"model"`
	ConvID    string `json:"conversation_id"`
	CreatedAt int64  `json:"created_at"`
}

type ContentDelta struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"delta"`
}

type MessageComplete struct {
	Type       string `json:"type"`
	StopReason string `json:"stop_reason"`
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

// StreamMessage handles streaming AI responses via Server-Sent Events (SSE)
func (h *MessageHandler) StreamMessage(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid conversation_id")
	}

	var req struct {
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "content and model are required")
	}

	stream, acc, err := h.service.StreamMessage(c.Context(), convID, req.Content, req.Model)
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
			if _, err := h.service.SaveAssistantMessage(context.Background(), convID, aiContent); err != nil {
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
func (h *MessageHandler) sendEvent(w *bufio.Writer, eventType string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, string(jsonData))
	w.Flush()
}

// Helper function to send error events
func (h *MessageHandler) sendErrorEvent(w *bufio.Writer, errorMsg string) {
	errorData := map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	}
	jsonData, _ := json.Marshal(errorData)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(jsonData))
	w.Flush()
}
