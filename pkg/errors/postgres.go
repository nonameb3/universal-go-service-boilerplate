package errors

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/universal-go-service/boilerplate/internal/domain"
)

// Database error codes for different database systems
const (
	// PostgreSQL error codes
	PostgreSQLUniqueViolation     = "23505"
	PostgreSQLForeignKeyViolation = "23503"
	PostgreSQLCheckViolation      = "23514"

	// MySQL error codes (for future support)
	MySQLDuplicateEntry       = 1062
	MySQLForeignKeyConstraint = 1452

	// SQLite error codes (for future support)
	SQLiteConstraintUnique     = 19
	SQLiteConstraintForeignKey = 787
)

// ErrorHandler provides database error detection and mapping utilities
type ErrorHandler struct{}

// NewErrorHandler creates a new database error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// IsUniqueConstraintViolation checks if the error is a unique constraint violation
func (eh *ErrorHandler) IsUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL unique constraint violation
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PostgreSQLUniqueViolation
	}

	// Check for common unique constraint keywords in error message
	errMsg := strings.ToLower(err.Error())
	uniqueKeywords := []string{
		"unique constraint",
		"duplicate key",
		"duplicate entry",
		"unique violation",
		"uniqueindex",
		"idx_items_name", // Our specific unique index name
	}

	for _, keyword := range uniqueKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// IsForeignKeyConstraintViolation checks if the error is a foreign key constraint violation
func (eh *ErrorHandler) IsForeignKeyConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL foreign key constraint violation
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PostgreSQLForeignKeyViolation
	}

	// Check for common foreign key keywords in error message
	errMsg := strings.ToLower(err.Error())
	foreignKeyKeywords := []string{
		"foreign key constraint",
		"foreign key violation",
		"violates foreign key",
		"fk constraint",
	}

	for _, keyword := range foreignKeyKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// MapDatabaseError converts database-specific errors to domain errors
func (eh *ErrorHandler) MapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Handle GORM-specific errors first
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrItemNotFound
	}

	// Handle constraint violations
	if eh.IsUniqueConstraintViolation(err) {
		return domain.ErrItemAlreadyExists
	}

	if eh.IsForeignKeyConstraintViolation(err) {
		// For future use when we have foreign key relationships
		return err // Return original error for now
	}

	// Return original error if not a recognized database error
	return err
}

// SafeDBOperation wraps a database operation with error mapping
func (eh *ErrorHandler) SafeDBOperation(operation func() error) error {
	err := operation()
	return eh.MapDatabaseError(err)
}

// SafeDBOperationWithResult wraps a database operation that returns a result with error mapping
func SafeDBOperationWithResult[T any](eh *ErrorHandler, operation func() (T, error)) (T, error) {
	result, err := operation()
	mappedErr := eh.MapDatabaseError(err)
	return result, mappedErr
}
