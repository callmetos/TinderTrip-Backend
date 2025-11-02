package service

import (
	"fmt"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FoodPreferenceService handles food preference business logic
type FoodPreferenceService struct {
}

// NewFoodPreferenceService creates a new food preference service
func NewFoodPreferenceService() *FoodPreferenceService {
	return &FoodPreferenceService{}
}

// GetFoodPreferences gets user's food preferences
func (s *FoodPreferenceService) GetFoodPreferences(userID string) ([]dto.FoodPreferenceResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get food preferences
	var preferences []models.FoodPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get food preferences: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.FoodPreferenceResponse, len(preferences))
	for i, pref := range preferences {
		responses[i] = dto.FoodPreferenceResponse{
			ID:              pref.ID.String(),
			UserID:          pref.UserID.String(),
			FoodCategory:    pref.FoodCategory,
			PreferenceLevel: pref.PreferenceLevel,
			CreatedAt:       pref.CreatedAt,
			UpdatedAt:       pref.UpdatedAt,
		}
	}

	return responses, nil
}

// UpdateFoodPreference updates a single food preference
func (s *FoodPreferenceService) UpdateFoodPreference(userID string, req dto.UpdateFoodPreferenceRequest) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Validate food category (check against database)
	if !s.isValidFoodCategory(req.FoodCategory) {
		return fmt.Errorf("invalid food category")
	}

	// Validate preference level
	if !models.IsValidPreferenceLevel(req.PreferenceLevel) {
		return fmt.Errorf("invalid preference level")
	}

	// Update or create food preference
	var preference models.FoodPreference
	err = database.GetDB().Where("user_id = ? AND food_category = ? AND deleted_at IS NULL", userUUID, req.FoodCategory).First(&preference).Error

	if err == gorm.ErrRecordNotFound {
		// Create new preference
		preference = models.FoodPreference{
			UserID:          userUUID,
			FoodCategory:    req.FoodCategory,
			PreferenceLevel: req.PreferenceLevel,
		}
		err = database.GetDB().Create(&preference).Error
	} else if err == nil {
		// Update existing preference
		preference.PreferenceLevel = req.PreferenceLevel
		err = database.GetDB().Save(&preference).Error
	}

	if err != nil {
		return fmt.Errorf("failed to update food preference: %w", err)
	}

	return nil
}

// UpdateAllFoodPreferences updates all food preferences
func (s *FoodPreferenceService) UpdateAllFoodPreferences(userID string, req dto.UpdateAllFoodPreferencesRequest) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Start transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update each preference
	for _, prefReq := range req.Preferences {
		// Validate food category (check against database)
		if !s.isValidFoodCategory(prefReq.FoodCategory) {
			tx.Rollback()
			return fmt.Errorf("invalid food category: %s", prefReq.FoodCategory)
		}

		// Validate preference level
		if !models.IsValidPreferenceLevel(prefReq.PreferenceLevel) {
			tx.Rollback()
			return fmt.Errorf("invalid preference level for category %s", prefReq.FoodCategory)
		}

		// Update or create preference
		var preference models.FoodPreference
		err = tx.Where("user_id = ? AND food_category = ? AND deleted_at IS NULL", userUUID, prefReq.FoodCategory).First(&preference).Error

		if err == gorm.ErrRecordNotFound {
			// Create new preference
			preference = models.FoodPreference{
				UserID:          userUUID,
				FoodCategory:    prefReq.FoodCategory,
				PreferenceLevel: prefReq.PreferenceLevel,
			}
			err = tx.Create(&preference).Error
		} else if err == nil {
			// Update existing preference
			preference.PreferenceLevel = prefReq.PreferenceLevel
			err = tx.Save(&preference).Error
		}

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update food preference for category %s: %w", prefReq.FoodCategory, err)
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit food preferences update: %w", err)
	}

	return nil
}

