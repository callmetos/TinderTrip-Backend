package dto

// CreateTagRequest represents a create tag request
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
	Kind string `json:"kind" binding:"required,min=1,max=20"`
}

// UpdateTagRequest represents an update tag request
type UpdateTagRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Kind *string `json:"kind,omitempty" binding:"omitempty,min=1,max=20"`
}

// GetTagsRequest represents a get tags request
type GetTagsRequest struct {
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	Kind   string `form:"kind"`
	Search string `form:"search"`
}

// GetTagsResponse represents a get tags response
type GetTagsResponse struct {
	Tags       []TagResponse `json:"tags"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// AddUserTagRequest represents an add user tag request
type AddUserTagRequest struct {
	TagID string `json:"tag_id" binding:"required"`
}

// RemoveUserTagRequest represents a remove user tag request
type RemoveUserTagRequest struct {
	TagID string `json:"tag_id" binding:"required"`
}

// GetUserTagsRequest represents a get user tags request
type GetUserTagsRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Kind   string `form:"kind"`
}

// GetUserTagsResponse represents a get user tags response
type GetUserTagsResponse struct {
	Tags []TagResponse `json:"tags"`
}
