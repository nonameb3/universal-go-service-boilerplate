package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// memoryCache is an in-memory cache implementation
type memoryCache struct {
	data   map[string]*cacheItem
	mutex  sync.RWMutex
	ticker *time.Ticker
	done   chan bool
}

// cacheItem represents a cached item with expiration
type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

// CacheProvider interface - defined locally to avoid import cycle
type CacheProvider interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Type        string        `yaml:"type"`
	Address     string        `yaml:"address"`
	Password    string        `yaml:"password"`
	Database    int           `yaml:"database"`
	MaxRetries  int           `yaml:"max_retries"`
	PoolSize    int           `yaml:"pool_size"`
	DefaultTTL  time.Duration `yaml:"default_ttl"`
}

// NewMemory creates a new in-memory cache
func NewMemory(config CacheConfig) (CacheProvider, error) {
	cache := &memoryCache{
		data: make(map[string]*cacheItem),
		done: make(chan bool),
	}

	// Start cleanup routine
	cache.ticker = time.NewTicker(5 * time.Minute)
	go cache.cleanup()

	return cache, nil
}

// Get retrieves a value from the cache
func (c *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		return nil, errors.New("key expired")
	}

	// Make a copy to prevent modification
	value := make([]byte, len(item.value))
	copy(value, item.value)

	return value, nil
}

// Set stores a value in the cache
func (c *memoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Make a copy to prevent external modification
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)

	expiresAt := time.Now().Add(ttl)
	if ttl <= 0 {
		// If no TTL specified, set to a default (1 hour)
		expiresAt = time.Now().Add(time.Hour)
	}

	c.data[key] = &cacheItem{
		value:     valueCopy,
		expiresAt: expiresAt,
	}

	return nil
}

// Delete removes a value from the cache
func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return nil
}

// Exists checks if a key exists in the cache
func (c *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		return false, nil
	}

	return true, nil
}

// Clear removes all values from the cache
func (c *memoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
	return nil
}

// cleanup removes expired items from the cache
func (c *memoryCache) cleanup() {
	for {
		select {
		case <-c.done:
			c.ticker.Stop()
			return
		case <-c.ticker.C:
			c.removeExpired()
		}
	}
}

// removeExpired removes expired items
func (c *memoryCache) removeExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiresAt) {
			delete(c.data, key)
		}
	}
}

// Close stops the cleanup routine
func (c *memoryCache) Close() {
	close(c.done)
}

// Stats returns cache statistics (useful for monitoring)
func (c *memoryCache) Stats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := CacheStats{
		Keys:    len(c.data),
		Expired: 0,
	}

	now := time.Now()
	for _, item := range c.data {
		if now.After(item.expiresAt) {
			stats.Expired++
		}
	}

	return stats
}

// CacheStats represents cache statistics
type CacheStats struct {
	Keys    int `json:"keys"`
	Expired int `json:"expired"`
}