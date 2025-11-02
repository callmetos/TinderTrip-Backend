package models

import (
	"time"

	"github.com/google/uuid"
)

// FoodCategoryMaster represents the food_categories master table
type FoodCategoryMaster struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code        string    `json:"code" gorm:"type:varchar(50);not null;uniqueIndex"`
	DisplayName string    `json:"display_name" gorm:"type:text;not null"`
	Icon        *string   `json:"icon,omitempty" gorm:"type:text"`
	Description *string   `json:"description,omitempty" gorm:"type:text"`
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}

// TableName returns the table name for FoodCategoryMaster
func (FoodCategoryMaster) TableName() string {
	return "food_categories"
}

