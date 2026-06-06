package websocket

import (
	"sync"
)

// Hub maintains the set of active clients and handles broadcasting messages.
type Hub struct {
	// Map of userID -> set of Clients (multiple connections per user allowed)
	clients map[int64]map[*Client]bool
	mu      sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]bool),
	}
}

// Register adds a client to the hub and returns true if this is the first connection for the user.
func (h *Hub) Register(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, exists := h.clients[c.userID]
	isFirst := !exists || len(userClients) == 0

	if !exists {
		h.clients[c.userID] = make(map[*Client]bool)
	}
	h.clients[c.userID][c] = true

	return isFirst
}

// Unregister removes a client from the hub and returns true if this was the last connection for the user.
func (h *Hub) Unregister(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, exists := h.clients[c.userID]
	if !exists {
		return false
	}

	delete(userClients, c)
	isLast := len(userClients) == 0
	if isLast {
		delete(h.clients, c.userID)
	}
	return isLast
}

// BroadcastToUsers sends a payload to all active connections of specified user IDs.
func (h *Hub) BroadcastToUsers(userIDs []int64, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range userIDs {
		if userClients, exists := h.clients[userID]; exists {
			for c := range userClients {
				select {
				case c.send <- payload:
				default:
					// Close channel or ignore. We let the client handle it.
				}
			}
		}
	}
}

// IsUserOnline checks if a user is online (has at least one active websocket connection).
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, online := h.clients[userID]
	return online
}
