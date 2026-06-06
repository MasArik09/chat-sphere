package tests

import (
	"context"
	"os"
	"testing"

	"chatsphere/internal/conversations"
	"chatsphere/internal/database"
	"chatsphere/internal/messages"
	"chatsphere/internal/users"
	"chatsphere/pkg/config"
)

func TestConversationAggregationIntegration(t *testing.T) {
	cfg := config.Load()
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

	// Clear tables
	_, _ = db.Exec("DELETE FROM messages")
	_, _ = db.Exec("DELETE FROM conversation_participants")
	_, _ = db.Exec("DELETE FROM conversations")
	_, _ = db.Exec("DELETE FROM users")

	// 1. Seed two test users
	userRepo := users.NewPostgresUserRepository(db)
	userAlice := &users.User{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
	}
	userBob := &users.User{
		Name:         "Bob",
		Email:        "bob@example.com",
		PasswordHash: "hash",
	}
	if err := userRepo.Create(ctx, userAlice); err != nil {
		t.Fatalf("failed to create Alice: %v", err)
	}
	if err := userRepo.Create(ctx, userBob); err != nil {
		t.Fatalf("failed to create Bob: %v", err)
	}

	// 2. Create conversation
	convRepo := conversations.NewPostgresConversationRepository(db)
	conv, err := convRepo.CreateConversation(ctx)
	if err != nil {
		t.Fatalf("failed to create conversation: %v", err)
	}

	// 3. Add participants
	partAlice := &conversations.ConversationParticipant{
		ConversationID: conv.ID,
		UserID:         userAlice.ID,
	}
	partBob := &conversations.ConversationParticipant{
		ConversationID: conv.ID,
		UserID:         userBob.ID,
	}
	if err := convRepo.AddParticipant(ctx, partAlice); err != nil {
		t.Fatalf("failed to add Alice to conversation: %v", err)
	}
	if err := convRepo.AddParticipant(ctx, partBob); err != nil {
		t.Fatalf("failed to add Bob to conversation: %v", err)
	}

	// 4. Seed multiple messages (to test lateral query and check duplication)
	msgRepo := messages.NewPostgresMessageRepository(db)
	msg1 := &messages.Message{
		ConversationID: conv.ID,
		SenderID:       userAlice.ID,
		Content:        "Hello Bob",
	}
	msg2 := &messages.Message{
		ConversationID: conv.ID,
		SenderID:       userBob.ID,
		Content:        "Hi Alice",
	}
	msg3 := &messages.Message{
		ConversationID: conv.ID,
		SenderID:       userBob.ID,
		Content:        "Are you there?",
	}
	if err := msgRepo.CreateMessage(ctx, msg1); err != nil {
		t.Fatalf("failed to create message 1: %v", err)
	}
	if err := msgRepo.CreateMessage(ctx, msg2); err != nil {
		t.Fatalf("failed to create message 2: %v", err)
	}
	if err := msgRepo.CreateMessage(ctx, msg3); err != nil {
		t.Fatalf("failed to create message 3: %v", err)
	}

	// 5. Query user conversations for Alice
	userConvs, err := convRepo.GetUserConversations(ctx, userAlice.ID, "")
	if err != nil {
		t.Fatalf("failed to retrieve conversations: %v", err)
	}

	if len(userConvs) != 1 {
		t.Fatalf("expected 1 conversation, got %d", len(userConvs))
	}

	c := userConvs[0]

	// Verify participant count is exactly 2
	if len(c.Participants) != 2 {
		t.Errorf("expected 2 participants, got %d", len(c.Participants))
	}

	// Verify participant aggregation does not contain duplicate entries
	seen := make(map[int64]bool)
	for _, p := range c.Participants {
		if seen[p.UserID] {
			t.Errorf("duplicate participant entry found for user ID %d", p.UserID)
		}
		seen[p.UserID] = true
	}

	// Verify unread count is correct (2 messages from Bob are unread by Alice)
	if c.UnreadCount != 2 {
		t.Errorf("expected 2 unread messages, got %d", c.UnreadCount)
	}

	// Verify latest message preview matches the last message (msg3)
	if c.LastMessage == nil {
		t.Fatal("expected last message preview, got nil")
	}
	if *c.LastMessage.ID != msg3.ID {
		t.Errorf("expected last message ID %d, got %d", msg3.ID, *c.LastMessage.ID)
	}
	if *c.LastMessage.Content != "Are you there?" {
		t.Errorf("expected last message content 'Are you there?', got '%s'", *c.LastMessage.Content)
	}
}
