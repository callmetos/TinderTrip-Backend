package models

import (
	"time"

	"github.com/google/uuid"
)

// EventInterest represents the event_interests table (many-to-many relationship)
// This connects events to the unified interests table
type EventInterest struct {
	EventID     uuid.UUID `json:"event_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	InterestID  uuid.UUID `json:"interest_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Event    *Event    `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Interest *Interest `json:"interest,omitempty" gorm:"foreignKey:InterestID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for EventInterest
func (EventInterest) TableName() string {
	return "event_interests"
}

