package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentGoalRepository struct {
	db *gorm.DB
}

func NewPaymentGoalRepository(db *gorm.DB) *PaymentGoalRepository {
	return &PaymentGoalRepository{db: db}
}

func (r *PaymentGoalRepository) FindAll(userID uint) ([]models.PaymentGoal, error) {
	var goals []models.PaymentGoal

	err := r.db.
		Preload("Status").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&goals).Error

	return goals, err
}

func (r *PaymentGoalRepository) FindByID(id int, userID uint) (*models.PaymentGoal, error) {
	var goal models.PaymentGoal

	err := r.db.
		Preload("Status").
		Where("user_id = ?", userID).
		First(&goal, id).Error

	return &goal, err
}
