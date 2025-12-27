package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) FindAllPaginated(page, limit int) ([]models.Payment, error) {
	var payments []models.Payment
	offset := (page - 1) * limit
	err := r.db.
		Select("payments.*, LOWER(payment_types.name) as type, (SELECT COUNT(*) FROM payment_item WHERE payment_item.payment_id = payments.id) as items_count, pa.id as account_id, pa.name as account_name, pa_to.id as account_to_id, pa_to.name as account_to_name").
		Joins("INNER JOIN payment_types ON payment_types.id = payments.type_id").
		Joins("INNER JOIN payment_accounts pa ON pa.id = payments.payment_account_id").
		Joins("LEFT JOIN payment_accounts pa_to ON pa_to.id = payments.payment_account_to_id").
		Offset(offset).Limit(limit).Order("updated_at desc").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Payment{}).Count(&count).Error
	return count, err
}
