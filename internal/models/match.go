package models

import (
	"time"
)

type Match struct {
	ID        int64     `json:"id" db:"id"`
	User1ID   int64     `json:"user1_id" db:"user1_id"`
	User2ID   int64     `json:"user2_id" db:"user2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Match with user details
type MatchWithUsers struct {
	ID        int64       `json:"id"`
	User1     UserProfile `json:"user1"`
	User2     UserProfile `json:"user2"`
	CreatedAt time.Time   `json:"created_at"`
}

// Match response for API
type MatchResponse struct {
	ID        int64       `json:"id"`
	User      UserProfile `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
} 