package conversations

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrUnauthorized is returned when a user is not a participant of the conversation.
	ErrUnauthorized = errors.New("unauthorized to access this conversation")
	// ErrInvalidParticipants is returned when conversation participant count is invalid.
	ErrInvalidParticipants = errors.New("a conversation must have exactly two participants")
)

// ConversationResponse represents basic conversation details for list query views.
type ConversationResponse struct {
	ID               int64     `json:"id"`
	ParticipantCount int       `json:"participant_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ConversationDetailResponse holds full details including the participant maps.
type ConversationDetailResponse struct {
	Conversation
	Participants []*ConversationParticipant `json:"participants"`
}

// ConversationService orchestrates operations on conversations.
type ConversationService interface {
	CreateConversation(ctx context.Context, creatorID int64, participantIDs []int64) (*Conversation, error)
	GetUserConversations(ctx context.Context, userID int64) ([]*ConversationResponse, error)
	GetConversationDetail(ctx context.Context, userID int64, conversationID int64) (*ConversationDetailResponse, error)
	AddParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error
	RemoveParticipant(ctx context.Context, requesterID int64, conversationID int64, userID int64) error
}

type conversationService struct {
	repo ConversationRepository
}

// NewConversationService creates a new ConversationService.
func NewConversationService(repo ConversationRepository) ConversationService {
	return &conversationService{repo: repo}
}

// CreateConversation inserts a new conversation record, enforcing duplicate checks.
func (s *conversationService) CreateConversation(ctx context.Context, creatorID int64, participantIDs []int64) (*Conversation, error) {
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

	// 1. Check duplicate conversation: Get all conversations for creatorID
	convs, err := s.repo.GetUserConversations(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	for _, c := range convs {
		participants, err := s.repo.GetParticipants(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		// Look for other participant in this conversation
		isOtherIn := false
		for _, p := range participants {
			if p.UserID == otherUserID {
				isOtherIn = true
				break
			}
		}

		if isOtherIn && len(participants) == 2 {
			// Duplicate private conversation found, return the existing conversation
			return c, nil
		}
	}

	// 2. Create new conversation
	c, err := s.repo.CreateConversation(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Add creator
	err = s.repo.AddParticipant(ctx, &ConversationParticipant{
		ConversationID: c.ID,
		UserID:         creatorID,
	})
	if err != nil {
		return nil, err
	}

	// 4. Add other participant
	err = s.repo.AddParticipant(ctx, &ConversationParticipant{
		ConversationID: c.ID,
		UserID:         otherUserID,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetUserConversations fetches a list of conversations for a user.
func (s *conversationService) GetUserConversations(ctx context.Context, userID int64) ([]*ConversationResponse, error) {
	convs, err := s.repo.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, err
	}

	var resp []*ConversationResponse
	for _, c := range convs {
		participants, err := s.repo.GetParticipants(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &ConversationResponse{
			ID:               c.ID,
			ParticipantCount: len(participants),
			CreatedAt:        c.CreatedAt,
			UpdatedAt:        c.UpdatedAt,
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
