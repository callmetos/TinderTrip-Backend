package service

import (
	"fmt"
	"log"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"
	"TinderTrip-Backend/pkg/email"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationService handles notifications
type NotificationService struct {
	emailService *EmailService
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		emailService: NewEmailService(),
	}
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

	// Save notification to database
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

	err = database.GetDB().Create(notification).Error
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Send email notification in background (don't block on error)
	go func() {
		// Get user email
		var user models.User
		err := database.GetDB().Where("id = ?", userUUID).First(&user).Error
		if err != nil {
			log.Printf("Failed to get user for email notification: %v", err)
			return
		}

		// Only send email if user has email address
		if user.Email == nil || *user.Email == "" {
			log.Printf("User %s has no email address, skipping email notification", userID)
			return
		}

		// Send email notification
		err = s.sendNotificationEmail(*user.Email, user.GetDisplayName(), title, body, data)
		if err != nil {
			log.Printf("Failed to send email notification to %s: %v", *user.Email, err)
		}
	}()

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
			"event_id":    eventID,
			"user_id":     userID,
			"type":        "user_joined",
			"event_title": event.Title,
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
			"event_id":    eventID,
			"user_id":     userID,
			"type":        "user_left",
			"event_title": event.Title,
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

// sendNotificationEmail sends an email notification
func (s *NotificationService) sendNotificationEmail(to, name, title, body string, data map[string]interface{}) error {
	subject := fmt.Sprintf("%s - TinderTrip", title)

	// Determine notification type from data if available
	notificationType := ""
	if data != nil {
		if nType, ok := data["type"].(string); ok {
			notificationType = nType
		}
	}

	// Create email body based on notification type
	var htmlBody string
	switch notificationType {
	case "user_joined", "user_left":
		htmlBody = s.createEventMemberChangeEmailHTML(name, title, body, data)
	case "event_reminder":
		htmlBody = s.createEventReminderEmailHTML(name, title, body, data)
	case "event_update":
		htmlBody = s.createEventUpdateEmailHTML(name, title, body, data)
	case "event_cancelled":
		htmlBody = s.createEventCancelledEmailHTML(name, title, body, data)
	case "event_completed":
		htmlBody = s.createEventCompletedEmailHTML(name, title, body, data)
	default:
		// Generic notification email
		htmlBody = s.createGenericNotificationEmailHTML(name, title, body, data)
	}

	// Send email using SMTP client directly
	message := &email.EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	smtpClient := email.NewSMTPClient()
	return smtpClient.SendEmail(message)
}

// createGenericNotificationEmailHTML creates HTML for generic notifications
func (s *NotificationService) createGenericNotificationEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.message { background-color: white; padding: 15px; border-radius: 4px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>TinderTrip</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<div class="message">
						<h3>%s</h3>
						<p>%s</p>
					</div>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, title, name, title, body)
}

// createEventMemberChangeEmailHTML creates HTML for event member join/leave notifications
func (s *NotificationService) createEventMemberChangeEmailHTML(name, title, body string, data map[string]interface{}) string {
	eventTitle := "an event"
	if data != nil {
		if eTitle, ok := data["event_title"].(string); ok && eTitle != "" {
			eventTitle = eTitle
		}
	}

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.event-info { background-color: white; padding: 15px; border-radius: 4px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>TinderTrip Notification</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<p>%s</p>
					<div class="event-info">
						<h3>%s</h3>
					</div>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, title, name, body, eventTitle)
}

// createEventReminderEmailHTML creates HTML for event reminder notifications
func (s *NotificationService) createEventReminderEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Event Reminder</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #FF9800; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.reminder-box { background-color: #fff3cd; border: 2px solid #FF9800; padding: 15px; border-radius: 4px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>⏰ Event Reminder</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<div class="reminder-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<p>Don't forget to prepare for your upcoming event!</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, title, body)
}

// createEventUpdateEmailHTML creates HTML for event update notifications
func (s *NotificationService) createEventUpdateEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Event Update</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #2196F3; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.update-box { background-color: white; border-left: 4px solid #2196F3; padding: 15px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>📢 Event Update</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<div class="update-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, title, body)
}

// createEventCancelledEmailHTML creates HTML for event cancellation notifications
func (s *NotificationService) createEventCancelledEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Event Cancelled</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #f44336; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.cancelled-box { background-color: #ffebee; border: 2px solid #f44336; padding: 15px; border-radius: 4px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>❌ Event Cancelled</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<div class="cancelled-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<p>We're sorry for any inconvenience. Please check for other events you might be interested in!</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, title, body)
}

// createEventCompletedEmailHTML creates HTML for event completion notifications
func (s *NotificationService) createEventCompletedEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Event Completed</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.completed-box { background-color: #e8f5e9; border: 2px solid #4CAF50; padding: 15px; border-radius: 4px; margin: 15px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>✅ Event Completed</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<div class="completed-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<p>Thanks for participating! We hope you had a great time.</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, title, body)
}
