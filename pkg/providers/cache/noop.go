package cache

import (
	"context"
	"errors"
	"time"
)

// noopCache is a cache that doesn't cache anything
type noopCache struct{}

// NewNoop creates a new no-op cache
func NewNoop(config CacheConfig) (CacheProvider, error) {
	return &noopCache{}, nil
}

// Get always returns key not found
func (c *noopCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, errors.New("key not found")
}

// Set does nothing
func (c *noopCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

// Delete does nothing
func (c *noopCache) Delete(ctx context.Context, key string) error {
	return nil
}

// Exists always returns false
func (c *noopCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// Clear does nothing
func (c *noopCache) Clear(ctx context.Context) error {
	return nil
}