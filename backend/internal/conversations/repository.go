package conversations

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	// ErrConversationNotFound is returned when a conversation is not found.
	ErrConversationNotFound = errors.New("conversation not found")
	// ErrParticipantNotFound is returned when a participant mapping is not found.
	ErrParticipantNotFound = errors.New("participant not found")
	// ErrParticipantConflict is returned when a user is already a participant of a conversation.
	ErrParticipantConflict = errors.New("participant already in conversation")
)

// ConversationRepository defines the persistence operations for Conversations.
type ConversationRepository interface {
	CreateConversation(ctx context.Context) (*Conversation, error)
	GetConversationByID(ctx context.Context, id int64) (*Conversation, error)
	GetUserConversations(ctx context.Context, userID int64) ([]*Conversation, error)
	AddParticipant(ctx context.Context, participant *ConversationParticipant) error
	RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error
	GetParticipants(ctx context.Context, conversationID int64) ([]*ConversationParticipant, error)
	UpdateConversationTimestamp(ctx context.Context, id int64) error
}

// PostgresConversationRepository implements ConversationRepository using PostgreSQL.
type PostgresConversationRepository struct {
	db *sql.DB
}

// NewPostgresConversationRepository creates a new PostgresConversationRepository.
func NewPostgresConversationRepository(db *sql.DB) ConversationRepository {
	return &PostgresConversationRepository{db: db}
}

// CreateConversation inserts a new conversation record.
func (r *PostgresConversationRepository) CreateConversation(ctx context.Context) (*Conversation, error) {
	query := `
		INSERT INTO conversations (created_at, updated_at)
		VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	var c Conversation
	err := r.db.QueryRowContext(ctx, query).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetConversationByID retrieves a conversation by its ID.
func (r *PostgresConversationRepository) GetConversationByID(ctx context.Context, id int64) (*Conversation, error) {
	query := `
		SELECT id, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	var c Conversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return &c, nil
}

// GetUserConversations retrieves all conversations associated with a specific user.
func (r *PostgresConversationRepository) GetUserConversations(ctx context.Context, userID int64) ([]*Conversation, error) {
	query := `
		SELECT c.id, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// AddParticipant adds a user participant to a conversation.
func (r *PostgresConversationRepository) AddParticipant(ctx context.Context, participant *ConversationParticipant) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		participant.ConversationID, participant.UserID,
	).Scan(&participant.ID, &participant.CreatedAt)

	if err != nil {
		// Unique key violation checking for uq_conversation_participants_conversation_user
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrParticipantConflict
		}
		return err
	}
	return nil
}

// RemoveParticipant removes a user participant from a conversation.
func (r *PostgresConversationRepository) RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error {
	query := `
		DELETE FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`
	res, err := r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrParticipantNotFound
	}
	return nil
}

// GetParticipants retrieves all participants for a specific conversation.
func (r *PostgresConversationRepository) GetParticipants(ctx context.Context, conversationID int64) ([]*ConversationParticipant, error) {
	query := `
		SELECT id, conversation_id, user_id, created_at
		FROM conversation_participants
		WHERE conversation_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ConversationParticipant
	for rows.Next() {
		var cp ConversationParticipant
		if err := rows.Scan(&cp.ID, &cp.ConversationID, &cp.UserID, &cp.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &cp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateConversationTimestamp updates the updated_at column to current time.
func (r *PostgresConversationRepository) UpdateConversationTimestamp(ctx context.Context, id int64) error {
	query := `
		UPDATE conversations
		SET updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrConversationNotFound
	}
	return nil
}
