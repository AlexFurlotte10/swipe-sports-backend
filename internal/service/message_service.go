package service

import (
	"encoding/json"
	"fmt"

	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/redis"
	"swipe-sports-backend/internal/repository"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	swipeRepo   *repository.SwipeRepository
}

func NewMessageService() *MessageService {
	return &MessageService{
		messageRepo: repository.NewMessageRepository(),
		swipeRepo:   repository.NewSwipeRepository(),
	}
}

func (s *MessageService) SendMessage(senderID int64, messageReq models.CreateMessageRequest) (*models.MessageWithSender, error) {
	// Verify user is part of the match
	isInMatch, err := s.swipeRepo.IsUserInMatch(senderID, messageReq.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user in match: %w", err)
	}

	if !isInMatch {
		return nil, fmt.Errorf("user not part of this match")
	}

	// Create message
	message := &models.Message{
		MatchID:     messageReq.MatchID,
		SenderID:    senderID,
		Content:     messageReq.Content,
		MessageType: messageReq.MessageType,
		MediaURL:    messageReq.MediaURL,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Get message with sender info
	messageWithSender, err := s.messageRepo.GetByID(message.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Clear cache for this match
	redis.DeleteMatchMessages(messageReq.MatchID)

	// Broadcast to WebSocket clients
	s.broadcastMessage(messageReq.MatchID, messageWithSender)

	return &models.MessageWithSender{
		ID:          messageWithSender.ID,
		MatchID:     messageWithSender.MatchID,
		SenderID:    messageWithSender.SenderID,
		Content:     messageWithSender.Content,
		MessageType: messageWithSender.MessageType,
		MediaURL:    messageWithSender.MediaURL,
		CreatedAt:   messageWithSender.CreatedAt,
	}, nil
}

func (s *MessageService) GetMessages(matchID, userID int64, page, limit int) ([]models.MessageWithSender, error) {
	// Verify user is part of the match
	isInMatch, err := s.swipeRepo.IsUserInMatch(userID, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user in match: %w", err)
	}

	if !isInMatch {
		return nil, fmt.Errorf("user not part of this match")
	}

	// Try to get from cache first
	cachedData, err := redis.GetMatchMessages(matchID)
	if err == nil {
		var messages []models.MessageWithSender
		if err := json.Unmarshal(cachedData, &messages); err == nil {
			// Apply pagination to cached data
			start := page * limit
			end := start + limit
			if start >= len(messages) {
				return []models.MessageWithSender{}, nil
			}
			if end > len(messages) {
				end = len(messages)
			}
			return messages[start:end], nil
		}
	}

	// Get from database
	offset := page * limit
	messages, err := s.messageRepo.GetByMatchID(matchID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Cache all messages for this match (without pagination)
	if page == 0 { // Only cache on first page request
		if data, err := json.Marshal(messages); err == nil {
			redis.SetMatchMessages(matchID, data)
		}
	}

	return messages, nil
}

func (s *MessageService) GetLatestMessage(matchID int64) (*models.MessageWithSender, error) {
	return s.messageRepo.GetLatestMessageByMatchID(matchID)
}

func (s *MessageService) DeleteMessage(messageID, userID int64) error {
	// Get message to verify ownership
	message, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	if message == nil {
		return fmt.Errorf("message not found")
	}

	// Only sender can delete their message
	if message.SenderID != userID {
		return fmt.Errorf("unauthorized to delete this message")
	}

	if err := s.messageRepo.DeleteByID(messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	// Clear cache for this match
	redis.DeleteMatchMessages(message.MatchID)

	return nil
}

func (s *MessageService) GetUnreadCount(userID, matchID int64) (int, error) {
	return s.messageRepo.GetUnreadCount(userID, matchID)
}

// WebSocket related methods
func (s *MessageService) broadcastMessage(matchID int64, message *models.Message) {
	// This would be implemented with a WebSocket manager
	// For now, we'll just log the broadcast
	fmt.Printf("Broadcasting message %d to match %d\n", message.ID, matchID)
}

func (s *MessageService) HandleTyping(matchID, userID int64, isTyping bool) {
	// Broadcast typing indicator to other users in the match
	typingMsg := models.WSTypingMessage{
		MatchID:  matchID,
		UserID:   userID,
		IsTyping: isTyping,
	}

	wsMessage := models.WSMessage{
		Type:    models.WSMessageTypeTyping,
		Payload: typingMsg,
	}

	// This would be sent via WebSocket
	fmt.Printf("Typing indicator: %+v\n", wsMessage)
}

func (s *MessageService) MarkUserOnline(userID int64) error {
	return redis.AddOnlineUser(userID)
}

func (s *MessageService) MarkUserOffline(userID int64) error {
	return redis.RemoveOnlineUser(userID)
}

func (s *MessageService) IsUserOnline(userID int64) (bool, error) {
	return redis.IsUserOnline(userID)
}

func (s *MessageService) GetOnlineUsers() ([]string, error) {
	return redis.GetOnlineUsers()
} 