package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections in development/production (controlled by auth token)
	},
}

// Handler handles websocket upgrade requests.
type Handler struct {
	manager *Manager
}

// NewHandler creates a new websocket Handler instance.
func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// Connect upgrades HTTP to WebSocket and registers the new client.
func (h *Handler) Connect(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "invalid user id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket handler: failed to upgrade connection for user %d: %v", userID, err)
		return
	}

	client := NewClient(userID, conn)

	// Register client and notify partners of online status
	h.manager.HandleUserConnect(context.Background(), client)

	// Start write pump in a separate goroutine
	go client.writePump()

	// Run read pump in connection handler goroutine (blocks until disconnect)
	client.readPump(func(disconnectedClient *Client) {
		h.manager.HandleUserDisconnect(context.Background(), disconnectedClient)
	}, func(client *Client, payload []byte) {
		h.manager.HandleIncomingMessage(context.Background(), client.userID, payload)
	})
}
