package websocket_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chatsphere/internal/conversations"
	"chatsphere/internal/messages"
	"chatsphere/internal/users"
	ws "chatsphere/internal/websocket"
)

// MockUserRepository implements users.UserRepository
type MockUserRepository struct {
	mu       sync.Mutex
	statuses map[int64]bool
	users    map[int64]*users.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		statuses: make(map[int64]bool),
		users:    make(map[int64]*users.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, u *users.User) error {
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*users.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, exists := m.users[id]
	if !exists {
		return &users.User{ID: id}, nil
	}
	return u, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}

func (m *MockUserRepository) UpdateOnlineStatus(ctx context.Context, id int64, isOnline bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statuses[id] = isOnline
	now := time.Now()
	m.users[id] = &users.User{
		ID:         id,
		IsOnline:   isOnline,
		LastSeenAt: &now,
	}
	return nil
}

// MockConversationRepository implements conversations.ConversationRepository
type MockConversationRepository struct {
	mu           sync.Mutex
	participants map[int64][]int64
	userConvs    map[int64][]int64
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{
		participants: make(map[int64][]int64),
		userConvs:    make(map[int64][]int64),
	}
}

func (m *MockConversationRepository) CreateConversation(ctx context.Context) (*conversations.Conversation, error) {
	return nil, nil
}

func (m *MockConversationRepository) GetConversationByID(ctx context.Context, id int64) (*conversations.Conversation, error) {
	return nil, nil
}

func (m *MockConversationRepository) GetUserConversations(ctx context.Context, userID int64, search string) ([]*conversations.Conversation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	convIDs := m.userConvs[userID]
	res := make([]*conversations.Conversation, len(convIDs))
	for i, cID := range convIDs {
		uIDs := m.participants[cID]
		var parts []*conversations.ConversationParticipant
		for _, uID := range uIDs {
			parts = append(parts, &conversations.ConversationParticipant{
				ConversationID: cID,
				UserID:         uID,
			})
		}
		res[i] = &conversations.Conversation{
			ID:           cID,
			Participants: parts,
		}
	}
	return res, nil
}

func (m *MockConversationRepository) UpdateLastReadMessage(ctx context.Context, conversationID int64, userID int64, messageID int64) error {
	return nil
}

func (m *MockConversationRepository) GetUnreadCount(ctx context.Context, conversationID int64, userID int64, lastReadMessageID int64) (int, error) {
	return 0, nil
}

func (m *MockConversationRepository) GetConversationPartners(ctx context.Context, userID int64) ([]int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	partnerSet := make(map[int64]bool)
	for _, uIDs := range m.participants {
		hasUser := false
		for _, uID := range uIDs {
			if uID == userID {
				hasUser = true
				break
			}
		}
		if hasUser {
			for _, uID := range uIDs {
				if uID != userID {
					partnerSet[uID] = true
				}
			}
		}
	}

	partners := make([]int64, 0, len(partnerSet))
	for pID := range partnerSet {
		partners = append(partners, pID)
	}
	return partners, nil
}

func (m *MockConversationRepository) AddParticipant(ctx context.Context, p *conversations.ConversationParticipant) error {
	return nil
}

func (m *MockConversationRepository) RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error {
	return nil
}

func (m *MockConversationRepository) GetParticipants(ctx context.Context, conversationID int64) ([]*conversations.ConversationParticipant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	uIDs := m.participants[conversationID]
	res := make([]*conversations.ConversationParticipant, len(uIDs))
	for i, uID := range uIDs {
		res[i] = &conversations.ConversationParticipant{
			ConversationID: conversationID,
			UserID:         uID,
		}
	}
	return res, nil
}

func (m *MockConversationRepository) UpdateConversationTimestamp(ctx context.Context, id int64) error {
	return nil
}

func (m *MockConversationRepository) GetConversationByParticipants(ctx context.Context, userID1, userID2 int64) (*conversations.Conversation, error) {
	return nil, conversations.ErrConversationNotFound
}

func TestWebSocketFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := NewMockUserRepository()
	convRepo := NewMockConversationRepository()

	// Alice (1) and Bob (2) share conversation 10. Charlie (3) is unrelated.
	convRepo.userConvs[1] = []int64{10}
	convRepo.userConvs[2] = []int64{10}
	convRepo.participants[10] = []int64{1, 2}

	manager := ws.NewManager(userRepo, convRepo)
	handler := ws.NewHandler(manager)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		token := c.Query("token")
		var userID int64
		switch token {
		case "alice":
			userID = 1
		case "bob":
			userID = 2
		case "charlie":
			userID = 3
		}
		if userID > 0 {
			c.Set("user_id", userID)
		}
		c.Next()
	})
	router.GET("/ws", handler.Connect)

	server := httptest.NewServer(router)
	defer server.Close()

	dialWS := func(token string) *websocket.Conn {
		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?token=" + token
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("failed to dial: %v", err)
		}
		return conn
	}

	// 1. Connect Bob
	bobConn := dialWS("bob")
	defer bobConn.Close()

	// 2. Connect Charlie
	charlieConn := dialWS("charlie")
	defer charlieConn.Close()

	// Wait a moment for register
	time.Sleep(50 * time.Millisecond)

	// Verify Bob's online status in repo
	userRepo.mu.Lock()
	bobOnline := userRepo.statuses[2]
	userRepo.mu.Unlock()
	if !bobOnline {
		t.Error("expected Bob to be online")
	}

	// 3. Connect Alice
	aliceConn := dialWS("alice")
	defer aliceConn.Close()

	// Alice and Bob share a conversation, Bob receives presence.online. Charlie does not.
	var event ws.WSEvent
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err := bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive presence online event, got error: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}
	if event.Event != "presence.online" {
		t.Errorf("expected presence.online event, got: %s", event.Event)
	}
	dataMap := event.Data.(map[string]interface{})
	if int64(dataMap["user_id"].(float64)) != 1 {
		t.Errorf("expected presence online from Alice (1), got %v", dataMap["user_id"])
	}

	// Verify Charlie didn't receive anything
	charlieConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = charlieConn.ReadMessage()
	if err == nil {
		t.Error("expected Charlie to NOT receive presence online event, but received one")
	}

	// 4. Test Message Broadcasting
	msg := &messages.Message{
		ID:             100,
		ConversationID: 10,
		SenderID:       1,
		Content:        "Hello Bob!",
		SentAt:         time.Now(),
	}

	manager.BroadcastMessage(messages.MessageEvent{
		ConversationID: 10,
		ParticipantIDs: []int64{1, 2},
		Message:        msg,
	})

	// Alice and Bob should receive it. Charlie should NOT.
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive message, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal message event: %v", err)
	}
	if event.Event != "message.received" {
		t.Errorf("expected message.received event, got: %s", event.Event)
	}

	aliceConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = aliceConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Alice to receive message copy, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal message event: %v", err)
	}
	if event.Event != "message.received" {
		t.Errorf("expected message.received event, got: %s", event.Event)
	}

	charlieConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = charlieConn.ReadMessage()
	if err == nil {
		t.Error("expected Charlie to NOT receive message event, but received one")
	}

	// 5. Test Disconnection
	aliceConn.Close()
	time.Sleep(50 * time.Millisecond)

	// Bob should receive presence.offline from Alice (preceded by automatic typing.stop cleanup)
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive Alice disconnect event, got error: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}
	if event.Event == "typing.stop" {
		_, p, err = bobConn.ReadMessage()
		if err != nil {
			t.Fatalf("expected Bob to receive presence offline event next, got error: %v", err)
		}
		if err := json.Unmarshal(p, &event); err != nil {
			t.Fatalf("failed to unmarshal event: %v", err)
		}
	}
	if event.Event != "presence.offline" {
		t.Errorf("expected presence.offline event, got: %s", event.Event)
	}
	dataMap = event.Data.(map[string]interface{})
	if int64(dataMap["user_id"].(float64)) != 1 {
		t.Errorf("expected presence offline from Alice (1), got %v", dataMap["user_id"])
	}

	// Charlie should NOT receive offline event
	charlieConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = charlieConn.ReadMessage()
	if err == nil {
		t.Error("expected Charlie to NOT receive presence offline event, but received one")
	}
}

func TestWebSocketTypingFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := NewMockUserRepository()
	convRepo := NewMockConversationRepository()

	// Alice (1) and Bob (2) share conversation 10.
	convRepo.userConvs[1] = []int64{10}
	convRepo.userConvs[2] = []int64{10}
	convRepo.participants[10] = []int64{1, 2}

	manager := ws.NewManager(userRepo, convRepo)
	handler := ws.NewHandler(manager)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		token := c.Query("token")
		var userID int64
		switch token {
		case "alice":
			userID = 1
		case "bob":
			userID = 2
		}
		if userID > 0 {
			c.Set("user_id", userID)
		}
		c.Next()
	})
	router.GET("/ws", handler.Connect)

	server := httptest.NewServer(router)
	defer server.Close()

	dialWS := func(token string) *websocket.Conn {
		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?token=" + token
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("failed to dial: %v", err)
		}
		return conn
	}

	bobConn := dialWS("bob")
	defer bobConn.Close()

	aliceConn := dialWS("alice")
	defer aliceConn.Close()

	// Wait a moment for register
	time.Sleep(50 * time.Millisecond)

	// Bob will receive online event from Alice (discard it first)
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err := bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive Alice's online event: %v", err)
	}

	// 1. Send typing.start from Alice
	startPayload := map[string]interface{}{
		"event": "typing.start",
		"data": map[string]interface{}{
			"conversation_id": 10,
		},
	}
	startBytes, _ := json.Marshal(startPayload)
	err = aliceConn.WriteMessage(websocket.TextMessage, startBytes)
	if err != nil {
		t.Fatalf("failed to send typing.start from Alice: %v", err)
	}

	// 2. Bob should receive typing.start event
	var event ws.WSEvent
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err := bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive typing.start, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if event.Event != "typing.start" {
		t.Errorf("expected typing.start event, got: %s", event.Event)
	}
	dataMap := event.Data.(map[string]interface{})
	if int64(dataMap["conversation_id"].(float64)) != 10 {
		t.Errorf("expected conversation_id 10, got %v", dataMap["conversation_id"])
	}
	if int64(dataMap["user_id"].(float64)) != 1 {
		t.Errorf("expected user_id 1 (Alice), got %v", dataMap["user_id"])
	}

	// 3. Send typing.stop from Alice
	stopPayload := map[string]interface{}{
		"event": "typing.stop",
		"data": map[string]interface{}{
			"conversation_id": 10,
		},
	}
	stopBytes, _ := json.Marshal(stopPayload)
	err = aliceConn.WriteMessage(websocket.TextMessage, stopBytes)
	if err != nil {
		t.Fatalf("failed to send typing.stop from Alice: %v", err)
	}

	// 4. Bob should receive typing.stop event
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive typing.stop, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if event.Event != "typing.stop" {
		t.Errorf("expected typing.stop event, got: %s", event.Event)
	}
	dataMap = event.Data.(map[string]interface{})
	if int64(dataMap["conversation_id"].(float64)) != 10 {
		t.Errorf("expected conversation_id 10, got %v", dataMap["conversation_id"])
	}
	if int64(dataMap["user_id"].(float64)) != 1 {
		t.Errorf("expected user_id 1 (Alice), got %v", dataMap["user_id"])
	}

	// 5. Send typing.start from Alice again
	err = aliceConn.WriteMessage(websocket.TextMessage, startBytes)
	if err != nil {
		t.Fatalf("failed to send typing.start again: %v", err)
	}

	// Bob receives typing.start
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive typing.start, got: %v", err)
	}

	// 6. Alice disconnects abruptly
	aliceConn.Close()
	time.Sleep(50 * time.Millisecond)

	// Bob should receive typing.stop first (automatic emission from server on disconnect)
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive typing.stop on Alice disconnect, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if event.Event != "typing.stop" {
		t.Errorf("expected typing.stop event first, got: %s", event.Event)
	}

	// Bob should receive presence.offline next
	bobConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, p, err = bobConn.ReadMessage()
	if err != nil {
		t.Fatalf("expected Bob to receive presence.offline on Alice disconnect, got: %v", err)
	}
	if err := json.Unmarshal(p, &event); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if event.Event != "presence.offline" {
		t.Errorf("expected presence.offline event, got: %s", event.Event)
	}
}
