package main

import (
	"flag"
	"fmt"
	"log"

	"TinderTrip-Backend/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, down, force, version")
		steps  = flag.Int("steps", 0, "Number of migration steps (for up/down)")
		ver    = flag.Int("version", 0, "Target version (for force)")
	)
	flag.Parse()

	// Load configuration
	config.LoadConfig()

	// Build database URL
	cfg := config.AppConfig.Database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	// Create migrate instance
	m, err := migrate.New(
		"file://pkg/database/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}
	defer m.Close()

	// Execute migration action
	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "force":
		if *ver < 0 {
			log.Fatal("Version must be non-negative for force action")
		}
		err = m.Force(*ver)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current version: %d, Dirty: %t\n", version, dirty)
		return
	default:
		log.Fatal("Invalid action. Use: up, down, force, version")
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			log.Fatal("Migration failed:", err)
		}
	} else {
		fmt.Printf("Migration %s completed successfully\n", *action)
	}
}
