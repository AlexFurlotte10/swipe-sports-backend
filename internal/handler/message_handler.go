package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"swipe-sports-backend/internal/auth"
	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/service"
	"swipe-sports-backend/internal/repository"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: service.NewMessageService(),
	}
}

// GET /messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	matchIDStr := c.Query("match_id")
	if matchIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match_id"})
		return
	}

	// Pagination parameters
	page := 0
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			page = p
		}
	}

	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get messages
	messages, err := h.messageService.GetMessages(matchID, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
		"page":     page,
		"limit":    limit,
		"has_more": len(messages) == limit,
	})
}

// POST /messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate message type
	if req.MessageType == "" {
		req.MessageType = models.MessageTypeText
	}

	if req.MessageType != models.MessageTypeText && 
	   req.MessageType != models.MessageTypeImage && 
	   req.MessageType != models.MessageTypeAudio {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message type"})
		return
	}

	// Validate content length
	if len(req.Content) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message too long. Maximum 1000 characters"})
		return
	}

	// Send message
	message, err := h.messageService.SendMessage(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// DELETE /messages/:id
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	err = h.messageService.DeleteMessage(messageID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// GET /messages/:match_id/latest
func (h *MessageHandler) GetLatestMessage(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	matchIDStr := c.Param("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Verify user is part of the match
	swipeRepo := repository.NewSwipeRepository()
	isInMatch, err := swipeRepo.IsUserInMatch(userID, matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isInMatch {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not part of this match"})
		return
	}

	message, err := h.messageService.GetLatestMessage(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// GET /messages/:match_id/unread-count
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	matchIDStr := c.Param("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Verify user is part of the match
	swipeRepo := repository.NewSwipeRepository()
	isInMatch, err := swipeRepo.IsUserInMatch(userID, matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isInMatch {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not part of this match"})
		return
	}

	count, err := h.messageService.GetUnreadCount(userID, matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// POST /messages/typing
func (h *MessageHandler) SendTypingIndicator(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		MatchID  int64 `json:"match_id" binding:"required"`
		IsTyping bool  `json:"is_typing"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user is part of the match
	swipeRepo := repository.NewSwipeRepository()
	isInMatch, err := swipeRepo.IsUserInMatch(userID, req.MatchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isInMatch {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not part of this match"})
		return
	}

	// Send typing indicator
	h.messageService.HandleTyping(req.MatchID, userID, req.IsTyping)

	c.JSON(http.StatusOK, gin.H{"message": "Typing indicator sent"})
} 