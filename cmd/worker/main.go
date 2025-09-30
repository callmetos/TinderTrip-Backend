package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to databases
	if err := database.ConnectPostgres(); err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer database.ClosePostgres()

	if err := database.ConnectRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()

	// Create worker service
	workerService := service.NewWorkerService()

	// Start workers
	go workerService.StartEmailWorker()
	go workerService.StartNotificationWorker()
	go workerService.StartCleanupWorker()

	log.Println("Worker started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// Give workers time to finish
	time.Sleep(5 * time.Second)

	log.Println("Worker stopped")
}
