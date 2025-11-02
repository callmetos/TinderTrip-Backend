package dto

import "time"

// FoodPreferenceResponse represents a food preference response
type FoodPreferenceResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	FoodCategory    string    `json:"food_category"`
	PreferenceLevel int       `json:"preference_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FoodPreferenceListResponse represents a food preference list response
type FoodPreferenceListResponse struct {
	Preferences []FoodPreferenceResponse `json:"preferences"`
}

// UpdateFoodPreferenceRequest represents an update food preference request
type UpdateFoodPreferenceRequest struct {
	FoodCategory    string `json:"food_category" binding:"required"` // Validation against database is done in service layer
	PreferenceLevel int    `json:"preference_level" binding:"required,min=1,max=3"`
}

// UpdateAllFoodPreferencesRequest represents an update all food preferences request
type UpdateAllFoodPreferencesRequest struct {
	Preferences []UpdateFoodPreferenceRequest `json:"preferences" binding:"required,min=1"`
}

// FoodPreferenceCategoryResponse represents a food preference category response
type FoodPreferenceCategoryResponse struct {
	Category        string `json:"category"`
	DisplayName     string `json:"display_name"`
	Icon            string `json:"icon"`
	PreferenceLevel int    `json:"preference_level,omitempty"`
}

// FoodPreferenceCategoriesResponse represents available food preference categories
type FoodPreferenceCategoriesResponse struct {
	Categories []FoodPreferenceCategoryResponse `json:"categories"`
}

// FoodPreferenceStatsResponse represents food preference statistics
type FoodPreferenceStatsResponse struct {
	TotalPreferences int `json:"total_preferences"`
	DislikeCount     int `json:"dislike_count"`
	NeutralCount     int `json:"neutral_count"`
	LoveCount        int `json:"love_count"`
}
