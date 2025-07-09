package models

import (
	"time"
)

type Message struct {
	ID          int64       `json:"id" db:"id"`
	MatchID     int64       `json:"match_id" db:"match_id"`
	SenderID    int64       `json:"sender_id" db:"sender_id"`
	Content     string      `json:"content" db:"content"`
	MessageType MessageType `json:"message_type" db:"message_type"`
	MediaURL    *string     `json:"media_url" db:"media_url"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeAudio MessageType = "audio"
)

// Message creation request
type CreateMessageRequest struct {
	MatchID     int64       `json:"match_id" binding:"required"`
	Content     string      `json:"content" binding:"required"`
	MessageType MessageType `json:"message_type"`
	MediaURL    *string     `json:"media_url"`
}

// Message with sender info
type MessageWithSender struct {
	ID          int64       `json:"id"`
	MatchID     int64       `json:"match_id"`
	SenderID    int64       `json:"sender_id"`
	SenderName  string      `json:"sender_name"`
	Content     string      `json:"content"`
	MessageType MessageType `json:"message_type"`
	MediaURL    *string     `json:"media_url"`
	CreatedAt   time.Time   `json:"created_at"`
}

// WebSocket message types
type WSMessageType string

const (
	WSMessageTypeChat     WSMessageType = "chat"
	WSMessageTypeTyping   WSMessageType = "typing"
	WSMessageTypeOnline   WSMessageType = "online"
	WSMessageTypeOffline  WSMessageType = "offline"
	WSMessageTypeMatch    WSMessageType = "match"
)

// WebSocket message
type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

// Chat message for WebSocket
type WSChatMessage struct {
	MatchID     int64       `json:"match_id"`
	SenderID    int64       `json:"sender_id"`
	Content     string      `json:"content"`
	MessageType MessageType `json:"message_type"`
	MediaURL    *string     `json:"media_url"`
	Timestamp   time.Time   `json:"timestamp"`
}

// Typing indicator
type WSTypingMessage struct {
	MatchID  int64 `json:"match_id"`
	UserID   int64 `json:"user_id"`
	IsTyping bool  `json:"is_typing"`
} 