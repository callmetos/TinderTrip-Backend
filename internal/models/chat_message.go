package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageType represents different types of chat messages
type MessageType string

const (
	MessageTypeText    MessageType = "text"
	MessageTypeImage   MessageType = "image"
	MessageTypeFile    MessageType = "file"
	MessageTypeSystem  MessageType = "system"
	MessageTypeJoin    MessageType = "join"
	MessageTypeLeave   MessageType = "leave"
	MessageTypeConfirm MessageType = "confirm"
)

// ChatMessage represents the chat_messages table
type ChatMessage struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID      uuid.UUID `json:"room_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	SenderID    uuid.UUID `json:"sender_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Body        *string   `json:"body" gorm:"type:text"`
	MessageType *string   `json:"message_type" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	Room   *ChatRoom `json:"room,omitempty" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
	Sender *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	// Used as confirmation message for event members
	EventMembers []EventMember `json:"event_members,omitempty" gorm:"foreignKey:ConfirmationMessageID;constraint:OnDelete:SET NULL"`
}

// TableName returns the table name for ChatMessage
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// BeforeCreate hook for ChatMessage
func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	return nil
}

// IsText checks if the message is a text message
func (cm *ChatMessage) IsText() bool {
	return cm.MessageType == nil || *cm.MessageType == string(MessageTypeText)
}

// IsImage checks if the message is an image message
func (cm *ChatMessage) IsImage() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeImage)
}

// IsFile checks if the message is a file message
func (cm *ChatMessage) IsFile() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeFile)
}

// IsSystem checks if the message is a system message
func (cm *ChatMessage) IsSystem() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeSystem)
}

// IsJoin checks if the message is a join message
func (cm *ChatMessage) IsJoin() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeJoin)
}

// IsLeave checks if the message is a leave message
func (cm *ChatMessage) IsLeave() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeLeave)
}

// IsConfirm checks if the message is a confirmation message
func (cm *ChatMessage) IsConfirm() bool {
	return cm.MessageType != nil && *cm.MessageType == string(MessageTypeConfirm)
}

// GetAge returns how long ago the message was sent
func (cm *ChatMessage) GetAge() time.Duration {
	return time.Since(cm.CreatedAt)
}

// IsRecent checks if the message was sent recently (within 1 hour)
func (cm *ChatMessage) IsRecent() bool {
	return cm.GetAge() < time.Hour
}

// GetDisplayBody returns the body text or a default message
func (cm *ChatMessage) GetDisplayBody() string {
	if cm.Body != nil && *cm.Body != "" {
		return *cm.Body
	}

	// Return default messages based on type
	switch {
	case cm.IsJoin():
		return "joined the event"
	case cm.IsLeave():
		return "left the event"
	case cm.IsConfirm():
		return "confirmed participation"
	case cm.IsSystem():
		return "system message"
	default:
		return "message"
	}
}
