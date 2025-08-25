package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/universal-go-service/boilerplate/pkg/providers/auth"
	"github.com/universal-go-service/boilerplate/pkg/providers/cache"
	"github.com/universal-go-service/boilerplate/pkg/providers/database"
)

// Placeholder implementations to avoid import cycles
// These can be replaced with proper implementations once interface alignment is fixed

// placeholderLogger is a simple logger that shows configuration worked
type placeholderLogger struct {
	loggerType  string
	serviceName string
}

func (p *placeholderLogger) Info(msg string, fields ...Field) {
	fmt.Printf("[%s] [%s] INFO: %s\n", p.loggerType, p.serviceName, msg)
}

func (p *placeholderLogger) Error(msg string, err error, fields ...Field) {
	fmt.Printf("[%s] [%s] ERROR: %s (err: %v)\n", p.loggerType, p.serviceName, msg, err)
}

func (p *placeholderLogger) Debug(msg string, fields ...Field) {
	fmt.Printf("[%s] [%s] DEBUG: %s\n", p.loggerType, p.serviceName, msg)
}

func (p *placeholderLogger) Warn(msg string, fields ...Field) {
	fmt.Printf("[%s] [%s] WARN: %s\n", p.loggerType, p.serviceName, msg)
}

func (p *placeholderLogger) WithContext(ctx context.Context) Logger {
	return p
}

func (p *placeholderLogger) WithCorrelationID(id string) Logger {
	return p
}

func (p *placeholderLogger) WithFields(fields ...Field) Logger {
	return p
}

// placeholderMetrics is a simple metrics collector that shows configuration worked
type placeholderMetrics struct {
	metricsType string
	serviceName string
}

func (p *placeholderMetrics) IncrementCounter(name string, labels map[string]string) {
	fmt.Printf("[%s] [%s] Counter incremented: %s\n", p.metricsType, p.serviceName, name)
}

func (p *placeholderMetrics) RecordHistogram(name string, value float64, labels map[string]string) {
	fmt.Printf("[%s] [%s] Histogram recorded: %s = %f\n", p.metricsType, p.serviceName, name, value)
}

func (p *placeholderMetrics) RecordGauge(name string, value float64, labels map[string]string) {
	fmt.Printf("[%s] [%s] Gauge recorded: %s = %f\n", p.metricsType, p.serviceName, name, value)
}

func (p *placeholderMetrics) StartTimer(name string) Timer {
	return &placeholderTimer{name: name, metricsType: p.metricsType, serviceName: p.serviceName, start: time.Now()}
}

// placeholderTimer is a simple timer implementation
type placeholderTimer struct {
	name        string
	metricsType string
	serviceName string
	start       time.Time
}

func (p *placeholderTimer) Stop(labels ...map[string]string) {
	duration := time.Since(p.start)
	fmt.Printf("[%s] [%s] Timer stopped: %s = %v\n", p.metricsType, p.serviceName, p.name, duration)
}

// Providers holds all provider instances
type Providers struct {
	Logger   Logger
	Metrics  MetricsCollector
	Auth     AuthProvider
	Cache    CacheProvider
	Database DatabaseProvider
	Health   HealthChecker
}

// ProviderRegistry maps provider types to their factory functions
type ProviderRegistry struct {
	loggerFactories   map[string]LoggerFactory
	metricsFactories  map[string]MetricsFactory
	authFactories     map[string]AuthFactory
	cacheFactories    map[string]CacheFactory
	databaseFactories map[string]DatabaseFactory
}

// Factory function types
type LoggerFactory func(config LoggerConfig) (Logger, error)
type MetricsFactory func(config MetricsConfig) (MetricsCollector, error)
type AuthFactory func(config AuthConfig) (AuthProvider, error)
type CacheFactory func(config CacheConfig) (CacheProvider, error)
type DatabaseFactory func(config DatabaseConfig) (DatabaseProvider, error)

// NewRegistry creates a new provider registry with default implementations
func NewRegistry() *ProviderRegistry {
	registry := &ProviderRegistry{
		loggerFactories:   make(map[string]LoggerFactory),
		metricsFactories:  make(map[string]MetricsFactory),
		authFactories:     make(map[string]AuthFactory),
		cacheFactories:    make(map[string]CacheFactory),
		databaseFactories: make(map[string]DatabaseFactory),
	}

	// Register default implementations
	registry.registerDefaults()
	return registry
}

