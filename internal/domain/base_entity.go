package domain

import "time"

type BaseEntity struct {
	Id        uint       `gorm:"primary_key;autoIncrement" json:"id"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index;default:null" json:"deleted_at,omitempty"`
}
