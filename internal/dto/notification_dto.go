package dto

import "time"

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"created_at"`
	ReadAt    *time.Time             `json:"read_at"`
}

// SendNotificationRequest represents a send notification request
type SendNotificationRequest struct {
	UserID string                 `json:"user_id" binding:"required"`
	Title  string                 `json:"title" binding:"required"`
	Body   string                 `json:"body" binding:"required"`
	Type   string                 `json:"type" binding:"required"`
	Data   map[string]interface{} `json:"data"`
}

// MarkNotificationReadRequest represents a mark notification read request
type MarkNotificationReadRequest struct {
	NotificationID string `json:"notification_id" binding:"required"`
}

// GetNotificationsRequest represents a get notifications request
type GetNotificationsRequest struct {
	Page  int   `form:"page" binding:"min=1"`
	Limit int   `form:"limit" binding:"min=1,max=100"`
	Read  *bool `form:"read"`
}

// GetNotificationsResponse represents a get notifications response
type GetNotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}
