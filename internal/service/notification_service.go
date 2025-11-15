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
		log.Printf("Attempting to send email notification to %s for user %s: %s - %s", *user.Email, userID, title, body)
		err = s.sendNotificationEmail(*user.Email, user.GetDisplayName(), title, body, data)
		if err != nil {
			log.Printf("Failed to send email notification to %s: %v", *user.Email, err)
		} else {
			log.Printf("Successfully sent email notification to %s for user %s", *user.Email, userID)
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

	// Send notification to event creator only
	if event.CreatorID != userUUID {
		title := "New Member Joined"
		body := fmt.Sprintf("%s joined your event: %s", user.GetDisplayName(), event.Title)
		data := map[string]interface{}{
			"event_id":    eventID,
			"user_id":     userID,
			"type":        "user_joined",
			"event_title": event.Title,
		}

		// Send push notification (includes email notification)
		err := s.SendPushNotification(event.CreatorID.String(), title, body, data)
		if err != nil {
			log.Printf("Error sending join notification to creator: %v", err)
			return err
		}
		log.Printf("Successfully sent join notification (push + email) to creator %s for event %s", event.CreatorID.String(), eventID)
	} else {
		log.Printf("User %s is the creator of event %s, skipping notification", userID, eventID)
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

		// Send push notification (includes email notification)
		err := s.SendPushNotification(event.CreatorID.String(), title, body, data)
		if err != nil {
			log.Printf("Error sending leave notification: %v", err)
		} else {
			log.Printf("Successfully sent leave notification (push + email) to creator %s for event %s", event.CreatorID.String(), eventID)
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
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>%s</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.message {
					background: linear-gradient(135deg, #f5f7fa 0%%, #e2e8f0 100%%);
					border: 2px solid #667eea;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
				}
				.message h3 {
					font-size: 22px;
					font-weight: 700;
					color: #1a202c;
					margin-bottom: 15px;
					text-align: center;
				}
				.message p {
					font-size: 16px;
					color: #4a5568;
					line-height: 1.7;
					margin: 0;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.message { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">📬</div>
					<h1>TinderTrip</h1>
					<p style="margin: 0; opacity: 0.9;">Notification</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<div class="message">
						<h3>%s</h3>
						<p>%s</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
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

	icon := "👥"
	if title == "New Member Joined" {
		icon = "🎉"
	} else if title == "Member Left" {
		icon = "👋"
	}

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>%s</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #10b981 0%%, #059669 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.event-info {
					background: linear-gradient(135deg, #f0fdf4 0%%, #dcfce7 100%%);
					border: 2px solid #10b981;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(16, 185, 129, 0.15);
					text-align: center;
				}
				.event-info h3 {
					font-size: 22px;
					font-weight: 700;
					color: #065f46;
					margin: 0;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.event-info { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">%s</div>
					<h1>TinderTrip</h1>
					<p style="margin: 0; opacity: 0.9;">%s</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<p>%s</p>
					<div class="event-info">
						<h3>%s</h3>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, title, icon, title, name, body, eventTitle)
}

// createEventReminderEmailHTML creates HTML for event reminder notifications
func (s *NotificationService) createEventReminderEmailHTML(name, title, body string, data map[string]interface{}) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Event Reminder</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.reminder-box {
					background: linear-gradient(135deg, #fffbeb 0%%, #fef3c7 100%%);
					border: 3px solid #f59e0b;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(245, 158, 11, 0.2);
					text-align: center;
				}
				.reminder-box h3 {
					font-size: 22px;
					font-weight: 700;
					color: #92400e;
					margin-bottom: 15px;
				}
				.reminder-box p {
					font-size: 16px;
					color: #78350f;
					margin: 0;
					line-height: 1.7;
				}
				.reminder-note {
					background: linear-gradient(135deg, #fef3c7 0%%, #fde68a 100%%);
					border-left: 4px solid #f59e0b;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					text-align: center;
				}
				.reminder-note p {
					margin: 0;
					color: #92400e;
					font-weight: 600;
					font-size: 16px;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.reminder-box { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">⏰</div>
					<h1>Event Reminder</h1>
					<p style="margin: 0; opacity: 0.9;">Don't miss out!</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<div class="reminder-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<div class="reminder-note">
						<p>📋 Don't forget to prepare for your upcoming event!</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
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
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Event Update</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #3b82f6 0%%, #2563eb 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.update-box {
					background: linear-gradient(135deg, #eff6ff 0%%, #dbeafe 100%%);
					border: 3px solid #3b82f6;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
				}
				.update-box h3 {
					font-size: 22px;
					font-weight: 700;
					color: #1e40af;
					margin-bottom: 15px;
					text-align: center;
				}
				.update-box p {
					font-size: 16px;
					color: #1e3a8a;
					margin: 0;
					line-height: 1.7;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.update-box { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">📢</div>
					<h1>Event Update</h1>
					<p style="margin: 0; opacity: 0.9;">Important information</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<div class="update-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
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
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Event Cancelled</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #ef4444 0%%, #dc2626 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.cancelled-box {
					background: linear-gradient(135deg, #fef2f2 0%%, #fee2e2 100%%);
					border: 3px solid #ef4444;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(239, 68, 68, 0.15);
					text-align: center;
				}
				.cancelled-box h3 {
					font-size: 22px;
					font-weight: 700;
					color: #991b1b;
					margin-bottom: 15px;
				}
				.cancelled-box p {
					font-size: 16px;
					color: #7f1d1d;
					margin: 0;
					line-height: 1.7;
				}
				.sorry-message {
					background: linear-gradient(135deg, #fef2f2 0%%, #fee2e2 100%%);
					border-left: 4px solid #ef4444;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					text-align: center;
				}
				.sorry-message p {
					margin: 0;
					color: #991b1b;
					font-weight: 600;
					font-size: 16px;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.cancelled-box { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">❌</div>
					<h1>Event Cancelled</h1>
					<p style="margin: 0; opacity: 0.9;">We're sorry</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<div class="cancelled-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<div class="sorry-message">
						<p>😔 We're sorry for any inconvenience. Please check for other events you might be interested in!</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
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
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Event Completed</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
					line-height: 1.6; 
					color: #333333;
					background-color: #f5f7fa;
					padding: 20px;
				}
				.email-wrapper {
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					border-radius: 12px;
					overflow: hidden;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				.header {
					background: linear-gradient(135deg, #10b981 0%%, #059669 100%%);
					color: white;
					padding: 40px 30px;
					text-align: center;
				}
				.header h1 {
					font-size: 32px;
					font-weight: 700;
					margin-bottom: 10px;
					letter-spacing: -0.5px;
				}
				.header-icon {
					font-size: 48px;
					margin-bottom: 15px;
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 24px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.completed-box {
					background: linear-gradient(135deg, #f0fdf4 0%%, #dcfce7 100%%);
					border: 3px solid #10b981;
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
					box-shadow: 0 4px 12px rgba(16, 185, 129, 0.15);
					text-align: center;
				}
				.completed-box h3 {
					font-size: 22px;
					font-weight: 700;
					color: #065f46;
					margin-bottom: 15px;
				}
				.completed-box p {
					font-size: 16px;
					color: #047857;
					margin: 0;
					line-height: 1.7;
				}
				.thanks-message {
					background: linear-gradient(135deg, #ecfdf5 0%%, #d1fae5 100%%);
					border-left: 4px solid #10b981;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					text-align: center;
				}
				.thanks-message p {
					margin: 0;
					color: #065f46;
					font-weight: 600;
					font-size: 16px;
				}
				.footer {
					text-align: center;
					padding: 30px;
					background-color: #f7fafc;
					border-top: 1px solid #e2e8f0;
				}
				.footer p {
					font-size: 13px;
					color: #718096;
					margin: 5px 0;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.completed-box { padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">✅</div>
					<h1>Event Completed</h1>
					<p style="margin: 0; opacity: 0.9;">Great job!</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<div class="completed-box">
						<h3>%s</h3>
						<p>%s</p>
					</div>
					<div class="thanks-message">
						<p>🎉 Thanks for participating! We hope you had a great time.</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, title, body)
}
