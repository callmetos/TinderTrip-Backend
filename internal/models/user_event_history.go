package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserEventHistory represents the user_event_history table
type UserEventHistory struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID     uuid.UUID  `json:"event_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Completed   bool       `json:"completed" gorm:"type:boolean;not null;default:false"`
	CompletedAt *time.Time `json:"completed_at" gorm:"type:timestamptz"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for UserEventHistory
func (UserEventHistory) TableName() string {
	return "user_event_history"
}

// BeforeCreate hook for UserEventHistory
func (ueh *UserEventHistory) BeforeCreate(tx *gorm.DB) error {
	if ueh.ID == uuid.Nil {
		ueh.ID = uuid.New()
	}
	return nil
}

// MarkCompleted marks the event as completed
func (ueh *UserEventHistory) MarkCompleted() {
	ueh.Completed = true
	now := time.Now()
	ueh.CompletedAt = &now
}

// MarkIncomplete marks the event as incomplete
func (ueh *UserEventHistory) MarkIncomplete() {
	ueh.Completed = false
	ueh.CompletedAt = nil
}

// IsCompleted checks if the event is completed
func (ueh *UserEventHistory) IsCompleted() bool {
	return ueh.Completed
}

// GetCompletionDuration returns how long it took to complete the event
func (ueh *UserEventHistory) GetCompletionDuration() *time.Duration {
	if !ueh.Completed || ueh.CompletedAt == nil {
		return nil
	}

	// This would need the event start time to calculate properly
	// For now, return nil
	return nil
}

// GetDaysSinceCompletion returns the number of days since completion
func (ueh *UserEventHistory) GetDaysSinceCompletion() *int {
	if !ueh.Completed || ueh.CompletedAt == nil {
		return nil
	}

	days := int(time.Since(*ueh.CompletedAt).Hours() / 24)
	return &days
}
