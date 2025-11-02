package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ImageHandler handles image serving requests
type ImageHandler struct {
	imageService *service.ImageService
	userService  *service.UserService
}

// NewImageHandler creates a new image handler
func NewImageHandler() (*ImageHandler, error) {
	imageService, err := service.NewImageService()
	if err != nil {
		return nil, err
	}

	return &ImageHandler{
		imageService: imageService,
		userService:  service.NewUserService(),
	}, nil
}

// ServeAvatar serves user avatar image
// @Summary Serve user avatar
// @Description Serves user avatar image by user ID
// @Tags images
// @Security BearerAuth
// @Produce image/*
// @Param user_id path string true "User ID"
// @Success 200 {file} file "Image file"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/avatars/{user_id} [get]
func (h *ImageHandler) ServeAvatar(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: "User ID is required",
		})
		return
	}

	// Ensure the requesting user is authorized to view this avatar
	_, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Authentication token is required",
		})
		return
	}

	// Parse user ID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	// Get user profile directly from database to get the actual storage key
	var profile models.UserProfile
	err = database.GetDB().Where("user_id = ?", userUUID).First(&profile).Error
	if err != nil {
		utils.Logger().WithField("error", err).Error("Failed to get user profile for avatar serving")
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "User not found",
			Message: "User profile could not be retrieved",
		})
		return
	}

	if profile.AvatarURL == nil || *profile.AvatarURL == "" {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Image not found",
			Message: "User does not have an avatar",
		})
		return
	}

	// The AvatarURL stored in the profile is the Nextcloud internal key
	imageKey := *profile.AvatarURL

	// Get image from storage using the actual storage key
	imageData, contentType, err := h.imageService.GetImageFromKey(c.Request.Context(), imageKey)
	if err != nil {
		utils.Logger().WithField("error", err).Error("Failed to get image from key")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Image retrieval failed",
			Message: "Could not retrieve image data",
		})
		return
	}

	// Generate ETag from image data for better caching
	etag := generateETag(imageData)
	c.Header("ETag", etag)

	// Check if client has cached version
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	// Set appropriate headers for browser caching
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	c.Header("Expires", time.Now().Add(24*time.Hour).UTC().Format(http.TimeFormat))
	c.Data(http.StatusOK, contentType, imageData)
}

// ServeEventImage serves event image
// @Summary Serve event image
// @Description Serves event image by event ID
// @Tags images
// @Security BearerAuth
// @Produce image/*
// @Param event_id path string true "Event ID"
// @Success 200 {file} file "Image file"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /images/events/{event_id} [get]
func (h *ImageHandler) ServeEventImage(c *gin.Context) {
	eventID := c.Param("event_id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: "Event ID is required",
		})
		return
	}

	// Ensure the requesting user is authorized to view this event image
	_, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Authentication token is required",
		})
		return
	}

	// Clean the event ID
	eventID = strings.TrimSpace(eventID)

	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid event ID",
			Message: "Event ID format is invalid",
		})
		return
	}

	// Get event from database to find the actual image URL
	var event models.Event
	err = database.GetDB().Where("id = ? AND deleted_at IS NULL", eventUUID).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Event not found",
				Message: "The requested event does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database error",
				Message: "Failed to retrieve event information",
			})
		}
		return
	}

	// Check if event has cover image
	if event.CoverImageURL == nil || *event.CoverImageURL == "" {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Event image not found",
			Message: "The requested event has no cover image",
		})
		return
	}

	// Use the stored URL directly - it's already a full Nextcloud URL
	coverURL := *event.CoverImageURL

	// Get image from storage using the full URL
	imageData, contentType, err := h.imageService.GetImageFromURL(c.Request.Context(), coverURL)
	if err != nil {
		utils.Logger().WithField("error", err).WithField("event_id", eventID).WithField("url", coverURL).Error("Failed to get event image from storage")
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Event image not found",
			Message: "The requested event image could not be found",
		})
		return
	}

	// Generate ETag from image data for better caching
	etag := generateETag(imageData)
	c.Header("ETag", etag)

	// Check if client has cached version
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	// Set appropriate headers for browser caching
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	c.Header("Expires", time.Now().Add(24*time.Hour).UTC().Format(http.TimeFormat))
	c.Data(http.StatusOK, contentType, imageData)
}

// generateETag generates an ETag from image data
func generateETag(imageData []byte) string {
	hash := sha256.Sum256(imageData)
	return fmt.Sprintf(`"%x"`, hash[:16]) // Use first 16 bytes for shorter ETag
}
