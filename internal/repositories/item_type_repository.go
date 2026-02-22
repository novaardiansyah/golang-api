package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type ItemTypeRepository struct {
	db *gorm.DB
}

func NewItemTypeRepository(db *gorm.DB) *ItemTypeRepository {
	return &ItemTypeRepository{db: db}
}

func (r *ItemTypeRepository) FindAll() ([]models.ItemType, error) {
	var itemTypes []models.ItemType
	err := r.db.Find(&itemTypes).Error
	return itemTypes, err
}

func (r *ItemTypeRepository) FindByID(id uint) (*models.ItemType, error) {
	var itemType models.ItemType
	err := r.db.First(&itemType, id).Error
	if err != nil {
		return nil, err
	}
	return &itemType, nil
}
