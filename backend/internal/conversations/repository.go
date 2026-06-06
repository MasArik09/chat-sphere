package conversations

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"chatsphere/internal/database"
	"chatsphere/internal/users"

	"github.com/lib/pq"
)

var (
	// ErrConversationNotFound is returned when a conversation is not found.
	ErrConversationNotFound = errors.New("conversation not found")
	// ErrParticipantNotFound is returned when a participant mapping is not found.
	ErrParticipantNotFound = errors.New("participant not found")
	// ErrParticipantConflict is returned when a user is already a participant of a conversation.
	ErrParticipantConflict = errors.New("participant already in conversation")
	// ErrInvalidMessage is returned when a message does not belong to the conversation.
	ErrInvalidMessage = errors.New("message does not belong to the conversation")
)

// ConversationRepository defines the persistence operations for Conversations.
type ConversationRepository interface {
	CreateConversation(ctx context.Context) (*Conversation, error)
	GetConversationByID(ctx context.Context, id int64) (*Conversation, error)
	GetUserConversations(ctx context.Context, userID int64, search string) ([]*Conversation, error)
	AddParticipant(ctx context.Context, participant *ConversationParticipant) error
	RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error
	GetParticipants(ctx context.Context, conversationID int64) ([]*ConversationParticipant, error)
	UpdateConversationTimestamp(ctx context.Context, id int64) error
	GetConversationByParticipants(ctx context.Context, userID1, userID2 int64) (*Conversation, error)
	UpdateLastReadMessage(ctx context.Context, conversationID int64, userID int64, messageID int64) error
	GetUnreadCount(ctx context.Context, conversationID int64, userID int64, lastReadMessageID int64) (int, error)
	GetConversationPartners(ctx context.Context, userID int64) ([]int64, error)
}

// PostgresConversationRepository implements ConversationRepository using PostgreSQL.
type PostgresConversationRepository struct {
	db *sql.DB
}

// NewPostgresConversationRepository creates a new PostgresConversationRepository.
func NewPostgresConversationRepository(db *sql.DB) ConversationRepository {
	return &PostgresConversationRepository{db: db}
}

