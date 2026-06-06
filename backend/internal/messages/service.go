package messages

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"chatsphere/internal/conversations"
)

var (
	// ErrUnauthorized is returned when a user is not a member of the conversation.
	ErrUnauthorized = errors.New("unauthorized to access this conversation")
	// ErrContentEmpty is returned when message content is empty after trimming.
	ErrContentEmpty = errors.New("message content cannot be empty")
	// ErrContentTooLong is returned when message content exceeds 2000 characters.
	ErrContentTooLong = errors.New("message content exceeds 2000 characters")
)

// MessageService handles business rules and verification for messages.
type MessageService interface {
	SendMessage(ctx context.Context, senderID int64, conversationID int64, content string) (*Message, error)
	GetConversationMessages(ctx context.Context, userID int64, conversationID int64, page int, limit int) ([]*Message, error)
}

type messageService struct {
	messageRepo      MessageRepository
	conversationRepo conversations.ConversationRepository
}

// NewMessageService creates a new MessageService.
func NewMessageService(messageRepo MessageRepository, conversationRepo conversations.ConversationRepository) MessageService {
	return &messageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
	}
}

// SendMessage validates content, validates membership, stores the message, and updates the conversation timestamp.
func (s *messageService) SendMessage(ctx context.Context, senderID int64, conversationID int64, content string) (*Message, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, ErrContentEmpty
	}
	if utf8.RuneCountInString(trimmed) > 2000 {
		return nil, ErrContentTooLong
	}

	// Verify membership
	participants, err := s.conversationRepo.GetParticipants(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return nil, conversations.ErrConversationNotFound
	}

	isMember := false
	for _, p := range participants {
		if p.UserID == senderID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, ErrUnauthorized
	}

	msg := &Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        trimmed,
	}

	err = s.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	err = s.conversationRepo.UpdateConversationTimestamp(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// GetConversationMessages verifies membership and retrieves message history with pagination.
func (s *messageService) GetConversationMessages(ctx context.Context, userID int64, conversationID int64, page int, limit int) ([]*Message, error) {
	// Verify membership
	participants, err := s.conversationRepo.GetParticipants(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return nil, conversations.ErrConversationNotFound
	}

	isMember := false
	for _, p := range participants {
		if p.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, ErrUnauthorized
	}

	// Pagination defaults and constraints
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	return s.messageRepo.GetMessagesByConversation(ctx, conversationID, limit, offset)
}
