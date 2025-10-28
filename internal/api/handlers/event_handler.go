package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// EventHandler handles event-related requests
type EventHandler struct {
	eventService *service.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler() *EventHandler {
	return &EventHandler{
		eventService: service.NewEventService(),
	}
}

// GetEvents gets events with pagination and filters
// @Summary Get events
// @Description Get events with pagination and filters
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param event_type query string false "Event type filter"
// @Param status query string false "Event status filter"
// @Success 200 {object} dto.EventListResponseWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	eventType := c.Query("event_type")
	status := c.Query("status")

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get user ID from context
	userID, _ := middleware.GetCurrentUserID(c)

	// Get events
	events, total, err := h.eventService.GetEvents(userID, page, limit, eventType, status)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get events", err)
		return
	}

	utils.PaginatedResponse(c, "Events retrieved successfully", events, int64(total), page, limit)
}

// GetJoinedEvents gets events that the user has joined
// @Summary Get joined events
// @Description Get events that the authenticated user has joined as a member
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Member status filter (pending, confirmed, declined)"
// @Success 200 {object} dto.EventListResponseWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /events/joined [get]
func (h *EventHandler) GetJoinedEvents(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	memberStatus := c.Query("status")

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get joined events
	events, total, err := h.eventService.GetJoinedEvents(userID, page, limit, memberStatus)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get joined events", err)
		return
	}

	utils.PaginatedResponse(c, "Joined events retrieved successfully", events, int64(total), page, limit)
}

// GetPublicEvents gets public events (no authentication required)
// @Summary Get public events
// @Description Get public events without authentication
// @Tags events
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param event_type query string false "Event type filter"
// @Success 200 {object} dto.EventListResponseWrapper
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /public/events [get]
func (h *EventHandler) GetPublicEvents(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	eventType := c.Query("event_type")

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get public events
	events, total, err := h.eventService.GetPublicEvents(page, limit, eventType)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get events", err)
		return
	}

	utils.PaginatedResponse(c, "Events retrieved successfully", events, int64(total), page, limit)
}

// GetEvent gets a specific event
// @Summary Get event
// @Description Get a specific event by ID
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 404 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, _ := middleware.GetCurrentUserID(c)

	// Get event
	event, err := h.eventService.GetEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "The requested event does not exist")
		} else {
			utils.InternalServerErrorResponse(c, "Failed to get event", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event retrieved successfully", event)
}

// GetPublicEvent gets a specific public event
// @Summary Get public event
// @Description Get a specific public event by ID without authentication
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 404 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /public/events/{id} [get]
func (h *EventHandler) GetPublicEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get public event
	event, err := h.eventService.GetPublicEvent(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "The requested event does not exist or is not public")
		} else {
			utils.InternalServerErrorResponse(c, "Failed to get event", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event retrieved successfully", event)
}

// CreateEvent creates a new event
// @Summary Create event
// @Description Create a new event with optional file uploads
// @Tags events
// @Security BearerAuth
// @Accept json,mpfd
// @Produce json
// @Param request body dto.CreateEventRequest true "Event data (JSON)"
// @Param title formData string false "Event title (multipart)"
// @Param description formData string false "Event description (multipart)"
// @Param event_type formData string false "Event type (multipart)"
// @Param address_text formData string false "Address text (multipart)"
// @Param lat formData number false "Latitude (multipart)"
// @Param lng formData number false "Longitude (multipart)"
// @Param start_at formData string false "Start time (multipart)"
// @Param end_at formData string false "End time (multipart)"
// @Param capacity formData int false "Capacity (multipart)"
// @Param category_ids formData string false "Category IDs comma separated (multipart)"
// @Param tag_ids formData string false "Tag IDs comma separated (multipart)"
// @Param file formData file false "Cover image file (multipart)"
// @Param files[] formData file false "Event photos (multipart)"
// @Success 201 {object} dto.EventResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Check content type to determine if it's multipart or JSON
	contentType := c.GetHeader("Content-Type")

	var req dto.CreateEventRequest
	var photoURLs []string

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data
		var coverImageURL *string
		var err error
		req, coverImageURL, photoURLs, err = h.parseCreateEventMultipart(c)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid request: "+err.Error())
			return
		}

		// Set cover image URL if provided
		if coverImageURL != nil {
			req.CoverImageURL = coverImageURL
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ValidationErrorResponse(c, "Invalid request", err.Error())
			return
		}
	}

	// Create event
	event, err := h.eventService.CreateEvent(userID, req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create event", err)
		return
	}

	// Add photos if provided in multipart
	if len(photoURLs) > 0 {
		if err := h.eventService.AppendEventPhotos(userID, event.ID, photoURLs); err != nil {
			// Log error but don't fail the event creation
			utils.Logger().WithField("error", err).Error("Failed to add event photos")
		}
	}

	utils.SuccessResponse(c, http.StatusCreated, "Event created successfully", event)
}

