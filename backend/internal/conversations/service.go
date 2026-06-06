package conversations

import (
	"context"
	"errors"
	"time"

	"chatsphere/internal/database"
)

var (
	// ErrUnauthorized is returned when a user is not a participant of the conversation.
	ErrUnauthorized = errors.New("unauthorized to access this conversation")
	// ErrInvalidParticipants is returned when conversation participant count is invalid.
	ErrInvalidParticipants = errors.New("a conversation must have exactly two participants")
)

type ReadEvent struct {
	ConversationID    int64
	UserID            int64
	LastReadMessageID int64
	ParticipantIDs    []int64
}

type ConversationEventHub interface {
	BroadcastReadReceipt(event ReadEvent)
}

// ConversationResponse represents basic conversation details for list query views.
type ConversationResponse struct {
	ID               int64                      `json:"id"`
	ParticipantCount int                        `json:"participant_count"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	Participants     []*ConversationParticipant `json:"participants"`
	LastMessage      *MessagePreview            `json:"last_message,omitempty"`
	UnreadCount      int                        `json:"unread_count"`
}

// ConversationDetailResponse holds full details including the participant maps.
type ConversationDetailResponse struct {
	Conversation
	Participants []*ConversationParticipant `json:"participants"`
}

// ConversationService orchestrates operations on conversations.
type ConversationService interface {
	CreateConversation(ctx context.Context, creatorID int64, participantIDs []int64) (*ConversationDetailResponse, error)
	GetUserConversations(ctx context.Context, userID int64, search string) ([]*ConversationResponse, error)
	GetConversationDetail(ctx context.Context, userID int64, conversationID int64) (*ConversationDetailResponse, error)
	AddParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error
	RemoveParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error
	MarkAsRead(ctx context.Context, userID int64, conversationID int64, messageID int64) error
}

type conversationService struct {
	repo     ConversationRepository
	tm       database.TransactionManager
	eventHub ConversationEventHub
}

// NewConversationService creates a new ConversationService.
func NewConversationService(repo ConversationRepository, tm database.TransactionManager, eventHub ConversationEventHub) ConversationService {
	return &conversationService{repo: repo, tm: tm, eventHub: eventHub}
}

// CreateConversation inserts a new conversation record, enforcing duplicate checks.
func (s *conversationService) CreateConversation(ctx context.Context, creatorID int64, participantIDs []int64) (*ConversationDetailResponse, error) {
	// Deduplicate and filter out creator ID to determine the other participant
	var otherUserID int64
	for _, pid := range participantIDs {
		id := int64(pid)
		if id != creatorID {
			otherUserID = id
			break
		}
	}

	if otherUserID == 0 {
		return nil, ErrInvalidParticipants
	}

	// 1. Check duplicate conversation using optimized single query
	conv, err := s.repo.GetConversationByParticipants(ctx, creatorID, otherUserID)
	if err == nil {
		participants, err := s.repo.GetParticipants(ctx, conv.ID)
		if err != nil {
			return nil, err
		}
		return &ConversationDetailResponse{
			Conversation: *conv,
			Participants: participants,
		}, nil
	}
	if !errors.Is(err, ErrConversationNotFound) {
		return nil, err
	}

	// 2. Create new conversation and add participants within a transaction
	var c *Conversation
	txErr := s.tm.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		c, err = s.repo.CreateConversation(txCtx)
		if err != nil {
			return err
		}

		err = s.repo.AddParticipant(txCtx, &ConversationParticipant{
			ConversationID: c.ID,
			UserID:         creatorID,
		})
		if err != nil {
			return err
		}

		err = s.repo.AddParticipant(txCtx, &ConversationParticipant{
			ConversationID: c.ID,
			UserID:         otherUserID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	participants, err := s.repo.GetParticipants(ctx, c.ID)
	if err != nil {
		return nil, err
	}

	return &ConversationDetailResponse{
		Conversation: *c,
		Participants: participants,
	}, nil
}

// GetUserConversations fetches a list of conversations for a user, with optional search filter.
func (s *conversationService) GetUserConversations(ctx context.Context, userID int64, search string) ([]*ConversationResponse, error) {
	convs, err := s.repo.GetUserConversations(ctx, userID, search)
	if err != nil {
		return nil, err
	}

	var resp []*ConversationResponse
	for _, c := range convs {
		resp = append(resp, &ConversationResponse{
			ID:               c.ID,
			ParticipantCount: len(c.Participants),
			CreatedAt:        c.CreatedAt,
			UpdatedAt:        c.UpdatedAt,
			Participants:     c.Participants,
			LastMessage:      c.LastMessage,
			UnreadCount:      c.UnreadCount,
		})
	}

	return resp, nil
}

// GetConversationDetail returns details if user is a member.
func (s *conversationService) GetConversationDetail(ctx context.Context, userID int64, conversationID int64) (*ConversationDetailResponse, error) {
	// 1. Verify membership
	participants, err := s.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return nil, err
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

	// 2. Fetch conversation
	c, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return &ConversationDetailResponse{
		Conversation: *c,
		Participants: participants,
	}, nil
}

// AddParticipant registers a user participant to a conversation.
func (s *conversationService) AddParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error {
	// 1. Verify requester belongs to conversation
	participants, err := s.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	isMember := false
	for _, p := range participants {
		if p.UserID == requesterID {
			isMember = true
			break
		}
	}

	if !isMember {
		return ErrUnauthorized
	}

	// 2. Add participant
	return s.repo.AddParticipant(ctx, &ConversationParticipant{
		ConversationID: conversationID,
		UserID:         userID,
	})
}

// RemoveParticipant deletes a user participant from a conversation.
func (s *conversationService) RemoveParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error {
	// 1. Verify requester belongs to conversation
	participants, err := s.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	isMember := false
	for _, p := range participants {
		if p.UserID == requesterID {
			isMember = true
			break
		}
	}

	if !isMember {
		return ErrUnauthorized
	}

	// 2. Remove participant
	return s.repo.RemoveParticipant(ctx, conversationID, userID)
}

// MarkAsRead updates the user's read cursor and broadcasts the receipt.
func (s *conversationService) MarkAsRead(ctx context.Context, userID int64, conversationID int64, messageID int64) error {
	// 1. Verify membership and get partner IDs
	participants, err := s.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	isMember := false
	var partnerIDs []int64
	for _, p := range participants {
		if p.UserID == userID {
			isMember = true
		} else {
			partnerIDs = append(partnerIDs, p.UserID)
		}
	}

	if !isMember {
		return ErrUnauthorized
	}

	// 2. Update database
	err = s.repo.UpdateLastReadMessage(ctx, conversationID, userID, messageID)
	if err != nil {
		return err
	}

	// 3. Broadcast read receipt to partners
	if s.eventHub != nil && len(partnerIDs) > 0 {
		s.eventHub.BroadcastReadReceipt(ReadEvent{
			ConversationID:    conversationID,
			UserID:            userID,
			LastReadMessageID: messageID,
			ParticipantIDs:    partnerIDs,
		})
	}

	return nil
}
