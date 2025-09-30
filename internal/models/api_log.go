package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APILog represents the api_logs table
type APILog struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RequestID  *string    `json:"request_id" gorm:"type:text"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:uuid;constraint:OnDelete:SET NULL"`
	Method     *string    `json:"method" gorm:"type:text"`
	Path       *string    `json:"path" gorm:"type:text"`
	Status     *int       `json:"status" gorm:"type:int"`
	DurationMs *int       `json:"duration_ms" gorm:"type:int"`
	IPAddress  *string    `json:"ip_address" gorm:"type:inet"`
	UserAgent  *string    `json:"user_agent" gorm:"type:text"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

// TableName returns the table name for APILog
func (APILog) TableName() string {
	return "api_logs"
}

// BeforeCreate hook for APILog
func (al *APILog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}
