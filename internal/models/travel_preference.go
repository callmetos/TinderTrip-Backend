package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TravelPreference represents a user's travel preference
type TravelPreference struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	TravelStyle string         `json:"travel_style" gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for TravelPreference
func (TravelPreference) TableName() string {
	return "travel_preferences"
}

// TravelStyle represents the available travel styles
type TravelStyle string

const (
	TravelStyleCafeDessert      TravelStyle = "cafe_dessert"
	TravelStyleBubbleTea        TravelStyle = "bubble_tea"
	TravelStyleBakeryCake       TravelStyle = "bakery_cake"
	TravelStyleBingsuIceCream   TravelStyle = "bingsu_ice_cream"
	TravelStyleCoffee           TravelStyle = "coffee"
	TravelStyleMatcha           TravelStyle = "matcha"
	TravelStylePancakes         TravelStyle = "pancakes"
	TravelStyleSocialActivity   TravelStyle = "social_activity"
	TravelStyleKaraoke          TravelStyle = "karaoke"
	TravelStyleGaming           TravelStyle = "gaming"
	TravelStyleMovie            TravelStyle = "movie"
	TravelStyleBoardGame        TravelStyle = "board_game"
	TravelStyleOutdoorActivity  TravelStyle = "outdoor_activity"
	TravelStylePartyCelebration TravelStyle = "party_celebration"
	TravelStyleSwimming         TravelStyle = "swimming"
	TravelStyleSkateboarding    TravelStyle = "skateboarding"
)

// IsValidTravelStyle checks if the travel style is valid
func IsValidTravelStyle(style string) bool {
	validStyles := []string{
		"cafe_dessert",
		"bubble_tea",
		"bakery_cake",
		"bingsu_ice_cream",
		"coffee",
		"matcha",
		"pancakes",
		"social_activity",
		"karaoke",
		"gaming",
		"movie",
		"board_game",
		"outdoor_activity",
		"party_celebration",
		"swimming",
		"skateboarding",
	}

	for _, validStyle := range validStyles {
		if style == validStyle {
			return true
		}
	}
	return false
}

// GetTravelStyleName returns the display name for travel style
func GetTravelStyleName(style string) string {
	switch style {
	case "cafe_dessert":
		return "Cafe & Dessert"
	case "bubble_tea":
		return "Bubble Tea"
	case "bakery_cake":
		return "Bakery / Cake"
	case "bingsu_ice_cream":
		return "Bingsu / Ice Cream"
	case "coffee":
		return "Coffee"
	case "matcha":
		return "Matcha"
	case "pancakes":
		return "Pancakes"
	case "social_activity":
		return "Social Activity"
	case "karaoke":
		return "Karaoke"
	case "gaming":
		return "Gaming"
	case "movie":
		return "Movie"
	case "board_game":
		return "Board Game"
	case "outdoor_activity":
		return "Outdoor Activity"
	case "party_celebration":
		return "Party / Celebration"
	case "swimming":
		return "Swimming"
	case "skateboarding":
		return "Skateboarding"
	default:
		return style
	}
}

// GetTravelStyleIcon returns the icon for travel style
func GetTravelStyleIcon(style string) string {
	switch style {
	case "cafe_dessert":
		return "🍰"
	case "bubble_tea":
		return "🧋"
	case "bakery_cake":
		return "🧁"
	case "bingsu_ice_cream":
		return "🍧"
	case "coffee":
		return "☕"
	case "matcha":
		return "🍵"
	case "pancakes":
		return "🥞"
	case "social_activity":
		return "👥"
	case "karaoke":
		return "🎤"
	case "gaming":
		return "🎮"
	case "movie":
		return "🎬"
	case "board_game":
		return "🎲"
	case "outdoor_activity":
		return "🏃"
	case "party_celebration":
		return "🎉"
	case "swimming":
		return "🏊"
	case "skateboarding":
		return "🛹"
	default:
		return "🎯"
	}
}
