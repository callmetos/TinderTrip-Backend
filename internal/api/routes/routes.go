package routes

import (
	"TinderTrip-Backend/internal/api/handlers"
	"TinderTrip-Backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes
func SetupRoutes(router *gin.Engine) {
	// Create API v1 group
	v1 := router.Group("/api/v1")

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"status":  "success",
			"message": "TinderTrip API is running",
		})
	})

	// Auth routes
	authHandler := handlers.NewAuthHandler()
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/google", authHandler.GoogleAuth)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		auth.POST("/refresh", middleware.AuthMiddleware(), authHandler.RefreshToken)
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
			events.POST("", eventHandler.CreateEvent)
			events.GET("/:id", eventHandler.GetEvent)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
			events.POST("/:id/join", eventHandler.JoinEvent)
			events.POST("/:id/leave", eventHandler.LeaveEvent)
			events.POST("/:id/swipe", eventHandler.SwipeEvent)
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
	}
}
