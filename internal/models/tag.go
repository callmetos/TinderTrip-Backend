package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TagKind represents different kinds of tags
type TagKind string

const (
	TagKindInterest      TagKind = "interest"
	TagKindCategory      TagKind = "category"
	TagKindActivity      TagKind = "activity"
	TagKindLocation      TagKind = "location"
	TagKindFood          TagKind = "food"
	TagKindTransport     TagKind = "transport"
	TagKindAccommodation TagKind = "accommodation"
)

// Tag represents the tags table
type Tag struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:citext;uniqueIndex;not null"`
	Kind      string    `json:"kind" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	UserTags        []UserTag       `json:"user_tags,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
	EventCategories []EventCategory `json:"event_categories,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
	EventTags       []EventTag      `json:"event_tags,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for Tag
func (Tag) TableName() string {
	return "tags"
}

// BeforeCreate hook for Tag
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsInterest checks if the tag is an interest tag
func (t *Tag) IsInterest() bool {
	return t.Kind == string(TagKindInterest)
}

// IsCategory checks if the tag is a category tag
func (t *Tag) IsCategory() bool {
	return t.Kind == string(TagKindCategory)
}

// IsActivity checks if the tag is an activity tag
func (t *Tag) IsActivity() bool {
	return t.Kind == string(TagKindActivity)
}

// IsLocation checks if the tag is a location tag
func (t *Tag) IsLocation() bool {
	return t.Kind == string(TagKindLocation)
}

// IsFood checks if the tag is a food tag
func (t *Tag) IsFood() bool {
	return t.Kind == string(TagKindFood)
}

// IsTransport checks if the tag is a transport tag
func (t *Tag) IsTransport() bool {
	return t.Kind == string(TagKindTransport)
}

// IsAccommodation checks if the tag is an accommodation tag
func (t *Tag) IsAccommodation() bool {
	return t.Kind == string(TagKindAccommodation)
}

// GetKindEnum returns the tag kind as an enum
func (t *Tag) GetKindEnum() TagKind {
	return TagKind(t.Kind)
}
