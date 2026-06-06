package messages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"chatsphere/internal/conversations"
)

// SendMessageRequest binding structure.
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MessageHandler holds HTTP controller logic.
type MessageHandler struct {
	service MessageService
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(service MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

// Send handles sending messages to a conversation.
func (h *MessageHandler) Send(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := userIDVal.(int64)

	convIDStr := c.Param("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid conversation ID"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	msg, err := h.service.SendMessage(c.Request.Context(), userID, convID, req.Content)
	if err != nil {
		if errors.Is(err, ErrContentEmpty) || errors.Is(err, ErrContentTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, conversations.ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": msg,
	})
}

// List handles listing messages within a conversation with pagination.
func (h *MessageHandler) List(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := userIDVal.(int64)

	convIDStr := c.Param("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid conversation ID"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	msgs, err := h.service.GetConversationMessages(c.Request.Context(), userID, convID, page, limit)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, conversations.ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to retrieve message history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": msgs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}
