package handlers

import (
	"net/http"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// FoodPreferenceHandler handles food preference-related requests
type FoodPreferenceHandler struct {
	foodPreferenceService *service.FoodPreferenceService
}

// NewFoodPreferenceHandler creates a new food preference handler
func NewFoodPreferenceHandler() *FoodPreferenceHandler {
	return &FoodPreferenceHandler{
		foodPreferenceService: service.NewFoodPreferenceService(),
	}
}

// GetFoodPreferences gets user's food preferences
// @Summary Get food preferences
// @Description Get the current user's food preferences
// @Tags food-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.FoodPreferenceListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences [get]
func (h *FoodPreferenceHandler) GetFoodPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get food preferences
	preferences, err := h.foodPreferenceService.GetFoodPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get food preferences",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.FoodPreferenceListResponse{
		Preferences: preferences,
	})
}

// UpdateFoodPreference updates a single food preference
// @Summary Update food preference
// @Description Update a single food preference for the current user
// @Tags food-preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateFoodPreferenceRequest true "Food preference data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences [put]
func (h *FoodPreferenceHandler) UpdateFoodPreference(c *gin.Context) {
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
	var req dto.UpdateFoodPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update food preference
	err := h.foodPreferenceService.UpdateFoodPreference(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update food preference",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Food preference updated successfully",
	})
}

// UpdateAllFoodPreferences updates all food preferences
// @Summary Update all food preferences
// @Description Update all food preferences for the current user
// @Tags food-preferences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateAllFoodPreferencesRequest true "All food preferences data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences/bulk [put]
func (h *FoodPreferenceHandler) UpdateAllFoodPreferences(c *gin.Context) {
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
	var req dto.UpdateAllFoodPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update all food preferences
	err := h.foodPreferenceService.UpdateAllFoodPreferences(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update food preferences",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Food preferences updated successfully",
	})
}

// GetFoodPreferenceCategories gets available food preference categories
// @Summary Get food preference categories
// @Description Get all available food preference categories from database (master data)
// @Tags food-preferences
// @Produce json
// @Success 200 {object} dto.FoodPreferenceCategoriesResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/food-preferences/categories [get]
func (h *FoodPreferenceHandler) GetFoodPreferenceCategories(c *gin.Context) {
	// Get categories
	categories := h.foodPreferenceService.GetFoodPreferenceCategories()

	// Send response
	c.JSON(http.StatusOK, dto.FoodPreferenceCategoriesResponse{
		Categories: categories,
	})
}

// GetFoodPreferenceCategoriesWithUserPreferences gets categories with user's current preferences
// @Summary Get food preference categories with user preferences
// @Description Get all available food preference categories with the current user's preferences
// @Tags food-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.FoodPreferenceCategoriesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences/categories [get]
func (h *FoodPreferenceHandler) GetFoodPreferenceCategoriesWithUserPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get categories with user preferences
	categories, err := h.foodPreferenceService.GetFoodPreferenceCategoriesWithUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get food preference categories",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.FoodPreferenceCategoriesResponse{
		Categories: categories,
	})
}

// GetFoodPreferenceStats gets food preference statistics
// @Summary Get food preference statistics
// @Description Get food preference statistics for the current user
// @Tags food-preferences
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.FoodPreferenceStatsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences/stats [get]
func (h *FoodPreferenceHandler) GetFoodPreferenceStats(c *gin.Context) {
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
	stats, err := h.foodPreferenceService.GetFoodPreferenceStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get food preference stats",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, stats)
}

// DeleteFoodPreference deletes a food preference
// @Summary Delete food preference
// @Description Delete a food preference for the current user
// @Tags food-preferences
// @Security BearerAuth
// @Produce json
// @Param category path string true "Food category"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/food-preferences/{category} [delete]
func (h *FoodPreferenceHandler) DeleteFoodPreference(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get food category from path
	foodCategory := c.Param("category")
	if foodCategory == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid food category",
			Message: "Food category is required",
		})
		return
	}

	// Delete food preference
	err := h.foodPreferenceService.DeleteFoodPreference(userID, foodCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete food preference",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Food preference deleted successfully",
	})
}
