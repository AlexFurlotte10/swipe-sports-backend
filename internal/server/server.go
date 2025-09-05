package server

import (
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"swipe-sports-backend/internal/auth"
	"swipe-sports-backend/internal/config"
	"swipe-sports-backend/internal/handler"
	"swipe-sports-backend/internal/redis"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
	redis  *redis.Client
}

func New(db *sql.DB, redisClient *redis.Client) *Server {
	// Set Gin mode
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: gin.DefaultWriter})

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
		UTC:    true,
		Format: "method={{.Method}} path={{.Path}} status={{.Status}} duration={{.Duration}} ip={{.ClientIP}} user_agent={{.UserAgent}}",
	}))

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.AppConfig.Server.CORSOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	server := &Server{
		router: router,
		db:     db,
		redis:  redisClient,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			authHandler := handler.NewAuthHandler()
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// Profile routes
			profile := protected.Group("/profile")
			{
				authHandler := handler.NewAuthHandler()
				profile.GET("/me", authHandler.GetMyProfile)
				profile.PUT("/me", authHandler.UpdateMyProfile)
				profile.PUT("/update", authHandler.UpdateProfileFromOnboarding)
				profile.POST("/picture", authHandler.UploadProfilePicture)
			}

			// Swipe routes
			swipe := protected.Group("")
			{
				swipeHandler := handler.NewSwipeHandler()
				swipe.GET("/profiles", swipeHandler.GetProfiles)
				swipe.POST("/swipe", swipeHandler.Swipe)
				swipe.GET("/matches", swipeHandler.GetMatches)
				swipe.GET("/matches/:id", swipeHandler.GetMatch)
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messageHandler := handler.NewMessageHandler()
				messages.GET("", messageHandler.GetMessages)
				messages.POST("", messageHandler.SendMessage)
				messages.DELETE("/:id", messageHandler.DeleteMessage)
				messages.GET("/:match_id/latest", messageHandler.GetLatestMessage)
				messages.GET("/:match_id/unread-count", messageHandler.GetUnreadCount)
				messages.POST("/typing", messageHandler.SendTypingIndicator)
			}
		}

		// WebSocket route (requires authentication in production)
		ws := v1.Group("/ws")
		{
			wsHandler := handler.NewWebSocketHandler()
			ws.GET("/chat", wsHandler.HandleWebSocket)
		}
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	if err := s.db.Ping(); err != nil {
		c.JSON(500, gin.H{
			"status": "unhealthy",
			"error":  "Database connection failed",
		})
		return
	}

	// Check Redis connection
	ctx := s.redis.Context()
	if err := s.redis.Ping(ctx).Err(); err != nil {
		c.JSON(500, gin.H{
			"status": "unhealthy",
			"error":  "Redis connection failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

func (s *Server) Run(addr string) error {
	log.Info().Msgf("Starting server on %s", addr)
	return s.router.Run(addr)
} 