// UpdateEvent updates an event
// @Summary Update event
// @Description Update an existing event
// @Tags events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body dto.UpdateEventRequest true "Event update data"
// @Success 200 {object} dto.EventResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	// Update event
	event, err := h.eventService.UpdateEvent(eventID, userID, req)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "unauthorized" {
			utils.ForbiddenResponse(c, "You don't have permission to update this event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to update event", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event updated successfully", event)
}

// DeleteEvent deletes an event
// @Summary Delete event
// @Description Delete an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Delete event
	err := h.eventService.DeleteEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "unauthorized" {
			utils.ForbiddenResponse(c, "You don't have permission to delete this event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete event", err)
		return
	}

	utils.SendSuccessResponse(c, "Event deleted successfully", nil)
}

// JoinEvent joins an event
// @Summary Join event
// @Description Join an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/join [post]
func (h *EventHandler) JoinEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Join event
	err := h.eventService.JoinEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "The requested event does not exist")
		} else if err.Error() == "user is already a member" {
			utils.ConflictResponse(c, "You are already a member of this event")
		} else {
			utils.InternalServerErrorResponse(c, "Failed to join event", err)
		}
		return
	}

	utils.SendSuccessResponse(c, "Successfully joined the event", nil)
}

// LeaveEvent leaves an event
// @Summary Leave event
// @Description Leave an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/leave [post]
func (h *EventHandler) LeaveEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Leave event
	err := h.eventService.LeaveEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "The requested event does not exist")
		} else if err.Error() == "user is not a member" {
			utils.NotFoundResponse(c, "You are not a member of this event")
		} else {
			utils.InternalServerErrorResponse(c, "Failed to leave event", err)
		}
		return
	}

	utils.SendSuccessResponse(c, "Successfully left the event", nil)
}

// ConfirmEvent confirms participation in an event
// @Summary Confirm event participation
// @Description Confirm participation in an event (change status from pending to confirmed)
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/confirm [post]
func (h *EventHandler) ConfirmEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Confirm event participation
	err := h.eventService.ConfirmEventParticipation(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "event is full" {
			utils.ConflictResponse(c, "Cannot confirm participation. Event has reached its capacity.")
			return
		}
		if err.Error() == "member not found" {
			utils.NotFoundResponse(c, "You are not a member of this event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to confirm event participation", err)
		return
	}

	utils.SendSuccessResponse(c, "Successfully confirmed participation in the event", nil)
}

// CancelEvent cancels participation in an event
// @Summary Cancel event participation
// @Description Cancel participation in an event (change status from pending to declined)
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/cancel [post]
func (h *EventHandler) CancelEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Cancel event participation
	err := h.eventService.CancelEventParticipation(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "member not found" {
			utils.NotFoundResponse(c, "You are not a member of this event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to cancel event participation", err)
		return
	}

	utils.SendSuccessResponse(c, "Successfully cancelled participation in the event", nil)
}

// CompleteEvent completes an event (creator only)
// @Summary Complete event
// @Description Complete an event (creator only)
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/complete [post]
func (h *EventHandler) CompleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Complete event
	err := h.eventService.CompleteEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "not authorized" {
			utils.ForbiddenResponse(c, "Only the event creator can complete the event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to complete event", err)
		return
	}

	utils.SendSuccessResponse(c, "Successfully completed the event", nil)
}

