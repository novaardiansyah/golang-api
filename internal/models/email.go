package models

import (
	"time"

	"gorm.io/gorm"
)

type Email struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UID             string         `gorm:"size:255;uniqueIndex;not null" json:"uid"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Email           string         `gorm:"size:255;not null" json:"email"`
	EmailSubject    string         `gorm:"column:subject;size:255;not null" json:"subject"`
	Message         string         `gorm:"type:text" json:"message"`
	Status          string         `gorm:"size:50" json:"status"`
	IsURLAttachment *bool          `gorm:"default:false" json:"is_url_attachment"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (Email) TableName() string {
	return "emails"
}
