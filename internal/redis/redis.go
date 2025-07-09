package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"swipe-sports-backend/internal/config"
)

var Client *redis.Client

func Init() (*redis.Client, error) {
	cfg := config.AppConfig.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	Client = client
	log.Println("Redis initialized successfully")
	return client, nil
}

// Cache keys
const (
	UserProfileKey     = "user:profile:%d"
	UserMatchesKey     = "user:matches:%d"
	MatchMessagesKey   = "match:messages:%d"
	OnlineUsersKey     = "online:users"
	RateLimitKey       = "rate_limit:%s"
	ProfileCacheKey    = "profiles:swipe:%d"
	ProfileCacheExpiry = 300 // 5 minutes
)

// Cache helper functions
func SetUserProfile(userID int64, data []byte) error {
	ctx := context.Background()
	key := fmt.Sprintf(UserProfileKey, userID)
	return Client.Set(ctx, key, data, 0).Err()
}

func GetUserProfile(userID int64) ([]byte, error) {
	ctx := context.Background()
	key := fmt.Sprintf(UserProfileKey, userID)
	return Client.Get(ctx, key).Bytes()
}

func DeleteUserProfile(userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(UserProfileKey, userID)
	return Client.Del(ctx, key).Err()
}

func SetUserMatches(userID int64, data []byte) error {
	ctx := context.Background()
	key := fmt.Sprintf(UserMatchesKey, userID)
	return Client.Set(ctx, key, data, 0).Err()
}

func GetUserMatches(userID int64) ([]byte, error) {
	ctx := context.Background()
	key := fmt.Sprintf(UserMatchesKey, userID)
	return Client.Get(ctx, key).Bytes()
}

func DeleteUserMatches(userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(UserMatchesKey, userID)
	return Client.Del(ctx, key).Err()
}

func SetMatchMessages(matchID int64, data []byte) error {
	ctx := context.Background()
	key := fmt.Sprintf(MatchMessagesKey, matchID)
	return Client.Set(ctx, key, data, 0).Err()
}

func GetMatchMessages(matchID int64) ([]byte, error) {
	ctx := context.Background()
	key := fmt.Sprintf(MatchMessagesKey, matchID)
	return Client.Get(ctx, key).Bytes()
}

func DeleteMatchMessages(matchID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(MatchMessagesKey, matchID)
	return Client.Del(ctx, key).Err()
}

func AddOnlineUser(userID int64) error {
	ctx := context.Background()
	return Client.SAdd(ctx, OnlineUsersKey, userID).Err()
}

func RemoveOnlineUser(userID int64) error {
	ctx := context.Background()
	return Client.SRem(ctx, OnlineUsersKey, userID).Err()
}

func GetOnlineUsers() ([]string, error) {
	ctx := context.Background()
	return Client.SMembers(ctx, OnlineUsersKey).Result()
}

func IsUserOnline(userID int64) (bool, error) {
	ctx := context.Background()
	return Client.SIsMember(ctx, OnlineUsersKey, userID).Result()
}

// Rate limiting
func CheckRateLimit(identifier string, limit int, window int) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf(RateLimitKey, identifier)
	
	count, err := Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	if count == 1 {
		Client.Expire(ctx, key, int64(window))
	}
	
	return count <= int64(limit), nil
} 