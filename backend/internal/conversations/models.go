package conversations

import "time"

// Conversation represents a private chat between users.
type Conversation struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConversationParticipant represents a member associated with a conversation.
type ConversationParticipant struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	UserID         int64     `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}
