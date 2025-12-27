package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

// ProductRepository handles data access for Product model
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates new product repository instance
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindAll retrieves all products with pagination
func (r *ProductRepository) FindAll(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * limit

	// Count total records
	r.db.Model(&models.Product{}).Count(&total)

	// Get paginated results
	err := r.db.Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

// FindByID retrieves product by ID
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create creates new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// Update updates existing product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
