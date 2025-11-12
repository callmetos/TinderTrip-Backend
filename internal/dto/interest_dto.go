package dto

import "time"

// InterestResponse represents an interest response
type InterestResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	DisplayName string    `json:"display_name"`
	Icon        *string   `json:"icon,omitempty"`
	Category    string    `json:"category"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	IsSelected  bool      `json:"is_selected,omitempty"` // For user-specific responses
}

// InterestListResponse represents an interest list response
type InterestListResponse struct {
	Interests []InterestResponse `json:"interests"`
	Total     int64              `json:"total"`
}

// UpdateUserInterestsRequest represents a bulk update user interests request
type UpdateUserInterestsRequest struct {
	InterestCodes []string `json:"interest_codes" binding:"required,min=1"`
}

// GetUserInterestsResponse represents user interests response with selection status
type GetUserInterestsResponse struct {
	Interests []InterestResponse `json:"interests"`
	Total     int64              `json:"total"`
}
