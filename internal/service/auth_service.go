package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/database"
	"TinderTrip-Backend/pkg/email"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	emailService *email.SMTPClient
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		emailService: email.NewSMTPClient(),
	}
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

	return &user, nil
}

// SendPasswordResetEmail sends a password reset email
func (s *AuthService) SendPasswordResetEmail(email string) error {
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

	// Generate reset token
	token, err := s.generateResetToken()
	if err != nil {
		return fmt.Errorf("token generation failed: %w", err)
	}

	// Create password reset record
	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	// Save password reset to database
	if err := database.GetDB().Create(passwordReset).Error; err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	// Send email
	resetURL := "http://localhost:3001/reset-password.html" // Frontend URL
	fmt.Printf("DEBUG: Attempting to send email to %s with token %s\n", email, token)
	err = s.emailService.SendPasswordResetEmail(email, token, resetURL)
	if err != nil {
		fmt.Printf("DEBUG: Email sending failed: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Printf("DEBUG: Email sent successfully to %s\n", email)
	return nil
}

// ResetPassword resets user password
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Find password reset record
	var passwordReset models.PasswordReset
	err := database.GetDB().Where("token = ? AND expires_at > ?", token, time.Now()).First(&passwordReset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("database error: %w", err)
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

	// Delete password reset record
	err = database.GetDB().Delete(&passwordReset).Error
	if err != nil {
		return fmt.Errorf("failed to delete password reset record: %w", err)
	}

	return nil
}

// generateResetToken generates a secure random token
func (s *AuthService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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
