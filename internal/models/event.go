package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventType represents the event type enum
type EventType string

const (
	EventTypeMeal       EventType = "meal"
	EventTypeOneDayTrip EventType = "one_day_trip"
	EventTypeOvernight  EventType = "overnight"
)

// EventStatus represents the event status enum
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusActive    EventStatus = "active"
	EventStatusClosed    EventStatus = "closed"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

// Event represents the events table
type Event struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatorID     uuid.UUID   `json:"creator_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Title         string      `json:"title" gorm:"type:text;not null"`
	Description   *string     `json:"description" gorm:"type:text"`
	EventType     EventType   `json:"event_type" gorm:"type:event_type;not null;default:'meal'"`
	AddressText   *string     `json:"address_text" gorm:"type:text"`
	Lat           *float64    `json:"lat" gorm:"type:double precision"`
	Lng           *float64    `json:"lng" gorm:"type:double precision"`
	StartAt       *time.Time  `json:"start_at" gorm:"type:timestamptz"`
	EndAt         *time.Time  `json:"end_at" gorm:"type:timestamptz"`
	Capacity      *int        `json:"capacity" gorm:"type:int;check:capacity IS NULL OR capacity >= 1"`
	Status        EventStatus `json:"status" gorm:"type:event_status;not null;default:'draft'"`
	CoverImageURL *string     `json:"cover_image_url" gorm:"type:text"`
	CreatedAt     time.Time   `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt     *time.Time  `json:"deleted_at" gorm:"type:timestamptz;index"`

	// Relationships
	Creator       *User              `json:"creator,omitempty" gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"`
	Photos        []EventPhoto       `json:"photos,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Categories    []EventCategory    `json:"categories,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Tags          []EventTag         `json:"tags,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Members       []EventMember      `json:"members,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Swipes        []EventSwipe       `json:"swipes,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	ChatRoom      *ChatRoom          `json:"chat_room,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	UserHistories []UserEventHistory `json:"user_histories,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for Event
func (Event) TableName() string {
	return "events"
}

// BeforeCreate hook for Event
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// IsActive checks if the event is active
func (e *Event) IsActive() bool {
	return e.Status == EventStatusActive
}

// IsDraft checks if the event is in draft status
func (e *Event) IsDraft() bool {
	return e.Status == EventStatusDraft
}

// IsCompleted checks if the event is completed
func (e *Event) IsCompleted() bool {
	return e.Status == EventStatusCompleted
}

// IsCancelled checks if the event is cancelled
func (e *Event) IsCancelled() bool {
	return e.Status == EventStatusCancelled
}

// IsClosed checks if the event is closed
func (e *Event) IsClosed() bool {
	return e.Status == EventStatusClosed
}

// HasLocation checks if the event has location data
func (e *Event) HasLocation() bool {
	return e.Lat != nil && e.Lng != nil
}

// GetDuration returns the duration of the event
func (e *Event) GetDuration() *time.Duration {
	if e.StartAt == nil || e.EndAt == nil {
		return nil
	}

	duration := e.EndAt.Sub(*e.StartAt)
	return &duration
}

// IsUpcoming checks if the event is in the future
func (e *Event) IsUpcoming() bool {
	if e.StartAt == nil {
		return false
	}
	return e.StartAt.After(time.Now())
}

// IsPast checks if the event is in the past
func (e *Event) IsPast() bool {
	if e.EndAt == nil {
		return false
	}
	return e.EndAt.Before(time.Now())
}
