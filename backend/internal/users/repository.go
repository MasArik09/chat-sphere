package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailConflict is returned when the email address is already in use.
	ErrEmailConflict = errors.New("email already in use")
)

// UserRepository defines the methods to persist and retrieve User domain entities.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateOnlineStatus(ctx context.Context, id int64, isOnline bool) error
}

// PostgresUserRepository implements UserRepository using PostgreSQL database.
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create inserts a new user record.
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, is_online, last_seen_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		user.Name, user.Email, user.PasswordHash, user.IsOnline, user.LastSeenAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		
		// Use driver type assertion for postgres errors
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailConflict
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by their ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, is_online, last_seen_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.IsOnline, &u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetByEmail retrieves a user by their email.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, is_online, last_seen_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.IsOnline, &u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// UpdateOnlineStatus updates user online presence status.
func (r *PostgresUserRepository) UpdateOnlineStatus(ctx context.Context, id int64, isOnline bool) error {
	query := `
		UPDATE users
		SET is_online = $1,
		    last_seen_at = CASE WHEN $1 = false THEN CURRENT_TIMESTAMP ELSE last_seen_at END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	res, err := r.db.ExecContext(ctx, query, isOnline, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
