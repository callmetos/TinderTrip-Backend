package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrefBudget represents the pref_budget table
type PrefBudget struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`
	MealMin      *int      `json:"meal_min" gorm:"type:int"`
	MealMax      *int      `json:"meal_max" gorm:"type:int"`
	DaytripMin   *int      `json:"daytrip_min" gorm:"type:int"`
	DaytripMax   *int      `json:"daytrip_max" gorm:"type:int"`
	OvernightMin *int      `json:"overnight_min" gorm:"type:int"`
	OvernightMax *int      `json:"overnight_max" gorm:"type:int"`
	Unlimited    bool      `json:"unlimited" gorm:"type:boolean;not null;default:false"`
	Currency     string    `json:"currency" gorm:"type:text;not null;default:'THB'"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for PrefBudget
func (PrefBudget) TableName() string {
	return "pref_budget"
}

// BeforeCreate hook for PrefBudget
func (pb *PrefBudget) BeforeCreate(tx *gorm.DB) error {
	if pb.ID == uuid.Nil {
		pb.ID = uuid.New()
	}
	return nil
}

// IsUnlimited checks if the budget is unlimited
func (pb *PrefBudget) IsUnlimited() bool {
	return pb.Unlimited
}

// GetMealBudget returns the meal budget range
func (pb *PrefBudget) GetMealBudget() (min, max *int) {
	return pb.MealMin, pb.MealMax
}

// GetDaytripBudget returns the daytrip budget range
func (pb *PrefBudget) GetDaytripBudget() (min, max *int) {
	return pb.DaytripMin, pb.DaytripMax
}

// GetOvernightBudget returns the overnight budget range
func (pb *PrefBudget) GetOvernightBudget() (min, max *int) {
	return pb.OvernightMin, pb.OvernightMax
}

// GetBudgetForEventType returns the budget range for a specific event type
func (pb *PrefBudget) GetBudgetForEventType(eventType EventType) (min, max *int) {
	switch eventType {
	case EventTypeMeal:
		return pb.GetMealBudget()
	case EventTypeDaytrip:
		return pb.GetDaytripBudget()
	case EventTypeOvernight:
		return pb.GetOvernightBudget()
	default:
		return nil, nil
	}
}

// IsWithinBudget checks if an amount is within the budget range
func (pb *PrefBudget) IsWithinBudget(amount int, eventType EventType) bool {
	if pb.IsUnlimited() {
		return true
	}

	min, max := pb.GetBudgetForEventType(eventType)

	if min != nil && amount < *min {
		return false
	}
	if max != nil && amount > *max {
		return false
	}

	return true
}

// GetBudgetRangeString returns a string representation of the budget range
func (pb *PrefBudget) GetBudgetRangeString(eventType EventType) string {
	if pb.IsUnlimited() {
		return "Unlimited"
	}

	min, max := pb.GetBudgetForEventType(eventType)

	if min == nil && max == nil {
		return "Not specified"
	}

	if min == nil {
		return fmt.Sprintf("Up to %d %s", *max, pb.Currency)
	}

	if max == nil {
		return fmt.Sprintf("From %d %s", *min, pb.Currency)
	}

	return fmt.Sprintf("%d - %d %s", *min, *max, pb.Currency)
}
