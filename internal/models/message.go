package models

import "github.com/google/uuid"

type Message struct {
	ID             uuid.UUID `db:"id"`
	ConversationID uuid.UUID `db:"conversation_id"`
	Role           string    `db:"role"` // "user" or "assistant"
	Content        string    `db:"content"`
	CreatedAt      string    `db:"created_at"`
}
