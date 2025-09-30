package models

import (
	"time"

	"github.com/google/uuid"
)

// SwipeDirection represents the swipe direction enum
type SwipeDirection string

const (
	SwipeDirectionLike SwipeDirection = "like"
	SwipeDirectionPass SwipeDirection = "pass"
)

// EventSwipe represents the event_swipes table
type EventSwipe struct {
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	EventID   uuid.UUID      `json:"event_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	Direction SwipeDirection `json:"direction" gorm:"type:swipe_direction;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for EventSwipe
func (EventSwipe) TableName() string {
	return "event_swipes"
}

// IsLike checks if the swipe is a like
func (es *EventSwipe) IsLike() bool {
	return es.Direction == SwipeDirectionLike
}

// IsPass checks if the swipe is a pass
func (es *EventSwipe) IsPass() bool {
	return es.Direction == SwipeDirectionPass
}

// GetSwipeAge returns how long ago the swipe was made
func (es *EventSwipe) GetSwipeAge() time.Duration {
	return time.Since(es.CreatedAt)
}

// IsRecent checks if the swipe was made recently (within 24 hours)
func (es *EventSwipe) IsRecent() bool {
	return es.GetSwipeAge() < 24*time.Hour
}