// registerDefaults registers all built-in provider implementations
func (r *ProviderRegistry) registerDefaults() {
	// Note: Provider implementations will be registered through wrapper functions
	// to avoid import cycle issues between the factory and individual provider packages.
	// For now, we'll register placeholder implementations that log their type.

	// Logger providers - using simple placeholders until interface alignment is fixed
	r.RegisterLogger("simple", func(config LoggerConfig) (Logger, error) {
		return &placeholderLogger{loggerType: "simple", serviceName: config.ServiceName}, nil
	})
	r.RegisterLogger("structured", func(config LoggerConfig) (Logger, error) {
		return &placeholderLogger{loggerType: "structured", serviceName: config.ServiceName}, nil
	})
	r.RegisterLogger("noop", func(config LoggerConfig) (Logger, error) {
		return &placeholderLogger{loggerType: "noop", serviceName: config.ServiceName}, nil
	})

	// Metrics providers - using simple placeholders
	r.RegisterMetrics("simple", func(config MetricsConfig) (MetricsCollector, error) {
		return &placeholderMetrics{metricsType: "simple", serviceName: config.ServiceName}, nil
	})
	r.RegisterMetrics("prometheus", func(config MetricsConfig) (MetricsCollector, error) {
		return &placeholderMetrics{metricsType: "prometheus", serviceName: config.ServiceName}, nil
	})
	r.RegisterMetrics("noop", func(config MetricsConfig) (MetricsCollector, error) {
		return &placeholderMetrics{metricsType: "noop", serviceName: config.ServiceName}, nil
	})

	// Auth providers (with type conversion adapters)
	r.RegisterAuth("simple", func(config AuthConfig) (AuthProvider, error) {
		authConfig := auth.AuthConfig{
			Type:         config.Type,
			Secret:       config.Secret,
			Issuer:       config.Issuer,
			Audience:     config.Audience,
			AccessTTL:    config.AccessTTL,
			RefreshTTL:   config.RefreshTTL,
			Algorithm:    config.Algorithm,
			PublicKeyURL: config.PublicKeyURL,
		}
		return auth.NewSimple(authConfig)
	})
	r.RegisterAuth("jwt", func(config AuthConfig) (AuthProvider, error) {
		authConfig := auth.AuthConfig{
			Type:         config.Type,
			Secret:       config.Secret,
			Issuer:       config.Issuer,
			Audience:     config.Audience,
			AccessTTL:    config.AccessTTL,
			RefreshTTL:   config.RefreshTTL,
			Algorithm:    config.Algorithm,
			PublicKeyURL: config.PublicKeyURL,
		}
		return auth.NewJWT(authConfig)
	})
	r.RegisterAuth("noop", func(config AuthConfig) (AuthProvider, error) {
		authConfig := auth.AuthConfig{
			Type:         config.Type,
			Secret:       config.Secret,
			Issuer:       config.Issuer,
			Audience:     config.Audience,
			AccessTTL:    config.AccessTTL,
			RefreshTTL:   config.RefreshTTL,
			Algorithm:    config.Algorithm,
			PublicKeyURL: config.PublicKeyURL,
		}
		return auth.NewNoop(authConfig)
	})

	// Cache providers (with type conversion adapters)
	r.RegisterCache("memory", func(config CacheConfig) (CacheProvider, error) {
		cacheConfig := cache.CacheConfig{
			Type:       config.Type,
			Address:    config.Address,
			Password:   config.Password,
			Database:   config.Database,
			MaxRetries: config.MaxRetries,
			PoolSize:   config.PoolSize,
			DefaultTTL: config.DefaultTTL,
		}
		return cache.NewMemory(cacheConfig)
	})
	r.RegisterCache("redis", func(config CacheConfig) (CacheProvider, error) {
		cacheConfig := cache.CacheConfig{
			Type:       config.Type,
			Address:    config.Address,
			Password:   config.Password,
			Database:   config.Database,
			MaxRetries: config.MaxRetries,
			PoolSize:   config.PoolSize,
			DefaultTTL: config.DefaultTTL,
		}
		return cache.NewRedis(cacheConfig)
	})
	r.RegisterCache("noop", func(config CacheConfig) (CacheProvider, error) {
		cacheConfig := cache.CacheConfig{
			Type:       config.Type,
			Address:    config.Address,
			Password:   config.Password,
			Database:   config.Database,
			MaxRetries: config.MaxRetries,
			PoolSize:   config.PoolSize,
			DefaultTTL: config.DefaultTTL,
		}
		return cache.NewNoop(cacheConfig)
	})

	// Database providers (with type conversion adapters)
	r.RegisterDatabase("postgres", func(config DatabaseConfig) (DatabaseProvider, error) {
		dbConfig := database.DatabaseConfig{
			Type:            config.Type,
			Host:            config.Host,
			Port:            config.Port,
			Username:        config.Username,
			Password:        config.Password,
			Database:        config.Database,
			SSLMode:         config.SSLMode,
			Timezone:        config.Timezone,
			MaxOpenConns:    config.MaxOpenConns,
			MaxIdleConns:    config.MaxIdleConns,
			ConnMaxLifetime: config.ConnMaxLifetime,
			ConnMaxIdleTime: config.ConnMaxIdleTime,
		}
		return database.NewPostgres(dbConfig)
	})
}

// Register methods for custom providers

// RegisterLogger registers a custom logger factory
func (r *ProviderRegistry) RegisterLogger(name string, factory LoggerFactory) {
	r.loggerFactories[name] = factory
}

// RegisterMetrics registers a custom metrics factory
func (r *ProviderRegistry) RegisterMetrics(name string, factory MetricsFactory) {
	r.metricsFactories[name] = factory
}

