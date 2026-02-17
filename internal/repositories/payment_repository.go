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
	UserID    uint
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
		Select("payments.*, (SELECT COUNT(*) FROM payment_item WHERE payment_item.payment_id = payments.id) as items_count").
		Preload("PaymentType").
		Preload("PaymentAccount").
		Preload("PaymentAccountTo")

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

	if filter.UserID > 0 {
		query = query.Where("payments.user_id = ?", filter.UserID)
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

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *PaymentRepository) FindByID(id int) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.
		Preload("PaymentType").
		Preload("PaymentAccount").
		Preload("PaymentAccountTo").
		First(&payment, id).Error

	return &payment, err
}

func (r *PaymentRepository) Create(tx *gorm.DB, payment *models.Payment) (*models.Payment, error) {
	if err := tx.Create(payment).Error; err != nil {
		return nil, err
	}

	if err := tx.
		Preload("PaymentType").
		Preload("PaymentAccount").
		Preload("PaymentAccountTo").
		First(payment, payment.ID).Error; err != nil {
		return nil, err
	}

	return payment, nil
}
