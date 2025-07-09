package service

import (
	"encoding/json"
	"fmt"

	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/redis"
	"swipe-sports-backend/internal/repository"
)

type SwipeService struct {
	userRepo   *repository.UserRepository
	swipeRepo  *repository.SwipeRepository
}

func NewSwipeService() *SwipeService {
	return &SwipeService{
		userRepo:  repository.NewUserRepository(),
		swipeRepo: repository.NewSwipeRepository(),
	}
}

func (s *SwipeService) Swipe(swiperID int64, swipeReq models.SwipeRequest) (*models.SwipeResponse, error) {
	// Check if user has already swiped on this profile
	existingSwipe, err := s.swipeRepo.GetBySwiperAndSwipee(swiperID, swipeReq.SwipeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing swipe: %w", err)
	}

	if existingSwipe != nil {
		return nil, fmt.Errorf("user has already swiped on this profile")
	}

	// Create the swipe
	swipe := &models.Swipe{
		SwiperID:  swiperID,
		SwipeeID:  swipeReq.SwipeeID,
		Direction: swipeReq.Direction,
	}

	if err := s.swipeRepo.Create(swipe); err != nil {
		return nil, fmt.Errorf("failed to create swipe: %w", err)
	}

	// Check for match only if swiped right
	if swipeReq.Direction == models.DirectionRight {
		isMatch, err := s.swipeRepo.CheckForMatch(swiperID, swipeReq.SwipeeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check for match: %w", err)
		}

		if isMatch {
			// Create match
			match := &models.Match{
				User1ID: swiperID,
				User2ID: swipeReq.SwipeeID,
			}

			if err := s.swipeRepo.CreateMatch(match); err != nil {
				return nil, fmt.Errorf("failed to create match: %w", err)
			}

			// Clear cache for both users
			redis.DeleteUserMatches(swiperID)
			redis.DeleteUserMatches(swipeReq.SwipeeID)

			return &models.SwipeResponse{
				IsMatch: true,
				Match:   match,
			}, nil
		}
	}

	return &models.SwipeResponse{
		IsMatch: false,
	}, nil
}

func (s *SwipeService) GetProfilesForSwipe(userID int64, filter models.ProfileFilter) ([]models.UserProfile, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("profiles:swipe:%d", userID)
	cachedData, err := redis.Client.Get(redis.Client.Context(), cacheKey).Bytes()
	if err == nil {
		var profiles []models.UserProfile
		if err := json.Unmarshal(cachedData, &profiles); err == nil {
			return profiles, nil
		}
	}

	// Get from database
	profiles, err := s.userRepo.GetProfilesForSwipe(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %w", err)
	}

	// Cache the results
	if data, err := json.Marshal(profiles); err == nil {
		redis.Client.Set(redis.Client.Context(), cacheKey, data, redis.ProfileCacheExpiry)
	}

	return profiles, nil
}

func (s *SwipeService) GetMatches(userID int64) ([]models.MatchResponse, error) {
	// Try to get from cache first
	cachedData, err := redis.GetUserMatches(userID)
	if err == nil {
		var matches []models.MatchResponse
		if err := json.Unmarshal(cachedData, &matches); err == nil {
			return matches, nil
		}
	}

	// Get matches from database
	matches, err := s.swipeRepo.GetMatchesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	// Convert to response format with user details
	var matchResponses []models.MatchResponse
	for _, match := range matches {
		// Get the other user's ID
		otherUserID := match.User1ID
		if match.User1ID == userID {
			otherUserID = match.User2ID
		}

		// Get user profile
		user, err := s.userRepo.GetByID(otherUserID)
		if err != nil {
			continue // Skip if user not found
		}

		// Convert to UserProfile
		userProfile := models.UserProfile{
			ID:              user.ID,
			Name:            user.Name,
			Gender:          user.Gender,
			Location:        user.Location,
			Rank:            user.Rank,
			ProfilePicURL:   user.ProfilePicURL,
			Bio:             user.Bio,
			SportPreferences: user.SportPreferences,
			SkillLevel:      user.SkillLevel,
			PlayStyle:       user.PlayStyle,
			Availability:    user.Availability,
			CreatedAt:       user.CreatedAt,
		}

		matchResponses = append(matchResponses, models.MatchResponse{
			ID:        match.ID,
			User:      userProfile,
			CreatedAt: match.CreatedAt,
		})
	}

	// Cache the results
	if data, err := json.Marshal(matchResponses); err == nil {
		redis.SetUserMatches(userID, data)
	}

	return matchResponses, nil
}

func (s *SwipeService) GetMatch(matchID, userID int64) (*models.MatchWithUsers, error) {
	// Check if user is part of the match
	isInMatch, err := s.swipeRepo.IsUserInMatch(userID, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user in match: %w", err)
	}

	if !isInMatch {
		return nil, fmt.Errorf("user not part of this match")
	}

	// Get match with user details
	matchWithUsers, err := s.swipeRepo.GetMatchWithUsers(matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return matchWithUsers, nil
} 