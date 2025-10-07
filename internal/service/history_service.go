package service

import (
	"fmt"
	"time"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HistoryService handles user event history business logic
type HistoryService struct {
}

// NewHistoryService creates a new history service
func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

// GetUserEventHistory gets event history for a user
func (s *HistoryService) GetUserEventHistory(userID string, page, limit int, completed *bool) ([]dto.UserEventHistoryResponse, int64, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	// Build query
	query := database.GetDB().Model(&models.UserEventHistory{}).Where("user_id = ?", userUUID)

	// Apply completed filter
	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	// Get total count
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count history: %w", err)
	}

	// Get history with pagination
	var history []models.UserEventHistory
	offset := (page - 1) * limit
	err = query.Preload("Event").Preload("Event.Creator").Preload("Event.Photos").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get history: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.UserEventHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = s.convertHistoryToResponse(h)
	}

	return responses, total, nil
}

// MarkEventAsComplete marks an event as complete
func (s *HistoryService) MarkEventAsComplete(eventID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if event exists
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Check if user is a member
	var member models.EventMember
	err = database.GetDB().Where("event_id = ? AND user_id = ? AND status = ?", eventUUID, userUUID, models.MemberStatusConfirmed).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user is not a member of this event")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Create or update history
	history := &models.UserEventHistory{
		EventID:     eventUUID,
		UserID:      userUUID,
		Completed:   true,
		CompletedAt: &time.Time{},
	}

	// Use upsert to create or update
	err = database.GetDB().Where("event_id = ? AND user_id = ?", eventUUID, userUUID).
		Assign(models.UserEventHistory{Completed: true, CompletedAt: &time.Time{}}).
		FirstOrCreate(history).Error
	if err != nil {
		return fmt.Errorf("failed to mark event as complete: %w", err)
	}

	return nil
}

// GetEventHistory gets history for a specific event
func (s *HistoryService) GetEventHistory(eventID string, page, limit int) ([]dto.UserEventHistoryResponse, int64, error) {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid event ID: %w", err)
	}

	// Check if event exists
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("event not found")
		}
		return nil, 0, fmt.Errorf("database error: %w", err)
	}

	// Get history
	var history []models.UserEventHistory
	offset := (page - 1) * limit
	err = database.GetDB().Preload("User").Preload("User.Profile").
		Where("event_id = ?", eventUUID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event history: %w", err)
	}

	// Get total count
	var total int64
	err = database.GetDB().Model(&models.UserEventHistory{}).Where("event_id = ?", eventUUID).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event history: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.UserEventHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = s.convertHistoryToResponse(h)
	}

	return responses, total, nil
}

// GetUserStats gets user statistics
func (s *HistoryService) GetUserStats(userID string) (*dto.UserStatsResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get total events
	var totalEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).Where("user_id = ?", userUUID).Count(&totalEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count total events: %w", err)
	}

	// Get completed events
	var completedEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).Where("user_id = ? AND completed = ?", userUUID, true).Count(&completedEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count completed events: %w", err)
	}

	// Get pending events
	var pendingEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).Where("user_id = ? AND completed = ?", userUUID, false).Count(&pendingEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count pending events: %w", err)
	}

	// Get events by type
	var mealEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).
		Joins("JOIN events ON user_event_history.event_id = events.id").
		Where("user_event_history.user_id = ? AND events.event_type = ?", userUUID, models.EventTypeMeal).
		Count(&mealEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count meal events: %w", err)
	}

	var dayTripEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).
		Joins("JOIN events ON user_event_history.event_id = events.id").
		Where("user_event_history.user_id = ? AND events.event_type = ?", userUUID, models.EventTypeDaytrip).
		Count(&dayTripEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count day trip events: %w", err)
	}

	var overnightEvents int64
	err = database.GetDB().Model(&models.UserEventHistory{}).
		Joins("JOIN events ON user_event_history.event_id = events.id").
		Where("user_event_history.user_id = ? AND events.event_type = ?", userUUID, models.EventTypeOvernight).
		Count(&overnightEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count overnight events: %w", err)
	}

	// Calculate completion rate
	var completionRate float64
	if totalEvents > 0 {
		completionRate = float64(completedEvents) / float64(totalEvents) * 100
	}

	response := &dto.UserStatsResponse{
		TotalEvents:     totalEvents,
		CompletedEvents: completedEvents,
		PendingEvents:   pendingEvents,
		MealEvents:      mealEvents,
		DayTripEvents:   dayTripEvents,
		OvernightEvents: overnightEvents,
		CompletionRate:  completionRate,
	}

	return response, nil
}

