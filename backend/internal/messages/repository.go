package messages

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chatsphere/internal/database"
)

var (
	// ErrMessageNotFound is returned when a message is not found.
	ErrMessageNotFound = errors.New("message not found")
)

// MessageRepository defines the persistence operations for Messages.
type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *Message) error
	GetMessagesByConversation(ctx context.Context, conversationID int64, limit int, offset int) ([]*Message, error)
	GetLatestMessage(ctx context.Context, conversationID int64) (*Message, error)
}

// PostgresMessageRepository implements MessageRepository using PostgreSQL.
type PostgresMessageRepository struct {
	db *sql.DB
}

// NewPostgresMessageRepository creates a new PostgresMessageRepository.
func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) getDB(ctx context.Context) database.DBTX {
	if tx := database.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

// CreateMessage inserts a new message record.
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, msg *Message) error {
	if msg.SentAt.IsZero() {
		msg.SentAt = time.Now()
	}

	query := `
		INSERT INTO messages (conversation_id, sender_id, content, sent_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.getDB(ctx).QueryRowContext(
		ctx, query,
		msg.ConversationID, msg.SenderID, msg.Content, msg.SentAt,
	).Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// GetMessagesByConversation retrieves message history for a conversation.
func (r *PostgresMessageRepository) GetMessagesByConversation(ctx context.Context, conversationID int64, limit int, offset int) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at ASC, id ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.getDB(ctx).QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Message
	for rows.Next() {
		var m Message
		err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.SentAt, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// GetLatestMessage retrieves the most recent message in a conversation.
func (r *PostgresMessageRepository) GetLatestMessage(ctx context.Context, conversationID int64) (*Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, sent_at, created_at, updated_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at DESC, id DESC
		LIMIT 1
	`
	var m Message
	err := r.getDB(ctx).QueryRowContext(ctx, query, conversationID).Scan(
		&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.SentAt, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	return &m, nil
}
