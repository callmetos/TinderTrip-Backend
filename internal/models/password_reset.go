package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordReset represents the password_resets table
type PasswordReset struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Token     string    `json:"token" gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for PasswordReset
func (PasswordReset) TableName() string {
	return "password_resets"
}

// BeforeCreate hook for PasswordReset
func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the password reset token is expired
func (pr *PasswordReset) IsExpired() bool {
	return time.Now().After(pr.ExpiresAt)
}

// IsValid checks if the password reset token is valid (not expired)
func (pr *PasswordReset) IsValid() bool {
	return !pr.IsExpired()
}
