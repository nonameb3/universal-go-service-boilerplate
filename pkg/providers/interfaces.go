package providers

import (
	"context"
	"database/sql"
	"io"
	"time"

	"gorm.io/gorm"
	
	"github.com/universal-go-service/boilerplate/pkg/types"
)

// Field represents a key-value pair for structured logging and metrics
type Field = types.Field

// Logger interface - universal logging abstraction
type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	WithContext(ctx context.Context) Logger
	WithCorrelationID(id string) Logger
	WithFields(fields ...Field) Logger
}

// MetricsCollector interface - universal metrics abstraction
type MetricsCollector interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
	StartTimer(name string) Timer
}

// Timer interface for measuring durations
type Timer interface {
	Stop(labels ...map[string]string)
}

// AuthProvider interface - universal authentication abstraction
type AuthProvider interface {
	ValidateToken(token string) (*types.UserClaims, error)
	GenerateToken(user *types.User) (string, error)
	RefreshToken(refreshToken string) (*types.TokenPair, error)
	RevokeToken(token string) error
}

// CacheProvider interface - universal caching abstraction
type CacheProvider interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// DatabaseProvider interface - universal database abstraction
type DatabaseProvider interface {
	GetDB() *gorm.DB
	GetSQLDB() *sql.DB
	Health() error
	Close() error
	Migrate(models ...interface{}) error
	Transaction(fn func(*gorm.DB) error) error
}

// HealthChecker interface - universal health checking
type HealthChecker interface {
	CheckHealth(ctx context.Context) types.HealthStatus
	RegisterCheck(name string, checker func(ctx context.Context) error)
}

// ConfigLoader interface - universal configuration loading
type ConfigLoader interface {
	Load(source string) error
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Reload() error
}

// Type aliases for convenience
type User = types.User
type UserClaims = types.UserClaims
type TokenPair = types.TokenPair
type HealthStatus = types.HealthStatus
type CheckResult = types.CheckResult
type LogLevel = types.LogLevel

// LogLevel constants
const (
	DebugLevel = types.DebugLevel
	InfoLevel  = types.InfoLevel
	WarnLevel  = types.WarnLevel
	ErrorLevel = types.ErrorLevel
	FatalLevel = types.FatalLevel
)

// Provider configurations

// ProviderConfig represents base provider configuration
type ProviderConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Type        string            `yaml:"type"`
	Level       LogLevel          `yaml:"level"`
	ServiceName string            `yaml:"service_name"`
	Format      string            `yaml:"format"` // json, text
	Output      io.Writer         `yaml:"-"`
	Fields      map[string]string `yaml:"fields"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Type        string `yaml:"type"`
	Enabled     bool   `yaml:"enabled"`
	Port        int    `yaml:"port"`
	Path        string `yaml:"path"`
	ServiceName string `yaml:"service_name"`
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

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type            string        `yaml:"type"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	SSLMode         string        `yaml:"ssl_mode"`
	Timezone        string        `yaml:"timezone"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}