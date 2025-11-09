package email

import (
	"TinderTrip-Backend/pkg/config"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
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

	// Connect to SMTP server (plain connection first)
	conn, err := net.Dial("tcp", addr)
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

	// Start TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.config.SMTPHost,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

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

// SendPasswordResetOTP sends a password reset OTP email
func (c *SMTPClient) SendPasswordResetOTP(to, otp string) error {
	subject := "Password Reset OTP - TinderTrip"

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Password Reset OTP</title>
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
				.otp-container {
					text-align: center;
					margin: 30px 0;
				}
				.otp-label {
					font-size: 14px;
					color: #718096;
					text-transform: uppercase;
					letter-spacing: 1px;
					margin-bottom: 10px;
					font-weight: 600;
				}
				.otp-code {
					font-size: 36px;
					font-weight: 700;
					color: #667eea;
					text-align: center;
					padding: 25px;
					background: linear-gradient(135deg, #f5f7fa 0%%, #c3cfe2 100%%);
					border: 3px solid #667eea;
					border-radius: 12px;
					margin: 15px auto;
					letter-spacing: 8px;
					display: inline-block;
					min-width: 280px;
					box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
				}
				.warning {
					background: linear-gradient(135deg, #fff5e6 0%%, #ffe0b2 100%%);
					border-left: 4px solid #ff9800;
					color: #e65100;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					box-shadow: 0 2px 8px rgba(255, 152, 0, 0.1);
				}
				.warning strong {
					display: block;
					margin-bottom: 8px;
					font-size: 16px;
				}
				.warning p {
					margin: 0;
					font-size: 14px;
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
				.footer a {
					color: #667eea;
					text-decoration: none;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.otp-code { font-size: 28px; letter-spacing: 4px; min-width: 240px; padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">🔐</div>
					<h1>TinderTrip</h1>
					<p style="margin: 0; opacity: 0.9;">Password Reset</p>
				</div>
				<div class="content">
					<h2>Password Reset Verification</h2>
					<p>Hello,</p>
					<p>We received a request to reset your password for your TinderTrip account. Please use the verification code below to complete the process.</p>
					
					<div class="otp-container">
						<div class="otp-label">Your Verification Code</div>
						<div class="otp-code">%s</div>
					</div>

					<div class="warning">
						<strong>⏱️ Important Security Notice</strong>
						<p>This verification code will expire in <strong>3 minutes</strong> for security reasons. Please use it promptly.</p>
					</div>

					<p style="margin-top: 25px; color: #718096; font-size: 14px;">
						<strong>Didn't request this?</strong> If you didn't request a password reset, please ignore this email. Your account remains secure.
					</p>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
					<p style="margin-top: 10px; font-size: 12px;">This is an automated message, please do not reply.</p>
				</div>
			</div>
		</body>
		</html>
	`, otp)

	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	return c.SendEmail(message)
}

// SendWelcomeEmail sends a welcome email to new users
func (c *SMTPClient) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to TinderTrip! 🎉"

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Welcome to TinderTrip</title>
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
					padding: 50px 30px;
					text-align: center;
					position: relative;
					overflow: hidden;
				}
				.header::before {
					content: '';
					position: absolute;
					top: -50%%;
					right: -50%%;
					width: 200%%;
					height: 200%%;
					background: radial-gradient(circle, rgba(255,255,255,0.1) 0%%, transparent 70%%);
					animation: pulse 3s ease-in-out infinite;
				}
				@keyframes pulse {
					0%%, 100%% { transform: scale(1); opacity: 0.5; }
					50%% { transform: scale(1.1); opacity: 0.8; }
				}
				.header h1 {
					font-size: 36px;
					font-weight: 700;
					margin-bottom: 15px;
					letter-spacing: -0.5px;
					position: relative;
					z-index: 1;
				}
				.header-icon {
					font-size: 64px;
					margin-bottom: 20px;
					position: relative;
					z-index: 1;
					display: inline-block;
					animation: bounce 2s ease-in-out infinite;
				}
				@keyframes bounce {
					0%%, 100%% { transform: translateY(0); }
					50%% { transform: translateY(-10px); }
				}
				.content {
					padding: 40px 30px;
					background-color: #ffffff;
				}
				.content h2 {
					font-size: 28px;
					font-weight: 600;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.welcome-name {
					color: #667eea;
					font-weight: 700;
				}
				.content p {
					font-size: 16px;
					color: #4a5568;
					margin-bottom: 15px;
					line-height: 1.7;
				}
				.features {
					background: linear-gradient(135deg, #f5f7fa 0%%, #e2e8f0 100%%);
					border-radius: 12px;
					padding: 30px;
					margin: 30px 0;
				}
				.features h3 {
					font-size: 20px;
					color: #1a202c;
					margin-bottom: 20px;
					text-align: center;
				}
				.features ul {
					list-style: none;
					padding: 0;
				}
				.features li {
					padding: 12px 0;
					padding-left: 35px;
					position: relative;
					font-size: 15px;
					color: #4a5568;
				}
				.features li::before {
					content: '✓';
					position: absolute;
					left: 0;
					width: 24px;
					height: 24px;
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					color: white;
					border-radius: 50%%;
					display: flex;
					align-items: center;
					justify-content: center;
					font-weight: bold;
					font-size: 14px;
				}
				.button-container {
					text-align: center;
					margin: 35px 0;
				}
				.button {
					display: inline-block;
					padding: 16px 40px;
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					color: white;
					text-decoration: none;
					border-radius: 8px;
					font-weight: 600;
					font-size: 16px;
					box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
					transition: transform 0.2s, box-shadow 0.2s;
				}
				.button:hover {
					transform: translateY(-2px);
					box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
				}
				.signature {
					margin-top: 30px;
					padding-top: 25px;
					border-top: 2px solid #e2e8f0;
					text-align: center;
				}
				.signature p {
					margin: 8px 0;
					color: #4a5568;
				}
				.signature .team-name {
					font-weight: 600;
					color: #667eea;
					font-size: 18px;
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
				.footer a {
					color: #667eea;
					text-decoration: none;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 40px 20px; }
					.header h1 { font-size: 28px; }
					.header-icon { font-size: 48px; }
					.content { padding: 30px 20px; }
					.content h2 { font-size: 24px; }
					.button { padding: 14px 30px; font-size: 15px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">✈️</div>
					<h1>Welcome to TinderTrip!</h1>
					<p style="margin: 0; opacity: 0.9; font-size: 18px; position: relative; z-index: 1;">Your journey begins here</p>
				</div>
				<div class="content">
					<h2>Hello <span class="welcome-name">%s</span>! 👋</h2>
					<p>We're absolutely thrilled to have you join our vibrant community of travelers and adventure seekers!</p>
					
					<div class="features">
						<h3>🌟 What you can do with TinderTrip:</h3>
						<ul>
							<li>Discover amazing travel events and activities</li>
							<li>Connect with like-minded travelers</li>
							<li>Create and join group trips</li>
							<li>Share your travel experiences</li>
						</ul>
					</div>

					<p>Ready to start your adventure? Get started by exploring events in your area or creating your first event!</p>

					<div class="button-container">
						<a href="http://localhost:3000/events" class="button">🚀 Explore Events</a>
					</div>

					<p>If you have any questions or need help getting started, feel free to reach out to our friendly support team. We're here to help!</p>

					<div class="signature">
						<p class="team-name">Happy travels! 🌍</p>
						<p>The TinderTrip Team</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
					<p style="margin-top: 10px; font-size: 12px;">Connect with us: <a href="#">Website</a> | <a href="#">Support</a></p>
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
	subject := fmt.Sprintf("Event Confirmation: %s - TinderTrip", eventTitle)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Event Confirmation</title>
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
				}
				.event-info h3 {
					font-size: 24px;
					font-weight: 700;
					color: #065f46;
					margin-bottom: 20px;
					text-align: center;
				}
				.event-detail {
					display: flex;
					align-items: center;
					justify-content: center;
					margin: 15px 0;
					padding: 15px;
					background-color: white;
					border-radius: 8px;
					box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
				}
				.event-detail-icon {
					font-size: 24px;
					margin-right: 15px;
				}
				.event-detail-text {
					font-size: 16px;
					color: #1a202c;
				}
				.event-detail-text strong {
					color: #065f46;
					display: block;
					margin-bottom: 5px;
					font-size: 12px;
					text-transform: uppercase;
					letter-spacing: 1px;
				}
				.success-message {
					background: linear-gradient(135deg, #ecfdf5 0%%, #d1fae5 100%%);
					border-left: 4px solid #10b981;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					text-align: center;
				}
				.success-message p {
					margin: 0;
					color: #065f46;
					font-weight: 600;
					font-size: 16px;
				}
				.signature {
					margin-top: 30px;
					padding-top: 25px;
					border-top: 2px solid #e2e8f0;
					text-align: center;
				}
				.signature p {
					margin: 8px 0;
					color: #4a5568;
				}
				.signature .team-name {
					font-weight: 600;
					color: #10b981;
					font-size: 18px;
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
					.event-info h3 { font-size: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">✅</div>
					<h1>Event Confirmed!</h1>
					<p style="margin: 0; opacity: 0.9;">You're all set</p>
				</div>
				<div class="content">
					<h2>Hello %s! 👋</h2>
					<p>Great news! You have successfully confirmed your participation in the following event:</p>
					
					<div class="event-info">
						<h3>%s</h3>
						<div class="event-detail">
							<div class="event-detail-icon">📅</div>
							<div class="event-detail-text">
								<strong>Event Date</strong>
								%s
							</div>
						</div>
					</div>

					<div class="success-message">
						<p>🎉 We're looking forward to seeing you there!</p>
					</div>

					<p>If you need to make any changes or have questions about the event, please don't hesitate to contact the event organizer through the TinderTrip app.</p>

					<div class="signature">
						<p class="team-name">Happy travels! ✈️</p>
						<p>The TinderTrip Team</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
					<p style="margin-top: 10px; font-size: 12px;">This is an automated confirmation email.</p>
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

// SendVerificationOTP sends an email verification OTP email
func (c *SMTPClient) SendVerificationOTP(to, otp string) error {
	subject := "Email Verification - TinderTrip"

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email Verification</title>
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
				.otp-container {
					text-align: center;
					margin: 30px 0;
				}
				.otp-label {
					font-size: 14px;
					color: #64748b;
					text-transform: uppercase;
					letter-spacing: 1px;
					margin-bottom: 10px;
					font-weight: 600;
				}
				.otp-code {
					font-size: 36px;
					font-weight: 700;
					color: #3b82f6;
					text-align: center;
					padding: 25px;
					background: linear-gradient(135deg, #eff6ff 0%%, #dbeafe 100%%);
					border: 3px solid #3b82f6;
					border-radius: 12px;
					margin: 15px auto;
					letter-spacing: 8px;
					display: inline-block;
					min-width: 280px;
					box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
				}
				.warning {
					background: linear-gradient(135deg, #fef3c7 0%%, #fde68a 100%%);
					border-left: 4px solid #f59e0b;
					color: #92400e;
					padding: 20px;
					border-radius: 8px;
					margin: 25px 0;
					box-shadow: 0 2px 8px rgba(245, 158, 11, 0.1);
				}
				.warning strong {
					display: block;
					margin-bottom: 8px;
					font-size: 16px;
				}
				.warning p {
					margin: 0;
					font-size: 14px;
				}
				.welcome-box {
					background: linear-gradient(135deg, #ecfdf5 0%%, #d1fae5 100%%);
					border: 2px solid #10b981;
					border-radius: 8px;
					padding: 20px;
					margin: 25px 0;
					text-align: center;
				}
				.welcome-box p {
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
				.footer a {
					color: #3b82f6;
					text-decoration: none;
				}
				@media only screen and (max-width: 600px) {
					.header { padding: 30px 20px; }
					.header h1 { font-size: 26px; }
					.content { padding: 30px 20px; }
					.otp-code { font-size: 28px; letter-spacing: 4px; min-width: 240px; padding: 20px; }
				}
			</style>
		</head>
		<body>
			<div class="email-wrapper">
				<div class="header">
					<div class="header-icon">✉️</div>
					<h1>TinderTrip</h1>
					<p style="margin: 0; opacity: 0.9;">Email Verification</p>
				</div>
				<div class="content">
					<h2>Verify Your Email Address</h2>
					<p>Hello,</p>
					<p>Thank you for registering with TinderTrip! We're excited to have you join our community. To complete your registration and secure your account, please verify your email address using the code below:</p>
					
					<div class="otp-container">
						<div class="otp-label">Your Verification Code</div>
						<div class="otp-code">%s</div>
					</div>

					<div class="warning">
						<strong>⏱️ Security Notice</strong>
						<p>This verification code will expire in <strong>10 minutes</strong> for security reasons. Please use it promptly to complete your registration.</p>
					</div>

					<p style="margin-top: 25px; color: #718096; font-size: 14px;">
						<strong>Didn't create an account?</strong> If you didn't register with TinderTrip, please ignore this email. No account will be created.
					</p>

					<div class="welcome-box">
						<p>🎉 Welcome to TinderTrip! We can't wait to see where your journey takes you.</p>
					</div>
				</div>
				<div class="footer">
					<p><strong>TinderTrip</strong></p>
					<p>&copy; 2024 TinderTrip. All rights reserved.</p>
					<p style="margin-top: 10px; font-size: 12px;">This is an automated message, please do not reply.</p>
				</div>
			</div>
		</body>
		</html>
	`, otp)

	message := &EmailMessage{
		To:      []string{to},
		Subject: subject,
		HTML:    htmlBody,
	}

	return c.SendEmail(message)
}
