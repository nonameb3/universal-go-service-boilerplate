package providers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// healthChecker implements the HealthChecker interface
type healthChecker struct {
	providers   *Providers
	checks      map[string]func(context.Context) error
	startTime   time.Time
	mutex       sync.RWMutex
}

// NewHealthChecker creates a new health checker with all providers
func NewHealthChecker(providers *Providers) HealthChecker {
	hc := &healthChecker{
		providers: providers,
		checks:    make(map[string]func(context.Context) error),
		startTime: time.Now(),
	}

	// Register default provider health checks
	hc.registerDefaultChecks()

	return hc
}

// registerDefaultChecks registers health checks for all providers
func (h *healthChecker) registerDefaultChecks() {
	// Database health check
	if h.providers.Database != nil {
		h.RegisterCheck("database", func(ctx context.Context) error {
			return h.providers.Database.Health()
		})
	}

	// Cache health check (if it supports health checking)
	if h.providers.Cache != nil {
		h.RegisterCheck("cache", func(ctx context.Context) error {
			// Try to set and get a test value
			testKey := "health_check"
			testValue := []byte("ok")
			
			if err := h.providers.Cache.Set(ctx, testKey, testValue, time.Second); err != nil {
				return fmt.Errorf("cache set failed: %w", err)
			}
			
			if _, err := h.providers.Cache.Get(ctx, testKey); err != nil {
				return fmt.Errorf("cache get failed: %w", err)
			}
			
			// Clean up
			_ = h.providers.Cache.Delete(ctx, testKey)
			return nil
		})
	}
}

// RegisterCheck registers a custom health check
func (h *healthChecker) RegisterCheck(name string, checker func(context.Context) error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.checks[name] = checker
}

// CheckHealth performs all registered health checks
func (h *healthChecker) CheckHealth(ctx context.Context) HealthStatus {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime),
		Checks:    make(map[string]CheckResult),
	}

	// Run all health checks concurrently
	type checkResult struct {
		name   string
		result CheckResult
	}

	resultsChan := make(chan checkResult, len(h.checks))
	var wg sync.WaitGroup

	for name, checker := range h.checks {
		wg.Add(1)
		go func(checkName string, checkFn func(context.Context) error) {
			defer wg.Done()
			
			start := time.Now()
			result := CheckResult{
				Status:  "pass",
				Latency: "",
			}

			if err := checkFn(ctx); err != nil {
				result.Status = "fail"
				result.Error = err.Error()
				result.Message = fmt.Sprintf("Health check failed: %v", err)
			}

			result.Latency = time.Since(start).String()
			resultsChan <- checkResult{name: checkName, result: result}
		}(name, checker)
	}

	// Wait for all checks to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	hasFailures := false
	hasWarnings := false

	for result := range resultsChan {
		status.Checks[result.name] = result.result
		
		if result.result.Status == "fail" {
			hasFailures = true
		} else if result.result.Status == "warn" {
			hasWarnings = true
		}
	}

	// Determine overall status
	if hasFailures {
		status.Status = "unhealthy"
	} else if hasWarnings {
		status.Status = "degraded"
	}

	return status
}

// Helper functions for common health checks

// DatabaseHealthCheck creates a standard database health check
func DatabaseHealthCheck(db DatabaseProvider) func(context.Context) error {
	return func(ctx context.Context) error {
		return db.Health()
	}
}

// CacheHealthCheck creates a standard cache health check
func CacheHealthCheck(cache CacheProvider) func(context.Context) error {
	return func(ctx context.Context) error {
		testKey := "health_check"
		testValue := []byte("ok")
		
		if err := cache.Set(ctx, testKey, testValue, time.Second); err != nil {
			return fmt.Errorf("cache set failed: %w", err)
		}
		
		if _, err := cache.Get(ctx, testKey); err != nil {
			return fmt.Errorf("cache get failed: %w", err)
		}
		
		// Clean up
		_ = cache.Delete(ctx, testKey)
		return nil
	}
}

// ExternalServiceHealthCheck creates a health check for external services
func ExternalServiceHealthCheck(name, url string) func(context.Context) error {
	return func(ctx context.Context) error {
		// This is a placeholder - in a real implementation, you'd make an HTTP request
		// to the external service endpoint
		return fmt.Errorf("external service health check not implemented for %s", name)
	}
}