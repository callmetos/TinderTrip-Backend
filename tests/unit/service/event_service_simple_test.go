package service_test

import (
	"testing"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/pkg/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSimpleEventServiceTest(t *testing.T) (*gorm.DB, *service.EventService) {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set global DB for testing
	database.DB = db

	eventService := service.NewEventService()
	return db, eventService
}

func TestEventService_Basic(t *testing.T) {
	_, eventService := setupSimpleEventServiceTest(t)

	t.Run("Create event service", func(t *testing.T) {
		assert.NotNil(t, eventService)
	})

	t.Run("Get public events with empty database", func(t *testing.T) {
		// This will fail gracefully since tables don't exist
		_, _, err := eventService.GetPublicEvents(1, 10, "")
		// We expect an error here since tables aren't created
		assert.Error(t, err)
	})

	t.Run("Get event with invalid ID", func(t *testing.T) {
		_, err := eventService.GetEvent("invalid-id", "user-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("Create event with invalid user ID", func(t *testing.T) {
		_, err := eventService.CreateEvent("invalid-id", dto.CreateEventRequest{
			Title:     "Test",
			EventType: string(models.EventTypeMeal),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("Update event with invalid ID", func(t *testing.T) {
		_, err := eventService.UpdateEvent("invalid-id", "user-id", dto.UpdateEventRequest{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("Delete event with invalid ID", func(t *testing.T) {
		err := eventService.DeleteEvent("invalid-id", "user-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("Join event with invalid ID", func(t *testing.T) {
		err := eventService.JoinEvent("invalid-id", "user-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("Leave event with invalid ID", func(t *testing.T) {
		err, _ := eventService.LeaveEvent("invalid-id", "user-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})

	t.Run("Complete event with invalid ID", func(t *testing.T) {
		err := eventService.CompleteEvent("invalid-id", "user-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event ID")
	})
}
