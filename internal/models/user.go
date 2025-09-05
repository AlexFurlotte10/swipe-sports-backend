package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type User struct {
	ID                int64       `json:"id" db:"id"`
	OAuthID           *string     `json:"oauth_id" db:"oauth_id"`
	OAuthProvider     *string     `json:"oauth_provider" db:"oauth_provider"`
	Name              string      `json:"name" db:"name"`
	FirstName         *string     `json:"first_name" db:"first_name"`
	LastName          *string     `json:"last_name" db:"last_name"`
	Age               *int        `json:"age" db:"age"`
	Email             *string     `json:"email" db:"email"`
	Gender            *Gender     `json:"gender" db:"gender"`
	Location          *string     `json:"location" db:"location"`
	Latitude          *float64    `json:"latitude" db:"latitude"`
	Longitude         *float64    `json:"longitude" db:"longitude"`
	Rank              int         `json:"rank" db:"rank"`
	ProfilePicURL     *string     `json:"profile_pic_url" db:"profile_pic_url"`
	Bio               *string     `json:"bio" db:"bio"`
	SportPreferences  SportPreferences `json:"sport_preferences" db:"sport_preferences"`
	SkillLevel        *string     `json:"skill_level" db:"skill_level"`
	NTRPRating        *float64    `json:"ntrp_rating" db:"ntrp_rating"`
	PlayStyle         *string     `json:"play_style" db:"play_style"`
	PreferredTimeslots *string    `json:"preferred_timeslots" db:"preferred_timeslots"`
	Availability      Availability `json:"availability" db:"availability"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type SportPreferences map[string]bool

func (sp SportPreferences) Value() (driver.Value, error) {
	return json.Marshal(sp)
}

func (sp *SportPreferences) Scan(value interface{}) error {
	if value == nil {
		*sp = make(SportPreferences)
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, sp)
}

type Availability map[string][]string

func (a Availability) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		*a = make(Availability)
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, a)
}

// User creation request
type CreateUserRequest struct {
	OAuthID         string            `json:"oauth_id" binding:"required"`
	OAuthProvider   string            `json:"oauth_provider" binding:"required"`
	Name            string            `json:"name" binding:"required"`
	Email           *string           `json:"email"`
	Gender          *Gender           `json:"gender"`
	Location        *string           `json:"location"`
	Latitude        *float64          `json:"latitude"`
	Longitude       *float64          `json:"longitude"`
	ProfilePicURL   *string           `json:"profile_pic_url"`
	Bio             *string           `json:"bio"`
	SportPreferences SportPreferences `json:"sport_preferences"`
	SkillLevel      *string           `json:"skill_level"`
	NTRPRating      *float64          `json:"ntrp_rating"`
	PlayStyle       *string           `json:"play_style"`
	Availability    Availability      `json:"availability"`
}

// User update request
type UpdateUserRequest struct {
	Name               *string           `json:"name"`
	FirstName          *string           `json:"first_name"`
	LastName           *string           `json:"last_name"`
	Age                *int              `json:"age"`
	Gender             *Gender           `json:"gender"`
	Location           *string           `json:"location"`
	Latitude           *float64          `json:"latitude"`
	Longitude          *float64          `json:"longitude"`
	ProfilePicURL      *string           `json:"profile_pic_url"`
	Bio                *string           `json:"bio"`
	SportPreferences   *SportPreferences `json:"sport_preferences"`
	SkillLevel         *string           `json:"skill_level"`
	NTRPRating         *float64          `json:"ntrp_rating"`
	PlayStyle          *string           `json:"play_style"`
	PreferredTimeslots *string           `json:"preferred_timeslots"`
	Availability       *Availability     `json:"availability"`
}

// Profile update request (matches frontend onboarding data structure)
type ProfileUpdateRequest struct {
	Name               string            `json:"name" binding:"required"`
	FirstName          string            `json:"first_name" binding:"required"`
	LastName           string            `json:"last_name" binding:"required"`
	Age                int               `json:"age" binding:"required,min=13,max=100"`
	Location           string            `json:"location" binding:"required"`
	Gender             Gender            `json:"gender" binding:"required"`
	PreferredTimeslots string            `json:"preferred_timeslots" binding:"required"`
	SportPreferences   SportPreferences  `json:"sport_preferences" binding:"required"`
	SkillLevel         string            `json:"skill_level" binding:"required"`
	NTRPRating         float64           `json:"ntrp_rating" binding:"required,min=1.0,max=5.5"`
	PlayStyle          string            `json:"play_style" binding:"required"`
	Bio                string            `json:"bio"`
	Availability       Availability      `json:"availability" binding:"required"`
}

// Profile filtering
type ProfileFilter struct {
	Gender    *Gender `json:"gender"`
	Location  *string `json:"location"`
	MinRank   *int    `json:"min_rank"`
	MaxRank   *int    `json:"max_rank"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Radius    *float64 `json:"radius"` // in kilometers
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
}

// User profile for swiping (excludes sensitive info)
type UserProfile struct {
	ID                int64            `json:"id"`
	Name              string           `json:"name"`
	Age               *int             `json:"age"`
	Gender            *Gender          `json:"gender"`
	Location          *string          `json:"location"`
	Rank              int              `json:"rank"`
	ProfilePicURL     *string          `json:"profile_pic_url"`
	Bio               *string          `json:"bio"`
	SportPreferences  SportPreferences `json:"sport_preferences"`
	SkillLevel        *string          `json:"skill_level"`
	NTRPRating        *float64         `json:"ntrp_rating"`
	PlayStyle         *string          `json:"play_style"`
	PreferredTimeslots *string         `json:"preferred_timeslots"`
	Availability      Availability     `json:"availability"`
	CreatedAt         time.Time        `json:"created_at"`
} 