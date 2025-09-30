package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Gender represents the gender enum
type Gender string

const (
	GenderMale         Gender = "male"
	GenderFemale       Gender = "female"
	GenderNonBinary    Gender = "nonbinary"
	GenderPreferNotSay Gender = "prefer_not_say"
)

// Smoking represents the smoking preference enum
type Smoking string

const (
	SmokingNo           Smoking = "no"
	SmokingYes          Smoking = "yes"
	SmokingOccasionally Smoking = "occasionally"
)

// UserProfile represents the user_profiles table
type UserProfile struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`
	Bio           *string    `json:"bio" gorm:"type:text"`
	Languages     *string    `json:"languages" gorm:"type:text"` // Comma-separated languages
	DateOfBirth   *time.Time `json:"date_of_birth" gorm:"type:date"`
	Gender        *Gender    `json:"gender" gorm:"type:gender"`
	JobTitle      *string    `json:"job_title" gorm:"type:text"`
	Smoking       *Smoking   `json:"smoking" gorm:"type:smoking"`
	InterestsNote *string    `json:"interests_note" gorm:"type:text"`
	AvatarURL     *string    `json:"avatar_url" gorm:"type:text"`
	HomeLocation  *string    `json:"home_location" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"type:timestamptz;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// BeforeCreate hook for UserProfile
func (up *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	return nil
}

// GetAge calculates age from date of birth
func (up *UserProfile) GetAge() *int {
	if up.DateOfBirth == nil {
		return nil
	}

	now := time.Now()
	age := now.Year() - up.DateOfBirth.Year()

	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < up.DateOfBirth.YearDay() {
		age--
	}

	return &age
}

// GetLanguagesArray returns languages as a slice
func (up *UserProfile) GetLanguagesArray() []string {
	if up.Languages == nil {
		return []string{}
	}

	// This would need JSON unmarshaling in practice
	// For now, return empty slice
	return []string{}
}
