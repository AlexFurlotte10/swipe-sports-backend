package repository

import (
	"database/sql"
	"fmt"

	"swipe-sports-backend/internal/database"
	"swipe-sports-backend/internal/models"
)

type SwipeRepository struct {
	db *sql.DB
}

func NewSwipeRepository() *SwipeRepository {
	return &SwipeRepository{db: database.DB}
}

func (r *SwipeRepository) Create(swipe *models.Swipe) error {
	query := `
		INSERT INTO swipes (swiper_id, swipee_id, direction)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, swipe.SwiperID, swipe.SwipeeID, swipe.Direction)
	if err != nil {
		return fmt.Errorf("failed to create swipe: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	swipe.ID = id
	return nil
}

func (r *SwipeRepository) GetBySwiperAndSwipee(swiperID, swipeeID int64) (*models.Swipe, error) {
	query := `SELECT * FROM swipes WHERE swiper_id = ? AND swipee_id = ?`
	
	var swipe models.Swipe
	err := r.db.QueryRow(query, swiperID, swipeeID).Scan(
		&swipe.ID, &swipe.SwiperID, &swipe.SwipeeID, &swipe.Direction, &swipe.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get swipe: %w", err)
	}

	return &swipe, nil
}

func (r *SwipeRepository) CheckForMatch(user1ID, user2ID int64) (bool, error) {
	// Check if both users have swiped right on each other
	query := `
		SELECT COUNT(*) FROM swipes 
		WHERE ((swiper_id = ? AND swipee_id = ? AND direction = 'right') OR 
		       (swiper_id = ? AND swipee_id = ? AND direction = 'right'))
		AND EXISTS (
			SELECT 1 FROM swipes 
			WHERE swiper_id = ? AND swipee_id = ? AND direction = 'right'
		)
		AND EXISTS (
			SELECT 1 FROM swipes 
			WHERE swiper_id = ? AND swipee_id = ? AND direction = 'right'
		)
	`

	var count int
	err := r.db.QueryRow(query, user1ID, user2ID, user2ID, user1ID, user1ID, user2ID, user2ID, user1ID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for match: %w", err)
	}

	return count >= 2, nil
}

func (r *SwipeRepository) CreateMatch(match *models.Match) error {
	query := `
		INSERT INTO matches (user1_id, user2_id)
		VALUES (?, ?)
	`

	result, err := r.db.Exec(query, match.User1ID, match.User2ID)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	match.ID = id
	return nil
}

func (r *SwipeRepository) GetMatchesByUserID(userID int64) ([]models.Match, error) {
	query := `
		SELECT * FROM matches 
		WHERE user1_id = ? OR user2_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		err := rows.Scan(&match.ID, &match.User1ID, &match.User2ID, &match.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *SwipeRepository) GetMatchByID(matchID int64) (*models.Match, error) {
	query := `SELECT * FROM matches WHERE id = ?`
	
	var match models.Match
	err := r.db.QueryRow(query, matchID).Scan(&match.ID, &match.User1ID, &match.User2ID, &match.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return &match, nil
}

func (r *SwipeRepository) IsUserInMatch(userID, matchID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM matches WHERE id = ? AND (user1_id = ? OR user2_id = ?)`
	
	var count int
	err := r.db.QueryRow(query, matchID, userID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user in match: %w", err)
	}

	return count > 0, nil
}

func (r *SwipeRepository) GetMatchWithUsers(matchID int64) (*models.MatchWithUsers, error) {
	query := `
		SELECT m.id, m.created_at,
		       u1.id, u1.name, u1.gender, u1.location, u1.rank, u1.profile_pic_url, 
		       u1.bio, u1.sport_preferences, u1.skill_level, u1.play_style, 
		       u1.availability, u1.created_at,
		       u2.id, u2.name, u2.gender, u2.location, u2.rank, u2.profile_pic_url, 
		       u2.bio, u2.sport_preferences, u2.skill_level, u2.play_style, 
		       u2.availability, u2.created_at
		FROM matches m
		JOIN users u1 ON m.user1_id = u1.id
		JOIN users u2 ON m.user2_id = u2.id
		WHERE m.id = ?
	`

	var matchWithUsers models.MatchWithUsers
	err := r.db.QueryRow(query, matchID).Scan(
		&matchWithUsers.ID, &matchWithUsers.CreatedAt,
		&matchWithUsers.User1.ID, &matchWithUsers.User1.Name, &matchWithUsers.User1.Gender,
		&matchWithUsers.User1.Location, &matchWithUsers.User1.Rank, &matchWithUsers.User1.ProfilePicURL,
		&matchWithUsers.User1.Bio, &matchWithUsers.User1.SportPreferences, &matchWithUsers.User1.SkillLevel,
		&matchWithUsers.User1.PlayStyle, &matchWithUsers.User1.Availability, &matchWithUsers.User1.CreatedAt,
		&matchWithUsers.User2.ID, &matchWithUsers.User2.Name, &matchWithUsers.User2.Gender,
		&matchWithUsers.User2.Location, &matchWithUsers.User2.Rank, &matchWithUsers.User2.ProfilePicURL,
		&matchWithUsers.User2.Bio, &matchWithUsers.User2.SportPreferences, &matchWithUsers.User2.SkillLevel,
		&matchWithUsers.User2.PlayStyle, &matchWithUsers.User2.Availability, &matchWithUsers.User2.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get match with users: %w", err)
	}

	return &matchWithUsers, nil
} 