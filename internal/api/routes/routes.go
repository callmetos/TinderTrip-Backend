package routes

import (
	"fmt"

	"TinderTrip-Backend/internal/api/handlers"
	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes
func SetupRoutes(router *gin.Engine) {
	// Create API v1 group
	v1 := router.Group("/api/v1")

	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SendSuccessResponse(c, "TinderTrip API is running", gin.H{
			"status": "healthy",
		})
	})

	// OTP monitoring for development
	otpHandler := handlers.NewOTPHandler()

	// Add dev/otp to v1 group (no auth required)
	v1.GET("/dev/otp", otpHandler.GetOTPs)

	// Image serving
	imageHandler, err := handlers.NewImageHandler()
	if err != nil {
		// Log error but don't fail startup
		fmt.Printf("Warning: Failed to initialize image handler: %v\n", err)
	} else {
		// Serve images with authentication
		imageGroup := router.Group("/images")
		imageGroup.Use(middleware.AuthMiddleware())
		{
			imageGroup.GET("/avatars/:user_id", imageHandler.ServeAvatar)
			imageGroup.GET("/events/:event_id", imageHandler.ServeEventImage)
		}
	}

	// OPTIONS handler for CORS preflight
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		c.Status(204)
	})

	// Auth routes
	authHandler := handlers.NewAuthHandler()
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/resend-verification", authHandler.ResendVerification)
		auth.POST("/login", authHandler.Login)
		auth.GET("/google", authHandler.GoogleAuth)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		auth.POST("/refresh", middleware.AuthMiddleware(), authHandler.RefreshToken)
		auth.GET("/check", middleware.AuthMiddleware(), authHandler.Check)
	}

	// Protected routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		userHandler := handlers.NewUserHandler()
		users := protected.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.DELETE("/profile", userHandler.DeleteProfile)
			users.GET("/setup-status", userHandler.GetSetupStatus)
		}

		// Preference routes
		preferenceHandler := handlers.NewPreferenceHandler()
		preferences := protected.Group("/users/preferences")
		{
			preferences.GET("/availability", preferenceHandler.GetAvailability)
			preferences.PUT("/availability", preferenceHandler.UpdateAvailability)
			preferences.GET("/budget", preferenceHandler.GetBudget)
			preferences.PUT("/budget", preferenceHandler.UpdateBudget)
		}

		// Event routes
		eventHandler := handlers.NewEventHandler()
		events := protected.Group("/events")
		{
			events.GET("", eventHandler.GetEvents)
			events.GET("/joined", eventHandler.GetJoinedEvents)
			events.GET("/suggestions", eventHandler.GetEventSuggestions)
			events.POST("", eventHandler.CreateEvent)
			events.GET("/:id", eventHandler.GetEvent)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
			events.POST("/:id/join", eventHandler.JoinEvent)
			events.POST("/:id/leave", eventHandler.LeaveEvent)
			events.POST("/:id/confirm", eventHandler.ConfirmEvent)
			events.POST("/:id/cancel", eventHandler.CancelEvent)
			events.POST("/:id/complete", eventHandler.CompleteEvent)
			events.POST("/:id/swipe", eventHandler.SwipeEvent)
			events.PUT("/:id/cover", eventHandler.UpdateCover)
			events.POST("/:id/photos", eventHandler.AddPhotos)
		}

		// Chat routes
		chatHandler := handlers.NewChatHandler()
		chat := protected.Group("/chat")
		{
			chat.GET("/rooms", chatHandler.GetRooms)
			chat.GET("/rooms/:id/messages", chatHandler.GetMessages)
			chat.POST("/rooms/:id/messages", chatHandler.SendMessage)
		}

		// History routes
		historyHandler := handlers.NewHistoryHandler()
		history := protected.Group("/history")
		{
			history.GET("", historyHandler.GetHistory)
			history.POST("/:id/complete", historyHandler.MarkComplete)
		}

		// Tag routes
		tagHandler := handlers.NewTagHandler()
		tags := protected.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
		}

		// User tag routes
		userTags := protected.Group("/users")
		{
			userTags.GET("/tags", tagHandler.GetUserTags)
			userTags.POST("/tags", tagHandler.AddUserTag)
			userTags.DELETE("/tags/:tag_id", tagHandler.RemoveUserTag)
		}

		// Food preference routes
		foodPreferenceHandler := handlers.NewFoodPreferenceHandler()
		foodPreferences := protected.Group("/users")
		{
			foodPreferences.GET("/food-preferences", foodPreferenceHandler.GetFoodPreferences)
			foodPreferences.PUT("/food-preferences", foodPreferenceHandler.UpdateFoodPreference)
			foodPreferences.PUT("/food-preferences/bulk", foodPreferenceHandler.UpdateAllFoodPreferences)
			foodPreferences.GET("/food-preferences/categories", foodPreferenceHandler.GetFoodPreferenceCategoriesWithUserPreferences)
			foodPreferences.GET("/food-preferences/stats", foodPreferenceHandler.GetFoodPreferenceStats)
			foodPreferences.DELETE("/food-preferences/:category", foodPreferenceHandler.DeleteFoodPreference)
		}

		// Travel preference routes
		travelPreferenceHandler := handlers.NewTravelPreferenceHandler()
		travelPreferences := protected.Group("/users")
		{
			travelPreferences.GET("/travel-preferences", travelPreferenceHandler.GetTravelPreferences)
			travelPreferences.POST("/travel-preferences", travelPreferenceHandler.AddTravelPreference)
			travelPreferences.PUT("/travel-preferences/bulk", travelPreferenceHandler.UpdateAllTravelPreferences)
			travelPreferences.GET("/travel-preferences/styles", travelPreferenceHandler.GetTravelPreferenceStylesWithUserPreferences)
			travelPreferences.GET("/travel-preferences/stats", travelPreferenceHandler.GetTravelPreferenceStats)
			travelPreferences.DELETE("/travel-preferences/:style", travelPreferenceHandler.DeleteTravelPreference)
		}

		// Event tag routes
		eventTags := protected.Group("/events")
		{
			eventTags.GET("/:id/tags", tagHandler.GetEventTags)
			eventTags.POST("/:id/tags", tagHandler.AddEventTag)
			eventTags.DELETE("/:id/tags/:tag_id", tagHandler.RemoveEventTag)
		}

		// Audit routes
		auditHandler := handlers.NewAuditHandler()
		audit := protected.Group("/audit")
		{
			audit.GET("/logs", auditHandler.GetAuditLogs)
			audit.GET("/entities/:entity_table/:entity_id", auditHandler.GetEntityAuditHistory)
		}
	}

	// Public routes (no authentication required)
	public := v1.Group("/public")
	{
		// Public event routes
		eventHandler := handlers.NewEventHandler()
		events := public.Group("/events")
		{
			events.GET("", eventHandler.GetPublicEvents)
			events.GET("/:id", eventHandler.GetPublicEvent)
		}

		// Public food preference routes
		foodPreferenceHandler := handlers.NewFoodPreferenceHandler()
		foodPreferences := public.Group("/food-preferences")
		{
			foodPreferences.GET("/categories", foodPreferenceHandler.GetFoodPreferenceCategories)
		}

		// Public travel preference routes
		travelPreferenceHandler := handlers.NewTravelPreferenceHandler()
		travelPreferences := public.Group("/travel-preferences")
		{
			travelPreferences.GET("/styles", travelPreferenceHandler.GetTravelPreferenceStyles)
		}
	}
}
