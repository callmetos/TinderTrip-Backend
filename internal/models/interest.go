package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interest represents the interests master table
type Interest struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code        string    `json:"code" gorm:"type:varchar(100);not null;uniqueIndex"`
	DisplayName string    `json:"display_name" gorm:"type:text;not null"`
	Icon        *string   `json:"icon,omitempty" gorm:"type:text"`
	Category    string    `json:"category" gorm:"type:varchar(50);not null"` // 'cafe', 'activity', 'pub_bar', 'sport'
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	UserInterests  []UserInterest  `json:"user_interests,omitempty" gorm:"foreignKey:InterestID;constraint:OnDelete:CASCADE"`
	EventInterests []EventInterest `json:"event_interests,omitempty" gorm:"foreignKey:InterestID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for Interest
func (Interest) TableName() string {
	return "interests"
}

// BeforeCreate hook for Interest
func (i *Interest) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// UserInterest represents the user_interests table (many-to-many relationship)
type UserInterest struct {
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	InterestID uuid.UUID `json:"interest_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Interest *Interest `json:"interest,omitempty" gorm:"foreignKey:InterestID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for UserInterest
func (UserInterest) TableName() string {
	return "user_interests"
}

