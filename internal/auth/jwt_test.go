package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := int64(123)
	email := "test@example.com"

	token, err := GenerateToken(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	userID := int64(123)
	email := "test@example.com"

	// Generate a token
	token, err := GenerateToken(userID, email)
	assert.NoError(t, err)

	// Validate the token
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Test with invalid token
	claims, err := ValidateToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestRefreshToken(t *testing.T) {
	userID := int64(123)
	email := "test@example.com"

	// Generate original token
	originalToken, err := GenerateToken(userID, email)
	assert.NoError(t, err)

	// Refresh the token
	newToken, err := RefreshToken(originalToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, originalToken, newToken)

	// Validate the new token
	claims, err := ValidateToken(newToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
} 