// GetFoodPreferenceCategories gets available food preference categories from database
func (s *FoodPreferenceService) GetFoodPreferenceCategories() []dto.FoodPreferenceCategoryResponse {
	// Get food categories from master table
	var masterCategories []models.FoodCategoryMaster
	err := database.GetDB().
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&masterCategories).Error

	if err != nil || len(masterCategories) == 0 {
		// Fallback to hardcoded categories if database is empty
		return s.getHardcodedFoodCategories()
	}

	// Convert master data to response DTOs
	responses := make([]dto.FoodPreferenceCategoryResponse, len(masterCategories))
	for i, category := range masterCategories {
		displayName := category.DisplayName
		if displayName == "" {
			displayName = models.GetFoodCategoryName(category.Code) // Fallback
		}

		icon := ""
		if category.Icon != nil && *category.Icon != "" {
			icon = *category.Icon
		} else {
			icon = models.GetFoodCategoryIcon(category.Code) // Fallback
		}

		responses[i] = dto.FoodPreferenceCategoryResponse{
			Category:    category.Code,
			DisplayName: displayName,
			Icon:        icon,
		}
	}

	return responses
}

// getHardcodedFoodCategories returns hardcoded food categories as fallback
func (s *FoodPreferenceService) getHardcodedFoodCategories() []dto.FoodPreferenceCategoryResponse {
	categories := []string{
		"thai_food",
		"japanese_food",
		"chinese_food",
		"international_food",
		"halal_food",
		"buffet",
		"bbq_grill",
	}

	responses := make([]dto.FoodPreferenceCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = dto.FoodPreferenceCategoryResponse{
			Category:    category,
			DisplayName: models.GetFoodCategoryName(category),
			Icon:        models.GetFoodCategoryIcon(category),
		}
	}

	return responses
}

// GetFoodPreferenceCategoriesWithUserPreferences gets categories with user's current preferences
func (s *FoodPreferenceService) GetFoodPreferenceCategoriesWithUserPreferences(userID string) ([]dto.FoodPreferenceCategoryResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get user's current preferences
	var preferences []models.FoodPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Create preference map
	preferenceMap := make(map[string]int)
	for _, pref := range preferences {
		preferenceMap[pref.FoodCategory] = pref.PreferenceLevel
	}

	// Get all categories from database
	masterCategories := s.GetFoodPreferenceCategories()

	responses := make([]dto.FoodPreferenceCategoryResponse, len(masterCategories))
	for i, category := range masterCategories {
		preferenceLevel := 2 // default to neutral
		if level, exists := preferenceMap[category.Category]; exists {
			preferenceLevel = level
		}

		responses[i] = dto.FoodPreferenceCategoryResponse{
			Category:        category.Category,
			DisplayName:     category.DisplayName,
			Icon:            category.Icon,
			PreferenceLevel: preferenceLevel,
		}
	}

	return responses, nil
}

// GetFoodPreferenceStats gets food preference statistics for a user
func (s *FoodPreferenceService) GetFoodPreferenceStats(userID string) (*dto.FoodPreferenceStatsResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get preferences
	var preferences []models.FoodPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	// Calculate stats
	stats := &dto.FoodPreferenceStatsResponse{
		TotalPreferences: len(preferences),
		DislikeCount:     0,
		NeutralCount:     0,
		LoveCount:        0,
	}

	for _, pref := range preferences {
		switch pref.PreferenceLevel {
		case 1:
			stats.DislikeCount++
		case 2:
			stats.NeutralCount++
		case 3:
			stats.LoveCount++
		}
	}

	return stats, nil
}

// DeleteFoodPreference deletes a food preference
func (s *FoodPreferenceService) DeleteFoodPreference(userID, foodCategory string) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Delete preference
	err = database.GetDB().Where("user_id = ? AND food_category = ? AND deleted_at IS NULL", userUUID, foodCategory).Delete(&models.FoodPreference{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete food preference: %w", err)
	}

	return nil
}

// isValidFoodCategory checks if food category exists in database
func (s *FoodPreferenceService) isValidFoodCategory(category string) bool {
	// Check against database first
	var masterCategory models.FoodCategoryMaster
	err := database.GetDB().
		Where("code = ? AND is_active = ?", category, true).
		First(&masterCategory).Error

	if err == nil {
		return true // Found in database
	}

	// Fallback to hardcoded validation (check against constants)
	validCategories := []string{
		"thai_food",
		"japanese_food",
		"chinese_food",
		"international_food",
		"halal_food",
		"buffet",
		"bbq_grill",
	}
	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}

	return false
}
