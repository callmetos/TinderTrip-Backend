package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatRoom represents the chat_rooms table
type ChatRoom struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID   uuid.UUID `json:"event_id" gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Event    *Event        `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	Messages []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for ChatRoom
func (ChatRoom) TableName() string {
	return "chat_rooms"
}

// BeforeCreate hook for ChatRoom
func (cr *ChatRoom) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == uuid.Nil {
		cr.ID = uuid.New()
	}
	return nil
}

// GetMessageCount returns the number of messages in the room
func (cr *ChatRoom) GetMessageCount() int64 {
	// This would typically be done with a count query
	// For now, return the length of the messages slice
	return int64(len(cr.Messages))
}

// GetLastMessage returns the last message in the room
func (cr *ChatRoom) GetLastMessage() *ChatMessage {
	if len(cr.Messages) == 0 {
		return nil
	}

	// Find the message with the latest created_at
	var lastMessage *ChatMessage
	for i := range cr.Messages {
		if lastMessage == nil || cr.Messages[i].CreatedAt.After(lastMessage.CreatedAt) {
			lastMessage = &cr.Messages[i]
		}
	}

	return lastMessage
}

// IsActive checks if the chat room is active (has recent messages)
func (cr *ChatRoom) IsActive() bool {
	lastMessage := cr.GetLastMessage()
	if lastMessage == nil {
		return false
	}

	// Consider active if there's a message within the last 7 days
	return time.Since(lastMessage.CreatedAt) < 7*24*time.Hour
}
