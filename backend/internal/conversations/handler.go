package conversations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateConversationRequest binding request object.
type CreateConversationRequest struct {
	ParticipantIDs []int64 `json:"participant_ids" binding:"required,min=1"`
}

// AddParticipantRequest binding request object.
type AddParticipantRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// ConversationHandler holds Gin controller methods.
type ConversationHandler struct {
	service ConversationService
}

// NewConversationHandler creates a new ConversationHandler.
func NewConversationHandler(service ConversationService) *ConversationHandler {
	return &ConversationHandler{service: service}
}

// Create inserts a conversation.
func (h *ConversationHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := userIDVal.(int64)

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	conv, err := h.service.CreateConversation(c.Request.Context(), userID, req.ParticipantIDs)
	if err != nil {
		if errors.Is(err, ErrInvalidParticipants) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":      true,
		"conversation": conv,
	})
}

// List lists user conversations.
func (h *ConversationHandler) List(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := userIDVal.(int64)

	convs, err := h.service.GetUserConversations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"conversations": convs,
	})
}

// Detail retrieves conversation info and participant lists.
func (h *ConversationHandler) Detail(c *gin.Context) {
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

	detail, err := h.service.GetConversationDetail(c.Request.Context(), userID, convID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get conversation details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"conversation": detail,
	})
}

// AddParticipant appends a user to participant lists.
func (h *ConversationHandler) AddParticipant(c *gin.Context) {
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

	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = h.service.AddParticipant(c.Request.Context(), userID, convID, req.UserID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, ErrParticipantConflict) {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Participant added successfully",
	})
}

// RemoveParticipant deletes a user from participant lists.
func (h *ConversationHandler) RemoveParticipant(c *gin.Context) {
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

	targetUserIDStr := c.Param("userId")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid target user ID"})
		return
	}

	err = h.service.RemoveParticipant(c.Request.Context(), userID, convID, targetUserID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		if errors.Is(err, ErrParticipantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to remove participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Participant removed successfully",
	})
}
