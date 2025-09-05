package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/typescript-any/llm-playground/internal/db"
)

type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

type MessageRepository struct{}

// NewMessageRepository constructor
func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

// SaveMessage inserts a message into conversation
func (r *MessageRepository) SaveMessage(ctx context.Context, conversationID uuid.UUID, role, content string) (*Message, error) {
	query := `
		INSERT INTO messages (id, conversation_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, conversation_id, role, content, created_at
	`

	id := uuid.New()
	createdAt := time.Now()

	row := db.GetPool().QueryRow(ctx, query, id, conversationID, role, content, createdAt)

	var m Message
	if err := row.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
		return nil, err
	}

	return &m, nil
}
