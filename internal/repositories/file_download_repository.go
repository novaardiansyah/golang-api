package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type FileDownloadRepository struct {
	db *gorm.DB
}

func NewFileDownloadRepository(db *gorm.DB) *FileDownloadRepository {
	return &FileDownloadRepository{db: db}
}

func (r *FileDownloadRepository) GetByUID(uid string) (*models.FileDownload, error) {
	var fileDownload models.FileDownload
	err := r.db.Preload("Files").Where("uid = ?", uid).First(&fileDownload).Error
	if err != nil {
		return nil, err
	}
	return &fileDownload, nil
}
