package models

import (
	"github.com/google/uuid"
)

// EventTag represents the event_tags table (many-to-many relationship)
type EventTag struct {
	EventID uuid.UUID `json:"event_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	TagID   uuid.UUID `json:"tag_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`

	// Relationships
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Tag   *Tag   `json:"tag,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for EventTag
func (EventTag) TableName() string {
	return "event_tags"
}
