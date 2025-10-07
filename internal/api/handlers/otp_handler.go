package handlers

import (
	"net/http"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// OTPHandler handles OTP monitoring requests
type OTPHandler struct{}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler() *OTPHandler {
	return &OTPHandler{}
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": err.Error(),
		})
		return
	}

	// Format response
	var otps []gin.H
	for _, ev := range emailVerifications {
		otps = append(otps, gin.H{
			"email":      ev.Email,
			"otp":        ev.OTP,
			"expires_at": ev.ExpiresAt,
			"created_at": ev.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Active OTPs retrieved successfully",
		"data": gin.H{
			"otps":  otps,
			"count": len(otps),
		},
	})
}
