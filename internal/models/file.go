package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Code                  string         `gorm:"size:255;uniqueIndex;not null" json:"code"`
	UserID                *uint          `gorm:"index" json:"user_id,omitempty"`
	FileDownloadID        *uint          `gorm:"index" json:"file_download_id,omitempty"`
	FileName              string         `gorm:"size:255;not null" json:"file_name"`
	FilePath              string         `gorm:"size:500;not null" json:"file_path"`
	FileSize              int64          `gorm:"not null" json:"file_size"`
	DownloadURL           string         `gorm:"size:500" json:"download_url"`
	FileAlias             string         `gorm:"size:255" json:"file_alias"`
	ScheduledDeletionTime *time.Time     `json:"scheduled_deletion_time,omitempty"`
	HasBeenDeleted        *bool          `gorm:"default:false" json:"has_been_deleted"`
	SubjectType           string         `gorm:"size:100;index:idx_subject" json:"subject_type"`
	SubjectID             uint           `gorm:"index:idx_subject" json:"subject_id"`
	FileDownload          *FileDownload  `gorm:"foreignKey:FileDownloadID" json:"file_download,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

func (File) TableName() string {
	return "files"
}
