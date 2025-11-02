package dto

import "time"

// EventResponse represents an event response
type EventResponse struct {
	ID            string                `json:"id"`
	CreatorID     string                `json:"creator_id"`
	Title         string                `json:"title"`
	Description   *string               `json:"description,omitempty"`
	EventType     string                `json:"event_type"`
	AddressText   *string               `json:"address_text,omitempty"`
	Lat           *float64              `json:"lat,omitempty"`
	Lng           *float64              `json:"lng,omitempty"`
	StartAt       *time.Time            `json:"start_at,omitempty"`
	EndAt         *time.Time            `json:"end_at,omitempty"`
	Capacity      *int                  `json:"capacity,omitempty"`
	BudgetMin     *int                  `json:"budget_min,omitempty"`
	BudgetMax     *int                  `json:"budget_max,omitempty"`
	Currency      *string               `json:"currency,omitempty"`
	Status        string                `json:"status"`
	CoverImageURL *string               `json:"cover_image_url,omitempty"`
	Creator       *UserResponse         `json:"creator,omitempty"`
	Photos        []EventPhotoResponse  `json:"photos,omitempty"`
	Categories    []TagResponse         `json:"categories,omitempty"`
	Tags          []TagResponse         `json:"tags,omitempty"`
	Members       []EventMemberResponse `json:"members,omitempty"`
	MemberCount   int                   `json:"member_count"`
	IsJoined      bool                  `json:"is_joined"`
	UserSwipe     *EventSwipeResponse   `json:"user_swipe,omitempty"`
	MatchScore    *float64              `json:"match_score,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// EventPhotoResponse represents an event photo response
type EventPhotoResponse struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	URL       string    `json:"url"`
	SortNo    *int      `json:"sort_no,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// EventMemberResponse represents an event member response
type EventMemberResponse struct {
	EventID     string     `json:"event_id"`
	UserID      string     `json:"user_id"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	JoinedAt    time.Time  `json:"joined_at"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	LeftAt      *time.Time `json:"left_at,omitempty"`
	Note        *string    `json:"note,omitempty"`
}

// EventSwipeResponse represents an event swipe response
type EventSwipeResponse struct {
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Direction string    `json:"direction"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateEventRequest represents a create event request
type CreateEventRequest struct {
	Title         string     `json:"title" binding:"required"`
	Description   *string    `json:"description,omitempty"`
	EventType     string     `json:"event_type" binding:"required,oneof=meal one_day_trip overnight"`
	AddressText   *string    `json:"address_text,omitempty"`
	Lat           *float64   `json:"lat,omitempty"`
	Lng           *float64   `json:"lng,omitempty"`
	StartAt       *time.Time `json:"start_at,omitempty"`
	EndAt         *time.Time `json:"end_at,omitempty"`
	Capacity      *int       `json:"capacity,omitempty"`
	BudgetMin     *int       `json:"budget_min,omitempty"`
	BudgetMax     *int       `json:"budget_max,omitempty"`
	Currency      *string    `json:"currency,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	CategoryIDs   []string   `json:"category_ids,omitempty"`
	TagIDs        []string   `json:"tag_ids,omitempty"`
}

// UpdateEventRequest represents an update event request
type UpdateEventRequest struct {
	Title         *string    `json:"title,omitempty"`
	Description   *string    `json:"description,omitempty"`
	EventType     *string    `json:"event_type,omitempty" binding:"omitempty,oneof=meal daytrip overnight activity other"`
	AddressText   *string    `json:"address_text,omitempty"`
	Lat           *float64   `json:"lat,omitempty"`
	Lng           *float64   `json:"lng,omitempty"`
	StartAt       *time.Time `json:"start_at,omitempty"`
	EndAt         *time.Time `json:"end_at,omitempty"`
	Capacity      *int       `json:"capacity,omitempty"`
	BudgetMin     *int       `json:"budget_min,omitempty"`
	BudgetMax     *int       `json:"budget_max,omitempty"`
	Currency      *string    `json:"currency,omitempty"`
	Status        *string    `json:"status,omitempty" binding:"omitempty,oneof=published cancelled completed"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	CategoryIDs   []string   `json:"category_ids,omitempty"`
	TagIDs        []string   `json:"tag_ids,omitempty"`
}

// JoinEventRequest represents a join event request
type JoinEventRequest struct {
	EventID string `json:"event_id" binding:"required"`
}

// LeaveEventRequest represents a leave event request
type LeaveEventRequest struct {
	EventID string `json:"event_id" binding:"required"`
}

// SwipeEventRequest represents a swipe event request
type SwipeEventRequest struct {
	EventID   string `json:"event_id" binding:"required"`
	Direction string `json:"direction" binding:"required,oneof=like pass"`
}

// EventListResponse represents an event list response
type EventListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

// GetEventsRequest represents a get events request
type GetEventsRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	EventType string `form:"event_type"`
	Status    string `form:"status"`
	Location  string `form:"location"`
	Search    string `form:"search"`
}

// GetEventsResponse represents a get events response
type GetEventsResponse struct {
	Events     []EventResponse `json:"events"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}
