package service

import (
	"fmt"
	"math"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TagService handles tag business logic
type TagService struct {
}

// NewTagService creates a new tag service
func NewTagService() *TagService {
	return &TagService{}
}

// GetTags gets tags with filtering
func (s *TagService) GetTags(page, limit int, kind string) ([]dto.TagResponse, int64, error) {
	// Build query
	query := database.GetDB().Model(&models.Tag{})

	// Apply kind filter
	if kind != "" {
		query = query.Where("kind = ?", kind)
	}

	// Get total count
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tags: %w", err)
	}

	// Get tags with pagination
	var tags []models.Tag
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tags: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = dto.TagResponse{
			ID:        tag.ID.String(),
			Name:      tag.Name,
			Kind:      tag.Kind,
			CreatedAt: tag.CreatedAt,
		}
	}

	return responses, total, nil
}

// GetUserTags gets user's tags
func (s *TagService) GetUserTags(userID string) ([]dto.TagResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get user tags
	var userTags []models.UserTag
	err = database.GetDB().Preload("Tag").Where("user_id = ?", userUUID).Find(&userTags).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user tags: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.TagResponse, len(userTags))
	for i, userTag := range userTags {
		if userTag.Tag != nil {
			responses[i] = dto.TagResponse{
				ID:        userTag.Tag.ID.String(),
				Name:      userTag.Tag.Name,
				Kind:      userTag.Tag.Kind,
				CreatedAt: userTag.Tag.CreatedAt,
			}
		}
	}

	return responses, nil
}

// AddUserTag adds a tag to user
func (s *TagService) AddUserTag(userID, tagID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID")
	}

	// Check if tag exists
	var tag models.Tag
	err = database.GetDB().Where("id = ?", tagUUID).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Check if user tag already exists
	var existingUserTag models.UserTag
	err = database.GetDB().Where("user_id = ? AND tag_id = ?", userUUID, tagUUID).First(&existingUserTag).Error
	if err == nil {
		return fmt.Errorf("tag already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user tag: %w", err)
	}

	// Create user tag
	userTag := &models.UserTag{
		UserID: userUUID,
		TagID:  tagUUID,
	}

	err = database.GetDB().Create(userTag).Error
	if err != nil {
		return fmt.Errorf("failed to add user tag: %w", err)
	}

	return nil
}

// RemoveUserTag removes a tag from user
func (s *TagService) RemoveUserTag(userID, tagID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID")
	}

	// Check if user tag exists
	var userTag models.UserTag
	err = database.GetDB().Where("user_id = ? AND tag_id = ?", userUUID, tagUUID).First(&userTag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to get user tag: %w", err)
	}

	// Delete user tag
	err = database.GetDB().Delete(&userTag).Error
	if err != nil {
		return fmt.Errorf("failed to remove user tag: %w", err)
	}

	return nil
}

// GetEventTags gets event's tags
func (s *TagService) GetEventTags(eventID string) ([]dto.TagResponse, error) {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID")
	}

	// Check if event exists
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Get event tags
	var eventTags []models.EventTag
	err = database.GetDB().Preload("Tag").Where("event_id = ?", eventUUID).Find(&eventTags).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get event tags: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.TagResponse, len(eventTags))
	for i, eventTag := range eventTags {
		if eventTag.Tag != nil {
			responses[i] = dto.TagResponse{
				ID:        eventTag.Tag.ID.String(),
				Name:      eventTag.Tag.Name,
				Kind:      eventTag.Tag.Kind,
				CreatedAt: eventTag.Tag.CreatedAt,
			}
		}
	}

	return responses, nil
}

