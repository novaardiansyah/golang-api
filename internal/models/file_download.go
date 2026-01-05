package models

import (
	"time"

	"gorm.io/gorm"
)

type FileDownloadStatus string

const (
	FileDownloadStatusActive   FileDownloadStatus = "active"
	FileDownloadStatusInactive FileDownloadStatus = "inactive"
)

type FileDownload struct {
	ID            uint               `gorm:"primaryKey" json:"id"`
	UID           string             `gorm:"size:255;uniqueIndex;not null" json:"uid"`
	Code          string             `gorm:"size:255;uniqueIndex;not null" json:"code"`
	Status        FileDownloadStatus `gorm:"size:50;default:active" json:"status"`
	DownloadCount int                `gorm:"default:0" json:"download_count"`
	AccessCount   int                `gorm:"default:0" json:"access_count"`
	Files         []File             `gorm:"foreignKey:FileDownloadID" json:"files,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (FileDownload) TableName() string {
	return "file_downloads"
}
