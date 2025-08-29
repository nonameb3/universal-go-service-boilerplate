package dto

import (
	"github.com/universal-go-service/boilerplate/internal/domain"
)

// PaginationRequest represents the business request for pagination
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Validate performs business validation and applies business rules for pagination
func (r *PaginationRequest) Validate() error {
	if r.Page < 0 {
		return domain.ErrInvalidPagination
	}
	
	if r.Limit < 0 {
		return domain.ErrInvalidPagination
	}
	
	// Business rule: Maximum limit is 100
	if r.Limit > 100 {
		return domain.ErrLimitTooLarge
	}
	
	return nil
}

// ApplyDefaults applies business default values
func (r *PaginationRequest) ApplyDefaults() {
	// Business rule: Default page is 1
	if r.Page <= 0 {
		r.Page = 1
	}
	
	// Business rule: Default limit is 10
	if r.Limit <= 0 {
		r.Limit = 10
	}
	
	// Business rule: Maximum limit is 100
	if r.Limit > 100 {
		r.Limit = 100
	}
}