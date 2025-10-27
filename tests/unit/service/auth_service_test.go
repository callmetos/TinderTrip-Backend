package service_test

import (
	"testing"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthServiceTest(t *testing.T) (*gorm.DB, *service.AuthService) {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create simplified tables for testing (SQLite compatible)
	// We manually create tables instead of AutoMigrate to avoid PostgreSQL-specific types
	sqlDB, _ := db.DB()

	// Users table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE,
			provider TEXT NOT NULL,
			password_hash TEXT,
			email_verified INTEGER NOT NULL DEFAULT 0,
			google_id TEXT,
			display_name TEXT,
			last_login_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatal("Failed to create users table:", err)
	}

	// Password resets table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS password_resets (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal("Failed to create password_resets table:", err)
	}

	// Email verifications table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS email_verifications (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			otp TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatal("Failed to create email_verifications table:", err)
	}

	// Audit logs table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			actor_user_id TEXT,
			entity_table TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			action TEXT NOT NULL,
			before_data TEXT,
			after_data TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal("Failed to create audit_logs table:", err)
	}

	// Set global DB for testing
	database.DB = db

	// Setup config
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpireHours: 24,
		},
	}

	authService := service.NewAuthService()
	return db, authService
}

func TestAuthService_GenerateOTP(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	// Since generateOTP is private, we test it indirectly through password reset
	email := "test@example.com"

	// Create a test user first
	password := "TestPass123!"
	hashedPass, _ := utils.HashPassword(password)
	user := &models.User{
		Email:        &email,
		Provider:     models.AuthProviderPassword,
		PasswordHash: &hashedPass,
	}
	database.DB.Create(user)

	// Test OTP generation through password reset
	// Note: This will fail to send email since SMTP is not configured, but OTP should still be generated
	_ = authService.SendPasswordResetOTP(email)

	// Verify OTP was created in database (even if email fails)
	var resetRecord models.PasswordReset
	err := database.DB.Where("user_id = ?", user.ID).First(&resetRecord).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, resetRecord.Token)
	assert.Len(t, resetRecord.Token, 6) // OTP should be 6 digits
}

func TestAuthService_VerifyOTP(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name    string
		setup   func() (email string, otp string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid OTP",
			setup: func() (string, string) {
				email := "valid@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:        &email,
					Provider:     models.AuthProviderPassword,
					PasswordHash: &hashedPass,
				}
				database.DB.Create(user)

				// Create reset record
				reset := &models.PasswordReset{
					UserID:    user.ID,
					Token:     "123456",
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				database.DB.Create(reset)

				return email, "123456"
			},
			wantErr: false,
		},
		{
			name: "Invalid OTP",
			setup: func() (string, string) {
				email := "invalid@example.com"
				return email, "999999"
			},
			wantErr: true,
			errMsg:  "invalid or expired OTP",
		},
		{
			name: "Expired OTP",
			setup: func() (string, string) {
				email := "expired@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:        &email,
					Provider:     models.AuthProviderPassword,
					PasswordHash: &hashedPass,
				}
				database.DB.Create(user)

				// Create expired reset record
				reset := &models.PasswordReset{
					UserID:    user.ID,
					Token:     "654321", // Use different OTP
					ExpiresAt: time.Now().Add(-10 * time.Minute),
				}
				database.DB.Create(reset)

				return email, "654321" // Use different OTP
			},
			wantErr: true,
			errMsg:  "invalid or expired OTP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, otp := tt.setup()
			err := authService.VerifyOTP(email, otp)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name    string
		setup   func() (email, password string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful login with password",
			setup: func() (string, string) {
				email := "login@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderPassword,
					PasswordHash:  &hashedPass,
					EmailVerified: true,
				}
				database.DB.Create(user)
				return email, password
			},
			wantErr: false,
		},
		{
			name: "Wrong password",
			setup: func() (string, string) {
				email := "wrong@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderPassword,
					PasswordHash:  &hashedPass,
					EmailVerified: true,
				}
				database.DB.Create(user)
				return email, "WrongPassword123!"
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "User not found",
			setup: func() (string, string) {
				return "notfound@example.com", "TestPass123!"
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "Google user cannot login with password",
			setup: func() (string, string) {
				email := "google@example.com"
				googleID := "google-123"
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderGoogle,
					GoogleID:      &googleID,
					EmailVerified: true,
				}
				database.DB.Create(user)
				return email, "anypassword"
			},
			wantErr: true,
			errMsg:  "Google",
		},
		{
			name: "Email not verified",
			setup: func() (string, string) {
				email := "unverified@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderPassword,
					PasswordHash:  &hashedPass,
					EmailVerified: false,
				}
				database.DB.Create(user)
				return email, password
			},
			wantErr: true,
			errMsg:  "email not verified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, password := tt.setup()
			user, err := authService.Login(email, password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, email, *user.Email)
				assert.NotNil(t, user.LastLoginAt)
			}
		})
	}
}

