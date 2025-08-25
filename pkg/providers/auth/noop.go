package auth

import (
	"errors"

	"github.com/universal-go-service/boilerplate/pkg/types"
)

// noopAuth is an auth provider that always fails - useful for testing
type noopAuth struct{}

// NewNoop creates a new no-op auth provider
func NewNoop(config AuthConfig) (AuthProvider, error) {
	return &noopAuth{}, nil
}

// ValidateToken always returns an error
func (a *noopAuth) ValidateToken(token string) (*types.UserClaims, error) {
	return nil, errors.New("no-op auth provider")
}

// GenerateToken always returns an error
func (a *noopAuth) GenerateToken(user *types.User) (string, error) {
	return "", errors.New("no-op auth provider")
}

// RefreshToken always returns an error
func (a *noopAuth) RefreshToken(refreshToken string) (*types.TokenPair, error) {
	return nil, errors.New("no-op auth provider")
}

// RevokeToken always returns an error
func (a *noopAuth) RevokeToken(token string) error {
	return errors.New("no-op auth provider")
}