package handlers

import (
	"net/http"
	"strconv"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

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
// @Success 200 {object} dto.EventListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	eventType := c.Query("event_type")
	status := c.Query("status")

	// Get user ID from context
	userID, _ := middleware.GetCurrentUserID(c)

	// Get events
	events, total, err := h.eventService.GetEvents(userID, page, limit, eventType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get events",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.EventListResponse{
		Events: events,
		Total:  total,
		Page:   page,
		Limit:  limit,
	})
}

// GetPublicEvents gets public events (no authentication required)
// @Summary Get public events
// @Description Get public events without authentication
// @Tags events
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param event_type query string false "Event type filter"
// @Success 200 {object} dto.EventListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/events [get]
func (h *EventHandler) GetPublicEvents(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	eventType := c.Query("event_type")

	// Get public events
	events, total, err := h.eventService.GetPublicEvents(page, limit, eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get events",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.EventListResponse{
		Events: events,
		Total:  total,
		Page:   page,
		Limit:  limit,
	})
}

// GetEvent gets a specific event
// @Summary Get event
// @Description Get a specific event by ID
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, _ := middleware.GetCurrentUserID(c)

	// Get event
	event, err := h.eventService.GetEvent(eventID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Event not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetPublicEvent gets a specific public event
// @Summary Get public event
// @Description Get a specific public event by ID without authentication
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/events/{id} [get]
func (h *EventHandler) GetPublicEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get public event
	event, err := h.eventService.GetPublicEvent(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Event not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// CreateEvent creates a new event
// @Summary Create event
// @Description Create a new event
// @Tags events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateEventRequest true "Event data"
// @Success 201 {object} dto.EventResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Create event
	event, err := h.eventService.CreateEvent(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, event)
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
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update event
	event, err := h.eventService.UpdateEvent(eventID, userID, req)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to update this event",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent deletes an event
// @Summary Delete event
// @Description Delete an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Delete event
	err := h.eventService.DeleteEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to delete this event",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Event deleted successfully",
	})
}

// JoinEvent joins an event
// @Summary Join event
// @Description Join an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/join [post]
func (h *EventHandler) JoinEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Join event
	err := h.eventService.JoinEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to join event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Successfully joined the event",
	})
}

// LeaveEvent leaves an event
// @Summary Leave event
// @Description Leave an existing event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/leave [post]
func (h *EventHandler) LeaveEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Leave event
	err := h.eventService.LeaveEvent(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to leave event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Successfully left the event",
	})
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
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/swipe [post]
func (h *EventHandler) SwipeEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.SwipeEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Swipe event
	err := h.eventService.SwipeEvent(eventID, userID, req.Direction)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to swipe event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Swipe recorded successfully",
	})
}
