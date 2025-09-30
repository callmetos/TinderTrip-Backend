package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"TinderTrip-Backend/docs"
	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/api/routes"
	"TinderTrip-Backend/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

func NewServer() *Server {
	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.Mode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.AuthContext()) // Extract user ID from JWT for API logging
	router.Use(middleware.APILogger())   // Add API logging to database
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())

	// Setup routes
	routes.SetupRoutes(router)

	// Add Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "TinderTrip API"
	docs.SwaggerInfo.Description = "A Tinder-like trip matching API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	return &Server{
		router: router,
	}
}

func (s *Server) Start() error {
	port := config.AppConfig.Server.Port
	host := config.AppConfig.Server.Host

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

	return s.httpServer.Shutdown(ctx)
}
