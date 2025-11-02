package dto

import "time"

// UserProfileResponse represents a user profile response
type UserProfileResponse struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	DisplayName   *string    `json:"display_name,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	Languages     *string    `json:"languages,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Age           *int       `json:"age,omitempty"`
	Gender        string     `json:"gender,omitempty"`
	JobTitle      *string    `json:"job_title,omitempty"`
	Smoking       string     `json:"smoking,omitempty"`
	InterestsNote *string    `json:"interests_note,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	HomeLocation  *string    `json:"home_location,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UpdateUserProfileRequest represents an update user profile request
type UpdateUserProfileRequest struct {
	DisplayName   *string    `json:"display_name,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	Languages     *string    `json:"languages,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Age           *int       `json:"age,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	JobTitle      *string    `json:"job_title,omitempty"`
	Smoking       *string    `json:"smoking,omitempty"`
	InterestsNote *string    `json:"interests_note,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	HomeLocation  *string    `json:"home_location,omitempty"`
}

// UpdateProfileRequest represents an update profile request (alias for compatibility)
type UpdateProfileRequest = UpdateUserProfileRequest

// GetUsersRequest represents a get users request
type GetUsersRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	Limit    int    `form:"limit" binding:"min=1,max=100"`
	Search   string `form:"search"`
	Gender   string `form:"gender"`
	Location string `form:"location"`
}

// GetUsersResponse represents a get users response
type GetUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalEvents     int64   `json:"total_events"`
	CompletedEvents int64   `json:"completed_events"`
	PendingEvents   int64   `json:"pending_events"`
	MealEvents      int64   `json:"meal_events"`
	DayTripEvents   int64   `json:"day_trip_events"`
	OvernightEvents int64   `json:"overnight_events"`
	CompletionRate  float64 `json:"completion_rate"`
}

// SetupStatusResponse represents user setup completion status
type SetupStatusResponse struct {
	SetupCompleted bool `json:"setup_completed"`
}
