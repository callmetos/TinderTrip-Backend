package models_test

import (
	"testing"
	"time"

	"TinderTrip-Backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvent_IsPublished(t *testing.T) {
	tests := []struct {
		name   string
		status models.EventStatus
		want   bool
	}{
		{
			name:   "Published status",
			status: models.EventStatusPublished,
			want:   true,
		},
		{
			name:   "Cancelled status",
			status: models.EventStatusCancelled,
			want:   false,
		},
		{
			name:   "Completed status",
			status: models.EventStatusCompleted,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{Status: tt.status}
			assert.Equal(t, tt.want, event.IsPublished())
		})
	}
}

func TestEvent_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status models.EventStatus
		want   bool
	}{
		{
			name:   "Completed status",
			status: models.EventStatusCompleted,
			want:   true,
		},
		{
			name:   "Published status",
			status: models.EventStatusPublished,
			want:   false,
		},
		{
			name:   "Cancelled status",
			status: models.EventStatusCancelled,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{Status: tt.status}
			assert.Equal(t, tt.want, event.IsCompleted())
		})
	}
}

func TestEvent_IsCancelled(t *testing.T) {
	tests := []struct {
		name   string
		status models.EventStatus
		want   bool
	}{
		{
			name:   "Cancelled status",
			status: models.EventStatusCancelled,
			want:   true,
		},
		{
			name:   "Published status",
			status: models.EventStatusPublished,
			want:   false,
		},
		{
			name:   "Completed status",
			status: models.EventStatusCompleted,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{Status: tt.status}
			assert.Equal(t, tt.want, event.IsCancelled())
		})
	}
}

func TestEvent_HasLocation(t *testing.T) {
	lat := 13.7563
	lng := 100.5018

	tests := []struct {
		name string
		lat  *float64
		lng  *float64
		want bool
	}{
		{
			name: "Has both lat and lng",
			lat:  &lat,
			lng:  &lng,
			want: true,
		},
		{
			name: "Has lat only",
			lat:  &lat,
			lng:  nil,
			want: false,
		},
		{
			name: "Has lng only",
			lat:  nil,
			lng:  &lng,
			want: false,
		},
		{
			name: "Has neither",
			lat:  nil,
			lng:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{Lat: tt.lat, Lng: tt.lng}
			assert.Equal(t, tt.want, event.HasLocation())
		})
	}
}

func TestEvent_GetDuration(t *testing.T) {
	now := time.Now()
	twoHoursLater := now.Add(2 * time.Hour)

	tests := []struct {
		name    string
		startAt *time.Time
		endAt   *time.Time
		want    *time.Duration
	}{
		{
			name:    "Valid duration",
			startAt: &now,
			endAt:   &twoHoursLater,
			want: func() *time.Duration {
				d := 2 * time.Hour
				return &d
			}(),
		},
		{
			name:    "No start time",
			startAt: nil,
			endAt:   &twoHoursLater,
			want:    nil,
		},
		{
			name:    "No end time",
			startAt: &now,
			endAt:   nil,
			want:    nil,
		},
		{
			name:    "Neither start nor end",
			startAt: nil,
			endAt:   nil,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{StartAt: tt.startAt, EndAt: tt.endAt}
			duration := event.GetDuration()

			if tt.want == nil {
				assert.Nil(t, duration)
			} else {
				assert.NotNil(t, duration)
				assert.Equal(t, *tt.want, *duration)
			}
		})
	}
}

func TestEvent_IsUpcoming(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	past := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name    string
		startAt *time.Time
		want    bool
	}{
		{
			name:    "Future event",
			startAt: &future,
			want:    true,
		},
		{
			name:    "Past event",
			startAt: &past,
			want:    false,
		},
		{
			name:    "No start time",
			startAt: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{StartAt: tt.startAt}
			assert.Equal(t, tt.want, event.IsUpcoming())
		})
	}
}

func TestEvent_IsPast(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	past := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name  string
		endAt *time.Time
		want  bool
	}{
		{
			name:  "Past event",
			endAt: &past,
			want:  true,
		},
		{
			name:  "Future event",
			endAt: &future,
			want:  false,
		},
		{
			name:  "No end time",
			endAt: nil,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.Event{EndAt: tt.endAt}
			assert.Equal(t, tt.want, event.IsPast())
		})
	}
}

func TestEvent_BeforeCreate(t *testing.T) {
	t.Run("Generates UUID if not set", func(t *testing.T) {
		event := &models.Event{}
		err := event.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, event.ID)
	})

	t.Run("Does not override existing UUID", func(t *testing.T) {
		existingID := uuid.New()
		event := &models.Event{ID: existingID}
		err := event.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.Equal(t, existingID, event.ID)
	})
}

func TestEvent_TableName(t *testing.T) {
	event := models.Event{}
	assert.Equal(t, "events", event.TableName())
}

func TestEventType_Constants(t *testing.T) {
	assert.Equal(t, models.EventType("meal"), models.EventTypeMeal)
	assert.Equal(t, models.EventType("daytrip"), models.EventTypeDaytrip)
	assert.Equal(t, models.EventType("overnight"), models.EventTypeOvernight)
	assert.Equal(t, models.EventType("activity"), models.EventTypeActivity)
	assert.Equal(t, models.EventType("other"), models.EventTypeOther)
}

func TestEventStatus_Constants(t *testing.T) {
	assert.Equal(t, models.EventStatus("published"), models.EventStatusPublished)
	assert.Equal(t, models.EventStatus("cancelled"), models.EventStatusCancelled)
	assert.Equal(t, models.EventStatus("completed"), models.EventStatusCompleted)
}

func TestEvent_CompleteScenario(t *testing.T) {
	t.Run("Complete meal event", func(t *testing.T) {
		now := time.Now()
		twoHoursLater := now.Add(2 * time.Hour)
		lat := 13.7563
		lng := 100.5018
		capacity := 5
		budgetMin := 300
		budgetMax := 500

		event := &models.Event{
			Title:       "Dinner at Italian Restaurant",
			Description: new(string),
			EventType:   models.EventTypeMeal,
			Lat:         &lat,
			Lng:         &lng,
			StartAt:     &now,
			EndAt:       &twoHoursLater,
			Capacity:    &capacity,
			BudgetMin:   &budgetMin,
			BudgetMax:   &budgetMax,
			Status:      models.EventStatusPublished,
		}

		assert.True(t, event.IsPublished())
		assert.False(t, event.IsCompleted())
		assert.False(t, event.IsCancelled())
		assert.True(t, event.HasLocation())
		assert.NotNil(t, event.GetDuration())
		assert.Equal(t, 2*time.Hour, *event.GetDuration())
	})

	t.Run("Upcoming daytrip event", func(t *testing.T) {
		future := time.Now().Add(7 * 24 * time.Hour)
		futureEnd := future.Add(8 * time.Hour)

		event := &models.Event{
			Title:     "Beach Day Trip",
			EventType: models.EventTypeDaytrip,
			StartAt:   &future,
			EndAt:     &futureEnd,
			Status:    models.EventStatusPublished,
		}

		assert.True(t, event.IsUpcoming())
		assert.False(t, event.IsPast())
	})

	t.Run("Completed event", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		pastEnd := past.Add(-20 * time.Hour)

		event := &models.Event{
			Title:     "Past Event",
			EventType: models.EventTypeActivity,
			StartAt:   &past,
			EndAt:     &pastEnd,
			Status:    models.EventStatusCompleted,
		}

		assert.False(t, event.IsPublished())
		assert.True(t, event.IsCompleted())
		assert.True(t, event.IsPast())
		assert.False(t, event.IsUpcoming())
	})
}
