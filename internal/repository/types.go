package repository

import "github.com/google/uuid"

// ConversationCreateParams holds parameters for creating a conversation
type ConversationCreateParams struct {
	UserID uuid.UUID
	Title  string
}

// ConversationListParams holds parameters for listing conversations
type ConversationListParams struct {
	UserID uuid.UUID
	Offset int
	Limit  int
}

// MessageSaveParams holds parameters for saving a message
type MessageSaveParams struct {
	ConversationID uuid.UUID
	Role           string
	Content        string
}

// MessageListParams holds parameters for listing messages by conversation
type MessageListParams struct {
	ConversationID uuid.UUID
	Limit          int
}
