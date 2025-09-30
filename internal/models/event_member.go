package models

import (
	"time"

	"github.com/google/uuid"
)

// MemberRole represents the member role enum
type MemberRole string

const (
	MemberRoleCreator     MemberRole = "creator"
	MemberRoleParticipant MemberRole = "participant"
)

// MemberStatus represents the member status enum
type MemberStatus string

const (
	MemberStatusPending   MemberStatus = "pending"
	MemberStatusConfirmed MemberStatus = "confirmed"
	MemberStatusDeclined  MemberStatus = "declined"
	MemberStatusKicked    MemberStatus = "kicked"
	MemberStatusLeft      MemberStatus = "left"
)

// EventMember represents the event_members table
type EventMember struct {
	EventID               uuid.UUID    `json:"event_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	UserID                uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;primaryKey;constraint:OnDelete:CASCADE"`
	Role                  MemberRole   `json:"role" gorm:"type:member_role;not null;default:'participant'"`
	Status                MemberStatus `json:"status" gorm:"type:member_status;not null;default:'pending'"`
	JoinedAt              time.Time    `json:"joined_at" gorm:"type:timestamptz;not null;default:now()"`
	ConfirmedAt           *time.Time   `json:"confirmed_at" gorm:"type:timestamptz"`
	LeftAt                *time.Time   `json:"left_at" gorm:"type:timestamptz"`
	Note                  *string      `json:"note" gorm:"type:text"`
	ConfirmationMessageID *uuid.UUID   `json:"confirmation_message_id" gorm:"type:uuid;constraint:OnDelete:SET NULL"`

	// Relationships
	Event               *Event       `json:"event,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	User                *User        `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ConfirmationMessage *ChatMessage `json:"confirmation_message,omitempty" gorm:"foreignKey:ConfirmationMessageID;constraint:OnDelete:SET NULL"`
}

// TableName returns the table name for EventMember
func (EventMember) TableName() string {
	return "event_members"
}

// IsCreator checks if the member is the creator
func (em *EventMember) IsCreator() bool {
	return em.Role == MemberRoleCreator
}

// IsParticipant checks if the member is a participant
func (em *EventMember) IsParticipant() bool {
	return em.Role == MemberRoleParticipant
}

// IsPending checks if the member status is pending
func (em *EventMember) IsPending() bool {
	return em.Status == MemberStatusPending
}

// IsConfirmed checks if the member status is confirmed
func (em *EventMember) IsConfirmed() bool {
	return em.Status == MemberStatusConfirmed
}

// IsDeclined checks if the member status is declined
func (em *EventMember) IsDeclined() bool {
	return em.Status == MemberStatusDeclined
}

// IsKicked checks if the member was kicked
func (em *EventMember) IsKicked() bool {
	return em.Status == MemberStatusKicked
}

// IsLeft checks if the member left
func (em *EventMember) IsLeft() bool {
	return em.Status == MemberStatusLeft
}

// IsActive checks if the member is active (confirmed and not left)
func (em *EventMember) IsActive() bool {
	return em.IsConfirmed() && !em.IsLeft() && !em.IsKicked()
}

// CanJoin checks if the member can join the event
func (em *EventMember) CanJoin() bool {
	return em.IsPending() || em.IsConfirmed()
}

// CanLeave checks if the member can leave the event
func (em *EventMember) CanLeave() bool {
	return em.IsConfirmed() && !em.IsLeft()
}

// GetDurationInEvent returns how long the member has been in the event
func (em *EventMember) GetDurationInEvent() time.Duration {
	endTime := time.Now()
	if em.LeftAt != nil {
		endTime = *em.LeftAt
	}
	return endTime.Sub(em.JoinedAt)
}
