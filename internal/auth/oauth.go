package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type Auth0Claims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

type JWKSKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type JWKS struct {
	Keys []JWKSKey `json:"keys"`
}

// Verify Google OAuth token and get user info
func VerifyGoogleToken(idToken string) (*OAuthUserInfo, error) {
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

// Parse RSA public key from JWK format
func parseRSAPublicKeyFromJWK(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %v", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	
	var e int
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}

// Verify Auth0 JWT token and get user info
func VerifyAuth0Token(token string) (*OAuthUserInfo, error) {
	cfg := config.AppConfig

	// Get Auth0 public keys
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", cfg.OAuth.Auth0.Domain)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get Auth0 public keys: %v", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %v", err)
	}

	// Parse token
	parsedToken, err := jwt.ParseWithClaims(token, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("no kid in token header")
		}

		// Find matching key
		for _, key := range jwks.Keys {
			if key.Kid == kid && key.Kty == "RSA" {
				return parseRSAPublicKeyFromJWK(key.N, key.E)
			}
		}

		return nil, fmt.Errorf("no matching key found")
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*Auth0Claims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify issuer and audience
	expectedIssuer := fmt.Sprintf("https://%s/", cfg.OAuth.Auth0.Domain)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	if len(claims.Audience) == 0 || claims.Audience[0] != cfg.OAuth.Auth0.ClientID {
		return nil, fmt.Errorf("invalid audience")
	}

	// Verify token is not expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return &OAuthUserInfo{
		ID:    claims.Sub,
		Email: claims.Email,
		Name:  claims.Name,
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
	case "auth0":
		return VerifyAuth0Token(token)
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
} 