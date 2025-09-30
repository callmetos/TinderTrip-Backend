package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"TinderTrip-Backend/pkg/config"
)

// SMTPClient represents an SMTP email client
type SMTPClient struct {
	config *config.EmailConfig
}

// EmailMessage represents an email message
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	HTML    string
}

// NewSMTPClient creates a new SMTP client
func NewSMTPClient() *SMTPClient {
	return &SMTPClient{
		config: &config.AppConfig.Email,
	}
}

// SendEmail sends an email using SMTP
func (c *SMTPClient) SendEmail(message *EmailMessage) error {
	// Create authentication
	auth := smtp.PlainAuth("", c.config.SMTPUsername, c.config.SMTPPassword, c.config.SMTPHost)

	// Create the email content
	msg := c.createMessage(message)

	// Connect to the server
	addr := fmt.Sprintf("%s:%d", c.config.SMTPHost, c.config.SMTPPort)

	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.config.SMTPHost,
	}

	// Connect to the server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, c.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err := client.Mail(c.config.SMTPUsername); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, to := range message.To {
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", to, err)
		}
	}

	// Send the email
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write email data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}

// createMessage creates the email message content
func (c *SMTPClient) createMessage(message *EmailMessage) []byte {
	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", c.config.SMTPFromName, c.config.SMTPUsername))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// Content type
	if message.HTML != "" {
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	msg.WriteString("\r\n")

	// Body
	if message.HTML != "" {
		msg.WriteString(message.HTML)
	} else {
		msg.WriteString(message.Body)
	}

	return []byte(msg.String())
}

// SendPasswordResetEmail sends a password reset email
func (c *SMTPClient) SendPasswordResetEmail(to, resetToken, resetURL string) error {
	subject := "Reset Your Password - TinderTrip"

	// Create reset URL with token
	fullResetURL := fmt.Sprintf("%s?token=%s", resetURL, resetToken)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Reset Your Password</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>TinderTrip</h1>
				</div>
				<div class="content">
					<h2>Reset Your Password</h2>
					<p>Hello,</p>
					<p>We received a request to reset your password for your TinderTrip account.</p>
					<p>Click the button below to reset your password:</p>
					<a href="%s" class="button">Reset Password</a>
					<p>If the button doesn't work, you can copy and paste this link into your browser:</p>
					<p><a href="%s">%s</a></p>
					<p>This link will expire in 24 hours for security reasons.</p>
					<p>If you didn't request this password reset, please ignore this email.</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, fullResetURL, fullResetURL, fullResetURL)

	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	return c.SendEmail(message)
}

// SendWelcomeEmail sends a welcome email to new users
func (c *SMTPClient) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to TinderTrip!"

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Welcome to TinderTrip</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Welcome to TinderTrip!</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<p>Welcome to TinderTrip! We're excited to have you join our community of travelers.</p>
					<p>With TinderTrip, you can:</p>
					<ul>
						<li>Discover amazing travel events and activities</li>
						<li>Connect with like-minded travelers</li>
						<li>Create and join group trips</li>
						<li>Share your travel experiences</li>
					</ul>
					<p>Get started by exploring events in your area or creating your first event!</p>
					<a href="http://localhost:3000/events" class="button">Explore Events</a>
					<p>If you have any questions, feel free to reach out to our support team.</p>
					<p>Happy travels!</p>
					<p>The TinderTrip Team</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name)

	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	return c.SendEmail(message)
}

// SendEventConfirmationEmail sends an event confirmation email
func (c *SMTPClient) SendEventConfirmationEmail(to, name, eventTitle, eventDate string) error {
	subject := fmt.Sprintf("Event Confirmation: %s", eventTitle)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Event Confirmation</title>
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
					<h1>Event Confirmation</h1>
				</div>
				<div class="content">
					<h2>Hello %s!</h2>
					<p>You have successfully confirmed your participation in the following event:</p>
					<div class="event-info">
						<h3>%s</h3>
						<p><strong>Date:</strong> %s</p>
					</div>
					<p>We're looking forward to seeing you there!</p>
					<p>If you need to make any changes or have questions, please contact the event organizer.</p>
					<p>Happy travels!</p>
					<p>The TinderTrip Team</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, name, eventTitle, eventDate)

	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	return c.SendEmail(message)
}
