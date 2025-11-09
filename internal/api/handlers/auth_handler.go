package handlers

import (
	"fmt"
	"net/http"
	"time"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService        *service.AuthService
	googleOAuthService *service.GoogleOAuthService
	emailService       *service.EmailService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService:        service.NewAuthService(),
		googleOAuthService: service.NewGoogleOAuthService(),
		emailService:       service.NewEmailService(),
	}
}

// StopCleanup stops background cleanup routines
func (h *AuthHandler) StopCleanup() {
	if h.authService != nil {
		h.authService.StopCleanup()
	}
}

// Register handles user registration - creates user and sends OTP for email verification
// @Summary Register a new user
// @Description Register a new user with email and password, creates user with unverified status and sends OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 409 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Send email verification OTP
	err := h.authService.SendEmailVerificationOTP(req.Email, req.DisplayName)
	if err != nil {
		// Check if error is due to user already exists
		if err.Error() == "user already exists" {
			utils.ConflictResponse(c, "An account with this email already exists")
			return
		}
		// Check if error is due to display name already taken
		if err.Error() == "display name already taken" {
			utils.ConflictResponse(c, "Display name already taken. Please choose a different name.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to send verification OTP", err)
		return
	}

	utils.SendSuccessResponse(c, "Verification OTP sent to your email", nil)
}

// Login handles user login
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} dto.AuthResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Authenticate user
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeAuthenticationFailed, "Email or password is incorrect", nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), *user.Email, string(user.Provider))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token", err)
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	database.GetDB().Save(user)

	// Create custom response with token and user at top level
	c.JSON(http.StatusOK, dto.AuthResponseWrapper{
		Success:   true,
		RequestID: utils.GetRequestID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "Login successful",
		Token:     token,
		User: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         *user.Email,
			DisplayName:   user.GetDisplayName(),
			Provider:      string(user.Provider),
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	})
}

// GoogleAuth handles Google OAuth authentication
// @Summary Google OAuth authentication
// @Description Get Google OAuth authorization URL
// @Tags auth
// @Produce json
// @Success 200 {object} dto.GoogleAuthResponseWrapper
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/google [get]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	// Generate state parameter
	state := uuid.New().String()

	// Store state in memory for validation (since Redis is optional)
	// In production, you should use Redis for this
	ctx := c.Request.Context()
	err := database.SetCache(ctx, "oauth_state:"+state, "valid", 10*time.Minute)
	if err != nil {
		// If Redis is not available, continue without state validation
		// This is not recommended for production
		fmt.Printf("Warning: Failed to store OAuth state in cache: %v\n", err)
	}

	// Check if Google OAuth service is properly configured
	if h.googleOAuthService == nil {
		utils.InternalServerErrorResponse(c, "Google OAuth service not initialized", nil)
		return
	}

	// Get Google OAuth URL
	authURL := h.googleOAuthService.GetAuthURL(state)

	// Create custom response with auth_url and state at top level
	c.JSON(http.StatusOK, dto.GoogleAuthResponseWrapper{
		Success:   true,
		RequestID: utils.GetRequestID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "OAuth URL generated successfully",
		AuthURL:   authURL,
		State:     state,
	})
}

// GoogleCallback handles Google OAuth callback
// @Summary Google OAuth callback
// @Description Handle Google OAuth callback and authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter"
// @Success 200 {object} dto.AuthResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// Helper function to redirect to frontend with error
	redirectToError := func(errorType, errorMessage string) {
		frontendURL := config.AppConfig.Server.FrontendURL
		redirectURL := fmt.Sprintf("%s/callbackerror?error=%s&message=%s",
			frontendURL,
			errorType,
			errorMessage,
		)
		c.Redirect(http.StatusFound, redirectURL)
	}

	if code == "" || state == "" {
		redirectToError("missing_parameters", "Authorization code and state are required")
		return
	}

	// Validate state parameter (skip if Redis is not available)
	ctx := c.Request.Context()
	validState, err := database.GetCache(ctx, "oauth_state:"+state)
	if err != nil {
		// If Redis is not available, skip state validation
		utils.Logger().WithField("error", err).Warn("Redis not available, skipping state validation")
	} else if validState != "valid" {
		redirectToError("invalid_state", "Invalid or expired state parameter")
		return
	}

	// Exchange code for token
	token, err := h.googleOAuthService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		redirectToError("token_exchange_failed", "Failed to exchange authorization code")
		return
	}

	// Get user info from Google
	userInfo, err := h.googleOAuthService.GetUserInfo(ctx, token)
	if err != nil {
		redirectToError("user_info_failed", "Failed to get user information from Google")
		return
	}

	// Create or update user
	user, isNewUser, err := h.googleOAuthService.CreateOrUpdateUser(ctx, userInfo)
	if err != nil {
		redirectToError("user_creation_failed", "Failed to create or update user account")
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateToken(user.ID.String(), userInfo.Email, string(user.Provider))
	if err != nil {
		redirectToError("token_generation_failed", "Failed to generate authentication token")
		return
	}

	// Clean up state
	database.DeleteCache(ctx, "oauth_state:"+state)

	// Send welcome email only for new Google OAuth users
	if isNewUser {
		go func() {
			if err := h.emailService.SendWelcomeEmail(userInfo.Email, userInfo.Name); err != nil {
				utils.Logger().WithFields(map[string]interface{}{
					"error":   err,
					"user_id": user.ID.String(),
					"email":   userInfo.Email,
				}).Error("Failed to send welcome email for Google OAuth user")
			}
		}()
	}

	// Redirect to frontend with token (success)
	frontendURL := config.AppConfig.Server.FrontendURL
	redirectURL := fmt.Sprintf("%s/callback?token=%s&user_id=%s&email=%s&display_name=%s&provider=%s&is_verified=%t",
		frontendURL,
		jwtToken,
		user.ID.String(),
		userInfo.Email,
		userInfo.Name,
		string(user.Provider),
		user.EmailVerified)

	c.Redirect(http.StatusFound, redirectURL)
}