// Helper function to convert history to response DTO
func (s *HistoryService) convertHistoryToResponse(history models.UserEventHistory) dto.UserEventHistoryResponse {
	response := dto.UserEventHistoryResponse{
		ID:          history.ID.String(),
		EventID:     history.EventID.String(),
		UserID:      history.UserID.String(),
		Completed:   history.Completed,
		CompletedAt: history.CompletedAt,
		CreatedAt:   history.CreatedAt,
	}

	// Add event info
	if history.Event != nil {
		// Convert cover image URL to public URL
		var publicCoverURL *string
		if history.Event.CoverImageURL != nil && *history.Event.CoverImageURL != "" {
			publicURL := fmt.Sprintf("https://api.tindertrip.phitik.com/images/events/%s", history.Event.ID.String())
			publicCoverURL = &publicURL
		}

		response.Event = &dto.EventResponse{
			ID:            history.Event.ID.String(),
			CreatorID:     history.Event.CreatorID.String(),
			Title:         history.Event.Title,
			Description:   history.Event.Description,
			EventType:     string(history.Event.EventType),
			AddressText:   history.Event.AddressText,
			Lat:           history.Event.Lat,
			Lng:           history.Event.Lng,
			StartAt:       history.Event.StartAt,
			EndAt:         history.Event.EndAt,
			Capacity:      history.Event.Capacity,
			Status:        string(history.Event.Status),
			CoverImageURL: publicCoverURL,
			CreatedAt:     history.Event.CreatedAt,
			UpdatedAt:     history.Event.UpdatedAt,
		}

		// Add creator info
		if history.Event.Creator != nil {
			response.Event.Creator = &dto.UserResponse{
				ID:          history.Event.Creator.ID.String(),
				Email:       *history.Event.Creator.Email,
				DisplayName: history.Event.Creator.GetDisplayName(),
				Provider:    string(history.Event.Creator.Provider),
				CreatedAt:   history.Event.Creator.CreatedAt,
			}
		}

		// Add photos
		response.Event.Photos = make([]dto.EventPhotoResponse, len(history.Event.Photos))
		for i, photo := range history.Event.Photos {
			// Convert photo URL to public URL
			var publicURL string
			if photo.URL != "" {
				publicURL = fmt.Sprintf("https://api.tindertrip.phitik.com/images/events/%s", history.Event.ID.String())
			}

			response.Event.Photos[i] = dto.EventPhotoResponse{
				ID:        photo.ID.String(),
				EventID:   photo.EventID.String(),
				URL:       publicURL,
				SortNo:    photo.SortNo,
				CreatedAt: photo.CreatedAt,
			}
		}
	}

	// Add user info
	if history.User != nil {
		response.User = &dto.UserResponse{
			ID:          history.User.ID.String(),
			Email:       *history.User.Email,
			DisplayName: history.User.GetDisplayName(),
			Provider:    string(history.User.Provider),
			CreatedAt:   history.User.CreatedAt,
		}

		// Add profile info
		if history.User != nil && history.User.Profile != nil {
			response.User.Profile = &dto.UserProfileResponse{
				ID:          history.User.Profile.ID.String(),
				UserID:      history.User.Profile.UserID.String(),
				Bio:         history.User.Profile.Bio,
				Languages:   history.User.Profile.Languages,
				DateOfBirth: history.User.Profile.DateOfBirth,
				Gender: func() string {
					if history.User.Profile.Gender != nil {
						return string(*history.User.Profile.Gender)
					} else {
						return ""
					}
				}(),
				JobTitle: history.User.Profile.JobTitle,
				Smoking: func() string {
					if history.User.Profile.Smoking != nil {
						return string(*history.User.Profile.Smoking)
					} else {
						return ""
					}
				}(),
				InterestsNote: history.User.Profile.InterestsNote,
				AvatarURL:     history.User.Profile.AvatarURL,
				HomeLocation:  history.User.Profile.HomeLocation,
				CreatedAt:     history.User.Profile.CreatedAt,
				UpdatedAt:     history.User.Profile.UpdatedAt,
			}
		}
	}

	return response
}

// GetHistory gets user event history
func (s *HistoryService) GetHistory(userID string, page, limit int, completed *bool) ([]dto.UserEventHistoryResponse, int64, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	// Build query
	query := database.GetDB().
		Preload("Event").
		Preload("Event.Creator").
		Preload("Event.Photos").
		Preload("User").
		Preload("User.Profile").
		Where("user_id = ?", userUUID)

	// Add completed filter if provided
	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	// Get total count
	var total int64
	err = query.Model(&models.UserEventHistory{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count history: %w", err)
	}

	// Get history with pagination
	var history []models.UserEventHistory
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get history: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.UserEventHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = s.convertHistoryToResponse(h)
	}

	return responses, total, nil
}

// MarkComplete marks an event as completed
func (s *HistoryService) MarkComplete(eventID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if history exists
	var history models.UserEventHistory
	err = database.GetDB().
		Where("event_id = ? AND user_id = ?", eventUUID, userUUID).
		First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("history not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Mark as completed
	history.MarkCompleted()

	// Update in database
	err = database.GetDB().Save(&history).Error
	if err != nil {
		return fmt.Errorf("failed to mark as complete: %w", err)
	}

	return nil
}