// RegisterAuth registers a custom auth factory
func (r *ProviderRegistry) RegisterAuth(name string, factory AuthFactory) {
	r.authFactories[name] = factory
}

// RegisterCache registers a custom cache factory
func (r *ProviderRegistry) RegisterCache(name string, factory CacheFactory) {
	r.cacheFactories[name] = factory
}

// RegisterDatabase registers a custom database factory
func (r *ProviderRegistry) RegisterDatabase(name string, factory DatabaseFactory) {
	r.databaseFactories[name] = factory
}

// Factory methods

// CreateLogger creates a logger instance based on configuration
func (r *ProviderRegistry) CreateLogger(config LoggerConfig) (Logger, error) {
	factory, exists := r.loggerFactories[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown logger type: %s", config.Type)
	}
	return factory(config)
}

// CreateMetrics creates a metrics collector instance based on configuration
func (r *ProviderRegistry) CreateMetrics(config MetricsConfig) (MetricsCollector, error) {
	factory, exists := r.metricsFactories[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown metrics type: %s", config.Type)
	}
	return factory(config)
}

// CreateAuth creates an auth provider instance based on configuration
func (r *ProviderRegistry) CreateAuth(config AuthConfig) (AuthProvider, error) {
	factory, exists := r.authFactories[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown auth type: %s", config.Type)
	}
	return factory(config)
}

// CreateCache creates a cache provider instance based on configuration
func (r *ProviderRegistry) CreateCache(config CacheConfig) (CacheProvider, error) {
	factory, exists := r.cacheFactories[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown cache type: %s", config.Type)
	}
	return factory(config)
}

// CreateDatabase creates a database provider instance based on configuration
func (r *ProviderRegistry) CreateDatabase(config DatabaseConfig) (DatabaseProvider, error) {
	factory, exists := r.databaseFactories[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown database type: %s", config.Type)
	}
	return factory(config)
}

// Default registry instance
var defaultRegistry = NewRegistry()

// Convenience functions using the default registry

// NewProviders creates all providers using the default registry
func NewProviders(config ProvidersConfig) (*Providers, error) {
	return NewProvidersWithRegistry(config, defaultRegistry)
}

// NewProvidersWithRegistry creates all providers using a custom registry
func NewProvidersWithRegistry(config ProvidersConfig, registry *ProviderRegistry) (*Providers, error) {
	providers := &Providers{}

	// Create logger
	loggerInstance, err := registry.CreateLogger(config.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	providers.Logger = loggerInstance

	// Create metrics
	metricsInstance, err := registry.CreateMetrics(config.Metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}
	providers.Metrics = metricsInstance

	// Create auth
	authInstance, err := registry.CreateAuth(config.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}
	providers.Auth = authInstance

	// Create cache
	cacheInstance, err := registry.CreateCache(config.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	providers.Cache = cacheInstance

	// Create database
	databaseInstance, err := registry.CreateDatabase(config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	providers.Database = databaseInstance

	// Create health checker with all providers
	providers.Health = NewHealthChecker(providers)

	return providers, nil
}

// RegisterCustomLogger registers a custom logger in the default registry
func RegisterCustomLogger(name string, factory LoggerFactory) {
	defaultRegistry.RegisterLogger(name, factory)
}

// RegisterCustomMetrics registers a custom metrics collector in the default registry
func RegisterCustomMetrics(name string, factory MetricsFactory) {
	defaultRegistry.RegisterMetrics(name, factory)
}

// RegisterCustomAuth registers a custom auth provider in the default registry
func RegisterCustomAuth(name string, factory AuthFactory) {
	defaultRegistry.RegisterAuth(name, factory)
}

// RegisterCustomCache registers a custom cache provider in the default registry
func RegisterCustomCache(name string, factory CacheFactory) {
	defaultRegistry.RegisterCache(name, factory)
}

// RegisterCustomDatabase registers a custom database provider in the default registry
func RegisterCustomDatabase(name string, factory DatabaseFactory) {
	defaultRegistry.RegisterDatabase(name, factory)
}

// ProvidersConfig holds configuration for all providers
type ProvidersConfig struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Auth     AuthConfig     `yaml:"auth"`
	Cache    CacheConfig    `yaml:"cache"`
	Database DatabaseConfig `yaml:"database"`
}

// GetDefaultProvidersConfig returns sensible default configuration
func GetDefaultProvidersConfig() ProvidersConfig {
	return ProvidersConfig{
		Logger: LoggerConfig{
			Type:        "simple",
			Level:       InfoLevel,
			ServiceName: "universal-service",
			Format:      "text",
		},
		Metrics: MetricsConfig{
			Type:        "simple",
			Enabled:     true,
			Port:        9090,
			Path:        "/metrics",
			ServiceName: "universal-service",
		},
		Auth: AuthConfig{
			Type:      "simple",
			Algorithm: "HS256",
		},
		Cache: CacheConfig{
			Type: "memory",
		},
		Database: DatabaseConfig{
			Type:         "postgres",
			Host:         "localhost",
			Port:         5432,
			MaxOpenConns: 25,
			MaxIdleConns: 10,
		},
	}
}
