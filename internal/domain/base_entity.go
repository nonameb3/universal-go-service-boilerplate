package domain

import "time"

type BaseEntity struct {
	Id        uint32     `gorm:"primary_key;autoIncrement"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"index;default:null"`
}
