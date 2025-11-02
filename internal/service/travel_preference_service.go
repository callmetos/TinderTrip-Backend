package service

import (
	"fmt"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TravelPreferenceService handles travel preference business logic
type TravelPreferenceService struct {
}

// NewTravelPreferenceService creates a new travel preference service
func NewTravelPreferenceService() *TravelPreferenceService {
	return &TravelPreferenceService{}
}

// GetTravelPreferences gets user's travel preferences
func (s *TravelPreferenceService) GetTravelPreferences(userID string) ([]dto.TravelPreferenceResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get travel preferences
	var preferences []models.TravelPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get travel preferences: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.TravelPreferenceResponse, len(preferences))
	for i, pref := range preferences {
		responses[i] = dto.TravelPreferenceResponse{
			ID:          pref.ID.String(),
			UserID:      pref.UserID.String(),
			TravelStyle: pref.TravelStyle,
			CreatedAt:   pref.CreatedAt,
			UpdatedAt:   pref.UpdatedAt,
		}
	}

	return responses, nil
}

// AddTravelPreference adds a travel preference
func (s *TravelPreferenceService) AddTravelPreference(userID string, req dto.AddTravelPreferenceRequest) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Validate travel style (check against database)
	if !s.isValidTravelStyle(req.TravelStyle) {
		return fmt.Errorf("invalid travel style")
	}

	// Check if preference already exists
	var existingPreference models.TravelPreference
	err = database.GetDB().Where("user_id = ? AND travel_style = ? AND deleted_at IS NULL", userUUID, req.TravelStyle).First(&existingPreference).Error

	if err == nil {
		return fmt.Errorf("travel preference already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing preference: %w", err)
	}

	// Create new preference
	preference := models.TravelPreference{
		UserID:      userUUID,
		TravelStyle: req.TravelStyle,
	}
	err = database.GetDB().Create(&preference).Error
	if err != nil {
		return fmt.Errorf("failed to add travel preference: %w", err)
	}

	return nil
}

// UpdateAllTravelPreferences updates all travel preferences (replace all)
func (s *TravelPreferenceService) UpdateAllTravelPreferences(userID string, req dto.UpdateAllTravelPreferencesRequest) error {
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

	// Delete all existing preferences
	err = tx.Where("user_id = ?", userUUID).Delete(&models.TravelPreference{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing preferences: %w", err)
	}

	// Add new preferences
	for _, travelStyle := range req.TravelStyles {
		// Validate travel style (check against database)
		if !s.isValidTravelStyle(travelStyle) {
			tx.Rollback()
			return fmt.Errorf("invalid travel style: %s", travelStyle)
		}

		// Create new preference
		preference := models.TravelPreference{
			UserID:      userUUID,
			TravelStyle: travelStyle,
		}
		err = tx.Create(&preference).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create travel preference for style %s: %w", travelStyle, err)
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit travel preferences update: %w", err)
	}

	return nil
}

// GetTravelPreferenceStyles gets available travel preference styles from database
func (s *TravelPreferenceService) GetTravelPreferenceStyles() []dto.TravelPreferenceStyleResponse {
	// Get travel styles from master table
	var masterStyles []models.TravelStyleMaster
	err := database.GetDB().
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&masterStyles).Error

	if err != nil || len(masterStyles) == 0 {
		// Fallback to hardcoded styles if database is empty
		return s.getHardcodedTravelStyles()
	}

	// Convert master data to response DTOs
	responses := make([]dto.TravelPreferenceStyleResponse, len(masterStyles))
	for i, style := range masterStyles {
		displayName := style.DisplayName
		if displayName == "" {
			displayName = models.GetTravelStyleName(style.Code) // Fallback
		}

		icon := ""
		if style.Icon != nil && *style.Icon != "" {
			icon = *style.Icon
		} else {
			icon = models.GetTravelStyleIcon(style.Code) // Fallback
		}

		responses[i] = dto.TravelPreferenceStyleResponse{
			Style:       style.Code,
			DisplayName: displayName,
			Icon:        icon,
		}
	}

	return responses
}

// getHardcodedTravelStyles returns hardcoded travel styles as fallback
func (s *TravelPreferenceService) getHardcodedTravelStyles() []dto.TravelPreferenceStyleResponse {
	styles := []string{
		"cafe_dessert",
		"bubble_tea",
		"bakery_cake",
		"bingsu_ice_cream",
		"coffee",
		"matcha",
		"pancakes",
		"social_activity",
		"karaoke",
		"gaming",
		"movie",
		"board_game",
		"outdoor_activity",
		"party_celebration",
		"swimming",
		"skateboarding",
	}

	responses := make([]dto.TravelPreferenceStyleResponse, len(styles))
	for i, style := range styles {
		responses[i] = dto.TravelPreferenceStyleResponse{
			Style:       style,
			DisplayName: models.GetTravelStyleName(style),
			Icon:        models.GetTravelStyleIcon(style),
		}
	}

	return responses
}

// GetTravelPreferenceStylesWithUserPreferences gets styles with user's current preferences
func (s *TravelPreferenceService) GetTravelPreferenceStylesWithUserPreferences(userID string) ([]dto.TravelPreferenceStyleResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get user's current preferences
	var preferences []models.TravelPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Create preference map (binary: true = selected, false = not selected)
	preferenceMap := make(map[string]bool)
	for _, pref := range preferences {
		preferenceMap[pref.TravelStyle] = true
	}

	// Get all styles from database
	masterStyles := s.GetTravelPreferenceStyles()

	responses := make([]dto.TravelPreferenceStyleResponse, len(masterStyles))
	for i, style := range masterStyles {
		isSelected := preferenceMap[style.Style] // true if selected, false if not

		responses[i] = dto.TravelPreferenceStyleResponse{
			Style:       style.Style,
			DisplayName: style.DisplayName,
			Icon:        style.Icon,
			IsSelected:  isSelected,
		}
	}

	return responses, nil
}

// GetTravelPreferenceStats gets travel preference statistics for a user
func (s *TravelPreferenceService) GetTravelPreferenceStats(userID string) (*dto.TravelPreferenceStatsResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get preferences
	var preferences []models.TravelPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&preferences).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	// Calculate stats
	stats := &dto.TravelPreferenceStatsResponse{
		TotalPreferences: len(preferences),
	}

	return stats, nil
}

// DeleteTravelPreference deletes a travel preference
func (s *TravelPreferenceService) DeleteTravelPreference(userID, travelStyle string) error {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Delete preference
	err = database.GetDB().Where("user_id = ? AND travel_style = ?", userUUID, travelStyle).Delete(&models.TravelPreference{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete travel preference: %w", err)
	}

	return nil
}

// isValidTravelStyle checks if travel style exists in database
func (s *TravelPreferenceService) isValidTravelStyle(style string) bool {
	// Check against database first
	var masterStyle models.TravelStyleMaster
	err := database.GetDB().
		Where("code = ? AND is_active = ?", style, true).
		First(&masterStyle).Error

	if err == nil {
		return true // Found in database
	}

	// Fallback to hardcoded validation
	return models.IsValidTravelStyle(style)
}
