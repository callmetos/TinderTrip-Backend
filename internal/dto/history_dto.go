package dto

import "time"

// UserEventHistoryResponse represents a user event history response
type UserEventHistoryResponse struct {
	ID          string         `json:"id"`
	EventID     string         `json:"event_id"`
	UserID      string         `json:"user_id"`
	Event       *EventResponse `json:"event,omitempty"`
	User        *UserResponse  `json:"user,omitempty"`
	Completed   bool           `json:"completed"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// MarkEventCompleteRequest represents a mark event complete request
type MarkEventCompleteRequest struct {
	EventID string `json:"event_id" binding:"required"`
}

// HistoryListResponse represents a history list response
type HistoryListResponse struct {
	History []UserEventHistoryResponse `json:"history"`
	Total   int64                      `json:"total"`
	Page    int                        `json:"page"`
	Limit   int                        `json:"limit"`
}

// GetEventHistoryRequest represents a get event history request
type GetEventHistoryRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	Completed *bool  `form:"completed"`
	EventType string `form:"event_type"`
}

// GetEventHistoryResponse represents a get event history response
type GetEventHistoryResponse struct {
	History    []UserEventHistoryResponse `json:"history"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// GetUserStatsRequest represents a get user stats request
type GetUserStatsRequest struct {
	UserID string `form:"user_id" binding:"required"`
}

// GetUserStatsResponse represents a get user stats response
type GetUserStatsResponse struct {
	TotalEvents     int64   `json:"total_events"`
	CompletedEvents int64   `json:"completed_events"`
	PendingEvents   int64   `json:"pending_events"`
	MealEvents      int64   `json:"meal_events"`
	DayTripEvents   int64   `json:"day_trip_events"`
	OvernightEvents int64   `json:"overnight_events"`
	CompletionRate  float64 `json:"completion_rate"`
}
