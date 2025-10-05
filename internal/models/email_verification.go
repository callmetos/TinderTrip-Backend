package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailVerification represents the email_verifications table
type EmailVerification struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"type:citext;not null;index"`
	OTP       string     `json:"otp" gorm:"type:varchar(6);not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"type:timestamptz;not null;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:timestamptz;index"`
}

// TableName returns the table name for EmailVerification
func (EmailVerification) TableName() string {
	return "email_verifications"
}

// BeforeCreate hook for EmailVerification
func (ev *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the verification is expired
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}
