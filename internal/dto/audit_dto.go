package dto

// AuditLogResponse represents an audit log response
type AuditLogResponse struct {
	ID          string  `json:"id"`
	ActorUserID *string `json:"actor_user_id,omitempty"`
	EntityTable string  `json:"entity_table"`
	EntityID    *string `json:"entity_id,omitempty"`
	Action      string  `json:"action"`
	BeforeData  *string `json:"before_data,omitempty"`
	AfterData   *string `json:"after_data,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// AuditLogListResponse represents a paginated audit log list response
type AuditLogListResponse struct {
	Logs       []AuditLogResponse `json:"logs"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
