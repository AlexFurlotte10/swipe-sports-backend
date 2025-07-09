package repository

import (
	"database/sql"
	"fmt"

	"swipe-sports-backend/internal/database"
	"swipe-sports-backend/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{db: database.DB}
}

func (r *MessageRepository) Create(message *models.Message) error {
	query := `
		INSERT INTO messages (match_id, sender_id, content, message_type, media_url)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, message.MatchID, message.SenderID, message.Content, message.MessageType, message.MediaURL)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	message.ID = id
	return nil
}

func (r *MessageRepository) GetByMatchID(matchID int64, limit, offset int) ([]models.MessageWithSender, error) {
	query := `
		SELECT m.id, m.match_id, m.sender_id, u.name as sender_name, 
		       m.content, m.message_type, m.media_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.match_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, matchID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageWithSender
	for rows.Next() {
		var message models.MessageWithSender
		err := rows.Scan(
			&message.ID, &message.MatchID, &message.SenderID, &message.SenderName,
			&message.Content, &message.MessageType, &message.MediaURL, &message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(messageID int64) (*models.Message, error) {
	query := `SELECT * FROM messages WHERE id = ?`
	
	var message models.Message
	err := r.db.QueryRow(query, messageID).Scan(
		&message.ID, &message.MatchID, &message.SenderID, &message.Content,
		&message.MessageType, &message.MediaURL, &message.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

func (r *MessageRepository) GetLatestMessageByMatchID(matchID int64) (*models.MessageWithSender, error) {
	query := `
		SELECT m.id, m.match_id, m.sender_id, u.name as sender_name, 
		       m.content, m.message_type, m.media_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.match_id = ?
		ORDER BY m.created_at DESC
		LIMIT 1
	`

	var message models.MessageWithSender
	err := r.db.QueryRow(query, matchID).Scan(
		&message.ID, &message.MatchID, &message.SenderID, &message.SenderName,
		&message.Content, &message.MessageType, &message.MediaURL, &message.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest message: %w", err)
	}

	return &message, nil
}

func (r *MessageRepository) GetUnreadCount(userID, matchID int64) (int, error) {
	// This is a simplified implementation
	// In a real app, you'd track read status per user
	query := `
		SELECT COUNT(*) FROM messages 
		WHERE match_id = ? AND sender_id != ?
	`

	var count int
	err := r.db.QueryRow(query, matchID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (r *MessageRepository) DeleteByID(messageID int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	
	_, err := r.db.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetMessagesByUserID(userID int64, limit, offset int) ([]models.MessageWithSender, error) {
	query := `
		SELECT m.id, m.match_id, m.sender_id, u.name as sender_name, 
		       m.content, m.message_type, m.media_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		JOIN matches mt ON m.match_id = mt.id
		WHERE mt.user1_id = ? OR mt.user2_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by user: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageWithSender
	for rows.Next() {
		var message models.MessageWithSender
		err := rows.Scan(
			&message.ID, &message.MatchID, &message.SenderID, &message.SenderName,
			&message.Content, &message.MessageType, &message.MediaURL, &message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
} 