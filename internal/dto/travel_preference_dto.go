package dto

import "time"

// TravelPreferenceResponse represents a travel preference response
type TravelPreferenceResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TravelStyle string    `json:"travel_style"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TravelPreferenceListResponse represents a travel preference list response
type TravelPreferenceListResponse struct {
	Preferences []TravelPreferenceResponse `json:"preferences"`
}

// AddTravelPreferenceRequest represents an add travel preference request
type AddTravelPreferenceRequest struct {
	TravelStyle string `json:"travel_style" binding:"required"` // Validation against database is done in service layer
}

// UpdateAllTravelPreferencesRequest represents an update all travel preferences request
type UpdateAllTravelPreferencesRequest struct {
	TravelStyles []string `json:"travel_styles" binding:"required,min=1"`
}

// TravelPreferenceStyleResponse represents a travel preference style response
type TravelPreferenceStyleResponse struct {
	Style       string `json:"style"`
	DisplayName string `json:"display_name"`
	Icon        string `json:"icon"`
	IsSelected  bool   `json:"is_selected,omitempty"`
}

// TravelPreferenceStylesResponse represents available travel preference styles
type TravelPreferenceStylesResponse struct {
	Styles []TravelPreferenceStyleResponse `json:"styles"`
}

// TravelPreferenceStatsResponse represents travel preference statistics
type TravelPreferenceStatsResponse struct {
	TotalPreferences int `json:"total_preferences"`
}
