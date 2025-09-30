package service

import (
	"fmt"
	"log"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationService handles notifications
type NotificationService struct {
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendPushNotification sends a push notification
func (s *NotificationService) SendPushNotification(userID, title, body string, data map[string]interface{}) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// TODO: Implement actual push notification sending
	// This would integrate with Firebase Cloud Messaging or similar service
	log.Printf("Sending push notification to user %s: %s - %s", userID, title, body)

	// For now, just log the notification
	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userUUID,
		Title:     title,
		Body:      body,
		Type:      "push",
		Data:      data,
		Read:      false,
		CreatedAt: time.Now(),
	}

	// Save notification to database
	err = database.GetDB().Create(notification).Error
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	return nil
}

// SendEventNotification sends an event-related notification
func (s *NotificationService) SendEventNotification(eventID, userID, title, body string, data map[string]interface{}) error {
	// Add event data
	if data == nil {
		data = make(map[string]interface{})
	}
	data["event_id"] = eventID

	// Send notification
	return s.SendPushNotification(userID, title, body, data)
}

// SendChatNotification sends a chat-related notification
func (s *NotificationService) SendChatNotification(roomID, userID, title, body string, data map[string]interface{}) error {
	// Add chat data
	if data == nil {
		data = make(map[string]interface{})
	}
	data["room_id"] = roomID

	// Send notification
	return s.SendPushNotification(userID, title, body, data)
}

// SendEventReminder sends an event reminder
func (s *NotificationService) SendEventReminder(eventID string) error {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Where("id = ?", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Get event members
	var members []models.EventMember
	err = database.GetDB().Preload("User").Where("event_id = ? AND status = ?", eventUUID, models.MemberStatusConfirmed).Find(&members).Error
	if err != nil {
		return fmt.Errorf("failed to get event members: %w", err)
	}

	// Send reminder to each member
	for _, member := range members {
		if member.User != nil {
			title := "Event Reminder"
			body := fmt.Sprintf("Don't forget! %s is starting soon.", event.Title)
			data := map[string]interface{}{
				"event_id": eventID,
				"type":     "event_reminder",
			}

			err := s.SendPushNotification(member.User.ID.String(), title, body, data)
			if err != nil {
				log.Printf("Error sending reminder to user %s: %v", member.User.ID, err)
			}
		}
	}

	return nil
}

// SendEventUpdate sends an event update notification
func (s *NotificationService) SendEventUpdate(eventID, title, body string) error {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event members
	var members []models.EventMember
	err = database.GetDB().Preload("User").Where("event_id = ? AND status = ?", eventUUID, models.MemberStatusConfirmed).Find(&members).Error
	if err != nil {
		return fmt.Errorf("failed to get event members: %w", err)
	}

	// Send update to each member
	for _, member := range members {
		if member.User != nil {
			data := map[string]interface{}{
				"event_id": eventID,
				"type":     "event_update",
			}

			err := s.SendPushNotification(member.User.ID.String(), title, body, data)
			if err != nil {
				log.Printf("Error sending update to user %s: %v", member.User.ID, err)
			}
		}
	}

	return nil
}

// SendWelcomeNotification sends a welcome notification
func (s *NotificationService) SendWelcomeNotification(userID string) error {
	title := "Welcome to TinderTrip!"
	body := "Thanks for joining! Start exploring events and meeting new people."
	data := map[string]interface{}{
		"type": "welcome",
	}

	return s.SendPushNotification(userID, title, body, data)
}

// SendUserJoinedEventNotification sends notification when user joins event
func (s *NotificationService) SendUserJoinedEventNotification(eventID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Where("id = ?", eventUUID).First(&event).Error
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Get user
	var user models.User
	err = database.GetDB().Where("id = ?", userUUID).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Send notification to event creator
	if event.CreatorID != userUUID {
		title := "New Member Joined"
		body := fmt.Sprintf("%s joined your event: %s", user.GetDisplayName(), event.Title)
		data := map[string]interface{}{
			"event_id": eventID,
			"user_id":  userID,
			"type":     "user_joined",
		}

		err := s.SendPushNotification(event.CreatorID.String(), title, body, data)
		if err != nil {
			log.Printf("Error sending join notification: %v", err)
		}
	}

	return nil
}

// SendUserLeftEventNotification sends notification when user leaves event
func (s *NotificationService) SendUserLeftEventNotification(eventID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Where("id = ?", eventUUID).First(&event).Error
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Get user
	var user models.User
	err = database.GetDB().Where("id = ?", userUUID).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Send notification to event creator
	if event.CreatorID != userUUID {
		title := "Member Left"
		body := fmt.Sprintf("%s left your event: %s", user.GetDisplayName(), event.Title)
		data := map[string]interface{}{
			"event_id": eventID,
			"user_id":  userID,
			"type":     "user_left",
		}

		err := s.SendPushNotification(event.CreatorID.String(), title, body, data)
		if err != nil {
			log.Printf("Error sending leave notification: %v", err)
		}
	}

	return nil
}

// SendEventCancelledNotification sends notification when event is cancelled
func (s *NotificationService) SendEventCancelledNotification(eventID string) error {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Where("id = ?", eventUUID).First(&event).Error
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Get event members
	var members []models.EventMember
	err = database.GetDB().Preload("User").Where("event_id = ? AND status = ?", eventUUID, models.MemberStatusConfirmed).Find(&members).Error
	if err != nil {
		return fmt.Errorf("failed to get event members: %w", err)
	}

	// Send cancellation notification to each member
	for _, member := range members {
		if member.User != nil {
			title := "Event Cancelled"
			body := fmt.Sprintf("The event '%s' has been cancelled.", event.Title)
			data := map[string]interface{}{
				"event_id": eventID,
				"type":     "event_cancelled",
			}

			err := s.SendPushNotification(member.User.ID.String(), title, body, data)
			if err != nil {
				log.Printf("Error sending cancellation notification to user %s: %v", member.User.ID, err)
			}
		}
	}

	return nil
}

// SendEventCompletedNotification sends notification when event is completed
func (s *NotificationService) SendEventCompletedNotification(eventID string) error {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Where("id = ?", eventUUID).First(&event).Error
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Get event members
	var members []models.EventMember
	err = database.GetDB().Preload("User").Where("event_id = ? AND status = ?", eventUUID, models.MemberStatusConfirmed).Find(&members).Error
	if err != nil {
		return fmt.Errorf("failed to get event members: %w", err)
	}

	// Send completion notification to each member
	for _, member := range members {
		if member.User != nil {
			title := "Event Completed"
			body := fmt.Sprintf("The event '%s' has been completed. Thanks for participating!", event.Title)
			data := map[string]interface{}{
				"event_id": eventID,
				"type":     "event_completed",
			}

			err := s.SendPushNotification(member.User.ID.String(), title, body, data)
			if err != nil {
				log.Printf("Error sending completion notification to user %s: %v", member.User.ID, err)
			}
		}
	}

	return nil
}
