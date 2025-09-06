package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/typescript-any/llm-playground/internal/models"
	"github.com/typescript-any/llm-playground/internal/repository"
)

type MessageService struct {
	repo   *repository.MessageRepo
	client *openai.Client
}

// Constructor function of MessageService
func NewMessageService(r *repository.MessageRepo, c *openai.Client) *MessageService {
	return &MessageService{
		repo:   r,
		client: c,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, convID uuid.UUID, content, model string) (*models.ChatMessage, error) {
	// 1. Save user message
	_, err := s.repo.SaveMessage(ctx, convID, models.RoleUser, content)
	if err != nil {
		return nil, err
	}

	// 2. Get conversation history
	history, err := s.repo.GetMessages(ctx, convID)
	if err != nil {
		return nil, err
	}

	// Convert history to openai message
	var messages []openai.ChatCompletionMessageParamUnion
	for _, m := range history {
		if m.Role == models.RoleUser {
			messages = append(messages, openai.UserMessage(m.Content))
		} else {
			messages = append(messages, openai.AssistantMessage(m.Content))
		}
	}

	// 3. Call OpenRouter via go-openai
	resp, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:  messages,
		Model:     defaultModel(model),
		MaxTokens: openai.Int(500),
	})
	if err != nil {
		return nil, err
	}

	reply := resp.Choices[0].Message.Content

	// 4. Save assistant reply
	if _, err := s.repo.SaveMessage(ctx, convID, models.RoleAssistant, reply); err != nil {
		return nil, err
	}

	return &models.ChatMessage{
		Role:    models.RoleAssistant,
		Content: reply,
	}, nil
}

func defaultModel(model string) string {
	if model == "" {
		return "gpt-4o"
	}
	return model
}
