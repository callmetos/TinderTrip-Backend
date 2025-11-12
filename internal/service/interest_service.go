package service

import (
	"fmt"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InterestService handles interest-related business logic
type InterestService struct{}

// NewInterestService creates a new interest service
func NewInterestService() *InterestService {
	return &InterestService{}
}

// GetAllInterests gets all active interests, optionally filtered by category
func (s *InterestService) GetAllInterests(category string) ([]dto.InterestResponse, int64, error) {
	var interests []models.Interest
	query := database.GetDB().Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query = query.Order("category ASC, sort_order ASC")

	if err := query.Find(&interests).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get interests: %w", err)
	}

	responses := make([]dto.InterestResponse, len(interests))
	for i, interest := range interests {
		responses[i] = dto.InterestResponse{
			ID:          interest.ID.String(),
			Code:        interest.Code,
			DisplayName: interest.DisplayName,
			Icon:        interest.Icon,
			Category:    interest.Category,
			SortOrder:   interest.SortOrder,
			IsActive:    interest.IsActive,
			CreatedAt:   interest.CreatedAt,
		}
	}

	return responses, int64(len(responses)), nil
}

// GetUserInterests gets all interests with user selection status
func (s *InterestService) GetUserInterests(userID uuid.UUID, category string) ([]dto.InterestResponse, int64, error) {
	var interests []models.Interest
	query := database.GetDB().Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query = query.Order("category ASC, sort_order ASC")

	if err := query.Find(&interests).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get interests: %w", err)
	}

	// Get user's selected interests
	var userInterests []models.UserInterest
	if err := database.GetDB().Where("user_id = ?", userID).Find(&userInterests).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get user interests: %w", err)
	}

	// Create a map of selected interest IDs
	selectedMap := make(map[uuid.UUID]bool)
	for _, ui := range userInterests {
		selectedMap[ui.InterestID] = true
	}

	responses := make([]dto.InterestResponse, len(interests))
	for i, interest := range interests {
		responses[i] = dto.InterestResponse{
			ID:          interest.ID.String(),
			Code:        interest.Code,
			DisplayName: interest.DisplayName,
			Icon:        interest.Icon,
			Category:    interest.Category,
			SortOrder:   interest.SortOrder,
			IsActive:    interest.IsActive,
			CreatedAt:   interest.CreatedAt,
			IsSelected:  selectedMap[interest.ID],
		}
	}

	return responses, int64(len(responses)), nil
}

// UpdateUserInterests updates user interests (bulk replace)
func (s *InterestService) UpdateUserInterests(userID uuid.UUID, interestCodes []string) error {
	// Validate that all interest codes exist
	var interests []models.Interest
	if err := database.GetDB().Where("code IN ? AND is_active = ?", interestCodes, true).Find(&interests).Error; err != nil {
		return fmt.Errorf("failed to validate interests: %w", err)
	}

	if len(interests) != len(interestCodes) {
		return fmt.Errorf("some interest codes are invalid or inactive")
	}

	// Start transaction
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		// Delete all existing user interests
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserInterest{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing interests: %w", err)
		}

		// Create new user interests
		userInterests := make([]models.UserInterest, len(interests))
		for i, interest := range interests {
			userInterests[i] = models.UserInterest{
				UserID:     userID,
				InterestID: interest.ID,
			}
		}

		if len(userInterests) > 0 {
			if err := tx.Create(&userInterests).Error; err != nil {
				return fmt.Errorf("failed to create user interests: %w", err)
			}
		}

		return nil
	})
}

// GetUserSelectedInterests gets only user's selected interests
func (s *InterestService) GetUserSelectedInterests(userID uuid.UUID) ([]dto.InterestResponse, error) {
	var userInterests []models.UserInterest
	if err := database.GetDB().
		Preload("Interest").
		Where("user_id = ?", userID).
		Find(&userInterests).Error; err != nil {
		return nil, fmt.Errorf("failed to get user interests: %w", err)
	}

	responses := make([]dto.InterestResponse, len(userInterests))
	for i, ui := range userInterests {
		if ui.Interest != nil {
			responses[i] = dto.InterestResponse{
				ID:          ui.Interest.ID.String(),
				Code:        ui.Interest.Code,
				DisplayName: ui.Interest.DisplayName,
				Icon:        ui.Interest.Icon,
				Category:    ui.Interest.Category,
				SortOrder:   ui.Interest.SortOrder,
				IsActive:    ui.Interest.IsActive,
				CreatedAt:   ui.Interest.CreatedAt,
				IsSelected:  true,
			}
		}
	}

	return responses, nil
}

