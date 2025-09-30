package models

import (
	"github.com/google/uuid"
)

// UserTag represents the user_tags table (many-to-many relationship)
type UserTag struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	TagID  uuid.UUID `json:"tag_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Tag  *Tag  `json:"tag,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for UserTag
func (UserTag) TableName() string {
	return "user_tags"
}
