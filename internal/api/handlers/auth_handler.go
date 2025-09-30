package handlers

import (
	"fmt"
	"net/http"
	"time"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"
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

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Create user
	user, err := h.authService.Register(req.Email, req.Password, req.DisplayName)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "User already exists",
				Message: "An account with this email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), *user.Email, string(user.Provider))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		})
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

	c.JSON(http.StatusCreated, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID.String(),
			Email:       *user.Email,
			DisplayName: user.GetDisplayName(),
			Provider:    string(user.Provider),
			CreatedAt:   user.CreatedAt,
		},
	})
}

// Login handles user login
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Authenticate user
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Email or password is incorrect",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), *user.Email, string(user.Provider))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	database.GetDB().Save(user)

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID.String(),
			Email:       *user.Email,
			DisplayName: user.GetDisplayName(),
			Provider:    string(user.Provider),
			CreatedAt:   user.CreatedAt,
		},
	})
}

// GoogleAuth handles Google OAuth authentication
// @Summary Google OAuth authentication
// @Description Get Google OAuth authorization URL
// @Tags auth
// @Produce json
// @Success 200 {object} dto.GoogleAuthResponse
// @Failure 500 {object} dto.ErrorResponse
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "OAuth service error",
			Message: "Google OAuth service not initialized",
		})
		return
	}

	// Get Google OAuth URL
	authURL := h.googleOAuthService.GetAuthURL(state)

	c.JSON(http.StatusOK, dto.GoogleAuthResponse{
		AuthURL: authURL,
		State:   state,
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
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Missing parameters",
			Message: "Authorization code and state are required",
		})
		return
	}

	// Validate state parameter
	ctx := c.Request.Context()
	validState, err := database.GetCache(ctx, "oauth_state:"+state)
	if err != nil || validState != "valid" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid state",
			Message: "Invalid or expired state parameter",
		})
		return
	}

	// Exchange code for token
	token, err := h.googleOAuthService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Token exchange failed",
			Message: "Failed to exchange authorization code for token",
		})
		return
	}

	// Get user info from Google
	userInfo, err := h.googleOAuthService.GetUserInfo(ctx, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "User info failed",
			Message: "Failed to get user information from Google",
		})
		return
	}

	// Create or update user
	user, err := h.googleOAuthService.CreateOrUpdateUser(ctx, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "User creation failed",
			Message: "Failed to create or update user",
		})
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateToken(user.ID.String(), userInfo.Email, string(user.Provider))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Token generation failed",
			Message: "Failed to generate authentication token",
		})
		return
	}

	// Clean up state
	database.DeleteCache(ctx, "oauth_state:"+state)

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: jwtToken,
		User: dto.UserResponse{
			ID:          user.ID.String(),
			Email:       userInfo.Email,
			DisplayName: userInfo.Name,
			Provider:    string(user.Provider),
			CreatedAt:   user.CreatedAt,
		},
	})
}

// ForgotPassword handles password reset request
// @Summary Request password reset
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Send password reset email
	err := h.authService.SendPasswordResetEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "If the email exists, a password reset link has been sent",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset password with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Reset password
	err := h.authService.ResetPassword(req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Password reset failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Password has been reset successfully",
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Logged out successfully",
	})
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get current token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Authorization header required",
			Message: "Please provide a valid token",
		})
		return
	}

	// Extract token
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Invalid token",
			Message: err.Error(),
		})
		return
	}

	// Refresh token
	newToken, err := utils.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Token refresh failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: newToken,
	})
}
