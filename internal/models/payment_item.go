package models

import (
	"time"
)

type PaymentItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ItemCode  string    `gorm:"size:255" json:"item_code"`
	PaymentID uint      `json:"payment_id"`
	ItemID    uint      `json:"item_id"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Total     int64     `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Payment *Payment `gorm:"foreignKey:PaymentID" json:"-"`
	Item    *Item    `gorm:"foreignKey:ItemID" json:"-"`
}

func (PaymentItem) TableName() string {
	return "payment_item"
}
