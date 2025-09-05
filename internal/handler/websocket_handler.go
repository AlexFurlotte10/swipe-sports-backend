package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/service"
	"swipe-sports-backend/internal/repository"
)

type WebSocketHandler struct {
	messageService *service.MessageService
	upgrader       websocket.Upgrader
	clients        map[int64]map[*websocket.Conn]bool // userID -> connections
	mutex          sync.RWMutex
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		messageService: service.NewMessageService(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, check against allowed origins
			},
		},
		clients: make(map[int64]map[*websocket.Conn]bool),
	}
}

// GET /ws/chat
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from query parameter (in production, use JWT token)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Get match ID from query parameter
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

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	h.registerClient(userID, conn)
	defer h.unregisterClient(userID, conn)

	// Mark user as online
	h.messageService.MarkUserOnline(userID)
	defer h.messageService.MarkUserOffline(userID)

	// Send welcome message
	welcomeMsg := models.WSMessage{
		Type: models.WSMessageTypeChat,
		Payload: gin.H{
			"message": "Connected to chat",
			"match_id": matchID,
			"user_id":  userID,
		},
	}
	
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
		return
	}

	// Handle incoming messages
	for {
		var wsMsg models.WSMessage
		if err := conn.ReadJSON(&wsMsg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle different message types
		switch wsMsg.Type {
		case models.WSMessageTypeChat:
			h.handleChatMessage(userID, matchID, wsMsg, conn)
		case models.WSMessageTypeTyping:
			h.handleTypingMessage(userID, matchID, wsMsg)
		default:
			log.Printf("Unknown message type: %s", wsMsg.Type)
		}
	}
}

func (h *WebSocketHandler) handleChatMessage(userID, matchID int64, wsMsg models.WSMessage, conn *websocket.Conn) {
	// Extract message content from payload
	payload, ok := wsMsg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message payload")
		return
	}

	content, ok := payload["content"].(string)
	if !ok || content == "" {
		log.Printf("Invalid message content")
		return
	}

	messageType := models.MessageTypeText
	if msgType, ok := payload["message_type"].(string); ok {
		messageType = models.MessageType(msgType)
	}

	// Create message request
	messageReq := models.CreateMessageRequest{
		MatchID:     matchID,
		Content:     content,
		MessageType: messageType,
	}

	// Send message
	message, err := h.messageService.SendMessage(userID, messageReq)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	// Broadcast to all clients in the match
	h.broadcastToMatch(matchID, models.WSMessage{
		Type: models.WSMessageTypeChat,
		Payload: models.WSChatMessage{
			MatchID:     message.MatchID,
			SenderID:    message.SenderID,
			Content:     message.Content,
			MessageType: message.MessageType,
			MediaURL:    message.MediaURL,
			Timestamp:   message.CreatedAt,
		},
	})
}

func (h *WebSocketHandler) handleTypingMessage(userID, matchID int64, wsMsg models.WSMessage) {
	// Extract typing info from payload
	payload, ok := wsMsg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid typing payload")
		return
	}

	isTyping, ok := payload["is_typing"].(bool)
	if !ok {
		log.Printf("Invalid typing status")
		return
	}

	// Broadcast typing indicator to other users in the match
	h.broadcastToMatch(matchID, models.WSMessage{
		Type: models.WSMessageTypeTyping,
		Payload: models.WSTypingMessage{
			MatchID:  matchID,
			UserID:   userID,
			IsTyping: isTyping,
		},
	})
}

func (h *WebSocketHandler) registerClient(userID int64, conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*websocket.Conn]bool)
	}
	h.clients[userID][conn] = true

	log.Printf("Client registered for user %d. Total connections: %d", userID, len(h.clients[userID]))
}

func (h *WebSocketHandler) unregisterClient(userID int64, conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.clients[userID] != nil {
		delete(h.clients[userID], conn)
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}

	log.Printf("Client unregistered for user %d", userID)
}

func (h *WebSocketHandler) broadcastToMatch(matchID int64, message models.WSMessage) {
	// Get all users in the match
	swipeRepo := repository.NewSwipeRepository()
	match, err := swipeRepo.GetMatchByID(matchID)
	if err != nil {
		log.Printf("Failed to get match: %v", err)
		return
	}

	if match == nil {
		log.Printf("Match not found: %d", matchID)
		return
	}

	// Broadcast to both users in the match
	userIDs := []int64{match.User1ID, match.User2ID}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, userID := range userIDs {
		if connections, exists := h.clients[userID]; exists {
			for conn := range connections {
				// Set write deadline
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("Failed to send message to user %d: %v", userID, err)
					// Remove the connection
					go h.unregisterClient(userID, conn)
				}
			}
		}
	}
}

func (h *WebSocketHandler) broadcastToUser(userID int64, message models.WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if connections, exists := h.clients[userID]; exists {
		for conn := range connections {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Failed to send message to user %d: %v", userID, err)
				// Remove the connection
				go h.unregisterClient(userID, conn)
			}
		}
	}
}

// Broadcast match notification
func (h *WebSocketHandler) BroadcastMatch(matchID int64, user1ID, user2ID int64) {
	matchMsg := models.WSMessage{
		Type: models.WSMessageTypeMatch,
		Payload: gin.H{
			"match_id": matchID,
			"user1_id": user1ID,
			"user2_id": user2ID,
			"message":  "It's a match!",
		},
	}

	// Send to both users
	h.broadcastToUser(user1ID, matchMsg)
	h.broadcastToUser(user2ID, matchMsg)
} 