package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FoodPreference represents a user's food preference
type FoodPreference struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	FoodCategory    string         `json:"food_category" gorm:"type:varchar(50);not null"`
	PreferenceLevel int            `json:"preference_level" gorm:"not null;check:preference_level IN (1,2,3)"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for FoodPreference
func (FoodPreference) TableName() string {
	return "food_preferences"
}

// FoodCategory represents the available food categories
type FoodCategory string

const (
	FoodCategoryThai          FoodCategory = "thai_food"
	FoodCategoryJapanese      FoodCategory = "japanese_food"
	FoodCategoryChinese       FoodCategory = "chinese_food"
	FoodCategoryInternational FoodCategory = "international_food"
	FoodCategoryHalal         FoodCategory = "halal_food"
	FoodCategoryBuffet        FoodCategory = "buffet"
	FoodCategoryBBQGrill      FoodCategory = "bbq_grill"
)

// PreferenceLevel represents the preference level
type PreferenceLevel int

const (
	PreferenceLevelDislike PreferenceLevel = 1 // 😱
	PreferenceLevelNeutral PreferenceLevel = 2 // 😃
	PreferenceLevelLove    PreferenceLevel = 3 // 🤩
)

// IsValidPreferenceLevel checks if the preference level is valid
func IsValidPreferenceLevel(level int) bool {
	return level >= 1 && level <= 3
}

// GetPreferenceLevelName returns the name of the preference level
func GetPreferenceLevelName(level int) string {
	switch level {
	case 1:
		return "dislike"
	case 2:
		return "neutral"
	case 3:
		return "love"
	default:
		return "unknown"
	}
}

// GetPreferenceLevelEmoji returns the emoji for the preference level
func GetPreferenceLevelEmoji(level int) string {
	switch level {
	case 1:
		return "😱"
	case 2:
		return "😃"
	case 3:
		return "🤩"
	default:
		return "❓"
	}
}

// GetFoodCategoryName returns the display name for food category
func GetFoodCategoryName(category string) string {
	switch category {
	case "thai_food":
		return "Thai Food"
	case "japanese_food":
		return "Japanese Food"
	case "chinese_food":
		return "Chinese Food"
	case "international_food":
		return "International Food"
	case "halal_food":
		return "Halal"
	case "buffet":
		return "Buffet"
	case "bbq_grill":
		return "BBQ / Grill"
	default:
		return category
	}
}

// GetFoodCategoryIcon returns the icon for food category
func GetFoodCategoryIcon(category string) string {
	switch category {
	case "thai_food":
		return "🍛"
	case "japanese_food":
		return "🍣"
	case "chinese_food":
		return "🥟"
	case "international_food":
		return "🌍"
	case "halal_food":
		return "☪️"
	case "buffet":
		return "🍽️"
	case "bbq_grill":
		return "🔥"
	default:
		return "🍴"
	}
}
