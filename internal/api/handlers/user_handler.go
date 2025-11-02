package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

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
// @Success 200 {object} dto.UserProfileResponseWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Get user profile
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Profile not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", profile)
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserProfileResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /users/profile [put]
// UpdateProfile handles both JSON and multipart form.
// If multipart and field "file" present -> upload to Nextcloud and update avatar_url.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
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
	DisplayName   *string    `json:"display_name,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	Languages     *string    `json:"languages,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Age           *int       `json:"age,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	JobTitle      *string    `json:"job_title,omitempty"`
	Smoking       *string    `json:"smoking,omitempty"`
	InterestsNote *string    `json:"interests_note,omitempty"`
	HomeLocation  *string    `json:"home_location,omitempty"`
}

// validateProfileUpdateRequest validates the profile update request
func validateProfileUpdateRequest(req updateProfileJSONReq) error {
	// Validate display name length
	if req.DisplayName != nil && utf8.RuneCountInString(*req.DisplayName) > 100 {
		return fmt.Errorf("display name must be 100 characters or less")
	}

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

	// Validate gender (align with models: male, female, nonbinary, prefer_not_say)
	if req.Gender != nil && *req.Gender != "" {
		validGenders := []string{"male", "female", "nonbinary", "prefer_not_say"}
		valid := false
		for _, g := range validGenders {
			if *req.Gender == g {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("gender must be one of: male, female, nonbinary, prefer_not_say")
		}
	}

	// Validate age (if provided)
	if req.Age != nil {
		if *req.Age < 0 || *req.Age > 120 {
			return fmt.Errorf("age must be between 0 and 120")
		}
	}

	// Validate job title length
	if req.JobTitle != nil && utf8.RuneCountInString(*req.JobTitle) > 100 {
		return fmt.Errorf("job title must be 100 characters or less")
	}

	// Validate smoking preference (align with models: no, yes, occasionally)
	if req.Smoking != nil && *req.Smoking != "" {
		validSmoking := []string{"no", "yes", "occasionally"}
		valid := false
		for _, s := range validSmoking {
			if *req.Smoking == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("smoking must be one of: no, yes, occasionally")
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
		utils.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	// Validate request data
	if err := validateProfileUpdateRequest(reqBody); err != nil {
		utils.ValidationErrorResponse(c, "Validation failed", err.Error())
		return
	}

	req := dto.UpdateProfileRequest{
		DisplayName:   reqBody.DisplayName,
		Bio:           reqBody.Bio,
		Languages:     reqBody.Languages,
		DateOfBirth:   reqBody.DateOfBirth,
		Age:           reqBody.Age,
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
		utils.InternalServerErrorResponse(c, "Update failed", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", profile)
}

func updateProfileMultipart(h *UserHandler, c *gin.Context, userID string) {
	// Text fields
	displayName := c.PostForm("display_name")
	bio := c.PostForm("bio")
	languages := c.PostForm("languages")
	dobStr := c.PostForm("date_of_birth")
	ageStr := c.PostForm("age")
	gender := c.PostForm("gender")
	jobTitle := c.PostForm("job_title")
	smoking := c.PostForm("smoking")
	interestsNote := c.PostForm("interests_note")
	homeLocation := c.PostForm("home_location")

	// Create request struct for validation
	reqBody := updateProfileJSONReq{
		DisplayName:   &displayName,
		Bio:           &bio,
		Languages:     &languages,
		DateOfBirth:   nil,
		Age:           nil,
		Gender:        &gender,
		JobTitle:      &jobTitle,
		Smoking:       &smoking,
		InterestsNote: &interestsNote,
		HomeLocation:  &homeLocation,
	}

	// Parse optional date_of_birth
	if dobStr != "" {
		if t, err := time.Parse(time.RFC3339, dobStr); err == nil {
			dob := t
			reqBody.DateOfBirth = &dob
		} else {
			utils.ValidationErrorResponse(c, "Invalid request", "date_of_birth must be RFC3339 format")
			return
		}
	}

	// Parse optional age
	if ageStr != "" {
		if n, err := strconv.Atoi(ageStr); err == nil {
			age := n
			reqBody.Age = &age
		} else {
			utils.ValidationErrorResponse(c, "Invalid request", "age must be an integer")
			return
		}
	}

	// Validate request data
	if err := validateProfileUpdateRequest(reqBody); err != nil {
		utils.ValidationErrorResponse(c, "Validation failed", err.Error())
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
			utils.BadRequestResponse(c, "Invalid file: "+err.Error())
			return
		}
		defer src.Close()

		fs, err := service.NewFileService()
		if err != nil {
			utils.InternalServerErrorResponse(c, "Storage initialization failed", err)
			return
		}
		_, url, _, _, _, err := fs.UploadImage(c, "avatars", fileHeader.Filename, src)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnsupportedMediaType, utils.ErrCodeInvalidInput, "Upload failed", err)
			return
		}
		// Store the full Nextcloud URL
		avatarURL = &url
	}

	req := dto.UpdateProfileRequest{
		DisplayName:   toPtr(displayName),
		Bio:           toPtr(bio),
		Languages:     toPtr(languages),
		DateOfBirth:   reqBody.DateOfBirth,
		Age:           reqBody.Age,
		Gender:        toPtr(gender),
		JobTitle:      toPtr(jobTitle),
		Smoking:       toPtr(smoking),
		InterestsNote: toPtr(interestsNote),
		HomeLocation:  toPtr(homeLocation),
		AvatarURL:     avatarURL, // only set if file uploaded
	}

	profile, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Update failed", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", profile)
}

// DeleteProfile deletes user profile
// @Summary Delete user profile
// @Description Delete current user's profile (soft delete)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /users/profile [delete]
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Delete profile
	err := h.userService.DeleteProfile(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Delete failed", err)
		return
	}

	utils.SendSuccessResponse(c, "Profile deleted successfully", nil)
}

// GetSetupStatus checks if user has completed initial setup
// @Summary Get user setup status
// @Description Check if current user has completed initial profile setup
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.SetupStatusResponseWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /users/setup-status [get]
func (h *UserHandler) GetSetupStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Check setup status
	setupCompleted, err := h.userService.CheckSetupStatus(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to check setup status", err)
		return
	}

	utils.SendSuccessResponse(c, "Setup status retrieved successfully", dto.SetupStatusResponse{
		SetupCompleted: setupCompleted,
	})
}
