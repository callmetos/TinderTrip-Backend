package dto

// APIResponse represents the standardized API response wrapper
type APIResponse struct {
	Success   bool        `json:"success" example:"true"`
	RequestID string      `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Code      string      `json:"code,omitempty" example:"SUCCESS"`
	Message   string      `json:"message" example:"Operation completed successfully"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Meta      *MetaData   `json:"meta,omitempty"`
}

// MetaData represents pagination and additional metadata
type MetaData struct {
	Page       *int   `json:"page,omitempty" example:"1"`
	Limit      *int   `json:"limit,omitempty" example:"10"`
	Total      *int64 `json:"total,omitempty" example:"100"`
	TotalPages *int   `json:"total_pages,omitempty" example:"10"`
}

// SuccessAPIResponse represents a successful API response
type SuccessAPIResponse struct {
	Success   bool        `json:"success" example:"true"`
	RequestID string      `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string      `json:"message" example:"Operation completed successfully"`
	Data      interface{} `json:"data,omitempty"`
	Meta      *MetaData   `json:"meta,omitempty"`
}

// ErrorAPIResponse represents an error API response
type ErrorAPIResponse struct {
	Success   bool        `json:"success" example:"false"`
	RequestID string      `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Code      string      `json:"code" example:"VALIDATION_ERROR"`
	Message   string      `json:"message" example:"Validation failed"`
	Errors    interface{} `json:"errors,omitempty"`
}

// AuthResponseWrapper wraps authentication response with token and user at top level
type AuthResponseWrapper struct {
	Success   bool         `json:"success" example:"true"`
	RequestID string       `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string       `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string       `json:"message" example:"Login successful"`
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User      UserResponse `json:"user"`
}

// TokenResponseWrapper wraps token-only response (for refresh token)
type TokenResponseWrapper struct {
	Success   bool   `json:"success" example:"true"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string `json:"message" example:"Token refreshed successfully"`
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserProfileResponseWrapper wraps UserProfileResponse in APIResponse format
type UserProfileResponseWrapper struct {
	Success   bool                `json:"success" example:"true"`
	RequestID string              `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string              `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string              `json:"message" example:"Profile retrieved successfully"`
	Data      UserProfileResponse `json:"data"`
}

// EventResponseWrapper wraps EventResponse in APIResponse format
type EventResponseWrapper struct {
	Success   bool          `json:"success" example:"true"`
	RequestID string        `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string        `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string        `json:"message" example:"Event retrieved successfully"`
	Data      EventResponse `json:"data"`
}

// EventListResponseWrapper wraps event list in APIResponse format
type EventListResponseWrapper struct {
	Success   bool            `json:"success" example:"true"`
	RequestID string          `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string          `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string          `json:"message" example:"Events retrieved successfully"`
	Data      []EventResponse `json:"data"`
	Meta      *MetaData       `json:"meta,omitempty"`
}

// GoogleAuthResponseWrapper wraps Google OAuth response with auth_url and state at top level
type GoogleAuthResponseWrapper struct {
	Success   bool   `json:"success" example:"true"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string `json:"message" example:"OAuth URL generated successfully"`
	AuthURL   string `json:"auth_url" example:"https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=xxx"`
	State     string `json:"state" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// SetupStatusResponseWrapper wraps SetupStatusResponse in APIResponse format
type SetupStatusResponseWrapper struct {
	Success   bool                `json:"success" example:"true"`
	RequestID string              `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string              `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string              `json:"message" example:"Setup status retrieved successfully"`
	Data      SetupStatusResponse `json:"data"`
}

// SuccessMessageWrapper wraps simple success message in APIResponse format
type SuccessMessageWrapper struct {
	Success   bool   `json:"success" example:"true"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string `json:"message" example:"Operation completed successfully"`
}

// FoodPreferenceListResponseWrapper wraps food preferences in APIResponse format
type FoodPreferenceListResponseWrapper struct {
	Success   bool                     `json:"success" example:"true"`
	RequestID string                   `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string                   `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string                   `json:"message" example:"Food preferences retrieved successfully"`
	Data      []FoodPreferenceResponse `json:"data"`
}

// TravelPreferenceListResponseWrapper wraps travel preferences in APIResponse format
type TravelPreferenceListResponseWrapper struct {
	Success   bool                       `json:"success" example:"true"`
	RequestID string                     `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string                     `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string                     `json:"message" example:"Travel preferences retrieved successfully"`
	Data      []TravelPreferenceResponse `json:"data"`
}

// TagListResponseWrapper wraps tags in APIResponse format
type TagListResponseWrapper struct {
	Success   bool          `json:"success" example:"true"`
	RequestID string        `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string        `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string        `json:"message" example:"Tags retrieved successfully"`
	Data      []TagResponse `json:"data"`
}

// ChatRoomListResponseWrapper wraps chat rooms in APIResponse format
type ChatRoomListResponseWrapper struct {
	Success   bool               `json:"success" example:"true"`
	RequestID string             `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string             `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string             `json:"message" example:"Chat rooms retrieved successfully"`
	Data      []ChatRoomResponse `json:"data"`
}

// ChatMessageListResponseWrapper wraps chat messages in APIResponse format
type ChatMessageListResponseWrapper struct {
	Success   bool                  `json:"success" example:"true"`
	RequestID string                `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string                `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string                `json:"message" example:"Messages retrieved successfully"`
	Data      []ChatMessageResponse `json:"data"`
}

// HistoryListResponseWrapper wraps event history in APIResponse format with pagination
type HistoryListResponseWrapper struct {
	Success   bool        `json:"success" example:"true"`
	RequestID string      `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string      `json:"message" example:"History retrieved successfully"`
	Data      interface{} `json:"data"`
	Meta      *MetaData   `json:"meta,omitempty"`
}

// AuditLogListResponseWrapper wraps audit logs in APIResponse format
type AuditLogListResponseWrapper struct {
	Success   bool               `json:"success" example:"true"`
	RequestID string             `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string             `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string             `json:"message" example:"Audit logs retrieved successfully"`
	Data      []AuditLogResponse `json:"data"`
	Meta      *MetaData          `json:"meta,omitempty"`
}

// EventSuggestionResponseWrapper wraps event suggestions in APIResponse format
type EventSuggestionResponseWrapper struct {
	Success   bool                    `json:"success" example:"true"`
	RequestID string                  `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp string                  `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Message   string                  `json:"message" example:"Event suggestions retrieved successfully"`
	Data      EventSuggestionResponse `json:"data"`
}
