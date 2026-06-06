package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"chatsphere/internal/conversations"
	"chatsphere/internal/messages"
	"chatsphere/internal/users"
)

// WSEvent defines the standard wrapper structure for WebSocket outgoing events.
type WSEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// PresenceOnlinePayload represents the data for a presence.online event.
type PresenceOnlinePayload struct {
	UserID int64 `json:"user_id"`
}

// PresenceOfflinePayload represents the data for a presence.offline event.
type PresenceOfflinePayload struct {
	UserID     int64     `json:"user_id"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

// Manager coordinates Hub operations, database updates, and event serialization.
type Manager struct {
	hub              *Hub
	userRepo         users.UserRepository
	conversationRepo conversations.ConversationRepository
}

// NewManager creates a new websocket Manager instance.
func NewManager(userRepo users.UserRepository, conversationRepo conversations.ConversationRepository) *Manager {
	return &Manager{
		hub:              NewHub(),
		userRepo:         userRepo,
		conversationRepo: conversationRepo,
	}
}

// GetHub returns the underlying websocket Hub.
func (m *Manager) GetHub() *Hub {
	return m.hub
}

// HandleUserConnect is called when a user establishes their first connection.
// It sets their status to online in the DB, registers the client,
// and broadcasts presence.online to their conversation partners.
func (m *Manager) HandleUserConnect(ctx context.Context, client *Client) {
	userID := client.userID

	// Register with hub and get if this is the first connection atomically
	isFirst := m.hub.Register(client)

	if isFirst {
		// First connection: update DB status and notify partners
		err := m.userRepo.UpdateOnlineStatus(ctx, userID, true)
		if err != nil {
			log.Printf("websocket manager: failed to update online status for user %d: %v", userID, err)
		}

		// Notify conversation partners
		partners, err := m.getConversationPartners(ctx, userID)
		if err != nil {
			log.Printf("websocket manager: failed to get conversation partners for user %d: %v", userID, err)
		} else if len(partners) > 0 {
			event := WSEvent{
				Event: "presence.online",
				Data: PresenceOnlinePayload{
					UserID: userID,
				},
			}
			payload, err := json.Marshal(event)
			if err == nil {
				m.hub.BroadcastToUsers(partners, payload)
			}
		}
	}
}

// HandleUserDisconnect is called when a client disconnects.
// If it was their last active connection, it updates their status to offline in the DB,
// unregisters the client, and broadcasts presence.offline to their conversation partners.
func (m *Manager) HandleUserDisconnect(ctx context.Context, client *Client) {
	userID := client.userID

	// Unregister from hub and get if this was the last connection atomically
	isLast := m.hub.Unregister(client)

	if isLast {
		// Last connection closed: update DB status and notify partners
		err := m.userRepo.UpdateOnlineStatus(ctx, userID, false)
		if err != nil {
			log.Printf("websocket manager: failed to update online status for user %d: %v", userID, err)
		}

		// Retrieve updated user to get accurate last_seen_at
		var lastSeen time.Time
		u, err := m.userRepo.GetByID(ctx, userID)
		if err == nil && u.LastSeenAt != nil {
			lastSeen = *u.LastSeenAt
		} else {
			lastSeen = time.Now()
		}

		// Broadcast typing.stop to partners in all user's conversations
		convs, err := m.conversationRepo.GetUserConversations(ctx, userID, "")
		if err == nil {
			for _, conv := range convs {
				var partners []int64
				for _, p := range conv.Participants {
					if p.UserID != userID {
						partners = append(partners, p.UserID)
					}
				}
				if len(partners) > 0 {
					stopEvent := WSEvent{
						Event: "typing.stop",
						Data: map[string]interface{}{
							"conversation_id": conv.ID,
							"user_id":         userID,
						},
					}
					payload, err := json.Marshal(stopEvent)
					if err == nil {
						m.hub.BroadcastToUsers(partners, payload)
					}
				}
			}
		}

		// Notify conversation partners of offline status
		partners, err := m.getConversationPartners(ctx, userID)
		if err != nil {
			log.Printf("websocket manager: failed to get conversation partners for user %d: %v", userID, err)
		} else if len(partners) > 0 {
			event := WSEvent{
				Event: "presence.offline",
				Data: PresenceOfflinePayload{
					UserID:     userID,
					LastSeenAt: lastSeen,
				},
			}
			payload, err := json.Marshal(event)
			if err == nil {
				m.hub.BroadcastToUsers(partners, payload)
			}
		}
	}
}

// BroadcastMessage implements messages.MessageEventHub interface.
// It receives MessageEvent, wraps it into a WSEvent and broadcasts it to participants.
func (m *Manager) BroadcastMessage(event messages.MessageEvent) {
	wsEvent := WSEvent{
		Event: "message.received",
		Data:  event.Message,
	}

	payload, err := json.Marshal(wsEvent)
	if err != nil {
		log.Printf("websocket manager: failed to marshal message payload: %v", err)
		return
	}

	m.hub.BroadcastToUsers(event.ParticipantIDs, payload)
}

// getConversationPartners returns list of unique user IDs who share a conversation with the target user.
func (m *Manager) getConversationPartners(ctx context.Context, userID int64) ([]int64, error) {
	return m.conversationRepo.GetConversationPartners(ctx, userID)
}

// WSIncomingMessage defines the structure of events sent by client.
type WSIncomingMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// TypingPayload holds the conversation context for typing events.
type TypingPayload struct {
	ConversationID int64 `json:"conversation_id"`
}

// HandleIncomingMessage parses and routes WebSocket events sent by clients.
func (m *Manager) HandleIncomingMessage(ctx context.Context, userID int64, payload []byte) {
	var msg WSIncomingMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return
	}

	switch msg.Event {
	case "typing.start", "typing.stop":
		var p TypingPayload
		if err := json.Unmarshal(msg.Data, &p); err != nil {
			return
		}

		// Verify membership
		participants, err := m.conversationRepo.GetParticipants(ctx, p.ConversationID)
		if err != nil || len(participants) == 0 {
			return
		}

		isMember := false
		var partners []int64
		for _, part := range participants {
			if part.UserID == userID {
				isMember = true
			} else {
				partners = append(partners, part.UserID)
			}
		}

		if !isMember {
			return
		}

		// Broadcast event to partners
		broadcastEvent := WSEvent{
			Event: msg.Event,
			Data: map[string]interface{}{
				"conversation_id": p.ConversationID,
				"user_id":         userID,
			},
		}

		broadcastPayload, err := json.Marshal(broadcastEvent)
		if err == nil {
			m.hub.BroadcastToUsers(partners, broadcastPayload)
		}
	}
}

// BroadcastReadReceipt implements conversations.ConversationEventHub interface.
func (m *Manager) BroadcastReadReceipt(event conversations.ReadEvent) {
	wsEvent := WSEvent{
		Event: "message.read",
		Data: map[string]interface{}{
			"conversation_id":      event.ConversationID,
			"user_id":              event.UserID,
			"last_read_message_id": event.LastReadMessageID,
		},
	}

	payload, err := json.Marshal(wsEvent)
	if err != nil {
		return
	}

	m.hub.BroadcastToUsers(event.ParticipantIDs, payload)
}
