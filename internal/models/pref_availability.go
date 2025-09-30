package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrefAvailability represents the pref_availability table
type PrefAvailability struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`
	Mon       bool      `json:"mon" gorm:"type:boolean;not null;default:false"`
	Tue       bool      `json:"tue" gorm:"type:boolean;not null;default:false"`
	Wed       bool      `json:"wed" gorm:"type:boolean;not null;default:false"`
	Thu       bool      `json:"thu" gorm:"type:boolean;not null;default:false"`
	Fri       bool      `json:"fri" gorm:"type:boolean;not null;default:false"`
	Sat       bool      `json:"sat" gorm:"type:boolean;not null;default:false"`
	Sun       bool      `json:"sun" gorm:"type:boolean;not null;default:false"`
	AllDay    bool      `json:"all_day" gorm:"type:boolean;not null;default:true"`
	Morning   bool      `json:"morning" gorm:"type:boolean;not null;default:false"`
	Afternoon bool      `json:"afternoon" gorm:"type:boolean;not null;default:false"`
	TimeRange *string   `json:"time_range" gorm:"type:tstzrange"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for PrefAvailability
func (PrefAvailability) TableName() string {
	return "pref_availability"
}

// BeforeCreate hook for PrefAvailability
func (pa *PrefAvailability) BeforeCreate(tx *gorm.DB) error {
	if pa.ID == uuid.Nil {
		pa.ID = uuid.New()
	}
	return nil
}

// IsAvailableOn checks if user is available on a specific day
func (pa *PrefAvailability) IsAvailableOn(day time.Weekday) bool {
	switch day {
	case time.Monday:
		return pa.Mon
	case time.Tuesday:
		return pa.Tue
	case time.Wednesday:
		return pa.Wed
	case time.Thursday:
		return pa.Thu
	case time.Friday:
		return pa.Fri
	case time.Saturday:
		return pa.Sat
	case time.Sunday:
		return pa.Sun
	default:
		return false
	}
}

// IsAvailableAllDay checks if user is available all day
func (pa *PrefAvailability) IsAvailableAllDay() bool {
	return pa.AllDay
}

// IsAvailableMorning checks if user is available in the morning
func (pa *PrefAvailability) IsAvailableMorning() bool {
	return pa.Morning
}

// IsAvailableAfternoon checks if user is available in the afternoon
func (pa *PrefAvailability) IsAvailableAfternoon() bool {
	return pa.Afternoon
}

// GetAvailableDays returns a slice of available days
func (pa *PrefAvailability) GetAvailableDays() []time.Weekday {
	var days []time.Weekday

	if pa.Mon {
		days = append(days, time.Monday)
	}
	if pa.Tue {
		days = append(days, time.Tuesday)
	}
	if pa.Wed {
		days = append(days, time.Wednesday)
	}
	if pa.Thu {
		days = append(days, time.Thursday)
	}
	if pa.Fri {
		days = append(days, time.Friday)
	}
	if pa.Sat {
		days = append(days, time.Saturday)
	}
	if pa.Sun {
		days = append(days, time.Sunday)
	}

	return days
}

// GetAvailableDaysCount returns the number of available days
func (pa *PrefAvailability) GetAvailableDaysCount() int {
	count := 0
	if pa.Mon {
		count++
	}
	if pa.Tue {
		count++
	}
	if pa.Wed {
		count++
	}
	if pa.Thu {
		count++
	}
	if pa.Fri {
		count++
	}
	if pa.Sat {
		count++
	}
	if pa.Sun {
		count++
	}
	return count
}

// IsAvailable checks if user is available on a specific day and time
func (pa *PrefAvailability) IsAvailable(day time.Weekday, hour int) bool {
	if !pa.IsAvailableOn(day) {
		return false
	}

	if pa.IsAvailableAllDay() {
		return true
	}

	// Check time preferences
	if hour >= 6 && hour < 12 && pa.IsAvailableMorning() {
		return true
	}
	if hour >= 12 && hour < 18 && pa.IsAvailableAfternoon() {
		return true
	}

	return false
}
