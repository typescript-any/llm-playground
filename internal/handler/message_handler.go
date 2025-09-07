package handler

import (
	"bufio"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	service "github.com/typescript-any/llm-playground/internal/services"
	"github.com/valyala/fasthttp"
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
		// Use the accumulator returned from service instead of creating new one
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta.Content
				if delta != "" {
					fmt.Fprintf(w, "data: %s\n\n", delta)
					w.Flush()
				}
			}
		}

		if stream.Err() != nil {
			fmt.Fprintf(w, "event: error\ndata: %v\n\n", stream.Err())
			w.Flush()
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

		fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
		w.Flush()
	}))

	return nil
}
