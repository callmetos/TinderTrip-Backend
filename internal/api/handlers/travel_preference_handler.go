package handlers

import (
	"net/http"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// TravelPreferenceHandler handles travel preference-related requests
type TravelPreferenceHandler struct {
	travelPreferenceService *service.TravelPreferenceService
}

// NewTravelPreferenceHandler creates a new travel preference handler
func NewTravelPreferenceHandler() *TravelPreferenceHandler {
	return &TravelPreferenceHandler{
		travelPreferenceService: service.NewTravelPreferenceService(),
	}
}

// GetTravelPreferences gets user's travel preferences
// @Summary Get travel preferences
// @Description Get the current user's travel preferences
// @Tags travel-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TravelPreferenceListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences [get]
func (h *TravelPreferenceHandler) GetTravelPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get travel preferences
	preferences, err := h.travelPreferenceService.GetTravelPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get travel preferences",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.TravelPreferenceListResponse{
		Preferences: preferences,
	})
}

// AddTravelPreference adds a travel preference
// @Summary Add travel preference
// @Description Add a travel preference for the current user
// @Tags travel-preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AddTravelPreferenceRequest true "Travel preference data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences [post]
func (h *TravelPreferenceHandler) AddTravelPreference(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Parse request
	var req dto.AddTravelPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Add travel preference
	err := h.travelPreferenceService.AddTravelPreference(userID, req)
	if err != nil {
		if err.Error() == "travel preference already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Travel preference already exists",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to add travel preference",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Travel preference added successfully",
	})
}

// UpdateAllTravelPreferences updates all travel preferences
// @Summary Update all travel preferences
// @Description Update all travel preferences for the current user (replace all)
// @Tags travel-preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateAllTravelPreferencesRequest true "All travel preferences data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences/bulk [put]
func (h *TravelPreferenceHandler) UpdateAllTravelPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Parse request
	var req dto.UpdateAllTravelPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update all travel preferences
	err := h.travelPreferenceService.UpdateAllTravelPreferences(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update travel preferences",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Travel preferences updated successfully",
	})
}

// GetTravelPreferenceStyles gets available travel preference styles
// @Summary Get travel preference styles
// @Description Get all available travel preference styles from database (master data)
// @Tags travel-preferences
// @Produce json
// @Success 200 {object} dto.TravelPreferenceStylesResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/travel-preferences/styles [get]
func (h *TravelPreferenceHandler) GetTravelPreferenceStyles(c *gin.Context) {
	// Get styles
	styles := h.travelPreferenceService.GetTravelPreferenceStyles()

	// Send response
	c.JSON(http.StatusOK, dto.TravelPreferenceStylesResponse{
		Styles: styles,
	})
}

// GetTravelPreferenceStylesWithUserPreferences gets styles with user's current preferences
// @Summary Get travel preference styles with user preferences
// @Description Get all available travel preference styles with the current user's preferences
// @Tags travel-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TravelPreferenceStylesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences/styles [get]
func (h *TravelPreferenceHandler) GetTravelPreferenceStylesWithUserPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get styles with user preferences
	styles, err := h.travelPreferenceService.GetTravelPreferenceStylesWithUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get travel preference styles",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.TravelPreferenceStylesResponse{
		Styles: styles,
	})
}

// GetTravelPreferenceStats gets travel preference statistics
// @Summary Get travel preference statistics
// @Description Get travel preference statistics for the current user
// @Tags travel-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.TravelPreferenceStatsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences/stats [get]
func (h *TravelPreferenceHandler) GetTravelPreferenceStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get stats
	stats, err := h.travelPreferenceService.GetTravelPreferenceStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get travel preference stats",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, stats)
}

// DeleteTravelPreference deletes a travel preference
// @Summary Delete travel preference
// @Description Delete a travel preference for the current user
// @Tags travel-preferences
// @Security BearerAuth
// @Produce json
// @Param style path string true "Travel style"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/travel-preferences/{style} [delete]
func (h *TravelPreferenceHandler) DeleteTravelPreference(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get travel style from path
	travelStyle := c.Param("style")
	if travelStyle == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid travel style",
			Message: "Travel style is required",
		})
		return
	}

	// Delete travel preference
	err := h.travelPreferenceService.DeleteTravelPreference(userID, travelStyle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete travel preference",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Travel preference deleted successfully",
	})
}
