package handlers

import (
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// OTPHandler handles OTP monitoring requests
type OTPHandler struct{}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler() *OTPHandler {
	return &OTPHandler{}
}

// OTPInfo represents OTP information
type OTPInfo struct {
	Email     string    `json:"email"`
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GetOTPs gets all active OTPs from email_verifications table
// @Summary Get OTPs for monitoring
// @Description Get all active OTPs from email_verifications table for development/testing
// @Tags development
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dev/otp [get]
func (h *OTPHandler) GetOTPs(c *gin.Context) {
	// Get all active OTPs (not expired)
	var emailVerifications []models.EmailVerification
	err := database.GetDB().Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		Find(&emailVerifications).Error

	if err != nil {
		utils.InternalServerErrorResponse(c, "Database error", err)
		return
	}

	// Format response
	var otps []OTPInfo
	for _, ev := range emailVerifications {
		otps = append(otps, OTPInfo{
			Email:     ev.Email,
			OTP:       ev.OTP,
			ExpiresAt: ev.ExpiresAt,
			CreatedAt: ev.CreatedAt,
		})
	}

	utils.SendSuccessResponse(c, "Active OTPs retrieved successfully", map[string]interface{}{
		"otps":  otps,
		"count": len(otps),
	})
}
