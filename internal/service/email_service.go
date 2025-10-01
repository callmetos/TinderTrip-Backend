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

// SendPasswordResetOTP sends a password reset OTP email
func (s *EmailService) SendPasswordResetOTP(to, otp string) error {
	return s.smtpClient.SendPasswordResetOTP(to, otp)
}

// SendEventConfirmationEmail sends an event confirmation email
func (s *EmailService) SendEventConfirmationEmail(to, name, eventTitle, eventDate string) error {
	return s.smtpClient.SendEventConfirmationEmail(to, name, eventTitle, eventDate)
}
