package main

import (
	"log"
	"os"

	"swipe-sports-backend/internal/config"
	"swipe-sports-backend/internal/database"
	"swipe-sports-backend/internal/server"
	"swipe-sports-backend/internal/redis"
)

func main() {
	// Load environment variables
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.Init()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	defer redisClient.Close()

	// Initialize and start server
	srv := server.New(db, redisClient)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Swipe Sports backend on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 