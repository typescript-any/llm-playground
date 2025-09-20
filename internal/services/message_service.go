package service

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
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

func defaultModel(model string) string {
	if model == "" {
		return "gpt-4o"
	}
	return model
}

func (s *MessageService) SendMessage(ctx context.Context, params MessageSendParams) (*models.ChatMessage, error) {
	// 1. Save user message
	_, err := s.repo.SaveMessage(ctx, repository.MessageSaveParams{
		ConversationID: params.ConversationID,
		Role:           models.RoleUser,
		Content:        params.Content,
	})
	if err != nil {
		return nil, err
	}

	// 2. Get conversation history
	history, err := s.repo.GetMessagesByConversation(ctx, repository.MessageListParams{
		ConversationID: params.ConversationID,
		Limit:          100, // Get all recent messages for context
	})
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
		Model:     defaultModel(params.Model),
		MaxTokens: openai.Int(500),
	})
	if err != nil {
		return nil, err
	}

	reply := resp.Choices[0].Message.Content

	// 4. Save assistant reply
	if _, err := s.repo.SaveMessage(ctx, repository.MessageSaveParams{
		ConversationID: params.ConversationID,
		Role:           models.RoleAssistant,
		Content:        reply,
	}); err != nil {
		return nil, err
	}

	return &models.ChatMessage{
		Role:    models.RoleAssistant,
		Content: reply,
	}, nil
}

func (s *MessageService) StreamMessage(ctx context.Context, params MessageStreamParams) (*ssestream.Stream[openai.ChatCompletionChunk], *openai.ChatCompletionAccumulator, error) {
	// 1. Save user message
	_, err := s.repo.SaveMessage(ctx, repository.MessageSaveParams{
		ConversationID: params.ConversationID,
		Role:           models.RoleUser,
		Content:        params.Content,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 2. Fetch recent history
	history, err := s.repo.GetMessagesByConversation(ctx, repository.MessageListParams{
		ConversationID: params.ConversationID,
		Limit:          20,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user history: %w", err)
	}

	// 3. Convert history into SDK messages
	var messages []openai.ChatCompletionMessageParamUnion
	for _, message := range history {
		switch message.Role {
		case models.RoleUser:
			messages = append(messages, openai.UserMessage(message.Content))

		case models.RoleAssistant:
			messages = append(messages, openai.AssistantMessage(message.Content))

		case models.RoleSystem:
			messages = append(messages, openai.SystemMessage(message.Content))
		default:
			// default to user if unknown
			messages = append(messages, openai.UserMessage(message.Content))
		}
	}

	// 4. Add the current user message to the end
	messages = append(messages, openai.UserMessage(params.Content))

	// 5. Create streaming request
	stream := s.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:       defaultModel(params.Model),
		Messages:    messages,
		MaxTokens:   openai.Int(500),
		Temperature: openai.Float(0.7),
		TopP:        openai.Float(1.0),
	})
	acc := openai.ChatCompletionAccumulator{}

	return stream, &acc, nil

}

// SaveAssistantMessage persists the assistant text after streaming completes.
func (s *MessageService) SaveAssistantMessage(ctx context.Context, params MessageSaveParams) (*models.Message, error) {
	return s.repo.SaveMessage(ctx, repository.MessageSaveParams{
		ConversationID: params.ConversationID,
		Role:           models.RoleAssistant,
		Content:        params.Content,
	})
}
