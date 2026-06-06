package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"chatsphere/pkg/config"
)

// DBTX defines the common query execution methods shared by *sql.DB and *sql.Tx.
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type txKey struct{}

// WithTx injects a database transaction into the context.
func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx extracts the transaction from context if present.
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// TransactionManager handles database transaction execution boundaries.
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type postgresTransactionManager struct {
	db *sql.DB
}

// NewTransactionManager instantiates a TransactionManager wrapper.
func NewTransactionManager(db *sql.DB) TransactionManager {
	return &postgresTransactionManager{db: db}
}

type noopTransactionManager struct{}

func (n *noopTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// NewNoopTransactionManager instantiates a no-op TransactionManager for testing.
func NewNoopTransactionManager() TransactionManager {
	return &noopTransactionManager{}
}

func (tm *postgresTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	txCtx := WithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Connection pool configurations
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Try pinging the database with retries to handle startup delay in Docker Compose
	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Println("Successfully connected to the database")
			return db, nil
		}
		log.Printf("Waiting for database connection... (attempt %d/5): %v", i+1, pingErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", pingErr)
}