// SwipeEvent swipes on an event
// @Summary Swipe event
// @Description Swipe on an event (like or pass)
// @Tags events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body dto.SwipeEventRequest true "Swipe data"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/swipe [post]
func (h *EventHandler) SwipeEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		utils.BadRequestResponse(c, "Event ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req dto.SwipeEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	// Swipe event
	err := h.eventService.SwipeEvent(eventID, userID, req.Direction)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "The requested event does not exist")
		} else {
			utils.InternalServerErrorResponse(c, "Failed to swipe event", err)
		}
		return
	}

	utils.SendSuccessResponse(c, "Swipe recorded successfully", nil)
}

// GetEventSuggestions gets event suggestions based on user interests
// @Summary Get event suggestions
// @Description Get event suggestions based on user's interests and tags
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.EventSuggestionResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/suggestions [get]
func (h *EventHandler) GetEventSuggestions(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get event suggestions
	suggestions, total, err := h.eventService.GetEventSuggestions(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get event suggestions", err)
		return
	}

	// Send response
	utils.SuccessResponse(c, http.StatusOK, "Event suggestions retrieved successfully", dto.EventSuggestionResponse{
		Events: suggestions,
		Total:  total,
		Page:   page,
		Limit:  limit,
	})
}

// UpdateCover updates event cover image (multipart: file)
// @Summary Update event cover image
// @Tags events
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param id path string true "Event ID"
// @Param file formData file true "Cover image file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/cover [put]
func (h *EventHandler) UpdateCover(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	eventID := c.Param("id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "File required")
		return
	}
	src, err := file.Open()
	if err != nil {
		utils.BadRequestResponse(c, "Invalid file")
		return
	}
	defer src.Close()

	fs, err := service.NewFileService()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Storage initialization failed", err)
		return
	}

	_, url, _, _, _, err := fs.UploadImage(c, "event_covers", file.Filename, src)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnsupportedMediaType, utils.ErrCodeInvalidInput, "Upload failed", err)
		return
	}

	if err := h.eventService.UpdateCoverImageURL(userID, eventID, &url); err != nil {
		utils.ForbiddenResponse(c, err.Error())
		return
	}
	utils.SendSuccessResponse(c, "Cover image updated successfully", gin.H{"cover_image_url": url})
}

// AddPhotos appends photos to event gallery (multipart: files[])
// @Summary Add event photos
// @Tags events
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param id path string true "Event ID"
// @Param files[] formData file true "Gallery images (multiple)"
// @Success 201 {object} map[string][]string
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /events/{id}/photos [post]
func (h *EventHandler) AddPhotos(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	eventID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil || form.File["files[]"] == nil {
		utils.BadRequestResponse(c, "files[] required")
		return
	}

	fs, err := service.NewFileService()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Storage initialization failed", err)
		return
	}

	var urls []string
	for _, f := range form.File["files[]"] {
		src, err := f.Open()
		if err != nil {
			utils.BadRequestResponse(c, "Invalid file")
			return
		}
		_, url, _, _, _, err := fs.UploadImage(c, "event_photos", f.Filename, src)
		src.Close()
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnsupportedMediaType, utils.ErrCodeInvalidInput, "Upload failed", err)
			return
		}
		urls = append(urls, url)
	}

	if err := h.eventService.AppendEventPhotos(userID, eventID, urls); err != nil {
		utils.ForbiddenResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Photos added successfully", gin.H{"urls": urls})
}

