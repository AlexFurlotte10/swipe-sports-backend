package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"swipe-sports-backend/internal/auth"
	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/service"
)

type SwipeHandler struct {
	swipeService *service.SwipeService
}

func NewSwipeHandler() *SwipeHandler {
	return &SwipeHandler{
		swipeService: service.NewSwipeService(),
	}
}

// GET /profiles
func (h *SwipeHandler) GetProfiles(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	filter := models.ProfileFilter{}

	// Gender filter
	if gender := c.Query("gender"); gender != "" {
		g := models.Gender(gender)
		if g == models.GenderMale || g == models.GenderFemale || g == models.GenderOther {
			filter.Gender = &g
		}
	}

	// Location filter
	if location := c.Query("location"); location != "" {
		filter.Location = &location
	}

	// Rank range filters
	if minRankStr := c.Query("min_rank"); minRankStr != "" {
		if minRank, err := strconv.Atoi(minRankStr); err == nil {
			filter.MinRank = &minRank
		}
	}

	if maxRankStr := c.Query("max_rank"); maxRankStr != "" {
		if maxRank, err := strconv.Atoi(maxRankStr); err == nil {
			filter.MaxRank = &maxRank
		}
	}

	// Location-based filters
	if latStr := c.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			filter.Latitude = &lat
		}
	}

	if lngStr := c.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			filter.Longitude = &lng
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			filter.Radius = &radius
		}
	}

	// Pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 20 // Default limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Get profiles
	profiles, err := h.swipeService.GetProfilesForSwipe(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"count":    len(profiles),
		"has_more": len(profiles) == filter.Limit,
	})
}

// POST /swipe
func (h *SwipeHandler) Swipe(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate swipe direction
	if req.Direction != models.DirectionLeft && req.Direction != models.DirectionRight {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swipe direction"})
		return
	}

	// Validate that user is not swiping on themselves
	if userID == req.SwipeeID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot swipe on your own profile"})
		return
	}

	// Process the swipe
	response, err := h.swipeService.Swipe(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GET /matches
func (h *SwipeHandler) GetMatches(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	matches, err := h.swipeService.GetMatches(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// GET /matches/:id
func (h *SwipeHandler) GetMatch(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	matchIDStr := c.Param("id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	match, err := h.swipeService.GetMatch(matchID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if match == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, match)
} 