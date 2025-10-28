package database

import (
	"fmt"
	"log"
	"time"

	"TinderTrip-Backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() error {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate models - DISABLED
	// All database tables already exist and are properly configured
	// AutoMigrate is disabled to prevent "insufficient arguments" errors
	// If you need to add new tables, uncomment and add only the new models
	/*
		err = DB.AutoMigrate(
			&models.User{},
			&models.UserProfile{},
			&models.PasswordReset{},
			&models.Tag{},
			&models.Event{},
			&models.EventPhoto{},
			&models.EventMember{},
			&models.EventSwipe{},
			&models.EventCategory{},
			&models.EventTag{},
			&models.ChatRoom{},
			&models.ChatMessage{},
			&models.UserEventHistory{},
			&models.AuditLog{},
			&models.APILog{},
			&models.PrefAvailability{},
			&models.PrefBudget{},
		)
		if err != nil {
			return fmt.Errorf("failed to auto-migrate models: %w", err)
		}
	*/

	log.Println("Successfully connected to PostgreSQL database")
	return nil
}

func ClosePostgres() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

// SetTestDB sets a test database instance (for testing purposes only)
func SetTestDB(testDB *gorm.DB) {
	DB = testDB
}
