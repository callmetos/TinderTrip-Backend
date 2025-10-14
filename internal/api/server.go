package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"TinderTrip-Backend/docs"
	"TinderTrip-Backend/internal/api/handlers"
	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/api/routes"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer        *http.Server
	router            *gin.Engine
	authHandler       *handlers.AuthHandler
	monitoringService *service.MonitoringService
}

func NewServer() *Server {
	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.Mode)

	// Create router
	router := gin.New()

	// Initialize monitoring service
	var monitoringService *service.MonitoringService
	if config.AppConfig.Monitoring.Enabled {
		monitoringService = service.NewMonitoringService(database.DB)
		middleware.SetMonitoringService(monitoringService)
	}

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.AuthContext()) // Extract user ID from JWT for API logging
	router.Use(middleware.APILogger())   // Add API logging to database
	router.Use(middleware.Recovery())
	router.Use(middleware.CustomCORS()) // CORS enabled for all responses
	// router.Use(middleware.RateLimit()) // Rate limit disabled for development

	// Add monitoring middleware
	if config.AppConfig.Monitoring.Enabled {
		router.Use(middleware.PrometheusMetrics())
		router.Use(middleware.BusinessMetrics())
	}

	// Setup routes
	routes.SetupRoutes(router)

	// Add Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "TinderTrip API"
	docs.SwaggerInfo.Description = "A Tinder-like trip matching API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "api.tindertrip.phitik.com"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Create auth handler for cleanup
	authHandler := handlers.NewAuthHandler()

	return &Server{
		router:            router,
		authHandler:       authHandler,
		monitoringService: monitoringService,
	}
}

func (s *Server) Start() error {
	// Start monitoring service
	if s.monitoringService != nil {
		if err := s.monitoringService.Start(); err != nil {
			log.Printf("Failed to start monitoring service: %v", err)
		}
	}

	port := config.AppConfig.Server.Port
	host := "0.0.0.0" // Hardcode to 0.0.0.0 for LAN access

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on %s:%s", host, port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop monitoring service
	if s.monitoringService != nil {
		if err := s.monitoringService.Stop(); err != nil {
			log.Printf("Error stopping monitoring service: %v", err)
		}
	}

	return s.httpServer.Shutdown(ctx)
}

// StopCleanup stops background cleanup routines
func (s *Server) StopCleanup() {
	if s.authHandler != nil {
		s.authHandler.StopCleanup()
	}
}
