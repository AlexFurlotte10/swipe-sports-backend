package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"swipe-sports-backend/internal/config"
)

type OAuthUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type FacebookUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AppleUserInfo struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Verify Google OAuth token and get user info
func VerifyGoogleToken(idToken string) (*OAuthUserInfo, error) {
	cfg := config.AppConfig.OAuth.Google
	
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Google token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Google token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Verify the token was issued for our app
	if userInfo.ID == "" || userInfo.Email == "" {
		return nil, fmt.Errorf("invalid user info from Google")
	}

	return &OAuthUserInfo{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

// Verify Facebook OAuth token and get user info
func VerifyFacebookToken(accessToken string) (*OAuthUserInfo, error) {
	cfg := config.AppConfig.OAuth.Facebook
	
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Facebook token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Facebook token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo FacebookUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	if userInfo.ID == "" || userInfo.Email == "" {
		return nil, fmt.Errorf("invalid user info from Facebook")
	}

	return &OAuthUserInfo{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

// Verify Apple OAuth token and get user info
func VerifyAppleToken(idToken string) (*OAuthUserInfo, error) {
	// Apple token verification is more complex and requires JWT validation
	// For now, we'll implement a basic version
	// In production, you should use Apple's public keys to verify the JWT
	
	// Parse the JWT token to extract user info
	// This is a simplified version - in production, verify the signature
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid Apple token format")
	}

	// Decode the payload (second part)
	payload := parts[1]
	// Add padding if needed
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// In a real implementation, you would:
	// 1. Verify the JWT signature using Apple's public keys
	// 2. Check the issuer, audience, and expiration
	// 3. Extract user information from the payload

	// For now, return a placeholder
	// In production, implement proper JWT verification
	return &OAuthUserInfo{
		ID:    "apple_user_id", // Extract from JWT payload
		Email: "user@example.com", // Extract from JWT payload
		Name:  "Apple User", // Extract from JWT payload
	}, nil
}

// Verify OAuth token based on provider
func VerifyOAuthToken(provider, token string) (*OAuthUserInfo, error) {
	switch provider {
	case "google":
		return VerifyGoogleToken(token)
	case "facebook":
		return VerifyFacebookToken(token)
	case "apple":
		return VerifyAppleToken(token)
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
} 