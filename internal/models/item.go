package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TypeID    uint           `json:"type_id"`
	Code      string         `gorm:"size:255" json:"code"`
	Name      string         `gorm:"size:255" json:"name"`
	Amount    int64          `json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (Item) TableName() string {
	return "items"
}
