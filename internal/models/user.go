package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthProvider represents the authentication provider
type AuthProvider string

const (
	AuthProviderPassword AuthProvider = "password"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderApple    AuthProvider = "apple"
	AuthProviderFacebook AuthProvider = "facebook"
)

// User represents the users table
type User struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         *string      `json:"email" gorm:"type:citext;uniqueIndex"`
	Provider      AuthProvider `json:"provider" gorm:"type:auth_provider;not null"`
	PasswordHash  *string      `json:"-" gorm:"type:text"`
	EmailVerified bool         `json:"email_verified" gorm:"type:boolean;not null;default:false"`
	GoogleID      *string      `json:"google_id" gorm:"type:text;uniqueIndex:ux_users_google_id,where:provider='google'"`
	DisplayName   *string      `json:"display_name" gorm:"type:text"`
	LastLoginAt   *time.Time   `json:"last_login_at" gorm:"type:timestamptz"`
	CreatedAt     time.Time    `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt     *time.Time   `json:"deleted_at" gorm:"type:timestamptz;index"`

	// Relationships
	Profile        *UserProfile      `json:"profile,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Availability   *PrefAvailability `json:"availability,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Budget         *PrefBudget       `json:"budget,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	UserTags       []UserTag         `json:"user_tags,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedEvents  []Event           `json:"created_events,omitempty" gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"`
	EventMembers   []EventMember     `json:"event_members,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	EventSwipes    []EventSwipe      `json:"event_swipes,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ChatMessages   []ChatMessage     `json:"chat_messages,omitempty" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	PasswordResets []PasswordReset   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// IsPasswordAuth checks if user uses password authentication
func (u *User) IsPasswordAuth() bool {
	return u.Provider == AuthProviderPassword
}

// IsGoogleAuth checks if user uses Google authentication
func (u *User) IsGoogleAuth() bool {
	return u.Provider == AuthProviderGoogle
}

// HasPassword checks if user has a password set
func (u *User) HasPassword() bool {
	return u.PasswordHash != nil && *u.PasswordHash != ""
}

// GetDisplayName returns display name or email as fallback
func (u *User) GetDisplayName() string {
	if u.DisplayName != nil && *u.DisplayName != "" {
		return *u.DisplayName
	}
	if u.Email != nil {
		return *u.Email
	}
	return "Unknown User"
}
