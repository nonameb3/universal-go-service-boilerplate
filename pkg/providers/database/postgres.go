package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewPostgres creates a new PostgreSQL database provider
func NewPostgres(config DatabaseConfig) (DatabaseProvider, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
		config.Timezone,
	)

	gormConfig := &gorm.Config{
		// Disable foreign key constraints for better performance and flexibility
		DisableForeignKeyConstraintWhenMigrating: true,

		// Enable prepared statements for better performance
		PrepareStmt: true,

		// Custom naming strategy (optional)
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // table name prefix
			SingularTable: false, // use singular table name, table for `User` would be `user` with this option enabled
		},

		// Logger configuration
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	return &postgresDatabase{db: db}, nil
}

// GetDB returns the GORM database instance
func (p *postgresDatabase) GetDB() *gorm.DB {
	return p.db
}

// GetSQLDB returns the underlying SQL database instance
func (p *postgresDatabase) GetSQLDB() *sql.DB {
	sqlDB, _ := p.db.DB()
	return sqlDB
}

// Health checks the database connection
func (p *postgresDatabase) Health() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Close closes the database connection
func (p *postgresDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Migrate runs auto-migration for the given models
func (p *postgresDatabase) Migrate(models ...interface{}) error {
	return p.db.AutoMigrate(models...)
}

// Transaction runs a function within a database transaction
func (p *postgresDatabase) Transaction(fn func(*gorm.DB) error) error {
	return p.db.Transaction(fn)
}
