package domain

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	Id        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (base *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if base.Id == uuid.Nil {
		base.Id = uuid.New()
	}
	return
}
