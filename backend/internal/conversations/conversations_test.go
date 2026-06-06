package conversations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"chatsphere/internal/conversations"
)

// MockConversationRepository mock db operations.
type MockConversationRepository struct {
	conversations map[int64]*conversations.Conversation
	participants  map[int64][]*conversations.ConversationParticipant
	nextConvID    int64
	nextPartID    int64
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{
		conversations: make(map[int64]*conversations.Conversation),
		participants:  make(map[int64][]*conversations.ConversationParticipant),
		nextConvID:    1,
		nextPartID:    1,
	}
}

func (m *MockConversationRepository) CreateConversation(ctx context.Context) (*conversations.Conversation, error) {
	c := &conversations.Conversation{
		ID:        m.nextConvID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.conversations[m.nextConvID] = c
	m.nextConvID++
	return c, nil
}

func (m *MockConversationRepository) GetConversationByID(ctx context.Context, id int64) (*conversations.Conversation, error) {
	c, exists := m.conversations[id]
	if !exists {
		return nil, conversations.ErrConversationNotFound
	}
	return c, nil
}

func (m *MockConversationRepository) GetUserConversations(ctx context.Context, userID int64) ([]*conversations.Conversation, error) {
	var list []*conversations.Conversation
	for cid, pList := range m.participants {
		for _, p := range pList {
			if p.UserID == userID {
				list = append(list, m.conversations[cid])
				break
			}
		}
	}
	return list, nil
}

func (m *MockConversationRepository) AddParticipant(ctx context.Context, p *conversations.ConversationParticipant) error {
	pList := m.participants[p.ConversationID]
	for _, exist := range pList {
		if exist.UserID == p.UserID {
			return conversations.ErrParticipantConflict
		}
	}

	p.ID = m.nextPartID
	m.nextPartID++
	p.CreatedAt = time.Now()
	m.participants[p.ConversationID] = append(pList, p)
	return nil
}

func (m *MockConversationRepository) RemoveParticipant(ctx context.Context, conversationID int64, userID int64) error {
	pList, exists := m.participants[conversationID]
	if !exists {
		return conversations.ErrParticipantNotFound
	}

	index := -1
	for i, p := range pList {
		if p.UserID == userID {
			index = i
			break
		}
	}

	if index == -1 {
		return conversations.ErrParticipantNotFound
	}

	m.participants[conversationID] = append(pList[:index], pList[index+1:]...)
	return nil
}

func (m *MockConversationRepository) GetParticipants(ctx context.Context, conversationID int64) ([]*conversations.ConversationParticipant, error) {
	pList, exists := m.participants[conversationID]
	if !exists {
		return []*conversations.ConversationParticipant{}, nil
	}
	return pList, nil
}

func TestConversationsFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockConversationRepository()
	service := conversations.NewConversationService(repo)
	handler := conversations.NewConversationHandler(service)

	router := gin.New()
	// Middleware mock helper to inject user contexts
	var mockUserID int64 = 1
	router.Use(func(c *gin.Context) {
		c.Set("user_id", mockUserID)
		c.Next()
	})

	router.POST("/conversations", handler.Create)
	router.GET("/conversations", handler.List)
	router.GET("/conversations/:id", handler.Detail)
	router.POST("/conversations/:id/participants", handler.AddParticipant)
	router.DELETE("/conversations/:id/participants/:userId", handler.RemoveParticipant)

	// 1. Create Conversation
	createReq := conversations.CreateConversationRequest{
		ParticipantIDs: []int64{2},
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/conversations", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var createResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &createResp)
	convMap := createResp["conversation"].(map[string]any)
	convID := int64(convMap["id"].(float64))

	// 2. Duplicate Check (returns same conversation)
	req, _ = http.NewRequest("POST", "/conversations", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	var dupResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &dupResp)
	dupConvMap := dupResp["conversation"].(map[string]any)
	if int64(dupConvMap["id"].(float64)) != convID {
		t.Fatal("expected duplicate check to return same conversation id")
	}

	// 3. List conversations
	req, _ = http.NewRequest("GET", "/conversations", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// 4. Detail success
	req, _ = http.NewRequest("GET", "/conversations/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// 5. Add participant
	addReq := conversations.AddParticipantRequest{
		UserID: 3,
	}
	addBody, _ := json.Marshal(addReq)
	req, _ = http.NewRequest("POST", "/conversations/1/participants", bytes.NewBuffer(addBody))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// 6. Remove participant
	req, _ = http.NewRequest("DELETE", "/conversations/1/participants/3", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// 7. Membership check / Unauthorized detail lookup
	mockUserID = 99 // Switch user
	req, _ = http.NewRequest("GET", "/conversations/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden for non-member, got %d", w.Code)
	}
}
