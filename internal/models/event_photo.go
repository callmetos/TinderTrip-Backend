package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventPhoto represents the event_photos table
type EventPhoto struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID   uuid.UUID `json:"event_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	URL       string    `json:"url" gorm:"type:text;not null"`
	SortNo    *int      `json:"sort_no" gorm:"type:int"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for EventPhoto
func (EventPhoto) TableName() string {
	return "event_photos"
}

// BeforeCreate hook for EventPhoto
func (ep *EventPhoto) BeforeCreate(tx *gorm.DB) error {
	if ep.ID == uuid.Nil {
		ep.ID = uuid.New()
	}
	return nil
}

// GetSortOrder returns the sort order or 0 if not set
func (ep *EventPhoto) GetSortOrder() int {
	if ep.SortNo == nil {
		return 0
	}
	return *ep.SortNo
}

// SetSortOrder sets the sort order
func (ep *EventPhoto) SetSortOrder(order int) {
	ep.SortNo = &order
}
