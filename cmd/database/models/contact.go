package models

import (
	"time"
)

type Contact struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	ContactID uint      `gorm:"not null"`
	Contact   User      `gorm:"foreignKey:ContactID;references:ID"`
	Status    string    `gorm:"size:20;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
