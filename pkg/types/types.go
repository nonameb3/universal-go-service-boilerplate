package types

import "time"

// User represents a user in the system
type User struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Roles    []string          `json:"roles"`
	Metadata map[string]string `json:"metadata"`
}

// UserClaims represents JWT token claims
type UserClaims struct {
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Roles     []string          `json:"roles"`
	ExpiresAt int64             `json:"exp"`
	IssuedAt  int64             `json:"iat"`
	Metadata  map[string]string `json:"metadata"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// HealthStatus represents overall system health
type HealthStatus struct {
	Status    string                 `json:"status"` // healthy, unhealthy, degraded
	Timestamp time.Time              `json:"timestamp"`
	Uptime    time.Duration          `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents individual health check result
type CheckResult struct {
	Status  string `json:"status"` // pass, fail, warn
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// Field represents a key-value pair for structured logging and metrics
type Field struct {
	Key   string
	Value interface{}
}

// LogLevel represents logging levels
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)