func (r *PostgresConversationRepository) getDB(ctx context.Context) database.DBTX {
	if tx := database.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

// CreateConversation inserts a new conversation record.
func (r *PostgresConversationRepository) CreateConversation(ctx context.Context) (*Conversation, error) {
	query := `
		INSERT INTO conversations (created_at, updated_at)
		VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	var c Conversation
	err := r.getDB(ctx).QueryRowContext(ctx, query).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
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
	err := r.getDB(ctx).QueryRowContext(ctx, query, id).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return &c, nil
}

// GetUserConversations retrieves all conversations associated with a specific user, matching search.
func (r *PostgresConversationRepository) GetUserConversations(ctx context.Context, userID int64, search string) ([]*Conversation, error) {
	query := `
		WITH user_convs AS (
			SELECT c.id, c.created_at, c.updated_at, cp_self.last_read_message_id
			FROM conversations c
			JOIN conversation_participants cp_self ON c.id = cp_self.conversation_id
	`

	var rows *sql.Rows
	var err error

	if search != "" {
		query += `
				JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id != cp_self.user_id
				JOIN users u2 ON cp2.user_id = u2.id
				WHERE cp_self.user_id = $1 AND u2.name ILIKE $2
			)
		`
	} else {
		query += `
				WHERE cp_self.user_id = $1
			)
		`
	}

	query += `
		SELECT 
			uc.id, uc.created_at, uc.updated_at,
			m.id, m.content, m.sender_id, m.sent_at,
			COALESCE(unr.unread_count, 0) AS unread_count,
			COALESCE(
				json_agg(
					json_build_object(
						'id', cp.id,
						'conversation_id', cp.conversation_id,
						'user_id', cp.user_id,
						'created_at', cp.created_at::timestamptz,
						'last_read_message_id', cp.last_read_message_id,
						'user', json_build_object(
							'id', u.id,
							'name', u.name,
							'email', u.email,
							'is_online', u.is_online,
							'last_seen_at', u.last_seen_at::timestamptz
						)
					)
				) FILTER (WHERE cp.id IS NOT NULL),
				'[]'::json
			) AS participants
		FROM user_convs uc
		LEFT JOIN conversation_participants cp ON uc.id = cp.conversation_id
		LEFT JOIN users u ON cp.user_id = u.id
		LEFT JOIN LATERAL (
			SELECT id, content, sender_id, sent_at
			FROM messages
			WHERE conversation_id = uc.id
			ORDER BY sent_at DESC, id DESC
			LIMIT 1
		) m ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) AS unread_count
			FROM messages
			WHERE conversation_id = uc.id
			  AND sender_id != $1
			  AND (uc.last_read_message_id IS NULL OR id > uc.last_read_message_id)
		) unr ON true
		GROUP BY uc.id, uc.created_at, uc.updated_at, m.id, m.content, m.sender_id, m.sent_at, unr.unread_count
		ORDER BY uc.updated_at DESC
	`

	if search != "" {
		rows, err = r.getDB(ctx).QueryContext(ctx, query, userID, "%"+search+"%")
	} else {
		rows, err = r.getDB(ctx).QueryContext(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Conversation
	for rows.Next() {
		var c Conversation
		var preview MessagePreview
		var participantsJSON []byte
		err := rows.Scan(
			&c.ID, &c.CreatedAt, &c.UpdatedAt,
			&preview.ID, &preview.Content, &preview.SenderID, &preview.SentAt,
			&c.UnreadCount, &participantsJSON,
		)
		if err != nil {
			return nil, err
		}
		if preview.ID != nil {
			c.LastMessage = &preview
		}
		if len(participantsJSON) > 0 {
			if err := json.Unmarshal(participantsJSON, &c.Participants); err != nil {
				return nil, err
			}
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
	err := r.getDB(ctx).QueryRowContext(
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
	res, err := r.getDB(ctx).ExecContext(ctx, query, conversationID, userID)
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

// GetParticipants retrieves all participants for a specific conversation, including user metadata.
func (r *PostgresConversationRepository) GetParticipants(ctx context.Context, conversationID int64) ([]*ConversationParticipant, error) {
	query := `
		SELECT cp.id, cp.conversation_id, cp.user_id, cp.created_at, cp.last_read_message_id,
		       u.id, u.name, u.email, u.is_online, u.last_seen_at, u.created_at, u.updated_at
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.conversation_id = $1
		ORDER BY cp.id ASC
	`
	rows, err := r.getDB(ctx).QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ConversationParticipant
	for rows.Next() {
		var cp ConversationParticipant
		var u users.User
		err := rows.Scan(
			&cp.ID, &cp.ConversationID, &cp.UserID, &cp.CreatedAt, &cp.LastReadMessageID,
			&u.ID, &u.Name, &u.Email, &u.IsOnline, &u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cp.User = &u
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
	res, err := r.getDB(ctx).ExecContext(ctx, query, id)
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

// GetConversationByParticipants finds a private conversation having exactly the two given participants.
func (r *PostgresConversationRepository) GetConversationByParticipants(ctx context.Context, userID1, userID2 int64) (*Conversation, error) {
	query := `
		SELECT c.id, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		GROUP BY c.id, c.created_at, c.updated_at
		HAVING COUNT(cp.user_id) = 2
		   AND SUM(CASE WHEN cp.user_id IN ($1, $2) THEN 1 ELSE 0 END) = 2
		LIMIT 1
	`
	var c Conversation
	err := r.getDB(ctx).QueryRowContext(ctx, query, userID1, userID2).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return &c, nil
}

// UpdateLastReadMessage updates the last_read_message_id column for a participant.
func (r *PostgresConversationRepository) UpdateLastReadMessage(ctx context.Context, conversationID int64, userID int64, messageID int64) error {
	// Verify that the message exists and belongs to this conversation
	var msgConvID int64
	err := r.getDB(ctx).QueryRowContext(ctx, "SELECT conversation_id FROM messages WHERE id = $1", messageID).Scan(&msgConvID)
	if err == sql.ErrNoRows {
		return ErrInvalidMessage
	} else if err != nil {
		return err
	}
	if msgConvID != conversationID {
		return ErrInvalidMessage
	}

	query := `
		UPDATE conversation_participants
		SET last_read_message_id = $3
		WHERE conversation_id = $1 AND user_id = $2
	`
	res, err := r.getDB(ctx).ExecContext(ctx, query, conversationID, userID, messageID)
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

// GetUnreadCount counts messages not sent by the user that are newer than their last_read_message_id.
func (r *PostgresConversationRepository) GetUnreadCount(ctx context.Context, conversationID int64, userID int64, lastReadMessageID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM messages
		WHERE conversation_id = $1 AND sender_id != $2 AND ($3 = 0 OR id > $3)
	`
	var count int
	err := r.getDB(ctx).QueryRowContext(ctx, query, conversationID, userID, lastReadMessageID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetConversationPartners returns all user IDs who share a conversation with the user in a single join query.
func (r *PostgresConversationRepository) GetConversationPartners(ctx context.Context, userID int64) ([]int64, error) {
	query := `
		SELECT DISTINCT cp2.user_id
		FROM conversation_participants cp1
		JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id AND cp2.user_id != cp1.user_id
		WHERE cp1.user_id = $1
	`
	rows, err := r.getDB(ctx).QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partners []int64
	for rows.Next() {
		var pID int64
		if err := rows.Scan(&pID); err != nil {
			return nil, err
		}
		partners = append(partners, pID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return partners, nil
}
