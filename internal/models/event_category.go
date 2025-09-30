package models

import (
	"github.com/google/uuid"
)

// EventCategory represents the event_categories table (many-to-many relationship)
type EventCategory struct {
	EventID uuid.UUID `json:"event_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	TagID   uuid.UUID `json:"tag_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`

	// Relationships
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Tag   *Tag   `json:"tag,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for EventCategory
func (EventCategory) TableName() string {
	return "event_categories"
}
