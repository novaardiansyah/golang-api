package models

import (
	"time"

	"gorm.io/gorm"
)

type ItemType struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (ItemType) TableName() string {
	return "item_types"
}
