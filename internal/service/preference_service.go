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

// PreferenceService handles user preference business logic
type PreferenceService struct{}

// NewPreferenceService creates a new preference service
func NewPreferenceService() *PreferenceService {
	return &PreferenceService{}
}

// GetAvailability gets user availability preferences
func (s *PreferenceService) GetAvailability(userID string) (*dto.PrefAvailabilityResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var availability models.PrefAvailability
	err = database.GetDB().Where("user_id = ?", userUUID).First(&availability).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default availability if not exists
			availability = models.PrefAvailability{
				UserID:    userUUID,
				AllDay:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = database.GetDB().Create(&availability).Error
			if err != nil {
				return nil, fmt.Errorf("failed to create default availability: %w", err)
			}
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	// Convert to response DTO
	response := &dto.PrefAvailabilityResponse{
		ID:        availability.ID.String(),
		UserID:    availability.UserID.String(),
		Mon:       availability.Mon,
		Tue:       availability.Tue,
		Wed:       availability.Wed,
		Thu:       availability.Thu,
		Fri:       availability.Fri,
		Sat:       availability.Sat,
		Sun:       availability.Sun,
		AllDay:    availability.AllDay,
		Morning:   availability.Morning,
		Afternoon: availability.Afternoon,
		TimeRange: availability.TimeRange,
		CreatedAt: availability.CreatedAt,
		UpdatedAt: availability.UpdatedAt,
	}

	return response, nil
}

// UpdateAvailability updates user availability preferences
func (s *PreferenceService) UpdateAvailability(userID string, req dto.UpdatePrefAvailabilityRequest) (*dto.PrefAvailabilityResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var availability models.PrefAvailability
	err = database.GetDB().Where("user_id = ?", userUUID).First(&availability).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new availability
			availability = models.PrefAvailability{
				UserID:    userUUID,
				CreatedAt: time.Now(),
			}
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	// Update fields
	if req.Mon != nil {
		availability.Mon = *req.Mon
	}
	if req.Tue != nil {
		availability.Tue = *req.Tue
	}
	if req.Wed != nil {
		availability.Wed = *req.Wed
	}
	if req.Thu != nil {
		availability.Thu = *req.Thu
	}
	if req.Fri != nil {
		availability.Fri = *req.Fri
	}
	if req.Sat != nil {
		availability.Sat = *req.Sat
	}
	if req.Sun != nil {
		availability.Sun = *req.Sun
	}
	if req.AllDay != nil {
		availability.AllDay = *req.AllDay
	}
	if req.Morning != nil {
		availability.Morning = *req.Morning
	}
	if req.Afternoon != nil {
		availability.Afternoon = *req.Afternoon
	}
	if req.TimeRange != nil {
		availability.TimeRange = req.TimeRange
	}

	availability.UpdatedAt = time.Now()

	// Save to database
	err = database.GetDB().Save(&availability).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update availability: %w", err)
	}

	// Convert to response DTO
	response := &dto.PrefAvailabilityResponse{
		ID:        availability.ID.String(),
		UserID:    availability.UserID.String(),
		Mon:       availability.Mon,
		Tue:       availability.Tue,
		Wed:       availability.Wed,
		Thu:       availability.Thu,
		Fri:       availability.Fri,
		Sat:       availability.Sat,
		Sun:       availability.Sun,
		AllDay:    availability.AllDay,
		Morning:   availability.Morning,
		Afternoon: availability.Afternoon,
		TimeRange: availability.TimeRange,
		CreatedAt: availability.CreatedAt,
		UpdatedAt: availability.UpdatedAt,
	}

	return response, nil
}

// GetBudget gets user budget preferences
func (s *PreferenceService) GetBudget(userID string) (*dto.PrefBudgetResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var budget models.PrefBudget
	err = database.GetDB().Where("user_id = ?", userUUID).First(&budget).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default budget if not exists
			budget = models.PrefBudget{
				UserID:    userUUID,
				Currency:  "THB",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = database.GetDB().Create(&budget).Error
			if err != nil {
				return nil, fmt.Errorf("failed to create default budget: %w", err)
			}
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	// Convert to response DTO
	response := &dto.PrefBudgetResponse{
		ID:           budget.ID.String(),
		UserID:       budget.UserID.String(),
		MealMin:      budget.MealMin,
		MealMax:      budget.MealMax,
		DaytripMin:   budget.DaytripMin,
		DaytripMax:   budget.DaytripMax,
		OvernightMin: budget.OvernightMin,
		OvernightMax: budget.OvernightMax,
		Unlimited:    budget.Unlimited,
		Currency:     budget.Currency,
		CreatedAt:    budget.CreatedAt,
		UpdatedAt:    budget.UpdatedAt,
	}

	return response, nil
}

// UpdateBudget updates user budget preferences
func (s *PreferenceService) UpdateBudget(userID string, req dto.UpdatePrefBudgetRequest) (*dto.PrefBudgetResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var budget models.PrefBudget
	err = database.GetDB().Where("user_id = ?", userUUID).First(&budget).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new budget
			budget = models.PrefBudget{
				UserID:    userUUID,
				Currency:  "THB",
				CreatedAt: time.Now(),
			}
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	// Update fields
	// For Meal budget
	if req.MealMin != nil {
		budget.MealMin = req.MealMin
	}
	if req.MealMax != nil {
		budget.MealMax = req.MealMax
		// If max is set but min is not provided in request, set min to 0
		if req.MealMin == nil {
			minValue := 0
			budget.MealMin = &minValue
		}
	}
	// For Daytrip budget
	if req.DaytripMin != nil {
		budget.DaytripMin = req.DaytripMin
	}
	if req.DaytripMax != nil {
		budget.DaytripMax = req.DaytripMax
		// If max is set but min is not provided in request, set min to 0
		if req.DaytripMin == nil {
			minValue := 0
			budget.DaytripMin = &minValue
		}
	}
	// For Overnight budget
	if req.OvernightMin != nil {
		budget.OvernightMin = req.OvernightMin
	}
	if req.OvernightMax != nil {
		budget.OvernightMax = req.OvernightMax
		// If max is set but min is not provided in request, set min to 0
		if req.OvernightMin == nil {
			minValue := 0
			budget.OvernightMin = &minValue
		}
	}
	if req.Unlimited != nil {
		budget.Unlimited = *req.Unlimited
	}
	if req.Currency != nil {
		budget.Currency = *req.Currency
	}

	budget.UpdatedAt = time.Now()

	// Save to database
	err = database.GetDB().Save(&budget).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	// Convert to response DTO
	response := &dto.PrefBudgetResponse{
		ID:           budget.ID.String(),
		UserID:       budget.UserID.String(),
		MealMin:      budget.MealMin,
		MealMax:      budget.MealMax,
		DaytripMin:   budget.DaytripMin,
		DaytripMax:   budget.DaytripMax,
		OvernightMin: budget.OvernightMin,
		OvernightMax: budget.OvernightMax,
		Unlimited:    budget.Unlimited,
		Currency:     budget.Currency,
		CreatedAt:    budget.CreatedAt,
		UpdatedAt:    budget.UpdatedAt,
	}

	return response, nil
}
