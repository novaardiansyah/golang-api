package repositories

import (
  "golang-api/internal/models"
  "gorm.io/gorm"
)

type GalleryRepository struct {
  db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
  return &GalleryRepository{db: db}
}

func (r *GalleryRepository) FindAllPaginated(page, limit int) ([]models.Gallery, error) {
  var galleries []models.Gallery
  offset := (page - 1) * limit
  err := r.db.Offset(offset).Limit(limit).Find(&galleries).Error
  return galleries, err
}

func (r *GalleryRepository) Count() (int64, error) {
  var count int64
  err := r.db.Model(&models.Gallery{}).Count(&count).Error
  return count, err
}
