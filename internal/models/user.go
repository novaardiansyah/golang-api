package models

import (
	"time"

	"gorm.io/gorm"
)

// User model - mirip dengan Eloquent Model di Laravel
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // "-" untuk hide password dari JSON response
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}
