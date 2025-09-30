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

// EventService handles event business logic
type EventService struct {
}

// NewEventService creates a new event service
func NewEventService() *EventService {
	return &EventService{}
}

// GetEvents gets events with pagination and filters
func (s *EventService) GetEvents(userID string, page, limit int, eventType, status string) ([]dto.EventResponse, int64, error) {
	// Build query
	query := database.GetDB().Model(&models.Event{}).Where("deleted_at IS NULL")

	// Apply filters
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	// Get events with pagination
	var events []models.Event
	offset := (page - 1) * limit
	err = query.Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = s.convertEventToResponse(event, userID)
	}

	return responses, total, nil
}

// GetPublicEvents gets public events (no authentication required)
func (s *EventService) GetPublicEvents(page, limit int, eventType string) ([]dto.EventResponse, int64, error) {
	// Build query for active events only
	query := database.GetDB().Model(&models.Event{}).Where("deleted_at IS NULL AND status = ?", models.EventStatusPublished)

	// Apply filters
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	// Get total count
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	// Get events with pagination
	var events []models.Event
	offset := (page - 1) * limit
	err = query.Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = s.convertEventToResponse(event, "")
	}

	return responses, total, nil
}

// GetEvent gets a specific event
func (s *EventService) GetEvent(eventID, userID string) (*dto.EventResponse, error) {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event
	var event models.Event
	err = database.GetDB().Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	response := s.convertEventToResponse(event, userID)
	return &response, nil
}

// GetPublicEvent gets a specific public event
func (s *EventService) GetPublicEvent(eventID string) (*dto.EventResponse, error) {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event (only active events)
	var event models.Event
	err = database.GetDB().Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Where("id = ? AND deleted_at IS NULL AND status = ?", eventUUID, models.EventStatusPublished).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	response := s.convertEventToResponse(event, "")
	return &response, nil
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(userID string, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create event
	event := &models.Event{
		CreatorID:     userUUID,
		Title:         req.Title,
		Description:   req.Description,
		EventType:     models.EventType(req.EventType),
		AddressText:   req.AddressText,
		Lat:           req.Lat,
		Lng:           req.Lng,
		StartAt:       req.StartAt,
		EndAt:         req.EndAt,
		Capacity:      req.Capacity,
		Status:        models.EventStatusDraft,
		CoverImageURL: req.CoverImageURL,
	}

	// Save event
	err = database.GetDB().Create(event).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Add creator as member
	member := &models.EventMember{
		EventID: event.ID,
		UserID:  userUUID,
		Role:    models.MemberRoleCreator,
		Status:  models.MemberStatusConfirmed,
	}
	err = database.GetDB().Create(member).Error
	if err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	// Create chat room
	chatRoom := &models.ChatRoom{
		EventID: event.ID,
	}
	err = database.GetDB().Create(chatRoom).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create chat room: %w", err)
	}

	// Load event with relationships
	err = database.GetDB().Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Where("id = ?", event.ID).First(event).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load event: %w", err)
	}

	response := s.convertEventToResponse(*event, userID)
	return &response, nil
}

// UpdateEvent updates an event
func (s *EventService) UpdateEvent(eventID, userID string, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if event exists and user is creator
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if event.CreatorID != userUUID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.EventType != nil {
		updates["event_type"] = *req.EventType
	}
	if req.AddressText != nil {
		updates["address_text"] = *req.AddressText
	}
	if req.Lat != nil {
		updates["lat"] = *req.Lat
	}
	if req.Lng != nil {
		updates["lng"] = *req.Lng
	}
	if req.StartAt != nil {
		updates["start_at"] = *req.StartAt
	}
	if req.EndAt != nil {
		updates["end_at"] = *req.EndAt
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.CoverImageURL != nil {
		updates["cover_image_url"] = *req.CoverImageURL
	}

	// Update event
	err = database.GetDB().Model(&event).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	// Load updated event with relationships
	err = database.GetDB().Preload("Creator").Preload("Photos").Preload("Categories.Tag").Preload("Tags.Tag").
		Where("id = ?", event.ID).First(&event).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load event: %w", err)
	}

	response := s.convertEventToResponse(event, userID)
	return &response, nil
}

