package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog represents the audit_logs table
type AuditLog struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ActorUserID *uuid.UUID `json:"actor_user_id" gorm:"type:uuid;constraint:OnDelete:SET NULL"`
	EntityTable string     `json:"entity_table" gorm:"type:text;not null"`
	EntityID    *uuid.UUID `json:"entity_id" gorm:"type:uuid"`
	Action      string     `json:"action" gorm:"type:text;not null"`
	BeforeData  *string    `json:"before_data" gorm:"type:jsonb"`
	AfterData   *string    `json:"after_data" gorm:"type:jsonb"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	ActorUser *User `json:"actor_user,omitempty" gorm:"foreignKey:ActorUserID;constraint:OnDelete:SET NULL"`
}

// TableName returns the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate hook for AuditLog
func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}
