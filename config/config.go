package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the complete application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
	Db     DbConfig     `yaml:"db"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	Environment  string        `yaml:"environment"`
}

// AppConfig represents application-specific configuration
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

type DbConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	SSLMode     string
	TimeZone    string
	AutoMigrate bool
}

// getConfig
func GetConfig(environment string) *Config {
	const (
		envDEV      = "dev"
		envFileName = ".env"
		envKey      = "GO_ENV"
	)

	env := os.Getenv(envKey)
	if env == "dev" || env == "development" || env == "local" {
		err := godotenv.Load(envFileName)
		if err != nil {
			panic(err.Error())
		}
	}

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("HOST", "0.0.0.0"),
			Port:         getEnvInt("PORT", 8080),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			Environment:  environment,
		},
		App: AppConfig{
			Name:    "universal-service",
			Version: "1.0.0",
			Debug:   environment == "development" || environment == "local",
		},
		Db: DbConfig{
			Host:        getEnv("DB_HOST", ""),
			Port:        getEnvInt("DB_PORT", 5432),
			User:        getEnv("DB_USERNAME", ""),
			Password:    getEnv("DB_PASSWORD", ""),
			DBName:      getEnv("DB_DATABASE", ""),
			SSLMode:     getEnv("DB_SSL_MODE", "require"),
			TimeZone:    getEnv("DB_TIMEZONE", "Asia/Bangkok"),
			AutoMigrate: StringToBoolean(getEnv("DB_AUTO_MIGRATE", "false")),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func StringToBoolean(source string) bool {
	booleanValue, err := strconv.ParseBool(source)
	if err != nil {
		panic(err)
	}
	return booleanValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}

// parseInt safely parses an integer from string
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// GetEnvironment returns the current environment
func GetEnvironment() string {
	return getEnv("GO_ENV", "development")
}

// IsDevelopment checks if running in development mode
func IsDevelopment() bool {
	env := GetEnvironment()
	return env == "development" || env == "dev" || env == "local"
}

// IsProduction checks if running in production mode
func IsProduction() bool {
	env := GetEnvironment()
	return env == "production" || env == "prod"
}

// IsTest checks if running in test mode
func IsTest() bool {
	env := GetEnvironment()
	return env == "test" || env == "testing"
}