// ForgotPassword handles password reset request
// @Summary Request password reset
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Send password reset OTP
	err := h.authService.SendPasswordResetOTP(req.Email)
	if err != nil {
		// Don't reveal if email exists or not (security best practice)
		utils.SendSuccessResponse(c, "If the email exists, a password reset OTP has been sent", nil)
		return
	}

	utils.SendSuccessResponse(c, "If the email exists, a password reset OTP has been sent", nil)
}

// ResetPassword handles password reset
// @Summary Reset password with OTP
// @Description Reset password with OTP verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password data with OTP"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Reset password with OTP
	err := h.authService.ResetPassword(req.Email, req.OTP, req.Password)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "Password has been reset successfully", nil)
}

// VerifyOTP handles OTP verification
// @Summary Verify OTP
// @Description Verify OTP for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "OTP verification data"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Verify OTP
	err := h.authService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, "OTP verified successfully", nil)
}

// VerifyEmail handles email verification with OTP
// @Summary Verify email with OTP
// @Description Verify email with OTP and complete user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterWithOTPRequest true "Email verification data"
// @Success 201 {object} dto.AuthResponseWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.RegisterWithOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request format", err.Error())
		return
	}

	// Verify OTP and create user
	user, err := h.authService.VerifyEmailOTP(req.Email, req.OTP, req.Password, req.DisplayName)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), *user.Email, string(user.Provider))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token", err)
		return
	}

	// Send welcome email
	go func() {
		if err := h.emailService.SendWelcomeEmail(*user.Email, user.GetDisplayName()); err != nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":   err,
				"user_id": user.ID.String(),
				"email":   *user.Email,
			}).Error("Failed to send welcome email")
		}
	}()

	// Create custom response with token and user at top level
	c.JSON(http.StatusCreated, dto.AuthResponseWrapper{
		Success:   true,
		RequestID: utils.GetRequestID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "Registration successful",
		Token:     token,
		User: dto.UserResponse{
			ID:            user.ID.String(),
			Email:         *user.Email,
			DisplayName:   user.GetDisplayName(),
			Provider:      string(user.Provider),
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	})
}

// ResendVerification handles resending verification OTP
// @Summary Resend verification OTP
// @Description Resend verification OTP to email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Email address"
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 409 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Resend verification OTP
	err := h.authService.ResendEmailVerificationOTP(req.Email)
	if err != nil {
		if err.Error() == "user already exists" {
			utils.ErrorResponse(c, http.StatusConflict, utils.ErrCodeResourceExists, "An account with this email already exists", nil)
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to resend verification", err)
		return
	}

	utils.SendSuccessResponse(c, "Verification OTP resent to your email. Please check your inbox.", nil)
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} dto.SuccessMessageWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	utils.SendSuccessResponse(c, "Logged out successfully", nil)
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} dto.TokenResponseWrapper
// @Failure 401 {object} dto.ErrorAPIResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get current token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Please provide a valid token")
		return
	}

	// Extract token
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeInvalidToken, err.Error(), nil)
		return
	}

	// Refresh token
	newToken, err := utils.RefreshToken(token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeExpiredToken, "Token refresh failed", err)
		return
	}

	// Create custom response with token at top level
	c.JSON(http.StatusOK, dto.TokenResponseWrapper{
		Success:   true,
		RequestID: utils.GetRequestID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   "Token refreshed successfully",
		Token:     newToken,
	})
}

// CheckResponse represents the check endpoint response
type CheckResponse struct {
	Status    string    `json:"status"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Check validates JWT token and returns user info
// @Summary Check JWT token
// @Description Validate JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} CheckResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /auth/check [get]
func (h *AuthHandler) Check(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Authentication token is required",
		})
		return
	}

	// Get user from database
	var user models.User
	err := database.GetDB().Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error
	if err != nil {
		utils.UnauthorizedResponse(c, "User account not found or has been deleted")
		return
	}

	// Return user info
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	username := ""
	if user.DisplayName != nil {
		username = *user.DisplayName
	}

	utils.SendSuccessResponse(c, "Token is valid", CheckResponse{
		Status:    "valid",
		Email:     email,
		Username:  username,
		CreatedAt: user.CreatedAt,
	})
}
