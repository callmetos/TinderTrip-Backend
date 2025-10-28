package handlers

import (
	"net/http"
	"strconv"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// TagHandler handles tag-related requests
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a new tag handler
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: service.NewTagService(),
	}
}

// GetTags gets all tags with filtering
// @Summary Get tags
// @Description Get all tags with optional filtering by kind
// @Tags tags
// @Produce json
// @Param kind query string false "Filter by tag kind (interest, category, activity, location, food, transport, accommodation)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	kind := c.Query("kind")

	// Validate pagination
	page, limit = utils.ValidatePagination(page, limit)

	// Get tags
	tags, total, err := h.tagService.GetTags(page, limit, kind)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get tags", err)
		return
	}

	// Send response with pagination
	utils.PaginatedResponse(c, "Tags retrieved successfully", tags, int64(total), page, limit)
}

// GetUserTags gets user's tags
// @Summary Get user tags
// @Description Get tags associated with the current user
// @Tags tags
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserTagListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/tags [get]
func (h *TagHandler) GetUserTags(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get user tags
	tags, err := h.tagService.GetUserTags(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get user tags",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.UserTagListResponse{
		Tags: tags,
	})
}

// AddUserTag adds a tag to user
// @Summary Add user tag
// @Description Add a tag to the current user's interests
// @Tags tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AddUserTagRequest true "Tag data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/tags [post]
func (h *TagHandler) AddUserTag(c *gin.Context) {
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
	var req dto.AddUserTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Add user tag
	err := h.tagService.AddUserTag(userID, req.TagID)
	if err != nil {
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Tag not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "tag already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Tag already exists",
				Message: "This tag is already associated with the user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to add user tag",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Tag added successfully",
	})
}

// RemoveUserTag removes a tag from user
// @Summary Remove user tag
// @Description Remove a tag from the current user's interests
// @Tags tags
// @Security BearerAuth
// @Produce json
// @Param tag_id path string true "Tag ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/tags/{tag_id} [delete]
func (h *TagHandler) RemoveUserTag(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get tag ID from path
	tagID := c.Param("tag_id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid tag ID",
			Message: "Tag ID is required",
		})
		return
	}

	// Remove user tag
	err := h.tagService.RemoveUserTag(userID, tagID)
	if err != nil {
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Tag not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to remove user tag",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Tag removed successfully",
	})
}

// GetEventTags gets event's tags
// @Summary Get event tags
// @Description Get tags associated with an event
// @Tags tags
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventTagListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/tags [get]
func (h *TagHandler) GetEventTags(c *gin.Context) {
	// Get event ID from path
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Get event tags
	tags, err := h.tagService.GetEventTags(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get event tags",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.EventTagListResponse{
		Tags: tags,
	})
}

// AddEventTag adds a tag to event
// @Summary Add event tag
// @Description Add a tag to an event
// @Tags tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body dto.AddEventTagRequest true "Tag data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/tags [post]
func (h *TagHandler) AddEventTag(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get event ID from path
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID is required",
		})
		return
	}

	// Parse request
	var req dto.AddEventTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Add event tag
	err := h.tagService.AddEventTag(eventID, req.TagID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Tag not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "not authorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Not authorized",
				Message: "Only the event creator can add tags",
			})
			return
		}
		if err.Error() == "tag already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Tag already exists",
				Message: "This tag is already associated with the event",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to add event tag",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Tag added successfully",
	})
}

// RemoveEventTag removes a tag from event
// @Summary Remove event tag
// @Description Remove a tag from an event
// @Tags tags
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Param tag_id path string true "Tag ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /events/{id}/tags/{tag_id} [delete]
func (h *TagHandler) RemoveEventTag(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get event ID and tag ID from path
	eventID := c.Param("id")
	tagID := c.Param("tag_id")
	if eventID == "" || tagID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid parameters",
			Message: "Event ID and Tag ID are required",
		})
		return
	}

	// Remove event tag
	err := h.tagService.RemoveEventTag(eventID, tagID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Tag not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "not authorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Not authorized",
				Message: "Only the event creator can remove tags",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to remove event tag",
			Message: err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Tag removed successfully",
	})
}
