package messages_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"chatsphere/internal/conversations"
	"chatsphere/internal/messages"
)

// MockConversationRepository mock db operations.
type MockConversationRepository struct {
	participants map[int64][]*conversations.ConversationParticipant
	updated      map[int64]time.Time
}

func (m *MockConversationRepository) CreateConversation(ctx context.Context) (*conversations.Conversation, error) {
	return nil, nil
}

func (m *MockConversationRepository) GetConversationByID(ctx context.Context, id int64) (*conversations.Conversation, error) {
	return &conversations.Conversation{ID: id}, nil
}

func (m *MockConversationRepository) GetUserConversations(ctx context.Context, userID int64) ([]*conversations.Conversation, error) {
	return nil, nil
}

func (m *MockConversationRepository) AddParticipant(ctx context.Context, p *conversations.ConversationParticipant) error {
	return nil
}

func (m *MockConversationRepository) RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error {
	return nil
}

func (m *MockConversationRepository) GetParticipants(ctx context.Context, conversationID int64) ([]*conversations.ConversationParticipant, error) {
	list, exists := m.participants[conversationID]
	if !exists {
		return []*conversations.ConversationParticipant{}, nil
	}
	return list, nil
}

func (m *MockConversationRepository) UpdateConversationTimestamp(ctx context.Context, id int64) error {
	m.updated[id] = time.Now()
	return nil
}

// MockMessageRepository mock db operations.
type MockMessageRepository struct {
	messages   map[int64][]*messages.Message
	nextMsgID  int64
	createdMsg *messages.Message
}

func (m *MockMessageRepository) CreateMessage(ctx context.Context, msg *messages.Message) error {
	msg.ID = m.nextMsgID
	m.nextMsgID++
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	m.createdMsg = msg
	m.messages[msg.ConversationID] = append(m.messages[msg.ConversationID], msg)
	return nil
}

func (m *MockMessageRepository) GetMessagesByConversation(ctx context.Context, conversationID int64, limit int, offset int) ([]*messages.Message, error) {
	list := m.messages[conversationID]
	if offset >= len(list) {
		return []*messages.Message{}, nil
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], nil
}

func (m *MockMessageRepository) GetLatestMessage(ctx context.Context, conversationID int64) (*messages.Message, error) {
	return nil, nil
}

func TestMessagesFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	convRepo := &MockConversationRepository{
		participants: map[int64][]*conversations.ConversationParticipant{
			1: {
				{ConversationID: 1, UserID: 1},
				{ConversationID: 1, UserID: 2},
			},
		},
		updated: make(map[int64]time.Time),
	}

	msgRepo := &MockMessageRepository{
		messages:  make(map[int64][]*messages.Message),
		nextMsgID: 1,
	}

	service := messages.NewMessageService(msgRepo, convRepo)
	handler := messages.NewMessageHandler(service)

	router := gin.New()
	var mockUserID int64 = 1
	router.Use(func(c *gin.Context) {
		c.Set("user_id", mockUserID)
		c.Next()
	})

	router.POST("/conversations/:id/messages", handler.Send)
	router.GET("/conversations/:id/messages", handler.List)

	// 1. Send Message Success
	sendReq := messages.SendMessageRequest{
		Content: "Hello Alice",
	}
	body, _ := json.Marshal(sendReq)
	req, _ := http.NewRequest("POST", "/conversations/1/messages", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	if msgRepo.createdMsg.Content != "Hello Alice" {
		t.Fatalf("expected message hello alice, got %s", msgRepo.createdMsg.Content)
	}
	if _, updated := convRepo.updated[1]; !updated {
		t.Fatal("expected conversation updated_at timestamp to be updated")
	}

	// 2. Empty Message Rejection
	sendReq.Content = "   "
	body, _ = json.Marshal(sendReq)
	req, _ = http.NewRequest("POST", "/conversations/1/messages", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	// 3. Message Length > 2000 Rejection
	sendReq.Content = strings.Repeat("a", 2001)
	body, _ = json.Marshal(sendReq)
	req, _ = http.NewRequest("POST", "/conversations/1/messages", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	// 4. Unauthorized Conversation Access (Send)
	mockUserID = 99 // Switch to user not in participant list
	sendReq.Content = "Intruder message"
	body, _ = json.Marshal(sendReq)
	req, _ = http.NewRequest("POST", "/conversations/1/messages", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden, got %d", w.Code)
	}

	// 5. Message History Retrieval
	mockUserID = 1 // Reset back to member
	// Seed multiple messages
	for i := 0; i < 5; i++ {
		msgRepo.messages[1] = append(msgRepo.messages[1], &messages.Message{
			ID:             int64(i + 2),
			ConversationID: 1,
			SenderID:       1,
			Content:        "Seed message",
			SentAt:         time.Now(),
		})
	}
	req, _ = http.NewRequest("GET", "/conversations/1/messages?page=1&limit=20", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var listResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &listResp)
	msgsList := listResp["messages"].([]any)
	// total = 1 sent + 5 seeded = 6 messages
	if len(msgsList) != 6 {
		t.Fatalf("expected 6 messages in history, got %d", len(msgsList))
	}

	// 6. Pagination Behavior
	req, _ = http.NewRequest("GET", "/conversations/1/messages?page=2&limit=2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var pageResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &pageResp)
	pageMsgs := pageResp["messages"].([]any)
	if len(pageMsgs) != 2 {
		t.Fatalf("expected 2 messages on page 2 with limit 2, got %d", len(pageMsgs))
	}
}
