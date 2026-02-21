package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentItemRepository struct {
	db *gorm.DB
}

func NewPaymentItemRepository(db *gorm.DB) *PaymentItemRepository {
	return &PaymentItemRepository{db: db}
}

func (r *PaymentItemRepository) FindAll() ([]models.PaymentItem, error) {
	var paymentItems []models.PaymentItem
	err := r.db.Find(&paymentItems).Error
	return paymentItems, err
}

func (r *PaymentItemRepository) FindByPaymentID(paymentID uint) ([]models.PaymentItem, error) {
	var paymentItems []models.PaymentItem
	err := r.db.Where("payment_id = ?", paymentID).Preload("Item").Find(&paymentItems).Error
	return paymentItems, err
}

func (r *PaymentItemRepository) FindByID(id uint) (*models.PaymentItem, error) {
	var paymentItem models.PaymentItem
	err := r.db.First(&paymentItem, id).Error
	if err != nil {
		return nil, err
	}
	return &paymentItem, nil
}

func (r *PaymentItemRepository) Create(tx *gorm.DB, paymentItem *models.PaymentItem) error {
	if tx != nil {
		return tx.Create(paymentItem).Error
	}
	return r.db.Create(paymentItem).Error
}

func (r *PaymentItemRepository) CreateBatch(tx *gorm.DB, paymentItems []models.PaymentItem) error {
	if tx != nil {
		return tx.Create(&paymentItems).Error
	}
	return r.db.Create(&paymentItems).Error
}

func (r *PaymentItemRepository) Update(paymentItem *models.PaymentItem) error {
	return r.db.Save(paymentItem).Error
}

func (r *PaymentItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.PaymentItem{}, id).Error
}

func (r *PaymentItemRepository) DeleteByPaymentID(tx *gorm.DB, paymentID uint) error {
	if tx != nil {
		return tx.Where("payment_id = ?", paymentID).Delete(&models.PaymentItem{}).Error
	}
	return r.db.Where("payment_id = ?", paymentID).Delete(&models.PaymentItem{}).Error
}
