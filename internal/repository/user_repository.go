package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"swipe-sports-backend/internal/database"
	"swipe-sports-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			oauth_id, oauth_provider, name, email, gender, location, 
			latitude, longitude, rank, profile_pic_url, bio, 
			sport_preferences, skill_level, play_style, availability
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		user.OAuthID, user.OAuthProvider, user.Name, user.Email, user.Gender,
		user.Location, user.Latitude, user.Longitude, user.Rank,
		user.ProfilePicURL, user.Bio, user.SportPreferences, user.SkillLevel,
		user.PlayStyle, user.Availability,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	
	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.OAuthID, &user.OAuthProvider, &user.Name, &user.Email,
		&user.Gender, &user.Location, &user.Latitude, &user.Longitude, &user.Rank,
		&user.ProfilePicURL, &user.Bio, &user.SportPreferences, &user.SkillLevel,
		&user.PlayStyle, &user.Availability, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByOAuthID(oauthID, provider string) (*models.User, error) {
	query := `SELECT * FROM users WHERE oauth_id = ? AND oauth_provider = ?`
	
	var user models.User
	err := r.db.QueryRow(query, oauthID, provider).Scan(
		&user.ID, &user.OAuthID, &user.OAuthProvider, &user.Name, &user.Email,
		&user.Gender, &user.Location, &user.Latitude, &user.Longitude, &user.Rank,
		&user.ProfilePicURL, &user.Bio, &user.SportPreferences, &user.SkillLevel,
		&user.PlayStyle, &user.Availability, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by oauth id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET 
			name = ?, email = ?, gender = ?, location = ?, latitude = ?, 
			longitude = ?, rank = ?, profile_pic_url = ?, bio = ?, 
			sport_preferences = ?, skill_level = ?, play_style = ?, 
			availability = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		user.Name, user.Email, user.Gender, user.Location, user.Latitude,
		user.Longitude, user.Rank, user.ProfilePicURL, user.Bio,
		user.SportPreferences, user.SkillLevel, user.PlayStyle,
		user.Availability, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetProfilesForSwipe(userID int64, filter models.ProfileFilter) ([]models.UserProfile, error) {
	var conditions []string
	var args []interface{}

	// Base condition: exclude users already swiped by current user
	conditions = append(conditions, "id NOT IN (SELECT swipee_id FROM swipes WHERE swiper_id = ?)")
	args = append(args, userID)

	// Exclude current user
	conditions = append(conditions, "id != ?")
	args = append(args, userID)

	// Add filter conditions
	if filter.Gender != nil {
		conditions = append(conditions, "gender = ?")
		args = append(args, *filter.Gender)
	}

	if filter.Location != nil {
		conditions = append(conditions, "location = ?")
		args = append(args, *filter.Location)
	}

	if filter.MinRank != nil {
		conditions = append(conditions, "rank >= ?")
		args = append(args, *filter.MinRank)
	}

	if filter.MaxRank != nil {
		conditions = append(conditions, "rank <= ?")
		args = append(args, *filter.MaxRank)
	}

	// Distance-based filtering (simplified - in production use PostGIS)
	if filter.Latitude != nil && filter.Longitude != nil && filter.Radius != nil {
		// This is a simplified distance calculation
		// In production, use PostGIS for proper geospatial queries
		conditions = append(conditions, `
			(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * 
			cos(radians(longitude) - radians(?)) + sin(radians(?)) * 
			sin(radians(latitude)))) <= ?
		`)
		args = append(args, *filter.Latitude, *filter.Longitude, *filter.Latitude, *filter.Radius)
	}

	// Set default limit if not provided
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	query := fmt.Sprintf(`
		SELECT id, name, gender, location, rank, profile_pic_url, bio, 
		       sport_preferences, skill_level, play_style, availability, created_at
		FROM users 
		WHERE %s
		ORDER BY RAND()
		LIMIT ? OFFSET ?
	`, strings.Join(conditions, " AND "))

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles for swipe: %w", err)
	}
	defer rows.Close()

	var profiles []models.UserProfile
	for rows.Next() {
		var profile models.UserProfile
		err := rows.Scan(
			&profile.ID, &profile.Name, &profile.Gender, &profile.Location,
			&profile.Rank, &profile.ProfilePicURL, &profile.Bio,
			&profile.SportPreferences, &profile.SkillLevel, &profile.PlayStyle,
			&profile.Availability, &profile.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
} 