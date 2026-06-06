package conversations

import (
	"time"

	"chatsphere/internal/users"
)

// MessagePreview holds a minimal summary of the last message in a conversation.
type MessagePreview struct {
	ID       *int64     `json:"id,omitempty"`
	Content  *string    `json:"content,omitempty"`
	SenderID *int64     `json:"sender_id,omitempty"`
	SentAt   *time.Time `json:"sent_at,omitempty"`
}

// Conversation represents a private chat between users.
type Conversation struct {
	ID           int64                      `json:"id"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
	LastMessage  *MessagePreview            `json:"last_message,omitempty"`
	Participants []*ConversationParticipant `json:"participants,omitempty"`
	UnreadCount  int                        `json:"unread_count,omitempty"`
}

// ConversationParticipant represents a member associated with a conversation.
type ConversationParticipant struct {
	ID                int64       `json:"id"`
	ConversationID    int64       `json:"conversation_id"`
	UserID            int64       `json:"user_id"`
	CreatedAt         time.Time   `json:"created_at"`
	User              *users.User `json:"user,omitempty"`
	LastReadMessageID *int64      `json:"last_read_message_id,omitempty"`
}
