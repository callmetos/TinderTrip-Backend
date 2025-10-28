package handlers

import (
	"strconv"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

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
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /history [get]
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	completed := c.Query("completed")

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
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
		utils.InternalServerErrorResponse(c, "Failed to get history", err)
		return
	}

	utils.PaginatedResponse(c, "History retrieved successfully", history, int64(total), page, limit)
}

// MarkComplete marks an event as completed
// @Summary Mark event as completed
// @Description Mark an event as completed in user's history
// @Tags history
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /history/{id}/complete [post]
func (h *HistoryHandler) MarkComplete(c *gin.Context) {
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

	// Mark as complete
	err := h.historyService.MarkComplete(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFoundResponse(c, "Event not found")
			return
		}
		if err.Error() == "history not found" {
			utils.NotFoundResponse(c, "You haven't participated in this event")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to mark as complete", err)
		return
	}

	utils.SendSuccessResponse(c, "Event marked as completed", nil)
}
