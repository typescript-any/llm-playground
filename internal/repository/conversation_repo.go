package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typescript-any/llm-playground/internal/models"
)

type ConversationRepo struct {
	db *pgxpool.Pool
}

// NewConversationRepository constructor
func NewConversationRepo(db *pgxpool.Pool) *ConversationRepo {
	return &ConversationRepo{
		db: db,
	}
}

func (r *ConversationRepo) CreateConversation(ctx context.Context, userID uuid.UUID, title string) (models.Conversation, error) {
	var conv models.Conversation
	query := `INSERT INTO conversations ( user_id, title)
			  VALUES ($1, $2)
		      RETURNING id, user_id, title, created_at`
	err := r.db.QueryRow(ctx, query, userID, title).Scan(
		&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt,
	)

	if err != nil {
		log.Printf("Error in creating conversation: %v", err)
		return models.Conversation{}, ErrInternal
	}

	return conv, nil
}

func (r *ConversationRepo) GetConversationsByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Conversation, error) {
	query := `SELECT id, user_id, title, created_at
			  FROM conversations
			  WHERE user_id = $1
			  ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt); err != nil {
			return nil, ErrInternal
		}
		conversations = append(conversations, conv)
	}

	if len(conversations) == 0 {
		return nil, ErrNotFound
	}
	return conversations, nil
}
