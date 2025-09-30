package handlers

import (
	"net/http"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// PreferenceHandler handles user preference requests
type PreferenceHandler struct {
	preferenceService *service.PreferenceService
}

// NewPreferenceHandler creates a new preference handler
func NewPreferenceHandler() *PreferenceHandler {
	return &PreferenceHandler{
		preferenceService: service.NewPreferenceService(),
	}
}

// GetAvailability gets user availability preferences
// @Summary Get user availability preferences
// @Description Get current user's availability preferences
// @Tags preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.PrefAvailabilityResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/preferences/availability [get]
func (h *PreferenceHandler) GetAvailability(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get availability preferences
	availability, err := h.preferenceService.GetAvailability(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Availability not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// UpdateAvailability updates user availability preferences
// @Summary Update user availability preferences
// @Description Update current user's availability preferences
// @Tags preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdatePrefAvailabilityRequest true "Availability preferences"
// @Success 200 {object} dto.PrefAvailabilityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/preferences/availability [put]
func (h *PreferenceHandler) UpdateAvailability(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.UpdatePrefAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update availability preferences
	availability, err := h.preferenceService.UpdateAvailability(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Update failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// GetBudget gets user budget preferences
// @Summary Get user budget preferences
// @Description Get current user's budget preferences
// @Tags preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.PrefBudgetResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/preferences/budget [get]
func (h *PreferenceHandler) GetBudget(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get budget preferences
	budget, err := h.preferenceService.GetBudget(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Budget not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// UpdateBudget updates user budget preferences
// @Summary Update user budget preferences
// @Description Update current user's budget preferences
// @Tags preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdatePrefBudgetRequest true "Budget preferences"
// @Success 200 {object} dto.PrefBudgetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/preferences/budget [put]
func (h *PreferenceHandler) UpdateBudget(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.UpdatePrefBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update budget preferences
	budget, err := h.preferenceService.UpdateBudget(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Update failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, budget)
}
