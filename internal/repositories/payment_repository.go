package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentFilter struct {
	DateFrom  string
	DateTo    string
	Type      int
	AccountID int
	Search    string
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) FindAllPaginated(page, limit int, filter PaymentFilter) ([]models.Payment, error) {
	var payments []models.Payment
	offset := (page - 1) * limit

	query := r.db.
		Select("payments.*, LOWER(payment_types.name) as type, (SELECT COUNT(*) FROM payment_item WHERE payment_item.payment_id = payments.id) as items_count, pa.id as account_id, pa.name as account_name, pa_to.id as account_to_id, pa_to.name as account_to_name").
		Joins("INNER JOIN payment_types ON payment_types.id = payments.type_id").
		Joins("INNER JOIN payment_accounts pa ON pa.id = payments.payment_account_id").
		Joins("LEFT JOIN payment_accounts pa_to ON pa_to.id = payments.payment_account_to_id")

	if filter.DateFrom != "" {
		query = query.Where("payments.date >= ?", filter.DateFrom)
	}

	if filter.DateTo != "" {
		query = query.Where("payments.date <= ?", filter.DateTo)
	}

	if filter.Type > 0 {
		query = query.Where("payments.type_id = ?", filter.Type)
	}

	if filter.AccountID > 0 {
		query = query.Where("payments.payment_account_id = ?", filter.AccountID)
	}

	if filter.Search != "" {
		query = query.Where("payments.name LIKE ?", "%"+filter.Search+"%")
	}

	err := query.Offset(offset).Limit(limit).Order("updated_at desc").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Count(filter PaymentFilter) (int64, error) {
	var count int64

	query := r.db.Model(&models.Payment{})

	if filter.DateFrom != "" {
		query = query.Where("date >= ?", filter.DateFrom)
	}

	if filter.DateTo != "" {
		query = query.Where("date <= ?", filter.DateTo)
	}

	if filter.Type > 0 {
		query = query.Where("type_id = ?", filter.Type)
	}

	if filter.AccountID > 0 {
		query = query.Where("payment_account_id = ?", filter.AccountID)
	}

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}

	err := query.Count(&count).Error
	return count, err
}
