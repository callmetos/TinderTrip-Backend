package dto

import "time"

// PrefAvailabilityResponse represents a preference availability response
type PrefAvailabilityResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Mon       bool      `json:"mon"`
	Tue       bool      `json:"tue"`
	Wed       bool      `json:"wed"`
	Thu       bool      `json:"thu"`
	Fri       bool      `json:"fri"`
	Sat       bool      `json:"sat"`
	Sun       bool      `json:"sun"`
	AllDay    bool      `json:"all_day"`
	Morning   bool      `json:"morning"`
	Afternoon bool      `json:"afternoon"`
	TimeRange *string   `json:"time_range,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PrefBudgetResponse represents a preference budget response
type PrefBudgetResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	MealMin      *int      `json:"meal_min,omitempty"`
	MealMax      *int      `json:"meal_max,omitempty"`
	DaytripMin   *int      `json:"daytrip_min,omitempty"`
	DaytripMax   *int      `json:"daytrip_max,omitempty"`
	OvernightMin *int      `json:"overnight_min,omitempty"`
	OvernightMax *int      `json:"overnight_max,omitempty"`
	Unlimited    bool      `json:"unlimited"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UpdatePrefAvailabilityRequest represents an update preference availability request
type UpdatePrefAvailabilityRequest struct {
	Mon       *bool   `json:"mon,omitempty"`
	Tue       *bool   `json:"tue,omitempty"`
	Wed       *bool   `json:"wed,omitempty"`
	Thu       *bool   `json:"thu,omitempty"`
	Fri       *bool   `json:"fri,omitempty"`
	Sat       *bool   `json:"sat,omitempty"`
	Sun       *bool   `json:"sun,omitempty"`
	AllDay    *bool   `json:"all_day,omitempty"`
	Morning   *bool   `json:"morning,omitempty"`
	Afternoon *bool   `json:"afternoon,omitempty"`
	TimeRange *string `json:"time_range,omitempty"`
}

// UpdatePrefBudgetRequest represents an update preference budget request
type UpdatePrefBudgetRequest struct {
	MealMin      *int    `json:"meal_min,omitempty"`
	MealMax      *int    `json:"meal_max,omitempty"`
	DaytripMin   *int    `json:"daytrip_min,omitempty"`
	DaytripMax   *int    `json:"daytrip_max,omitempty"`
	OvernightMin *int    `json:"overnight_min,omitempty"`
	OvernightMax *int    `json:"overnight_max,omitempty"`
	Unlimited    *bool   `json:"unlimited,omitempty"`
	Currency     *string `json:"currency,omitempty"`
}

// GetPreferencesRequest represents a get preferences request
type GetPreferencesRequest struct {
	UserID string `form:"user_id" binding:"required"`
}

// GetPreferencesResponse represents a get preferences response
type GetPreferencesResponse struct {
	Availability *PrefAvailabilityResponse `json:"availability,omitempty"`
	Budget       *PrefBudgetResponse       `json:"budget,omitempty"`
}
