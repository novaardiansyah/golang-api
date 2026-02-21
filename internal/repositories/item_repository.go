package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *ItemRepository) FindAllPaginated(page, limit int) ([]models.Item, error) {
	var items []models.Item
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, err
}

func (r *ItemRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Item{}).Count(&count).Error
	return count, err
}

func (r *ItemRepository) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) FindByCode(code string) (*models.Item, error) {
	var item models.Item
	err := r.db.Where("code = ?", code).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *ItemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *ItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.Item{}, id).Error
}
