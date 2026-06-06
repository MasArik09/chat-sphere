package tests

import (
	"context"
	"errors"
	"os"
	"testing"

	"chatsphere/internal/conversations"
	"chatsphere/internal/database"
	"chatsphere/internal/messages"
	"chatsphere/internal/users"
	"chatsphere/pkg/config"
)

type FailingConversationRepository struct {
	conversations.ConversationRepository
	failAddParticipant  bool
	failUpdateTimestamp bool
	failUserID          int64
}

func (r *FailingConversationRepository) AddParticipant(ctx context.Context, participant *conversations.ConversationParticipant) error {
	if r.failAddParticipant && participant.UserID == r.failUserID {
		return errors.New("simulated error on participant insert")
	}
	return r.ConversationRepository.AddParticipant(ctx, participant)
}

func (r *FailingConversationRepository) UpdateConversationTimestamp(ctx context.Context, id int64) error {
	if r.failUpdateTimestamp {
		return errors.New("simulated error on conversation timestamp update")
	}
	return r.ConversationRepository.UpdateConversationTimestamp(ctx, id)
}

func TestTransactionRollbackIntegration(t *testing.T) {
	// Read configuration
	cfg := config.Load()
	
	// Force the host to postgres if we are running in docker loopback/bridge
	if os.Getenv("DB_HOST") == "" {
		cfg.DBHost = "postgres"
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Skipf("Skipping integration test; database connection failed: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	// Clear test tables to avoid conflicts
	_, _ = db.Exec("DELETE FROM messages")
	_, _ = db.Exec("DELETE FROM conversation_participants")
	_, _ = db.Exec("DELETE FROM conversations")
	_, _ = db.Exec("DELETE FROM users")

	// 1. Seed two test users
	userRepo := users.NewPostgresUserRepository(db)
	user1 := &users.User{
		Name:         "User One",
		Email:        "user1@example.com",
		PasswordHash: "hash",
	}
	user2 := &users.User{
		Name:         "User Two",
		Email:        "user2@example.com",
		PasswordHash: "hash",
	}
	if err := userRepo.Create(ctx, user1); err != nil {
		t.Fatalf("failed to create user 1: %v", err)
	}
	if err := userRepo.Create(ctx, user2); err != nil {
		t.Fatalf("failed to create user 2: %v", err)
	}

	txManager := database.NewTransactionManager(db)
	realConvRepo := conversations.NewPostgresConversationRepository(db)

	t.Run("Conversation creation rollback when participant insertion fails", func(t *testing.T) {
		failingRepo := &FailingConversationRepository{
			ConversationRepository: realConvRepo,
			failAddParticipant:     true,
			failUserID:             user2.ID,
		}

		convService := conversations.NewConversationService(failingRepo, txManager, nil)

		// Get initial counts
		var initialConvs, initialParts int
		_ = db.QueryRow("SELECT COUNT(*) FROM conversations").Scan(&initialConvs)
		_ = db.QueryRow("SELECT COUNT(*) FROM conversation_participants").Scan(&initialParts)

		// Create conversation should fail
		_, err := convService.CreateConversation(ctx, user1.ID, []int64{user1.ID, user2.ID})
		if err == nil {
			t.Fatal("expected CreateConversation to fail due to simulated error, but it succeeded")
		}

		// Verify rollback: count remains same, no partial writes
		var finalConvs, finalParts int
		_ = db.QueryRow("SELECT COUNT(*) FROM conversations").Scan(&finalConvs)
		_ = db.QueryRow("SELECT COUNT(*) FROM conversation_participants").Scan(&finalParts)

		if finalConvs != initialConvs {
			t.Errorf("expected conversations count to remain %d, got %d (rollback failed)", initialConvs, finalConvs)
		}
		if finalParts != initialParts {
			t.Errorf("expected participants count to remain %d, got %d (rollback failed)", initialParts, finalParts)
		}
	})

	t.Run("Message sending rollback when conversation timestamp update fails", func(t *testing.T) {
		// Create a valid conversation first
		convService := conversations.NewConversationService(realConvRepo, txManager, nil)
		conv, err := convService.CreateConversation(ctx, user1.ID, []int64{user1.ID, user2.ID})
		if err != nil {
			t.Fatalf("failed to create base conversation: %v", err)
		}

		failingRepo := &FailingConversationRepository{
			ConversationRepository: realConvRepo,
			failUpdateTimestamp:    true,
		}

		msgRepo := messages.NewPostgresMessageRepository(db)
		msgService := messages.NewMessageService(msgRepo, failingRepo, nil, txManager)

		// Get initial messages count
		var initialMsgs int
		_ = db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&initialMsgs)

		// Send message should fail due to timestamp update failure
		_, err = msgService.SendMessage(ctx, user1.ID, conv.ID, "this should roll back")
		if err == nil {
			t.Fatal("expected SendMessage to fail due to simulated error, but it succeeded")
		}

		// Verify rollback: message count remains same
		var finalMsgs int
		_ = db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&finalMsgs)

		if finalMsgs != initialMsgs {
			t.Errorf("expected messages count to remain %d, got %d (rollback failed)", initialMsgs, finalMsgs)
		}
	})
}
