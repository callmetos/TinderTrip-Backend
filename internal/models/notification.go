package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a user notification
type Notification struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"type:uuid;not null"`
	Title     string                 `json:"title" gorm:"not null"`
	Body      string                 `json:"body" gorm:"not null"`
	Type      string                 `json:"type" gorm:"not null"` // push, email, sms
	Data      map[string]interface{} `json:"data" gorm:"type:jsonb"`
	Read      bool                   `json:"read" gorm:"default:false"`
	CreatedAt time.Time              `json:"created_at" gorm:"not null;default:now()"`
	ReadAt    *time.Time             `json:"read_at"`
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}