func TestAuthService_ResetPassword(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name        string
		setup       func() (email, otp, newPassword string)
		wantErr     bool
		errMsg      string
		checkResult func(t *testing.T, email, newPassword string)
	}{
		{
			name: "Successful password reset",
			setup: func() (string, string, string) {
				email := "resetpass123@example.com"
				oldPassword := "OldPass123!"
				hashedPass, _ := utils.HashPassword(oldPassword)
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderPassword,
					PasswordHash:  &hashedPass,
					EmailVerified: true,
				}
				database.DB.Create(user)

				// Create reset record
				otp := "555555"
				reset := &models.PasswordReset{
					UserID:    user.ID,
					Token:     otp,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				database.DB.Create(reset)

				newPassword := "NewPass123!"
				return email, otp, newPassword
			},
			wantErr: false,
			checkResult: func(t *testing.T, email, newPassword string) {
				// Try to login with new password
				var user models.User
				err := database.DB.Where("email = ?", email).First(&user).Error
				assert.NoError(t, err)

				if user.PasswordHash != nil {
					valid, err := utils.VerifyPassword(newPassword, *user.PasswordHash)
					assert.NoError(t, err)
					assert.True(t, valid)
				}
			},
		},
		{
			name: "Invalid OTP",
			setup: func() (string, string, string) {
				email := "invalid@example.com"
				return email, "999999", "NewPass123!"
			},
			wantErr: true,
			errMsg:  "invalid or expired OTP",
		},
		{
			name: "Email mismatch",
			setup: func() (string, string, string) {
				email := "mismatch@example.com"
				password := "TestPass123!"
				hashedPass, _ := utils.HashPassword(password)
				user := &models.User{
					Email:         &email,
					Provider:      models.AuthProviderPassword,
					PasswordHash:  &hashedPass,
					EmailVerified: true,
				}
				database.DB.Create(user)

				// Create reset record for this user
				otp := "123456"
				reset := &models.PasswordReset{
					UserID:    user.ID,
					Token:     otp,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				database.DB.Create(reset)

				// But try to reset with different email
				return "different@example.com", otp, "NewPass123!"
			},
			wantErr: true,
			errMsg:  "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, otp, newPassword := tt.setup()
			err := authService.ResetPassword(email, otp, newPassword)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, email, newPassword)
				}

				// Note: The reset record may or may not be deleted depending on implementation
				// So we don't assert here, just verify password was changed in checkResult
			}
		})
	}
}

func TestAuthService_GetUserByEmail(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "User exists",
			setup: func() string {
				email := "exists@example.com"
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				database.DB.Create(user)
				return email
			},
			wantErr: false,
		},
		{
			name: "User not found",
			setup: func() string {
				return "notfound@example.com"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email := tt.setup()
			user, err := authService.GetUserByEmail(email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, email, *user.Email)
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "User exists",
			setup: func() string {
				email := "byid@example.com"
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				database.DB.Create(user)
				return user.ID.String()
			},
			wantErr: false,
		},
		{
			name: "User not found",
			setup: func() string {
				return "00000000-0000-0000-0000-000000000000"
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID",
			setup: func() string {
				return "invalid-uuid"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			user, err := authService.GetUserByID(userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, userID, user.ID.String())
			}
		})
	}
}

func TestAuthService_UpdateUser(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	t.Run("Update user display name", func(t *testing.T) {
		email := "update@example.com"
		user := &models.User{
			Email:    &email,
			Provider: models.AuthProviderPassword,
		}
		database.DB.Create(user)

		newName := "Updated Name"
		err := authService.UpdateUser(user.ID.String(), map[string]interface{}{
			"display_name": newName,
		})

		assert.NoError(t, err)

		// Verify update
		var updatedUser models.User
		database.DB.First(&updatedUser, user.ID)
		assert.NotNil(t, updatedUser.DisplayName)
		assert.Equal(t, newName, *updatedUser.DisplayName)
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		err := authService.UpdateUser("00000000-0000-0000-0000-000000000000", map[string]interface{}{
			"display_name": "Test",
		})

		// Should not error (GORM doesn't error on updates that affect 0 rows)
		assert.NoError(t, err)
	})
}

func TestAuthService_DeleteUser(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	t.Run("Soft delete user", func(t *testing.T) {
		email := "delete@example.com"
		user := &models.User{
			Email:    &email,
			Provider: models.AuthProviderPassword,
		}
		database.DB.Create(user)

		err := authService.DeleteUser(user.ID.String())
		assert.NoError(t, err)

		// Verify soft delete (deleted_at is set)
		// Note: Need to check if the field is actually populated after delete
		sqlDB, _ := database.DB.DB()
		var deletedAt *string
		err = sqlDB.QueryRow("SELECT deleted_at FROM users WHERE id = ?", user.ID.String()).Scan(&deletedAt)
		assert.NoError(t, err)
		assert.NotNil(t, deletedAt, "deleted_at should be set")
	})
}

func TestAuthService_VerifyEmailOTP(t *testing.T) {
	_, authService := setupAuthServiceTest(t)

	tests := []struct {
		name    string
		setup   func() (email, otp, password, displayName string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid email verification",
			setup: func() (string, string, string, string) {
				email := "verify@example.com"
				otp := "123456"

				// Create email verification record
				verification := &models.EmailVerification{
					Email:     email,
					OTP:       otp,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				database.DB.Create(verification)

				return email, otp, "TestPass123!", "Test User"
			},
			wantErr: false,
		},
		{
			name: "Invalid OTP",
			setup: func() (string, string, string, string) {
				return "invalid@example.com", "999999", "TestPass123!", "Test"
			},
			wantErr: true,
			errMsg:  "invalid or expired OTP",
		},
		{
			name: "User already exists",
			setup: func() (string, string, string, string) {
				email := "existing@example.com"
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				database.DB.Create(user)

				otp := "123456"
				verification := &models.EmailVerification{
					Email:     email,
					OTP:       otp,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				database.DB.Create(verification)

				return email, otp, "TestPass123!", "Test"
			},
			wantErr: true,
			errMsg:  "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, otp, password, displayName := tt.setup()
			user, err := authService.VerifyEmailOTP(email, otp, password, displayName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, email, *user.Email)
				assert.Equal(t, displayName, *user.DisplayName)
				assert.True(t, user.EmailVerified)
				assert.True(t, user.HasPassword())
			}
		})
	}
}
