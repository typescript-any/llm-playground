package service

import (
	"context"

	"github.com/google/uuid"
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

func (s *ConversationService) ListConversations(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	return s.repo.GetConversationsByUser(ctx, userID)
}