// AddEventTag adds a tag to event
func (s *TagService) AddEventTag(eventID, tagID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID")
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Check if event exists and user is creator
	var event models.Event
	err = database.GetDB().Where("id = ? AND creator_id = ? AND deleted_at IS NULL", eventUUID, userUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Check if tag exists
	var tag models.Tag
	err = database.GetDB().Where("id = ?", tagUUID).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Check if event tag already exists
	var existingEventTag models.EventTag
	err = database.GetDB().Where("event_id = ? AND tag_id = ?", eventUUID, tagUUID).First(&existingEventTag).Error
	if err == nil {
		return fmt.Errorf("tag already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing event tag: %w", err)
	}

	// Create event tag
	eventTag := &models.EventTag{
		EventID: eventUUID,
		TagID:   tagUUID,
	}

	err = database.GetDB().Create(eventTag).Error
	if err != nil {
		return fmt.Errorf("failed to add event tag: %w", err)
	}

	return nil
}

// RemoveEventTag removes a tag from event
func (s *TagService) RemoveEventTag(eventID, tagID, userID string) error {
	// Parse IDs
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID")
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Check if event exists and user is creator
	var event models.Event
	err = database.GetDB().Where("id = ? AND creator_id = ? AND deleted_at IS NULL", eventUUID, userUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Check if event tag exists
	var eventTag models.EventTag
	err = database.GetDB().Where("event_id = ? AND tag_id = ?", eventUUID, tagUUID).First(&eventTag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to get event tag: %w", err)
	}

	// Delete event tag
	err = database.GetDB().Delete(&eventTag).Error
	if err != nil {
		return fmt.Errorf("failed to remove event tag: %w", err)
	}

	return nil
}

// GetEventSuggestions gets event suggestions based on user interests
func (s *TagService) GetEventSuggestions(userID string, page, limit int) ([]dto.EventSuggestionItem, int64, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID")
	}

	// Get user tags
	var userTags []models.UserTag
	err = database.GetDB().Preload("Tag").Where("user_id = ?", userUUID).Find(&userTags).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user tags: %w", err)
	}

	// If user has no tags, return empty result
	if len(userTags) == 0 {
		return []dto.EventSuggestionItem{}, 0, nil
	}

	// Get user tag IDs
	userTagIDs := make([]uuid.UUID, len(userTags))
	for i, userTag := range userTags {
		userTagIDs[i] = userTag.TagID
	}

	// Get events with tags that match user interests
	var events []models.Event
	err = database.GetDB().
		Preload("Creator").
		Preload("Photos").
		Preload("Categories.Tag").
		Preload("Tags.Tag").
		Preload("Members").
		Joins("JOIN event_tags ON events.id = event_tags.event_id").
		Where("events.deleted_at IS NULL AND events.status = ? AND event_tags.tag_id IN ?",
			models.EventStatusPublished, userTagIDs).
		Group("events.id").
		Order("events.created_at DESC").
		Find(&events).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	// Calculate match scores
	suggestions := make([]dto.EventSuggestionItem, len(events))
	for i, event := range events {
		// Calculate match score
		matchScore, matchedTags := s.calculateMatchScore(userTags, event.Tags)

		// Convert event to response
		eventResponse := s.convertEventToResponse(event, userID)

		suggestions[i] = dto.EventSuggestionItem{
			Event:       eventResponse,
			MatchScore:  matchScore,
			MatchedTags: matchedTags,
		}
	}

	// Sort by match score (highest first)
	for i := 0; i < len(suggestions)-1; i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[i].MatchScore < suggestions[j].MatchScore {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	// Apply pagination
	total := int64(len(suggestions))
	offset := (page - 1) * limit
	if offset >= len(suggestions) {
		return []dto.EventSuggestionItem{}, total, nil
	}

	end := offset + limit
	if end > len(suggestions) {
		end = len(suggestions)
	}

	return suggestions[offset:end], total, nil
}

// calculateMatchScore calculates match score between user tags and event tags
func (s *TagService) calculateMatchScore(userTags []models.UserTag, eventTags []models.EventTag) (float64, []dto.TagResponse) {
	// Tag kind weights (higher = more important)
	kindWeights := map[string]float64{
		"interest":      1.0, // User interests are most important
		"activity":      0.8, // Activities are very important
		"location":      0.6, // Location is important
		"food":          0.7, // Food preferences are important
		"category":      0.5, // Categories are somewhat important
		"transport":     0.3, // Transport is less important
		"accommodation": 0.4, // Accommodation is less important
	}

	// Create maps for quick lookup
	userTagMap := make(map[uuid.UUID]models.Tag)
	for _, userTag := range userTags {
		if userTag.Tag != nil {
			userTagMap[userTag.TagID] = *userTag.Tag
		}
	}

	eventTagMap := make(map[uuid.UUID]models.Tag)
	for _, eventTag := range eventTags {
		if eventTag.Tag != nil {
			eventTagMap[eventTag.TagID] = *eventTag.Tag
		}
	}

	// Calculate matches
	var totalScore float64
	var matchedTags []dto.TagResponse

	// Check for exact matches
	for userTagID, userTag := range userTagMap {
		if eventTag, exists := eventTagMap[userTagID]; exists {
			// Exact match - full weight
			weight := kindWeights[userTag.Kind]
			totalScore += weight

			matchedTags = append(matchedTags, dto.TagResponse{
				ID:        eventTag.ID.String(),
				Name:      eventTag.Name,
				Kind:      eventTag.Kind,
				CreatedAt: eventTag.CreatedAt,
			})
		}
	}

	// Normalize score to 0-100
	maxPossibleScore := 0.0
	for _, userTag := range userTagMap {
		weight := kindWeights[userTag.Kind]
		maxPossibleScore += weight
	}

	var normalizedScore float64
	if maxPossibleScore > 0 {
		normalizedScore = (totalScore / maxPossibleScore) * 100
	}

	return math.Round(normalizedScore*100) / 100, matchedTags
}

// convertEventToResponse converts event model to response DTO
func (s *TagService) convertEventToResponse(event models.Event, userID string) dto.EventResponse {
	// Convert cover image URL to public URL
	var publicCoverURL *string
	if event.CoverImageURL != nil && *event.CoverImageURL != "" {
		publicURL := fmt.Sprintf("https://api.tindertrip.phitik.com/images/events/%s", event.ID.String())
		publicCoverURL = &publicURL
	}

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
		CoverImageURL: publicCoverURL,
		CreatedAt:     event.CreatedAt,
		UpdatedAt:     event.UpdatedAt,
	}

	// Add creator
	if event.Creator != nil {
		response.Creator = &dto.UserResponse{
			ID:          event.Creator.ID.String(),
			Email:       *event.Creator.Email,
			DisplayName: *event.Creator.DisplayName,
			Provider:    string(event.Creator.Provider),
			CreatedAt:   event.Creator.CreatedAt,
		}
	}

	// Add photos
	response.Photos = make([]dto.EventPhotoResponse, len(event.Photos))
	for i, photo := range event.Photos {
		// Convert photo URL to public URL
		var publicURL string
		if photo.URL != "" {
			publicURL = fmt.Sprintf("https://api.tindertrip.phitik.com/images/events/%s", event.ID.String())
		}

		response.Photos[i] = dto.EventPhotoResponse{
			ID:        photo.ID.String(),
			EventID:   photo.EventID.String(),
			URL:       publicURL,
			SortNo:    photo.SortNo,
			CreatedAt: photo.CreatedAt,
		}
	}

	// Add categories
	response.Categories = make([]dto.TagResponse, len(event.Categories))
	for i, category := range event.Categories {
		if category.Tag != nil {
			response.Categories[i] = dto.TagResponse{
				ID:        category.Tag.ID.String(),
				Name:      category.Tag.Name,
				Kind:      category.Tag.Kind,
				CreatedAt: category.Tag.CreatedAt,
			}
		}
	}

	// Add tags
	response.Tags = make([]dto.TagResponse, len(event.Tags))
	for i, tag := range event.Tags {
		if tag.Tag != nil {
			response.Tags[i] = dto.TagResponse{
				ID:        tag.Tag.ID.String(),
				Name:      tag.Tag.Name,
				Kind:      tag.Tag.Kind,
				CreatedAt: tag.Tag.CreatedAt,
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

	// Count confirmed members
	confirmedCount := 0
	for _, member := range event.Members {
		if member.Status == models.MemberStatusConfirmed {
			confirmedCount++
		}
	}
	response.MemberCount = confirmedCount

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

	return response
}
