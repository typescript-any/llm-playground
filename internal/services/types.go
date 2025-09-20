package service

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

// ConversationNewParams holds parameters for creating a new conversation with AI-generated title
type ConversationNewParams struct {
	UserID  uuid.UUID
	Content string
	Model   string
}

// MessageSendParams holds parameters for sending a message
type MessageSendParams struct {
	ConversationID uuid.UUID
	Content        string
	Model          string
}

// MessageStreamParams holds parameters for streaming a message
type MessageStreamParams struct {
	ConversationID uuid.UUID
	Content        string
	Model          string
}

// MessageSaveParams holds parameters for saving an assistant message
type MessageSaveParams struct {
	ConversationID uuid.UUID
	Content        string
}
