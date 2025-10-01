package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"TinderTrip-Backend/internal/api"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"
)

// @title TinderTrip API
// @version 1.0
// @description A Tinder-like trip matching API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api.tindertrip.phitik.com
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to databases
	if err := database.ConnectPostgres(); err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer database.ClosePostgres()

	if err := database.ConnectRedis(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (continuing without Redis)", err)
	} else {
		defer database.CloseRedis()
	}

	// Create and start server
	srv := api.NewServer()

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop background cleanup routines
	srv.StopCleanup()

	// Graceful shutdown
	if err := srv.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
