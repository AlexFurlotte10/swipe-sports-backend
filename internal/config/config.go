package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	AWS      AWSConfig
	Server   ServerConfig
	RateLimit RateLimitConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type OAuthConfig struct {
	Google   OAuthProvider
	Apple    OAuthProvider
	Facebook OAuthProvider
	Auth0    Auth0Provider
}

type OAuthProvider struct {
	ClientID     string
	ClientSecret string
}

type Auth0Provider struct {
	Domain   string
	ClientID string
}

type AWSConfig struct {
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	S3Bucket         string
	CloudFrontDomain string
}

type ServerConfig struct {
	Port        string
	Environment string
	CORSOrigin  string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

var AppConfig Config

func Load() error {
	// Load .env file if it exists
	godotenv.Load()

	// Database config
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	AppConfig.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "swipe_sports"),
	}

	// Redis config
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	AppConfig.Redis = RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     redisPort,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       redisDB,
	}

	// JWT config
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	AppConfig.JWT = JWTConfig{
		Secret: getEnv("JWT_SECRET", "default-secret-change-in-production"),
		Expiry: jwtExpiry,
	}

	// OAuth config
	AppConfig.OAuth = OAuthConfig{
		Google: OAuthProvider{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		Apple: OAuthProvider{
			ClientID:     getEnv("APPLE_CLIENT_ID", ""),
			ClientSecret: getEnv("APPLE_CLIENT_SECRET", ""),
		},
		Facebook: OAuthProvider{
			ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
			ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
		},
		Auth0: Auth0Provider{
			Domain:   getEnv("AUTH0_DOMAIN", ""),
			ClientID: getEnv("AUTH0_CLIENT_ID", ""),
		},
	}

	// AWS config
	AppConfig.AWS = AWSConfig{
		Region:           getEnv("AWS_REGION", "us-east-1"),
		AccessKeyID:      getEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
		S3Bucket:         getEnv("AWS_S3_BUCKET", ""),
		CloudFrontDomain: getEnv("AWS_CLOUDFRONT_DOMAIN", ""),
	}

	// Server config
	AppConfig.Server = ServerConfig{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENV", "development"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}

	// Rate limit config
	rateLimitRequests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))
	AppConfig.RateLimit = RateLimitConfig{
		Requests: rateLimitRequests,
		Window:   rateLimitWindow,
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 