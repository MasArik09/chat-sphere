package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"chatsphere/internal/auth"
	"chatsphere/internal/users"
	"chatsphere/pkg/config"
)

// MockUserRepository implements users.UserRepository for testing.
type MockUserRepository struct {
	users  map[int64]*users.User
	emails map[string]*users.User
	nextID int64
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int64]*users.User),
		emails: make(map[string]*users.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(ctx context.Context, u *users.User) error {
	if _, exists := m.emails[u.Email]; exists {
		return users.ErrEmailConflict
	}
	u.ID = m.nextID
	m.nextID++
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	m.users[u.ID] = u
	m.emails[u.Email] = u
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*users.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, users.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	u, exists := m.emails[email]
	if !exists {
		return nil, users.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserRepository) UpdateOnlineStatus(ctx context.Context, id int64, isOnline bool) error {
	u, exists := m.users[id]
	if !exists {
		return users.ErrUserNotFound
	}
	u.IsOnline = isOnline
	return nil
}

func TestAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTSecret:    "test_jwt_secret_key_chatsphere_v1_2026",
		JWTExpiresIn: "1h",
	}

	userRepo := NewMockUserRepository()
	authService := auth.NewAuthService(userRepo, cfg)
	authHandler := auth.NewAuthHandler(authService, userRepo)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/me", auth.AuthMiddleware(cfg), authHandler.Me)

	// 1. Register Success
	regReq := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}
	body, _ := json.Marshal(regReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// 2. Register Duplicate Email
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}

	// 3. Login Success
	loginReq := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "secret123",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var loginResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("failed to extract token from login response")
	}

	// 4. Login Wrong Password
	loginReq.Password = "wrongpassword"
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}

	// 5. Access Protected Route Without Token
	req, _ = http.NewRequest("GET", "/me", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}

	// 6. Access Protected Route With Token
	req, _ = http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var meResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &meResp)
	userMap, ok := meResp["user"].(map[string]any)
	if !ok {
		t.Fatal("failed to get user object from profile endpoint response")
	}
	if userMap["email"] != "john@example.com" {
		t.Fatalf("expected user email john@example.com, got %v", userMap["email"])
	}
}