// DeleteEvent deletes an event
func (s *EventService) DeleteEvent(eventID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if event exists and user is creator
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if event.CreatorID != userUUID {
		return fmt.Errorf("unauthorized")
	}

	// Soft delete event
	now := time.Now()
	err = database.GetDB().Model(&event).Update("deleted_at", now).Error
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// JoinEvent joins an event
func (s *EventService) JoinEvent(eventID, userID string) error {
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

	// Check if user is already a member
	var existingMember models.EventMember
	err = database.GetDB().Where("event_id = ? AND user_id = ?", eventUUID, userUUID).First(&existingMember).Error
	if err == nil {
		return fmt.Errorf("user is already a member")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("database error: %w", err)
	}

	// Create member
	member := &models.EventMember{
		EventID: eventUUID,
		UserID:  userUUID,
		Role:    models.MemberRoleParticipant,
		Status:  models.MemberStatusPending,
	}

	err = database.GetDB().Create(member).Error
	if err != nil {
		return fmt.Errorf("failed to join event: %w", err)
	}

	return nil
}

// LeaveEvent leaves an event
func (s *EventService) LeaveEvent(eventID, userID string) error {
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
	err = database.GetDB().Where("event_id = ? AND user_id = ?", eventUUID, userUUID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user is not a member")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Update member status
	now := time.Now()
	err = database.GetDB().Model(&member).Updates(map[string]interface{}{
		"status":  models.MemberStatusLeft,
		"left_at": now,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to leave event: %w", err)
	}

	return nil
}

// SwipeEvent swipes on an event
func (s *EventService) SwipeEvent(eventID, userID, direction string) error {
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

	// Create or update swipe
	swipe := &models.EventSwipe{
		UserID:    userUUID,
		EventID:   eventUUID,
		Direction: models.SwipeDirection(direction),
	}

	// Use upsert to create or update
	err = database.GetDB().Where("user_id = ? AND event_id = ?", userUUID, eventUUID).
		Assign(models.EventSwipe{Direction: models.SwipeDirection(direction)}).
		FirstOrCreate(swipe).Error
	if err != nil {
		return fmt.Errorf("failed to swipe event: %w", err)
	}

	return nil
}

// Helper function to convert event to response DTO
func (s *EventService) convertEventToResponse(event models.Event, userID string) dto.EventResponse {
	response := dto.EventResponse{
		ID:            event.ID.String(),
		CreatorID:     event.CreatorID.String(),
		Title:         event.Title,
		Description:   event.Description,
		EventType:     string(event.EventType),
		AddressText:   event.AddressText,
		Lat:           event.Lat,
		Lng:           event.Lng,
		StartAt:       event.StartAt,
		EndAt:         event.EndAt,
		Capacity:      event.Capacity,
		Status:        string(event.Status),
		CoverImageURL: event.CoverImageURL,
		CreatedAt:     event.CreatedAt,
		UpdatedAt:     event.UpdatedAt,
	}

	// Add creator info
	if event.Creator != nil {
		response.Creator = &dto.UserResponse{
			ID:          event.Creator.ID.String(),
			Email:       *event.Creator.Email,
			DisplayName: event.Creator.GetDisplayName(),
			Provider:    string(event.Creator.Provider),
			CreatedAt:   event.Creator.CreatedAt,
		}
	}

	// Add photos
	response.Photos = make([]dto.EventPhotoResponse, len(event.Photos))
	for i, photo := range event.Photos {
		response.Photos[i] = dto.EventPhotoResponse{
			ID:        photo.ID.String(),
			EventID:   photo.EventID.String(),
			URL:       photo.URL,
			SortNo:    photo.SortNo,
			CreatedAt: photo.CreatedAt,
		}
	}

	// Add categories
	response.Categories = make([]dto.TagResponse, len(event.Categories))
	for i, category := range event.Categories {
		if category.Tag != nil {
			response.Categories[i] = dto.TagResponse{
				ID:   category.Tag.ID.String(),
				Name: category.Tag.Name,
				Kind: category.Tag.Kind,
			}
		}
	}

	// Add tags
	response.Tags = make([]dto.TagResponse, len(event.Tags))
	for i, tag := range event.Tags {
		if tag.Tag != nil {
			response.Tags[i] = dto.TagResponse{
				ID:   tag.Tag.ID.String(),
				Name: tag.Tag.Name,
				Kind: tag.Tag.Kind,
			}
		}
	}

	// Add members
	response.Members = make([]dto.EventMemberResponse, len(event.Members))
	for i, member := range event.Members {
		response.Members[i] = dto.EventMemberResponse{
			EventID:     member.EventID.String(),
			UserID:      member.UserID.String(),
			Role:        string(member.Role),
			Status:      string(member.Status),
			JoinedAt:    member.JoinedAt,
			ConfirmedAt: member.ConfirmedAt,
			LeftAt:      member.LeftAt,
			Note:        member.Note,
		}
	}

	response.MemberCount = len(event.Members)

	// Check if user is joined
	if userID != "" {
		userUUID, err := uuid.Parse(userID)
		if err == nil {
			for _, member := range event.Members {
				if member.UserID == userUUID && member.Status == models.MemberStatusConfirmed {
					response.IsJoined = true
					break
				}
			}
		}
	}

	// Add user swipe
	if userID != "" {
		userUUID, err := uuid.Parse(userID)
		if err == nil {
			for _, swipe := range event.Swipes {
				if swipe.UserID == userUUID {
					response.UserSwipe = &dto.EventSwipeResponse{
						UserID:    swipe.UserID.String(),
						EventID:   swipe.EventID.String(),
						Direction: string(swipe.Direction),
						CreatedAt: swipe.CreatedAt,
					}
					break
				}
			}
		}
	}

	return response
}
