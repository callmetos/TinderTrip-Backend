package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/audit"
	"TinderTrip-Backend/pkg/database"
	"TinderTrip-Backend/pkg/email"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	emailService *email.SMTPClient
	auditLogger  *audit.AuditLogger
	stopCleanup  chan bool
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	service := &AuthService{
		emailService: email.NewSMTPClient(),
		auditLogger:  audit.NewAuditLogger(),
		stopCleanup:  make(chan bool),
	}

	// Start background cleanup
	go service.startCleanupRoutine()

	return service
}

// Register registers a new user
func (s *AuthService) Register(email, password, displayName string) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	err := database.GetDB().Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        &email,
		Provider:     models.AuthProviderPassword,
		PasswordHash: &hashedPassword,
		DisplayName:  &displayName,
	}

	// Save user to database
	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Log user registration
	userIDStr := user.ID.String()
	s.auditLogger.LogCreate(&userIDStr, "users", &userIDStr, user)

	return user, nil
}

// Login authenticates a user
func (s *AuthService) Login(email, password string) (*models.User, error) {
	// Find user by email
	var user models.User
	err := database.GetDB().Where("email = ? AND provider = ?", email, models.AuthProviderPassword).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user has password
	if user.PasswordHash == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	valid, err := utils.VerifyPassword(password, *user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("password verification failed: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := database.GetDB().Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update last login time: %w", err)
	}

	// Log login action
	userIDStr := user.ID.String()
	s.auditLogger.LogLogin(&userIDStr, map[string]interface{}{
		"email":      email,
		"login_time": now,
	})

	return &user, nil
}

// SendPasswordResetOTP sends a password reset OTP email
func (s *AuthService) SendPasswordResetOTP(email string) error {
	// Find user by email
	var user models.User
	err := database.GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if user exists or not
			return nil
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Generate 6-digit OTP
	otp := s.generateOTP()

	// Delete existing reset tokens for this user
	database.GetDB().Where("user_id = ?", user.ID).Delete(&models.PasswordReset{})

	// Clean up expired tokens
	database.GetDB().Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{})

	// Create password reset record with OTP
	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     otp,
		ExpiresAt: time.Now().Add(3 * time.Minute), // 3 minutes expiry
	}

	// Save password reset to database
	if err := database.GetDB().Create(passwordReset).Error; err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	// Send OTP email
	fmt.Printf("DEBUG: Attempting to send OTP to %s: %s\n", email, otp)
	err = s.emailService.SendPasswordResetOTP(email, otp)
	if err != nil {
		fmt.Printf("DEBUG: OTP email sending failed: %v\n", err)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}
	fmt.Printf("DEBUG: OTP email sent successfully to %s\n", email)
	return nil
}

// ResetPassword resets user password with OTP
func (s *AuthService) ResetPassword(email, otp, newPassword string) error {
	// Find password reset record
	var passwordReset models.PasswordReset
	err := database.GetDB().Where("token = ? AND expires_at > ?", otp, time.Now()).First(&passwordReset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired OTP")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Verify the email matches the OTP
	var user models.User
	err = database.GetDB().Where("id = ? AND email = ?", passwordReset.UserID, email).First(&user).Error
	if err != nil {
		return fmt.Errorf("invalid email for this OTP")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	// Update user password
	err = database.GetDB().Model(&models.User{}).Where("id = ?", passwordReset.UserID).Update("password_hash", hashedPassword).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Log password reset
	userIDStr := passwordReset.UserID.String()
	s.auditLogger.LogPasswordReset(&userIDStr)

	// Delete password reset record
	err = database.GetDB().Delete(&passwordReset).Error
	if err != nil {
		return fmt.Errorf("failed to delete password reset record: %w", err)
	}

	// Clean up expired tokens
	database.GetDB().Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{})

	return nil
}

// VerifyOTP verifies OTP for password reset
func (s *AuthService) VerifyOTP(email, otp string) error {
	// Find password reset record
	var passwordReset models.PasswordReset
	err := database.GetDB().Where("token = ? AND expires_at > ?", otp, time.Now()).First(&passwordReset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired OTP")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Verify the email matches the OTP
	var user models.User
	err = database.GetDB().Where("id = ? AND email = ?", passwordReset.UserID, email).First(&user).Error
	if err != nil {
		return fmt.Errorf("invalid email for this OTP")
	}

	// Clean up expired tokens
	database.GetDB().Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{})

	return nil
}

// startCleanupRoutine starts background cleanup for expired OTPs
func (s *AuthService) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute) // Run every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Clean up expired OTPs
			database.GetDB().Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{})
		case <-s.stopCleanup:
			return
		}
	}
}

// StopCleanup stops the background cleanup routine
func (s *AuthService) StopCleanup() {
	close(s.stopCleanup)
}

// generateResetToken generates a secure random token
func (s *AuthService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateOTP generates a 6-digit OTP
func (s *AuthService) generateOTP() string {
	// Generate a random 6-digit number
	otp := mathrand.Intn(900000) + 100000 // 100000 to 999999
	return fmt.Sprintf("%06d", otp)
}

// ValidateToken validates a password reset token
func (s *AuthService) ValidateToken(token string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	err := database.GetDB().Where("token = ? AND expires_at > ?", token, time.Now()).First(&passwordReset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &passwordReset, nil
}

// GetUserByID gets a user by ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := database.GetDB().Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// GetUserByEmail gets a user by email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(userID string, updates map[string]interface{}) error {
	err := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user (soft delete)
func (s *AuthService) DeleteUser(userID string) error {
	now := time.Now()
	err := database.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("deleted_at", now).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
