package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/go-sql-driver/mysql"
	"swipe-sports-backend/internal/models"
)

// Simple config from environment variables
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

func getConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "swipe_sports"),
		JWTSecret:  getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// JWT Claims
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Simple database operations
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			oauth_id, oauth_provider, name, first_name, last_name, age, email, gender, location, 
			latitude, longitude, rank, profile_pic_url, bio, 
			sport_preferences, skill_level, ntrp_rating, play_style, preferred_timeslots, availability
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		user.OAuthID, user.OAuthProvider, user.Name, user.FirstName, user.LastName, user.Age, user.Email, user.Gender,
		user.Location, user.Latitude, user.Longitude, user.Rank,
		user.ProfilePicURL, user.Bio, user.SportPreferences, user.SkillLevel,
		user.NTRPRating, user.PlayStyle, user.PreferredTimeslots, user.Availability,
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

func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	
	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.OAuthID, &user.OAuthProvider, &user.Name, &user.FirstName, &user.LastName, &user.Age, &user.Email,
		&user.Gender, &user.Location, &user.Latitude, &user.Longitude, &user.Rank,
		&user.ProfilePicURL, &user.Bio, &user.SportPreferences, &user.SkillLevel,
		&user.NTRPRating, &user.PlayStyle, &user.PreferredTimeslots, &user.Availability, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByOAuthID(oauthID, provider string) (*models.User, error) {
	query := `SELECT * FROM users WHERE oauth_id = ? AND oauth_provider = ?`
	
	var user models.User
	err := r.db.QueryRow(query, oauthID, provider).Scan(
		&user.ID, &user.OAuthID, &user.OAuthProvider, &user.Name, &user.FirstName, &user.LastName, &user.Age, &user.Email,
		&user.Gender, &user.Location, &user.Latitude, &user.Longitude, &user.Rank,
		&user.ProfilePicURL, &user.Bio, &user.SportPreferences, &user.SkillLevel,
		&user.NTRPRating, &user.PlayStyle, &user.PreferredTimeslots, &user.Availability, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by oauth id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET 
			name = ?, first_name = ?, last_name = ?, age = ?, email = ?, gender = ?, location = ?, latitude = ?, 
			longitude = ?, rank = ?, profile_pic_url = ?, bio = ?, 
			sport_preferences = ?, skill_level = ?, ntrp_rating = ?, play_style = ?, preferred_timeslots = ?,
			availability = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		user.Name, user.FirstName, user.LastName, user.Age, user.Email, user.Gender, user.Location, user.Latitude,
		user.Longitude, user.Rank, user.ProfilePicURL, user.Bio,
		user.SportPreferences, user.SkillLevel, user.NTRPRating, user.PlayStyle, user.PreferredTimeslots,
		user.Availability, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// JWT functions
func generateToken(userID int64, email, jwtSecret string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func validateToken(tokenString, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Auth middleware
func authMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := validateToken(tokenParts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// Response types
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type AuthRequest struct {
	Provider string `json:"provider" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// Simple Auth0 token verification (mock for now)
func verifyAuth0Token(token string) (string, string, string, error) {
	// In a real implementation, you would verify the Auth0 JWT token
	// For now, we'll just extract a mock user ID from the token
	// This is just for testing - you need proper Auth0 verification
	
	// Mock extraction - in reality you'd decode the JWT properly
	userID := "auth0|" + token[0:10] // Simple mock
	email := "user@example.com"
	name := "Test User"
	
	return userID, email, name, nil
}

// Initialize database schema
func initSchema(db *sql.DB) error {
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
			` + "`rank`" + ` INT DEFAULT 1000,
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
			INDEX idx_rank (` + "`rank`" + `),
			INDEX idx_oauth (oauth_id, oauth_provider),
			INDEX idx_age (age),
			INDEX idx_skill_level (skill_level)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

func main() {
	cfg := getConfig()

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	// Initialize repository
	userRepo := NewUserRepository(db)

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://swipesports.co"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Auth routes (no auth required) - match frontend expectations
	auth := r.Group("/auth")
	{
		auth.POST("/signup", func(c *gin.Context) {
			var req AuthRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.Provider != "auth0" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only auth0 provider supported"})
				return
			}

			// Verify Auth0 token
			oauthID, email, name, err := verifyAuth0Token(req.Token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// Check if user exists
			user, err := userRepo.GetUserByOAuthID(oauthID, req.Provider)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			// If user doesn't exist, create new user
			if user == nil {
				user = &models.User{
					OAuthID:         &oauthID,
					OAuthProvider:   &req.Provider,
					Name:            name,
					Email:           &email,
					Rank:            1000,
					SportPreferences: make(models.SportPreferences),
					Availability:     make(models.Availability),
				}

				if err := userRepo.CreateUser(user); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}
			}

			// Generate JWT token
			jwtToken, err := generateToken(user.ID, email, cfg.JWTSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, AuthResponse{
				Token: jwtToken,
				User:  *user,
			})
		})

		auth.POST("/login", func(c *gin.Context) {
			var req AuthRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.Provider != "auth0" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only auth0 provider supported"})
				return
			}

			// Verify Auth0 token
			oauthID, email, name, err := verifyAuth0Token(req.Token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// Get user
			user, err := userRepo.GetUserByOAuthID(oauthID, req.Provider)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			if user == nil {
				// Create user if doesn't exist
				user = &models.User{
					OAuthID:         &oauthID,
					OAuthProvider:   &req.Provider,
					Name:            name,
					Email:           &email,
					Rank:            1000,
					SportPreferences: make(models.SportPreferences),
					Availability:     make(models.Availability),
				}

				if err := userRepo.CreateUser(user); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}
			}

			// Generate JWT token
			jwtToken, err := generateToken(user.ID, email, cfg.JWTSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, AuthResponse{
				Token: jwtToken,
				User:  *user,
			})
		})
	}

	// Protected routes  
	protected := r.Group("")
	protected.Use(authMiddleware(cfg.JWTSecret))
	{
		// Profile routes
		profile := protected.Group("/profile")
		{
			profile.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				
				user, err := userRepo.GetUserByID(userID.(int64))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
					return
				}

				if user == nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				c.JSON(http.StatusOK, user)
			})

			profile.PUT("/update", func(c *gin.Context) {
				userID, _ := c.Get("user_id")

				var req models.ProfileUpdateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				user, err := userRepo.GetUserByID(userID.(int64))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
					return
				}

				if user == nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				// Validate and update user fields
				if req.Gender != models.GenderMale && req.Gender != models.GenderFemale && req.Gender != models.GenderOther {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
					return
				}

				validSkillLevels := []string{"beginner", "intermediate", "advanced"}
				isValidSkill := false
				for _, level := range validSkillLevels {
					if req.SkillLevel == level {
						isValidSkill = true
						break
					}
				}
				if !isValidSkill {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill level"})
					return
				}

				if req.PlayStyle != "ranked" && req.PlayStyle != "fun" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid play style"})
					return
				}

				validTimeslots := []string{"weekends-evenings", "anytime-anywhere", "weekends-only", "weekdays-only"}
				isValidTimeslot := false
				for _, timeslot := range validTimeslots {
					if req.PreferredTimeslots == timeslot {
						isValidTimeslot = true
						break
					}
				}
				if !isValidTimeslot {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preferred timeslots"})
					return
				}

				// Update user fields
				user.Name = req.Name
				user.FirstName = &req.FirstName
				user.LastName = &req.LastName
				user.Age = &req.Age
				user.Gender = &req.Gender
				user.Location = &req.Location
				user.SkillLevel = &req.SkillLevel
				user.NTRPRating = &req.NTRPRating
				user.PlayStyle = &req.PlayStyle
				user.PreferredTimeslots = &req.PreferredTimeslots
				user.Bio = &req.Bio
				user.SportPreferences = req.SportPreferences
				user.Availability = req.Availability

				// Set default coordinates (Halifax)
				if user.Latitude == nil || user.Longitude == nil {
					defaultLat := 44.6488
					defaultLng := -63.5752
					user.Latitude = &defaultLat
					user.Longitude = &defaultLng
				}

				// Save to database
				if err := userRepo.UpdateUser(user); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
					return
				}

				// Generate new token
				var email string
				if user.Email != nil {
					email = *user.Email
				}
				jwtToken, err := generateToken(user.ID, email, cfg.JWTSecret)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				c.JSON(http.StatusOK, AuthResponse{
					Token: jwtToken,
					User:  *user,
				})
			})
		}
	}

	// Start server
	port := ":" + cfg.Port
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Health check: http://localhost%s/health", port)
	log.Printf("Auth signup: http://localhost%s/auth/signup", port)
	log.Printf("Auth login: http://localhost%s/auth/login", port)
	log.Printf("Profile endpoints: http://localhost%s/profile/*", port)
	log.Fatal(r.Run(port))
}
