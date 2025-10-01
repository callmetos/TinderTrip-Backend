package dto

import "time"

// TagResponse represents a tag response
type TagResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
}

// TagListResponse represents a tag list response
type TagListResponse struct {
	Tags  []TagResponse `json:"tags"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

// UserTagListResponse represents a user tag list response
type UserTagListResponse struct {
	Tags []TagResponse `json:"tags"`
}

// EventTagListResponse represents an event tag list response
type EventTagListResponse struct {
	Tags []TagResponse `json:"tags"`
}

// AddUserTagRequest represents an add user tag request
type AddUserTagRequest struct {
	TagID string `json:"tag_id" binding:"required"`
}

// AddEventTagRequest represents an add event tag request
type AddEventTagRequest struct {
	TagID string `json:"tag_id" binding:"required"`
}

// CreateTagRequest represents a create tag request (admin only)
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
	Kind string `json:"kind" binding:"required,oneof=interest category activity location food transport accommodation"`
}

// EventSuggestionRequest represents an event suggestion request
type EventSuggestionRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// EventSuggestionResponse represents an event suggestion response
type EventSuggestionResponse struct {
	Events []EventSuggestionItem `json:"events"`
	Total  int64                 `json:"total"`
	Page   int                   `json:"page"`
	Limit  int                   `json:"limit"`
}

// EventSuggestionItem represents an event suggestion item with match score
type EventSuggestionItem struct {
	Event       EventResponse `json:"event"`
	MatchScore  float64       `json:"match_score"`
	MatchedTags []TagResponse `json:"matched_tags"`
}
