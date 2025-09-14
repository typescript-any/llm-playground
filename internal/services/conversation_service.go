package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/typescript-any/llm-playground/internal/llm"
	"github.com/typescript-any/llm-playground/internal/models"
	"github.com/typescript-any/llm-playground/internal/repository"
)

type ConversationService struct {
	repo *repository.ConversationRepo
}

func NewConversationService(repo *repository.ConversationRepo) *ConversationService {
	return &ConversationService{repo: repo}
}

func (s *ConversationService) CreateConversation(ctx context.Context, userID uuid.UUID, title string) (models.Conversation, error) {
	return s.repo.CreateConversation(ctx, userID, title)
}

func (s *ConversationService) ListConversations(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Conversation, error) {
	return s.repo.GetConversationsByUser(ctx, userID, offset, limit)
}

func (s *ConversationService) CreateNewConversation(ctx context.Context, userID uuid.UUID, content string, model string) (models.Conversation, error) {
	client := llm.NewClient()

	prompt := `
		You are a system that generates very short conversation titles.

		Example:
		User: How can I format dates in React?
		Title: Formatting Dates in React

		User: ` + content + `
		Title:
		`

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:     defaultModel(model),
		MaxTokens: openai.Int(500),
	})

	if err != nil || len(resp.Choices) == 0 {
		log.Info("Error generating title:", err)
		return s.repo.CreateConversation(ctx, userID, "New Conversation")
	}

	title := resp.Choices[0].Message.Content
	log.Info("Generated title:", title)
	return s.repo.CreateConversation(ctx, userID, title)

}
