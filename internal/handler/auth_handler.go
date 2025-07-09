package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"swipe-sports-backend/internal/auth"
	"swipe-sports-backend/internal/models"
	"swipe-sports-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

type OAuthRequest struct {
	Provider string `json:"provider" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate provider
	if req.Provider != "google" && req.Provider != "apple" && req.Provider != "facebook" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported OAuth provider"})
		return
	}

	// Authenticate with OAuth
	authResponse, err := h.authService.AuthenticateOAuth(req.Provider, req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate provider
	if req.Provider != "google" && req.Provider != "apple" && req.Provider != "facebook" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported OAuth provider"})
		return
	}

	// Authenticate with OAuth
	authResponse, err := h.authService.AuthenticateOAuth(req.Provider, req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GET /profile/me
func (h *AuthHandler) GetMyProfile(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /profile/me
func (h *AuthHandler) UpdateMyProfile(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.UpdateUser(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /profile/picture
func (h *AuthHandler) UploadProfilePicture(c *gin.Context) {
	userID, exists := auth.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	if file.Header.Get("Content-Type") != "image/jpeg" && 
	   file.Header.Get("Content-Type") != "image/png" && 
	   file.Header.Get("Content-Type") != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, and WebP are allowed"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 5MB"})
		return
	}

	// In a real implementation, you would:
	// 1. Upload the file to S3
	// 2. Get the URL
	// 3. Update the user's profile_pic_url

	// For now, we'll just return a placeholder URL
	profilePicURL := "https://example.com/profile-pictures/" + strconv.FormatInt(userID, 10) + ".jpg"

	// Update user profile
	updateReq := models.UpdateUserRequest{
		ProfilePicURL: &profilePicURL,
	}

	user, err := h.authService.UpdateUser(userID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile_pic_url": profilePicURL,
		"user":            user,
	})
} 