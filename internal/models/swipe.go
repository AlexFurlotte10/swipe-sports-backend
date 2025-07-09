package models

import (
	"time"
)

type Swipe struct {
	ID        int64     `json:"id" db:"id"`
	SwiperID  int64     `json:"swiper_id" db:"swiper_id"`
	SwipeeID  int64     `json:"swipee_id" db:"swipee_id"`
	Direction Direction `json:"direction" db:"direction"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Direction string

const (
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

// Swipe request
type SwipeRequest struct {
	SwipeeID  int64     `json:"swipee_id" binding:"required"`
	Direction Direction `json:"direction" binding:"required"`
}

// Swipe response
type SwipeResponse struct {
	IsMatch bool    `json:"is_match"`
	Match   *Match  `json:"match,omitempty"`
} 