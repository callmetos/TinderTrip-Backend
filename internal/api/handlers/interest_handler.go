package handlers

import (
	"net/http"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InterestHandler handles interest-related requests
type InterestHandler struct {
	interestService *service.InterestService
}

// NewInterestHandler creates a new interest handler
func NewInterestHandler() *InterestHandler {
	return &InterestHandler{
		interestService: service.NewInterestService(),
	}
}

// GetAllInterests gets all active interests
// @Summary Get all interests
// @Description Get all active interests, optionally filtered by category
// @Tags interests
// @Produce json
// @Param category query string false "Filter by category (cafe, activity, pub_bar, sport)"
// @Success 200 {object} dto.InterestListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /interests [get]
func (h *InterestHandler) GetAllInterests(c *gin.Context) {
	category := c.Query("category")

	interests, total, err := h.interestService.GetAllInterests(category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get interests", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Interests retrieved successfully", dto.InterestListResponse{
		Interests: interests,
		Total:     total,
	})
}

// GetUserInterests gets all interests with user selection status
// @Summary Get user interests
// @Description Get all interests with user selection status
// @Tags interests
// @Security BearerAuth
// @Produce json
// @Param category query string false "Filter by category (cafe, activity, pub_bar, sport)"
// @Success 200 {object} dto.GetUserInterestsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/interests [get]
func (h *InterestHandler) GetUserInterests(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	category := c.Query("category")

	interests, total, err := h.interestService.GetUserInterests(userID, category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get user interests", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User interests retrieved successfully", dto.GetUserInterestsResponse{
		Interests: interests,
		Total:     total,
	})
}

// UpdateUserInterests updates user interests (bulk replace)
// @Summary Update user interests
// @Description Update user interests (bulk replace - replaces all existing selections)
// @Tags interests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserInterestsRequest true "Interest codes"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/interests [put]
func (h *InterestHandler) UpdateUserInterests(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserInterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request")
		return
	}

	if err := h.interestService.UpdateUserInterests(userID, req.InterestCodes); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user interests", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User interests updated successfully", nil)
}

// GetUserSelectedInterests gets only user's selected interests
// @Summary Get user selected interests
// @Description Get only user's selected interests
// @Tags interests
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.InterestListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/interests/selected [get]
func (h *InterestHandler) GetUserSelectedInterests(c *gin.Context) {
	userIDStr, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	interests, err := h.interestService.GetUserSelectedInterests(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get selected interests", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Selected interests retrieved successfully", dto.InterestListResponse{
		Interests: interests,
		Total:     int64(len(interests)),
	})
}
