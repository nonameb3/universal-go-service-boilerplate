package domain

import "errors"

// Domain-specific errors for business rules
var (
	// Item validation errors
	ErrItemNameRequired    = errors.New("item name is required")
	ErrItemNameTooLong     = errors.New("item name cannot exceed 100 characters")
	ErrItemAmountTooLarge  = errors.New("item amount cannot exceed 999999")
	
	// Item business logic errors
	ErrItemNotFound        = errors.New("item not found")
	ErrItemAlreadyExists   = errors.New("item already exists")
	ErrItemCannotBeDeleted = errors.New("item cannot be deleted")
	
	// Pagination errors
	ErrInvalidPagination   = errors.New("invalid pagination parameters")
	ErrPageTooLarge        = errors.New("page number too large")
	ErrLimitTooLarge       = errors.New("limit too large")
)