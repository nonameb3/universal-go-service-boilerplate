package cache

import (
	"fmt"
)

// NewRedis creates a new Redis cache provider
func NewRedis(config CacheConfig) (CacheProvider, error) {
	return nil, fmt.Errorf("Redis provider not implemented yet")
}