// parseCreateEventMultipart parses multipart form data for event creation
func (h *EventHandler) parseCreateEventMultipart(c *gin.Context) (dto.CreateEventRequest, *string, []string, error) {
	var req dto.CreateEventRequest
	var coverImageURL *string
	var photoURLs []string

	// Parse text fields
	title := c.PostForm("title")
	description := c.PostForm("description")
	eventType := c.PostForm("event_type")
	addressText := c.PostForm("address_text")
	latStr := c.PostForm("lat")
	lngStr := c.PostForm("lng")
	startAtStr := c.PostForm("start_at")
	endAtStr := c.PostForm("end_at")
	capacityStr := c.PostForm("capacity")
	budgetMinStr := c.PostForm("budget_min")
	budgetMaxStr := c.PostForm("budget_max")
	currency := c.PostForm("currency")
	categoryIDsStr := c.PostForm("category_ids")
	tagIDsStr := c.PostForm("tag_ids")

	// Validate required fields
	if title == "" {
		return req, nil, nil, fmt.Errorf("title is required")
	}
	if eventType == "" {
		return req, nil, nil, fmt.Errorf("event_type is required")
	}

	req.Title = title
	if description != "" {
		req.Description = &description
	}
	req.EventType = eventType
	if addressText != "" {
		req.AddressText = &addressText
	}

	// Parse latitude
	if latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			req.Lat = &lat
		}
	}

	// Parse longitude
	if lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			req.Lng = &lng
		}
	}

	// Parse start time
	if startAtStr != "" {
		if startAt, err := time.Parse(time.RFC3339, startAtStr); err == nil {
			req.StartAt = &startAt
		}
	}

	// Parse end time
	if endAtStr != "" {
		if endAt, err := time.Parse(time.RFC3339, endAtStr); err == nil {
			req.EndAt = &endAt
		}
	}

	// Parse capacity
	if capacityStr != "" {
		if capacity, err := strconv.Atoi(capacityStr); err == nil {
			req.Capacity = &capacity
		}
	}

	// Parse budget min
	if budgetMinStr != "" {
		if budgetMin, err := strconv.Atoi(budgetMinStr); err == nil {
			req.BudgetMin = &budgetMin
		}
	}

	// Parse budget max
	if budgetMaxStr != "" {
		if budgetMax, err := strconv.Atoi(budgetMaxStr); err == nil {
			req.BudgetMax = &budgetMax
		}
	}

	// Parse currency
	if currency != "" {
		req.Currency = &currency
	}

	// Parse category IDs
	if categoryIDsStr != "" {
		categoryIDs := strings.Split(categoryIDsStr, ",")
		for i, id := range categoryIDs {
			categoryIDs[i] = strings.TrimSpace(id)
		}
		req.CategoryIDs = categoryIDs
	}

	// Parse tag IDs
	if tagIDsStr != "" {
		tagIDs := strings.Split(tagIDsStr, ",")
		for i, id := range tagIDs {
			tagIDs[i] = strings.TrimSpace(id)
		}
		req.TagIDs = tagIDs
	}

	// Handle cover image upload
	if fileHeader, err := c.FormFile("file"); err == nil && fileHeader != nil {
		src, err := fileHeader.Open()
		if err != nil {
			return req, nil, nil, fmt.Errorf("invalid cover image file: %w", err)
		}
		defer src.Close()

		fs, err := service.NewFileService()
		if err != nil {
			return req, nil, nil, fmt.Errorf("storage init failed: %w", err)
		}

		_, url, _, _, _, err := fs.UploadImage(c, "event_covers", fileHeader.Filename, src)
		if err != nil {
			return req, nil, nil, fmt.Errorf("cover image upload failed: %w", err)
		}
		coverImageURL = &url
	}

	// Handle multiple photos upload
	if form, err := c.MultipartForm(); err == nil && form.File["files[]"] != nil {
		fs, err := service.NewFileService()
		if err != nil {
			return req, nil, nil, fmt.Errorf("storage init failed: %w", err)
		}

		for _, f := range form.File["files[]"] {
			src, err := f.Open()
			if err != nil {
				return req, nil, nil, fmt.Errorf("invalid photo file: %w", err)
			}

			_, url, _, _, _, err := fs.UploadImage(c, "event_photos", f.Filename, src)
			src.Close()
			if err != nil {
				return req, nil, nil, fmt.Errorf("photo upload failed: %w", err)
			}
			photoURLs = append(photoURLs, url)
		}
	}

	return req, coverImageURL, photoURLs, nil
}
