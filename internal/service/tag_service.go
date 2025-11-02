package service

import (
	"fmt"
	"math"
	"strings"

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

// GetEventSuggestions gets event suggestions based on user preferences (travel, food, budget, event type)
func (s *TagService) GetEventSuggestions(userID string, page, limit int) ([]dto.EventSuggestionItem, int64, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID")
	}

	// Get travel preferences
	var travelPrefs []models.TravelPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&travelPrefs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get travel preferences: %w", err)
	}

	// Get food preferences
	var foodPrefs []models.FoodPreference
	err = database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&foodPrefs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get food preferences: %w", err)
	}

	// Get user budget preference
	var userBudget models.PrefBudget
	err = database.GetDB().Where("user_id = ?", userUUID).First(&userBudget).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, fmt.Errorf("failed to get user budget: %w", err)
	}
	hasBudget := err == nil

	// Get all published events (will be sorted by match score later)
	var events []models.Event
	err = database.GetDB().
		Preload("Creator").
		Preload("Photos").
		Preload("Categories.Tag").
		Preload("Tags.Tag").
		Preload("Members").
		Where("deleted_at IS NULL AND status = ?", models.EventStatusPublished).
		Order("created_at DESC").
		Find(&events).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	// Calculate match scores
	suggestions := make([]dto.EventSuggestionItem, len(events))
	for i, event := range events {
		// Calculate travel preference score
		travelScore := s.calculateTravelPreferenceScore(travelPrefs, event)

		// Calculate food preference score
		foodScore := s.calculateFoodPreferenceScore(foodPrefs, event)

		// Calculate event type (trip duration) score
		eventTypeScore := s.calculateEventTypeScore(userUUID, event)

		// Calculate budget match score
		var budgetScore float64
		if hasBudget {
			budgetScore = s.calculateBudgetMatchScore(userBudget, event)
		} else {
			// No budget preference = neutral score (50)
			budgetScore = 50.0
		}

		// Combined score: 30% travel + 30% food + 30% budget + 10% event type
		// Weights can be adjusted based on importance
		combinedScore := (travelScore * 0.3) + (foodScore * 0.3) + (budgetScore * 0.3) + (eventTypeScore * 0.1)

		// Convert event to response
		eventResponse := s.convertEventToResponse(event, userID)

		// Get matched tags for display (from event tags)
		matchedTags := make([]dto.TagResponse, len(event.Tags))
		for j, eventTag := range event.Tags {
			if eventTag.Tag != nil {
				matchedTags[j] = dto.TagResponse{
					ID:        eventTag.Tag.ID.String(),
					Name:      eventTag.Tag.Name,
					Kind:      eventTag.Tag.Kind,
					CreatedAt: eventTag.Tag.CreatedAt,
				}
			}
		}

		suggestions[i] = dto.EventSuggestionItem{
			Event:       eventResponse,
			MatchScore:  math.Round(combinedScore*100) / 100,
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

// calculateTagMatchScore calculates match score between user tags and event tags
func (s *TagService) calculateTagMatchScore(userTags []models.UserTag, eventTags []models.EventTag) (float64, []dto.TagResponse) {
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

// calculateBudgetMatchScore calculates match score between user budget preference and event budget
func (s *TagService) calculateBudgetMatchScore(userBudget models.PrefBudget, event models.Event) float64 {
	// If user has unlimited budget, match all events
	if userBudget.Unlimited {
		return 100.0
	}

	// If event has no budget, give neutral score
	if event.BudgetMin == nil && event.BudgetMax == nil {
		return 50.0
	}

	// Get user budget range for event type
	userMin, userMax := userBudget.GetBudgetForEventType(event.EventType)

	// If user has no budget preference for this event type, give neutral score
	if userMin == nil && userMax == nil {
		return 50.0
	}

	// Calculate event budget range
	eventMin := 0
	eventMax := 0
	if event.BudgetMin != nil {
		eventMin = *event.BudgetMin
	}
	if event.BudgetMax != nil {
		eventMax = *event.BudgetMax
	} else {
		// If event has only min, assume reasonable max (2x min or 100k)
		eventMax = eventMin * 2
		if eventMax < 100000 {
			eventMax = 100000
		}
	}

	// Normalize to same currency (for simplicity, assume same currency)
	// In production, you might want to add currency conversion here

	// Calculate overlap and match score
	var userMinVal, userMaxVal int
	if userMin != nil {
		userMinVal = *userMin
	} else {
		userMinVal = 0
	}
	if userMax != nil {
		userMaxVal = *userMax
	} else {
		userMaxVal = 1000000 // Large number for unlimited max
	}

	// Check if ranges overlap
	overlapMin := math.Max(float64(userMinVal), float64(eventMin))
	overlapMax := math.Min(float64(userMaxVal), float64(eventMax))

	if overlapMin > overlapMax {
		// No overlap - calculate distance-based score
		distance := 0.0
		if overlapMin > float64(eventMax) {
			distance = overlapMin - float64(eventMax)
		} else {
			distance = float64(eventMin) - overlapMax
		}

		// Calculate score based on distance (penalty)
		// Use percentage of user budget range as reference
		budgetRange := float64(userMaxVal - userMinVal)
		if budgetRange <= 0 {
			budgetRange = 10000 // Default range if user has no range
		}

		penaltyPercent := (distance / budgetRange) * 100
		if penaltyPercent > 50 {
			return 0.0 // Too far outside budget
		}
		return 50.0 - (penaltyPercent * 0.6) // Penalty: 0-30 points
	}

	// Ranges overlap - calculate overlap percentage
	overlapRange := overlapMax - overlapMin
	eventRange := float64(eventMax - eventMin)
	userRange := float64(userMaxVal - userMinVal)

	if eventRange <= 0 || userRange <= 0 {
		return 100.0 // Exact match or edge case
	}

	// Calculate how much of the event budget overlaps with user budget
	overlapPercent := (overlapRange / eventRange) * 100

	// Calculate score: base score from overlap, bonus for being within user range
	score := 70.0 + (overlapPercent * 0.3) // 70-100 points for overlap

	// Bonus if event budget is completely within user budget
	if float64(eventMin) >= float64(userMinVal) && float64(eventMax) <= float64(userMaxVal) {
		score = 100.0 // Perfect match
	}

	return math.Round(score*100) / 100
}

// calculateTravelPreferenceScore calculates match score between travel preferences and event tags
func (s *TagService) calculateTravelPreferenceScore(travelPrefs []models.TravelPreference, event models.Event) float64 {
	if len(travelPrefs) == 0 {
		// No travel preferences = neutral score (50)
		return 50.0
	}

	// Map travel styles to potential tag names/activities
	travelStyleMap := map[string][]string{
		"outdoor_activity": {"fitness", "camping", "hiking", "outdoor", "sports"},
		"social_activity":  {"social", "meetup", "gathering", "party"},
		"karaoke":          {"karaoke", "singing", "music"},
		"gaming":           {"gaming", "games", "esports"},
		"movie":            {"movie", "cinema", "film"},
		"board_game":       {"board game", "games", "tabletop"},
		"swimming":         {"swimming", "pool", "water"},
		"skateboarding":    {"skateboarding", "skate", "extreme"},
		"cafe_dessert":     {"cafe", "dessert", "coffee", "bakery"},
		"bubble_tea":       {"bubble tea", "tea", "drinks"},
		"bakery_cake":      {"bakery", "cake", "pastry"},
		"bingsu_ice_cream": {"bingsu", "ice cream", "dessert"},
		"coffee":           {"coffee", "cafe"},
		"matcha":           {"matcha", "tea"},
		"pancakes":         {"pancakes", "breakfast", "brunch"},
		"party_celebration": {"party", "celebration", "event"},
	}

	// Create map of user travel styles
	userTravelStyles := make(map[string]bool)
	for _, pref := range travelPrefs {
		userTravelStyles[pref.TravelStyle] = true
	}

	// Check if event tags match any travel preferences
	var matches int
	var totalChecks int

	for _, eventTag := range event.Tags {
		if eventTag.Tag == nil {
			continue
		}

		tagNameLower := strings.ToLower(eventTag.Tag.Name)

		// Check each travel style mapping
		for travelStyle, keywords := range travelStyleMap {
			if userTravelStyles[travelStyle] {
				totalChecks++
				// Check if tag name contains any keyword
				for _, keyword := range keywords {
					if strings.Contains(tagNameLower, strings.ToLower(keyword)) {
						matches++
						break
					}
				}
			}
		}
	}

	// Calculate score: matches / total checks * 100
	if totalChecks == 0 {
		return 50.0 // No preferences to check
	}

	score := (float64(matches) / float64(totalChecks)) * 100
	if score < 50 {
		// Penalize if no matches found
		score = 50.0 - ((50.0 - score) * 0.5) // Reduce penalty
	}

	return math.Round(score*100) / 100
}

// calculateFoodPreferenceScore calculates match score between food preferences and event tags
func (s *TagService) calculateFoodPreferenceScore(foodPrefs []models.FoodPreference, event models.Event) float64 {
	if len(foodPrefs) == 0 {
		// No food preferences = neutral score (50)
		return 50.0
	}

	// Map food categories to potential tag names
	foodCategoryMap := map[string][]string{
		"thai_food":          {"thai", "thailand", "pad thai", "tom yum"},
		"japanese_food":      {"japanese", "japan", "sushi", "ramen"},
		"chinese_food":       {"chinese", "china", "dim sum", "dumpling"},
		"international_food": {"international", "western", "global"},
		"halal_food":         {"halal", "muslim", "islamic"},
		"buffet":             {"buffet", "all you can eat"},
		"bbq_grill":          {"bbq", "barbecue", "grill", "grilled"},
	}

	// Create map of user food preferences with levels
	userFoodPrefs := make(map[string]int) // food_category -> preference_level
	for _, pref := range foodPrefs {
		userFoodPrefs[pref.FoodCategory] = pref.PreferenceLevel
	}

	// Check if event tags match any food preferences
	var totalScore float64
	var matchCount int

	for _, eventTag := range event.Tags {
		if eventTag.Tag == nil {
			continue
		}

		// Check if tag is food-related
		if eventTag.Tag.Kind != "food" {
			continue
		}

		tagNameLower := strings.ToLower(eventTag.Tag.Name)

		// Check each food category mapping
		for foodCategory, keywords := range foodCategoryMap {
			if prefLevel, exists := userFoodPrefs[foodCategory]; exists {
				// Check if tag name contains any keyword
				for _, keyword := range keywords {
					if strings.Contains(tagNameLower, strings.ToLower(keyword)) {
						matchCount++
						// Calculate score based on preference level
						// 1 = dislike (0-20 points), 2 = neutral (40-60 points), 3 = love (80-100 points)
						switch prefLevel {
						case 1: // Dislike
							totalScore += 10.0 // Low score
						case 2: // Neutral
							totalScore += 50.0 // Neutral score
						case 3: // Love
							totalScore += 90.0 // High score
						default:
							totalScore += 50.0
						}
						break
					}
				}
			}
		}
	}

	// If no matches, check if user dislikes everything (negative preference)
	if matchCount == 0 {
		// Check if user has mostly dislike preferences
		dislikeCount := 0
		loveCount := 0
		for _, pref := range foodPrefs {
			if pref.PreferenceLevel == 1 {
				dislikeCount++
			} else if pref.PreferenceLevel == 3 {
				loveCount++
			}
		}
		if dislikeCount > loveCount {
			// User dislikes more than loves, neutral score for events without food
			return 50.0
		}
	}

	// Normalize score
	if matchCount == 0 {
		return 50.0 // No food tags found or no matches
	}

	avgScore := totalScore / float64(matchCount)
	return math.Round(avgScore*100) / 100
}

// calculateEventTypeScore calculates match score based on event type (trip duration preference)
func (s *TagService) calculateEventTypeScore(userUUID uuid.UUID, event models.Event) float64 {
	// For now, we use neutral score (50) as we don't have explicit event type preferences
	// In the future, we could add:
	// - User's preferred event types
	// - Historical event types user joined
	// - Event type matching with budget preferences (already handled in budget score)

	// Check if user has budget preference for this event type (indirect preference)
	var userBudget models.PrefBudget
	err := database.GetDB().Where("user_id = ?", userUUID).First(&userBudget).Error
	if err == nil {
		// If user has budget preference for this event type, it's a positive signal
		_, max := userBudget.GetBudgetForEventType(event.EventType)
		if max != nil {
			return 70.0 // User has budget preference for this event type = preference
		}
	}

	// Default: neutral score
	return 50.0
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
		BudgetMin:     event.BudgetMin,
		BudgetMax:     event.BudgetMax,
		Currency:      event.Currency,
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
