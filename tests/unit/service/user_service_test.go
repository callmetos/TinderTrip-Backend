package service_test

import (
	"fmt"
	"testing"
	"time"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserServiceTest(t *testing.T) (*gorm.DB, *service.UserService) {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create simplified tables for testing (SQLite compatible)
	sqlDB, _ := db.DB()

	// Users table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE,
			provider TEXT NOT NULL,
			password_hash TEXT,
			email_verified BOOLEAN NOT NULL DEFAULT 0,
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

	// User profiles table
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS user_profiles (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			bio TEXT,
			languages TEXT,
			date_of_birth DATE,
			gender TEXT,
			job_title TEXT,
			smoking TEXT,
			interests_note TEXT,
			avatar_url TEXT,
			home_location TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatal("Failed to create user_profiles table:", err)
	}

	// Set global DB for testing
	database.DB = db

	userService := service.NewUserService()
	return db, userService
}

func TestUserService_GetProfile(t *testing.T) {
	db, userService := setupUserServiceTest(t)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
		errMsg  string
	}{
		{
			name: "Profile exists",
			setup: func() string {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("user-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				bio := "Test bio"
				profile := &models.UserProfile{
					UserID: user.ID,
					Bio:    &bio,
				}
				if err := db.Create(profile).Error; err != nil {
					t.Fatalf("Failed to create profile: %v", err)
				}

				return user.ID.String()
			},
			wantErr: false,
		},
		{
			name: "Profile not found",
			setup: func() string {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("noprofile-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				return user.ID.String()
			},
			wantErr: true,
			errMsg:  "profile not found",
		},
		{
			name: "Invalid user ID",
			setup: func() string {
				return "invalid-uuid"
			},
			wantErr: true,
			errMsg:  "invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			profile, err := userService.GetProfile(userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, profile)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, userID, profile.UserID)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	db, userService := setupUserServiceTest(t)

	tests := []struct {
		name    string
		setup   func() (string, dto.UpdateProfileRequest)
		wantErr bool
		check   func(t *testing.T, profile *dto.UserProfileResponse)
	}{
		{
			name: "Update existing profile",
			setup: func() (string, dto.UpdateProfileRequest) {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("update-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				bio := "Old bio"
				profile := &models.UserProfile{
					UserID: user.ID,
					Bio:    &bio,
				}
				if err := db.Create(profile).Error; err != nil {
					t.Fatalf("Failed to create profile: %v", err)
				}

				newBio := "New bio"
				return user.ID.String(), dto.UpdateProfileRequest{
					Bio: &newBio,
				}
			},
			wantErr: false,
			check: func(t *testing.T, profile *dto.UserProfileResponse) {
				assert.NotNil(t, profile.Bio)
				assert.Equal(t, "New bio", *profile.Bio)
			},
		},
		{
			name: "Create new profile",
			setup: func() (string, dto.UpdateProfileRequest) {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("new-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				bio := "New profile bio"
				languages := "English, Thai"
				return user.ID.String(), dto.UpdateProfileRequest{
					Bio:       &bio,
					Languages: &languages,
				}
			},
			wantErr: false,
			check: func(t *testing.T, profile *dto.UserProfileResponse) {
				assert.NotNil(t, profile.Bio)
				assert.Equal(t, "New profile bio", *profile.Bio)
				assert.NotNil(t, profile.Languages)
				assert.Equal(t, "English, Thai", *profile.Languages)
			},
		},
		{
			name: "Update with all fields",
			setup: func() (string, dto.UpdateProfileRequest) {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("allfields-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				bio := "Full bio"
				languages := "English"
				dob := time.Now().AddDate(-25, 0, 0)
				gender := "male"
				jobTitle := "Engineer"
				smoking := "non_smoker"
				interests := "Travel"
				homeLocation := "Bangkok"

				return user.ID.String(), dto.UpdateProfileRequest{
					Bio:           &bio,
					Languages:     &languages,
					DateOfBirth:   &dob,
					Gender:        &gender,
					JobTitle:      &jobTitle,
					Smoking:       &smoking,
					InterestsNote: &interests,
					HomeLocation:  &homeLocation,
				}
			},
			wantErr: false,
			check: func(t *testing.T, profile *dto.UserProfileResponse) {
				assert.NotNil(t, profile.Bio)
				assert.Equal(t, "Full bio", *profile.Bio)
				assert.NotNil(t, profile.Languages)
				assert.Equal(t, "English", *profile.Languages)
				assert.NotNil(t, profile.JobTitle)
				assert.Equal(t, "Engineer", *profile.JobTitle)
				assert.NotNil(t, profile.HomeLocation)
				assert.Equal(t, "Bangkok", *profile.HomeLocation)
			},
		},
		{
			name: "Invalid user ID",
			setup: func() (string, dto.UpdateProfileRequest) {
				bio := "Test"
				return "invalid-uuid", dto.UpdateProfileRequest{
					Bio: &bio,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, req := tt.setup()
			profile, err := userService.UpdateProfile(userID, req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				if tt.check != nil {
					tt.check(t, profile)
				}
			}
		})
	}
}

func TestUserService_DeleteProfile(t *testing.T) {
	db, userService := setupUserServiceTest(t)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "Delete existing profile",
			setup: func() string {
				// Use unique email with timestamp to avoid UNIQUE constraint issues
				email := fmt.Sprintf("delete-%d@example.com", time.Now().UnixNano())
				user := &models.User{
					Email:    &email,
					Provider: models.AuthProviderPassword,
				}
				if err := db.Create(user).Error; err != nil {
					t.Fatalf("Failed to create user: %v", err)
				}

				bio := "To be deleted"
				profile := &models.UserProfile{
					UserID: user.ID,
					Bio:    &bio,
				}
				if err := db.Create(profile).Error; err != nil {
					t.Fatalf("Failed to create profile: %v", err)
				}

				return user.ID.String()
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent profile",
			setup: func() string {
				// Valid UUID but no profile
				return uuid.New().String()
			},
			wantErr: false, // GORM doesn't error on updates that affect 0 rows
		},
		{
			name: "Invalid user ID",
			setup: func() string {
				return "invalid-uuid"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			err := userService.DeleteProfile(userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ProfileWithFullData(t *testing.T) {
	db, userService := setupUserServiceTest(t)

	// Create user with complete profile
	// Use unique email with timestamp to avoid UNIQUE constraint issues
	email := fmt.Sprintf("complete-%d@example.com", time.Now().UnixNano())
	displayName := "Complete User"
	user := &models.User{
		Email:       &email,
		Provider:    models.AuthProviderPassword,
		DisplayName: &displayName,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	bio := "Complete bio"
	languages := "English, Thai, Japanese"
	dob := time.Now().AddDate(-30, 0, 0)
	gender := models.GenderMale
	jobTitle := "Senior Engineer"
	smoking := models.SmokingNo
	interests := "Travel, Food, Photography"
	avatarURL := "https://example.com/avatar.jpg"
	homeLocation := "Bangkok, Thailand"

	profile := &models.UserProfile{
		UserID:        user.ID,
		Bio:           &bio,
		Languages:     &languages,
		DateOfBirth:   &dob,
		Gender:        &gender,
		JobTitle:      &jobTitle,
		Smoking:       &smoking,
		InterestsNote: &interests,
		AvatarURL:     &avatarURL,
		HomeLocation:  &homeLocation,
	}
	if err := db.Create(profile).Error; err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Retrieve and verify
	result, err := userService.GetProfile(user.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, user.ID.String(), result.UserID)
	assert.NotNil(t, result.Bio)
	assert.Equal(t, bio, *result.Bio)
	assert.NotNil(t, result.Languages)
	assert.Equal(t, languages, *result.Languages)
	assert.NotNil(t, result.Gender)
	assert.Equal(t, "male", result.Gender)
	assert.NotNil(t, result.JobTitle)
	assert.Equal(t, jobTitle, *result.JobTitle)
	assert.NotNil(t, result.Smoking)
	assert.Equal(t, "no", result.Smoking)
	assert.NotNil(t, result.HomeLocation)
	assert.Equal(t, homeLocation, *result.HomeLocation)
}

func TestUserService_UpdateProfile_PartialUpdate(t *testing.T) {
	db, userService := setupUserServiceTest(t)

	// Create user with profile
	// Use unique email with timestamp to avoid UNIQUE constraint issues
	email := fmt.Sprintf("partial-%d@example.com", time.Now().UnixNano())
	user := &models.User{
		Email:    &email,
		Provider: models.AuthProviderPassword,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	bio := "Original bio"
	languages := "English"
	profile := &models.UserProfile{
		UserID:    user.ID,
		Bio:       &bio,
		Languages: &languages,
	}
	if err := db.Create(profile).Error; err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Update only bio
	newBio := "Updated bio"
	result, err := userService.UpdateProfile(user.ID.String(), dto.UpdateProfileRequest{
		Bio: &newBio,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newBio, *result.Bio)
	// Languages should still be the same
	assert.NotNil(t, result.Languages)
	assert.Equal(t, languages, *result.Languages)
}
