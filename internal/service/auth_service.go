package service

import (
	"encoding/json"
	"fmt"

	"swipe-sports-backend/internal/auth"
	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/redis"
	"swipe-sports-backend/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *AuthService) AuthenticateOAuth(provider, token string) (*AuthResponse, error) {
	// Verify OAuth token
	oauthUser, err := auth.VerifyOAuthToken(provider, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OAuth token: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByOAuthID(oauthUser.ID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// If user doesn't exist, create new user
	if user == nil {
		user = &models.User{
			OAuthID:       &oauthUser.ID,
			OAuthProvider: &provider,
			Name:          oauthUser.Name,
			Email:         &oauthUser.Email,
			Rank:          1000, // Default rank
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWT token
	jwtToken, err := auth.GenerateToken(user.ID, oauthUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Cache user profile
	if err := s.cacheUserProfile(user); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache user profile: %v\n", err)
	}

	return &AuthResponse{
		Token: jwtToken,
		User:  *user,
	}, nil
}

func (s *AuthService) RefreshToken(token string) (string, error) {
	return auth.RefreshToken(token)
}

func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	// Try to get from cache first
	cachedData, err := redis.GetUserProfile(userID)
	if err == nil {
		var user models.User
		if err := json.Unmarshal(cachedData, &user); err == nil {
			return &user, nil
		}
	}

	// Get from database
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user != nil {
		// Cache the user profile
		if err := s.cacheUserProfile(user); err != nil {
			fmt.Printf("Failed to cache user profile: %v\n", err)
		}
	}

	return user, nil
}

func (s *AuthService) UpdateUser(userID int64, updateReq models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields if provided
	if updateReq.Name != nil {
		user.Name = *updateReq.Name
	}
	if updateReq.Gender != nil {
		user.Gender = updateReq.Gender
	}
	if updateReq.Location != nil {
		user.Location = updateReq.Location
	}
	if updateReq.Latitude != nil {
		user.Latitude = updateReq.Latitude
	}
	if updateReq.Longitude != nil {
		user.Longitude = updateReq.Longitude
	}
	if updateReq.ProfilePicURL != nil {
		user.ProfilePicURL = updateReq.ProfilePicURL
	}
	if updateReq.Bio != nil {
		user.Bio = updateReq.Bio
	}
	if updateReq.SportPreferences != nil {
		user.SportPreferences = *updateReq.SportPreferences
	}
	if updateReq.SkillLevel != nil {
		user.SkillLevel = updateReq.SkillLevel
	}
	if updateReq.PlayStyle != nil {
		user.PlayStyle = updateReq.PlayStyle
	}
	if updateReq.Availability != nil {
		user.Availability = *updateReq.Availability
	}

	// Save to database
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update cache
	if err := s.cacheUserProfile(user); err != nil {
		fmt.Printf("Failed to cache user profile: %v\n", err)
	}

	return user, nil
}

func (s *AuthService) cacheUserProfile(user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	return redis.SetUserProfile(user.ID, data)
} 