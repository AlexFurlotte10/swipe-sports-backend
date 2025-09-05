package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"swipe-sports-backend/internal/config"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	cfg := config.AppConfig.Database
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db

	// Initialize database schema
	if err := initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			oauth_id VARCHAR(255) UNIQUE,
			oauth_provider VARCHAR(50),
			name VARCHAR(255) NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			age INT,
			email VARCHAR(255) UNIQUE,
			gender ENUM('male', 'female', 'other'),
			location VARCHAR(255),
			latitude DECIMAL(10, 8),
			longitude DECIMAL(11, 8),
			rank INT DEFAULT 1000,
			profile_pic_url VARCHAR(500),
			bio TEXT,
			sport_preferences JSON,
			skill_level VARCHAR(50),
			ntrp_rating DECIMAL(3, 1),
			play_style VARCHAR(100),
			preferred_timeslots VARCHAR(100),
			availability JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_location (location),
			INDEX idx_gender (gender),
			INDEX idx_rank (rank),
			INDEX idx_oauth (oauth_id, oauth_provider),
			INDEX idx_age (age),
			INDEX idx_skill_level (skill_level)
		)`,
		`CREATE TABLE IF NOT EXISTS swipes (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			swiper_id BIGINT NOT NULL,
			swipee_id BIGINT NOT NULL,
			direction ENUM('left', 'right') NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (swiper_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (swipee_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE KEY unique_swipe (swiper_id, swipee_id),
			INDEX idx_swiper (swiper_id),
			INDEX idx_swipee (swipee_id)
		)`,
		`CREATE TABLE IF NOT EXISTS matches (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user1_id BIGINT NOT NULL,
			user2_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE KEY unique_match (user1_id, user2_id),
			INDEX idx_user1 (user1_id),
			INDEX idx_user2 (user2_id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			match_id BIGINT NOT NULL,
			sender_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			message_type ENUM('text', 'image', 'audio') DEFAULT 'text',
			media_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
			FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_match_created (match_id, created_at),
			INDEX idx_sender (sender_id)
		)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	// Run migrations to add new columns if they don't exist
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	migrations := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS first_name VARCHAR(255)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name VARCHAR(255)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS age INT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS ntrp_rating DECIMAL(3, 1)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS preferred_timeslots VARCHAR(100)`,
		`ALTER TABLE users ADD INDEX IF NOT EXISTS idx_age (age)`,
		`ALTER TABLE users ADD INDEX IF NOT EXISTS idx_skill_level (skill_level)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			// Ignore errors for existing columns/indexes, but log them
			log.Printf("Migration warning (likely column/index already exists): %v", err)
		}
	}

	return nil
} 