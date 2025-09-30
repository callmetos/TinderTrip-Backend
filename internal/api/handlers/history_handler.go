package handlers

import (
	"net/http"
	"strconv"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// HistoryHandler handles history-related requests
type HistoryHandler struct {
	historyService *service.HistoryService
}

// NewHistoryHandler creates a new history handler
func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{
		historyService: service.NewHistoryService(),
	}
}

// GetHistory gets user event history
// @Summary Get event history
// @Description Get current user's event history
// @Tags history
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param completed query bool false "Filter by completion status"
// @Success 200 {object} dto.HistoryListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /history [get]
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	completed := c.Query("completed")

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Parse completed filter
	var completedFilter *bool
	if completed != "" {
		if completed == "true" {
			completedFilter = &[]bool{true}[0]
		} else if completed == "false" {
			completedFilter = &[]bool{false}[0]
		}
	}

	// Get history
	history, total, err := h.historyService.GetHistory(userID, page, limit, completedFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get history",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.HistoryListResponse{
		History: history,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// MarkComplete marks an event as completed
// @Summary Mark event as completed
// @Description Mark an event as completed in user's history
// @Tags history
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /history/{id}/complete [post]
func (h *HistoryHandler) MarkComplete(c *gin.Context) {
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

	// Mark as complete
	err := h.historyService.MarkComplete(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "history not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "History not found",
				Message: "You haven't participated in this event",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to mark as complete",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Event marked as completed",
	})
}
