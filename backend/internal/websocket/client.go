package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client represents a single client connection.
type Client struct {
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

// NewClient instantiates a new Client wrapper.
func NewClient(userID int64, conn *websocket.Conn) *Client {
	return &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

// readPump reads control frames and messages from the connection.
func (c *Client) readPump(onDisconnect func(*Client), onMessage func(*Client, []byte)) {
	defer func() {
		onDisconnect(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if onMessage != nil {
			onMessage(c, payload)
		}
	}
}

// writePump writes messages to the websocket client.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
