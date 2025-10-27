package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// GetProfile gets user profile
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserProfileResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get user profile
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Profile not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [put]
// UpdateProfile handles both JSON and multipart form.
// If multipart and field "file" present -> upload to Nextcloud and update avatar_url.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") || strings.Contains(c.GetHeader("Content-Type"), "boundary=") {
		updateProfileMultipart(h, c, userID)
		return
	}
	updateProfileJSON(h, c, userID)
}

type updateProfileJSONReq struct {
	Bio           *string `json:"bio,omitempty"`
	Languages     *string `json:"languages,omitempty"`
	Gender        *string `json:"gender,omitempty"`
	JobTitle      *string `json:"job_title,omitempty"`
	Smoking       *string `json:"smoking,omitempty"`
	InterestsNote *string `json:"interests_note,omitempty"`
	HomeLocation  *string `json:"home_location,omitempty"`
}

// validateProfileUpdateRequest validates the profile update request
func validateProfileUpdateRequest(req updateProfileJSONReq) error {
	// Validate bio length
	if req.Bio != nil && utf8.RuneCountInString(*req.Bio) > 500 {
		return fmt.Errorf("bio must be 500 characters or less")
	}

	// Validate languages format (comma-separated)
	if req.Languages != nil && *req.Languages != "" {
		languages := strings.Split(*req.Languages, ",")
		for _, lang := range languages {
			lang = strings.TrimSpace(lang)
			if utf8.RuneCountInString(lang) > 50 {
				return fmt.Errorf("each language must be 50 characters or less")
			}
		}
	}

	// Validate gender
	if req.Gender != nil && *req.Gender != "" {
		validGenders := []string{"male", "female", "other", "prefer_not_to_say"}
		valid := false
		for _, g := range validGenders {
			if *req.Gender == g {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("gender must be one of: male, female, other, prefer_not_to_say")
		}
	}

	// Validate job title length
	if req.JobTitle != nil && utf8.RuneCountInString(*req.JobTitle) > 100 {
		return fmt.Errorf("job title must be 100 characters or less")
	}

	// Validate smoking preference
	if req.Smoking != nil && *req.Smoking != "" {
		validSmoking := []string{"yes", "no", "occasionally", "prefer_not_to_say"}
		valid := false
		for _, s := range validSmoking {
			if *req.Smoking == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("smoking preference must be one of: yes, no, occasionally, prefer_not_to_say")
		}
	}

	// Validate interests note length
	if req.InterestsNote != nil && utf8.RuneCountInString(*req.InterestsNote) > 1000 {
		return fmt.Errorf("interests note must be 1000 characters or less")
	}

	// Validate home location length
	if req.HomeLocation != nil && utf8.RuneCountInString(*req.HomeLocation) > 200 {
		return fmt.Errorf("home location must be 200 characters or less")
	}

	return nil
}

func updateProfileJSON(h *UserHandler, c *gin.Context, userID string) {
	var reqBody updateProfileJSONReq
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}

	// Validate request data
	if err := validateProfileUpdateRequest(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	req := dto.UpdateProfileRequest{
		Bio:           reqBody.Bio,
		Languages:     reqBody.Languages,
		Gender:        reqBody.Gender,
		JobTitle:      reqBody.JobTitle,
		Smoking:       reqBody.Smoking,
		InterestsNote: reqBody.InterestsNote,
		HomeLocation:  reqBody.HomeLocation,
		AvatarURL:     nil, // JSON mode: do not touch avatar
	}

	// Update profile
	profile, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Update failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func updateProfileMultipart(h *UserHandler, c *gin.Context, userID string) {
	// Text fields
	bio := c.PostForm("bio")
	languages := c.PostForm("languages")
	gender := c.PostForm("gender")
	jobTitle := c.PostForm("job_title")
	smoking := c.PostForm("smoking")
	interestsNote := c.PostForm("interests_note")
	homeLocation := c.PostForm("home_location")

	// Create request struct for validation
	reqBody := updateProfileJSONReq{
		Bio:           &bio,
		Languages:     &languages,
		Gender:        &gender,
		JobTitle:      &jobTitle,
		Smoking:       &smoking,
		InterestsNote: &interestsNote,
		HomeLocation:  &homeLocation,
	}

	// Validate request data
	if err := validateProfileUpdateRequest(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
		return
	}

	toPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		v := s
		return &v
	}

	var avatarURL *string
	if fileHeader, err := c.FormFile("file"); err == nil && fileHeader != nil {
		src, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid file", Message: err.Error()})
			return
		}
		defer src.Close()

		fs, err := service.NewFileService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Storage init failed", Message: err.Error()})
			return
		}
		_, url, _, _, _, err := fs.UploadImage(c, "avatars", fileHeader.Filename, src)
		if err != nil {
			c.JSON(http.StatusUnsupportedMediaType, dto.ErrorResponse{Error: "Upload failed", Message: err.Error()})
			return
		}
		// Store the full Nextcloud URL
		avatarURL = &url
	}

	req := dto.UpdateProfileRequest{
		Bio:           toPtr(bio),
		Languages:     toPtr(languages),
		Gender:        toPtr(gender),
		JobTitle:      toPtr(jobTitle),
		Smoking:       toPtr(smoking),
		InterestsNote: toPtr(interestsNote),
		HomeLocation:  toPtr(homeLocation),
		AvatarURL:     avatarURL, // only set if file uploaded
	}

	profile, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Update failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// DeleteProfile deletes user profile
// @Summary Delete user profile
// @Description Delete current user's profile (soft delete)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [delete]
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Delete profile
	err := h.userService.DeleteProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Delete failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Profile deleted successfully",
	})
}

// GetSetupStatus checks if user has completed initial setup
// @Summary Get user setup status
// @Description Check if current user has completed initial profile setup
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.SetupStatusResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/setup-status [get]
func (h *UserHandler) GetSetupStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Check setup status
	setupCompleted, err := h.userService.CheckSetupStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to check setup status",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.SetupStatusResponse{
			SetupCompleted: setupCompleted,
		},
		"message": "Setup status retrieved successfully",
	})
}
