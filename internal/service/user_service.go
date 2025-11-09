package service

import (
	"fmt"
	"strings"
	"time"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct {
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// GetProfile gets user profile
func (s *UserService) GetProfile(userID string) (*dto.UserProfileResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get user profile
	var profile models.UserProfile
	err = database.GetDB().Where("user_id = ?", userUUID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Get user to get display_name
	var user models.User
	err = database.GetDB().Where("id = ?", userUUID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert to response DTO
	var gender, smoking string
	if profile.Gender != nil {
		gender = string(*profile.Gender)
	}
	if profile.Smoking != nil {
		smoking = string(*profile.Smoking)
	}

	response := &dto.UserProfileResponse{
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
		DisplayName:   user.DisplayName,
		Bio:           profile.Bio,
		Languages:     profile.Languages,
		DateOfBirth:   profile.DateOfBirth,
		Age:           profile.GetAge(),
		Gender:        gender,
		JobTitle:      profile.JobTitle,
		Smoking:       smoking,
		InterestsNote: profile.InterestsNote,
		AvatarURL: func() *string {
			if profile.AvatarURL != nil && *profile.AvatarURL != "" {
				userID := profile.UserID.String()
				publicURL := fmt.Sprintf("https://api.tindertrip.phitik.com/images/avatars/%s", userID)
				return &publicURL
			}
			return nil
		}(),
		HomeLocation: profile.HomeLocation,
		CreatedAt:    profile.CreatedAt,
		UpdatedAt:    profile.UpdatedAt,
	}

	return response, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if profile exists
	var profile models.UserProfile
	err = database.GetDB().Where("user_id = ?", userUUID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new profile
			profile = models.UserProfile{
				UserID: userUUID,
			}
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	// Update display_name in users table if provided
	if req.DisplayName != nil {
		// Check if display_name already exists (excluding current user)
		var existingDisplayName models.User
		err = database.GetDB().Where("display_name = ? AND id != ? AND deleted_at IS NULL", *req.DisplayName, userUUID).First(&existingDisplayName).Error
		if err == nil {
			return nil, fmt.Errorf("display name already taken")
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("database error: %w", err)
		}
		
		// Update display_name
		err = database.GetDB().Model(&models.User{}).Where("id = ?", userUUID).Update("display_name", *req.DisplayName).Error
		if err != nil {
			// Check if error is due to unique constraint violation for display_name
			errStr := strings.ToLower(err.Error())
			// PostgreSQL unique constraint violation error codes and messages
			// Error format: "ERROR: duplicate key value violates unique constraint \"ux_users_display_name\"\nSQL state: 23505"
			if errStr != "" && (
				strings.Contains(errStr, "ux_users_display_name") || 
				strings.Contains(errStr, "duplicate key value") || 
				strings.Contains(errStr, "unique constraint") ||
				strings.Contains(errStr, "23505") || // PostgreSQL unique violation error code
				strings.Contains(errStr, "violates unique constraint") ||
				strings.Contains(errStr, "already exists") ||
				strings.Contains(errStr, "key (display_name)") ||
				strings.Contains(errStr, "sql state: 23505")) {
				return nil, fmt.Errorf("display name already taken")
			}
			return nil, fmt.Errorf("failed to update display_name: %w", err)
		}
	}

	// Update profile fields
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Languages != nil {
		profile.Languages = req.Languages
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != nil {
		gender := models.Gender(*req.Gender)
		profile.Gender = &gender
	}
	if req.JobTitle != nil {
		profile.JobTitle = req.JobTitle
	}
	if req.Smoking != nil {
		smoking := models.Smoking(*req.Smoking)
		profile.Smoking = &smoking
	}
	if req.InterestsNote != nil {
		profile.InterestsNote = req.InterestsNote
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL
	}
	if req.HomeLocation != nil {
		profile.HomeLocation = req.HomeLocation
	}

	// Save profile
	if profile.ID == uuid.Nil {
		err = database.GetDB().Create(&profile).Error
	} else {
		err = database.GetDB().Save(&profile).Error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to save profile: %w", err)
	}

	// Get user to get display_name (after update)
	var user models.User
	err = database.GetDB().Where("id = ?", userUUID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert to response DTO
	var gender, smoking string
	if profile.Gender != nil {
		gender = string(*profile.Gender)
	}
	if profile.Smoking != nil {
		smoking = string(*profile.Smoking)
	}

	response := &dto.UserProfileResponse{
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
		DisplayName:   user.DisplayName,
		Bio:           profile.Bio,
		Languages:     profile.Languages,
		DateOfBirth:   profile.DateOfBirth,
		Age:           req.Age, // Use age from request instead of calculating
		Gender:        gender,
		JobTitle:      profile.JobTitle,
		Smoking:       smoking,
		InterestsNote: profile.InterestsNote,
		AvatarURL: func() *string {
			if profile.AvatarURL != nil && *profile.AvatarURL != "" {
				userID := profile.UserID.String()
				publicURL := fmt.Sprintf("https://api.tindertrip.phitik.com/images/avatars/%s", userID)
				return &publicURL
			}
			return nil
		}(),
		HomeLocation: profile.HomeLocation,
		CreatedAt:    profile.CreatedAt,
		UpdatedAt:    profile.UpdatedAt,
	}

	return response, nil
}

// DeleteProfile deletes user profile
func (s *UserService) DeleteProfile(userID string) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Soft delete profile
	now := time.Now()
	err = database.GetDB().Model(&models.UserProfile{}).Where("user_id = ?", userUUID).Update("deleted_at", now).Error
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// CheckSetupStatus checks if user has completed initial setup
func (s *UserService) CheckSetupStatus(userID string) (bool, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get user profile
	var profile models.UserProfile
	err = database.GetDB().Where("user_id = ?", userUUID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No profile created yet = setup not completed
			return false, nil
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	// Check if essential fields are filled
	// Consider setup completed if user has at least bio OR gender OR languages
	hasEssentialInfo := (profile.Bio != nil && *profile.Bio != "") ||
		(profile.Gender != nil) ||
		(profile.Languages != nil && *profile.Languages != "")

	return hasEssentialInfo, nil
}

// Helper functions
func convertTimeToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.Format("2006-01-02")
	return &str
}

func convertGenderToString(g *models.Gender) *string {
	if g == nil {
		return nil
	}
	str := string(*g)
	return &str
}

func convertSmokingToString(s *models.Smoking) *string {
	if s == nil {
		return nil
	}
	str := string(*s)
	return &str
}
