package auth

import (
	"fmt"
)

// NewJWT creates a new JWT auth provider
func NewJWT(config AuthConfig) (AuthProvider, error) {
	return nil, fmt.Errorf("JWT provider not implemented yet")
}