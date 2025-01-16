package db

import "time"

// Users
type User struct {
	ID        uint      `gorm:"primaryKey;column:user_id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	PublicKey string    `gorm:"column:public_key;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// Sessions
type Session struct {
	ID        uint      `gorm:"primaryKey;column:session_id"`
	UserID    uint      `gorm:"column:user_id;not null"`
	Token     string    `gorm:"column:token;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// Contacts
type Contact struct {
	ID            uint      `gorm:"primaryKey;column:contact_id"`
	UserID        uint      `gorm:"column:user_id;not null"`
	ContactUserID uint      `gorm:"column:contact_user_id;not null"`
	Status        string    `gorm:"size:20;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// Messages
type Message struct {
	ID          uint      `gorm:"primaryKey;column:message_id"`
	SenderID    uint      `gorm:"column:sender_id;not null"`
	RecipientID uint      `gorm:"column:recipient_id;not null"`
	Content     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// Groups
type Group struct {
	ID                      uint      `gorm:"primaryKey;column:group_id"`
	Name                    string    `gorm:"size:100;not null"`
	CreatedBy               uint      `gorm:"column:created_by;not null"`
	CreatedAt               time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	RequiresAgeVerification bool      `gorm:"column:requires_age_verification;default:false"`
}

// Group Members
type GroupMember struct {
	GroupID  uint      `gorm:"column:group_id;not null;primaryKey"`
	UserID   uint      `gorm:"column:user_id;not null;primaryKey"`
	JoinedAt time.Time `gorm:"column:joined_at;default:CURRENT_TIMESTAMP"`
}

// Age Proofs
type AgeProof struct {
	ID        uint      `gorm:"primaryKey;column:proof_id"`
	UserID    uint      `gorm:"column:user_id;not null"`
	ProofData string    `gorm:"column:proof_data;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}
