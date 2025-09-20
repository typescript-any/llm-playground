package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typescript-any/llm-playground/internal/models"
)

type MessageRepo struct {
	db *pgxpool.Pool
}

// NewMessageRepo constructor
func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

// SaveMessage inserts a message into conversation
func (r *MessageRepo) SaveMessage(ctx context.Context, params MessageSaveParams) (*models.Message, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM conversations WHERE id=$1)", params.ConversationID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("conversation %s does not exist", params.ConversationID)
	}

	query := `
		INSERT INTO messages (id, conversation_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, conversation_id, role, content, created_at
	`

	id := uuid.New()
	createdAt := time.Now()

	row := r.db.QueryRow(ctx, query, id, params.ConversationID, params.Role, params.Content, createdAt)

	var m models.Message
	if err := row.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
		fmt.Printf("Failed to save message %v", err)
		return nil, ErrInternal
	}

	return &m, nil
}

// List messages
func (r *MessageRepo) GetMessages(ctx context.Context, convID uuid.UUID) ([]models.Message, error) {
	query := `SELECT id, conversation_id, role, content, created_at from messages`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		fmt.Printf("Failed to select messages %v", err)
		return nil, ErrInternal
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.Role, &message.Content, &message.CreatedAt); err != nil {
			return nil, ErrInternal
		}
		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, ErrNotFound
	}

	return messages, nil
}

// Get messages by conversation
func (r *MessageRepo) GetMessagesByConversation(ctx context.Context, params MessageListParams) ([]models.Message, error) {
	query := `SELECT id, conversation_id, role, content, created_at
			  FROM messages
			  WHERE conversation_id = $1
			  ORDER BY created_at ASC
			  LIMIT $2 
			  `
	rows, err := r.db.Query(ctx, query, params.ConversationID, params.Limit)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.Role, &message.Content, &message.CreatedAt); err != nil {
			return nil, ErrInternal
		}
		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, ErrNotFound
	}
	return messages, nil

}
