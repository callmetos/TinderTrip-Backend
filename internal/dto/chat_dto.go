package dto

import "time"

// ChatRoomResponse represents a chat room response
type ChatRoomResponse struct {
	ID        string         `json:"id"`
	EventID   string         `json:"event_id"`
	Event     *EventResponse `json:"event,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// ChatMessageResponse represents a chat message response
type ChatMessageResponse struct {
	ID          string        `json:"id"`
	RoomID      string        `json:"room_id"`
	SenderID    string        `json:"sender_id"`
	Sender      *UserResponse `json:"sender,omitempty"`
	Body        *string       `json:"body,omitempty"`
	MessageType string        `json:"message_type"`
	ImageURL    *string       `json:"image_url,omitempty"`
	FileURL     *string       `json:"file_url,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

// SendMessageRequest represents a send message request (JSON)
type SendMessageRequest struct {
	RoomID      string `json:"room_id" binding:"required"`
	Body        string `json:"body"`
	MessageType string `json:"message_type" binding:"required"`
}

// SendMessageMultipartRequest represents a send message request (multipart form data)
type SendMessageMultipartRequest struct {
	RoomID      string `form:"room_id" binding:"required"`
	Body        string `form:"body"`
	MessageType string `form:"message_type" binding:"required"`
	File        string `form:"file"` // Will be handled as multipart.FileHeader
}

// GetMessagesRequest represents a get messages request
type GetMessagesRequest struct {
	RoomID string `form:"room_id" binding:"required"`
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
}

// GetMessagesResponse represents a get messages response
type GetMessagesResponse struct {
	Messages   []ChatMessageResponse `json:"messages"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// GetChatRoomsRequest represents a get chat rooms request
type GetChatRoomsRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// GetChatRoomsResponse represents a get chat rooms response
type GetChatRoomsResponse struct {
	Rooms      []ChatRoomResponse `json:"rooms"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ChatRoomListResponse represents a chat room list response
type ChatRoomListResponse struct {
	Rooms []ChatRoomResponse `json:"rooms"`
}

// ChatMessageListResponse represents a chat message list response
type ChatMessageListResponse struct {
	Messages []ChatMessageResponse `json:"messages"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	Limit    int                   `json:"limit"`
}
