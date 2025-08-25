package database

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// postgresDatabase implements DatabaseProvider for PostgreSQL
type postgresDatabase struct {
	db *gorm.DB
}

// DatabaseProvider interface - defined locally to avoid import cycle
type DatabaseProvider interface {
	GetDB() *gorm.DB
	GetSQLDB() *sql.DB
	Health() error
	Close() error
	Migrate(models ...interface{}) error
	Transaction(fn func(*gorm.DB) error) error
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
