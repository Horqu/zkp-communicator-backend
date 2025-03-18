package models

import (
	"time"
)

type Message struct {
	ID          uint      `gorm:"primaryKey"`
	ReaderID    uint      `gorm:"not null"`
	Reader      User      `gorm:"foreignKey:ReaderID;references:ID"`
	SenderID    uint      `gorm:"not null"`
	Sender      User      `gorm:"foreignKey:SenderID;references:ID"`
	RecipientID uint      `gorm:"not null"`
	Recipient   User      `gorm:"foreignKey:RecipientID;references:ID"`
	Content     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
