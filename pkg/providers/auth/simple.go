package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// simpleAuth is a basic in-memory auth provider
type simpleAuth struct {
	tokens map[string]*tokenInfo
	users  map[string]*types.User
	mutex  sync.RWMutex
}

// tokenInfo holds token metadata
type tokenInfo struct {
	userID    string
	expiresAt time.Time
	tokenType string // "access" or "refresh"
}

// AuthProvider interface - defined locally to avoid import cycle
type AuthProvider interface {
	ValidateToken(token string) (*types.UserClaims, error)
	GenerateToken(user *types.User) (string, error)
	RefreshToken(refreshToken string) (*types.TokenPair, error)
	RevokeToken(token string) error
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type         string        `yaml:"type"`
	Secret       string        `yaml:"secret"`
	Issuer       string        `yaml:"issuer"`
	Audience     string        `yaml:"audience"`
	AccessTTL    time.Duration `yaml:"access_ttl"`
	RefreshTTL   time.Duration `yaml:"refresh_ttl"`
	Algorithm    string        `yaml:"algorithm"`
	PublicKeyURL string        `yaml:"public_key_url"`
}

// NewSimple creates a new simple auth provider
func NewSimple(config AuthConfig) (AuthProvider, error) {
	auth := &simpleAuth{
		tokens: make(map[string]*tokenInfo),
		users:  make(map[string]*types.User),
	}

	// Add a default test user
	testUser := &types.User{
		ID:       "test-user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user"},
		Metadata: map[string]string{
			"created_at": time.Now().Format(time.RFC3339),
		},
	}
	auth.users[testUser.ID] = testUser

	return auth, nil
}

// ValidateToken validates a token and returns user claims
func (a *simpleAuth) ValidateToken(token string) (*types.UserClaims, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	tokenData, exists := a.tokens[token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	if time.Now().After(tokenData.expiresAt) {
		return nil, errors.New("token expired")
	}

	// Get user
	user, exists := a.users[tokenData.userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Create user claims
	claims := &types.UserClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     user.Roles,
		ExpiresAt: tokenData.expiresAt.Unix(),
		IssuedAt:  time.Now().Unix(),
		Metadata:  user.Metadata,
	}

	return claims, nil
}

// GenerateToken generates tokens for a user
func (a *simpleAuth) GenerateToken(user *types.User) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token info
	a.tokens[token] = &tokenInfo{
		userID:    user.ID,
		expiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
		tokenType: "access",
	}

	// Ensure user exists in our store
	a.users[user.ID] = user

	return token, nil
}

// RefreshToken refreshes an existing token
func (a *simpleAuth) RefreshToken(refreshToken string) (*types.TokenPair, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	tokenData, exists := a.tokens[refreshToken]
	if !exists {
		return nil, errors.New("invalid refresh token")
	}

	if tokenData.tokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	// Check if refresh token is expired
	if time.Now().After(tokenData.expiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// Get user
	user, exists := a.users[tokenData.userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	accessTokenBytes := make([]byte, 32)
	if _, err := rand.Read(accessTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	accessToken := hex.EncodeToString(accessTokenBytes)

	// Generate new refresh token
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	// Store new tokens
	a.tokens[accessToken] = &tokenInfo{
		userID:    user.ID,
		expiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
		tokenType: "access",
	}

	a.tokens[newRefreshToken] = &tokenInfo{
		userID:    user.ID,
		expiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 day expiry
		tokenType: "refresh",
	}

	// Remove old refresh token
	delete(a.tokens, refreshToken)

	return &types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
		TokenType:    "Bearer",
	}, nil
}

// RevokeToken revokes a token
func (a *simpleAuth) RevokeToken(token string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if _, exists := a.tokens[token]; !exists {
		return errors.New("token not found")
	}

	delete(a.tokens, token)
	return nil
}

// Helper methods for testing and user management

// AddUser adds a user to the auth provider (useful for testing)
func (a *simpleAuth) AddUser(user *types.User) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.users[user.ID] = user
}

// GetUser gets a user by ID (useful for testing)
func (a *simpleAuth) GetUser(userID string) (*types.User, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	user, exists := a.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetTokenCount returns the number of active tokens (useful for testing)
func (a *simpleAuth) GetTokenCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.tokens)
}