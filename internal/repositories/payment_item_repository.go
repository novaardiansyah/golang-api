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

func (r *PaymentItemRepository) CountByPaymentID(paymentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.PaymentItem{}).Where("payment_id = ?", paymentID).Count(&count).Error
	return count, err
}

func (r *PaymentItemRepository) FindByPaymentIDPaginated(paymentID uint, page, limit int) ([]models.PaymentItem, error) {
	var paymentItems []models.PaymentItem
	offset := (page - 1) * limit
	err := r.db.Where("payment_id = ?", paymentID).Preload("Item").Offset(offset).Limit(limit).Order("updated_at desc").Find(&paymentItems).Error
	return paymentItems, err
}

func (r *PaymentItemRepository) GetSummaryByPaymentID(paymentID uint) (*PaymentItemSummary, error) {
	var summary PaymentItemSummary
	err := r.db.Model(&models.PaymentItem{}).
		Select("payment_id, COUNT(*) as total_items, SUM(quantity) as total_qty, SUM(total) as total_amount").
		Where("payment_id = ?", paymentID).
		Group("payment_id").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

type PaymentItemSummary struct {
	PaymentID   uint  `gorm:"column:payment_id"`
	TotalItems  int64 `gorm:"column:total_items"`
	TotalQty    int64 `gorm:"column:total_qty"`
	TotalAmount int64 `gorm:"column:total_amount"`
}
