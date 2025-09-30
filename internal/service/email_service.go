package service

import (
	"TinderTrip-Backend/pkg/email"
)

// EmailService handles email operations
type EmailService struct {
	smtpClient *email.SMTPClient
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		smtpClient: email.NewSMTPClient(),
	}
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	return s.smtpClient.SendWelcomeEmail(to, name)
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, token, resetURL string) error {
	return s.smtpClient.SendPasswordResetEmail(to, token, resetURL)
}

// SendEventConfirmationEmail sends an event confirmation email
func (s *EmailService) SendEventConfirmationEmail(to, name, eventTitle, eventDate string) error {
	return s.smtpClient.SendEventConfirmationEmail(to, name, eventTitle, eventDate)
}
