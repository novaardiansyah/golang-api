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

func (r *PaymentGoalRepository) GetOverview(userID uint) (int64, int64, error) {
	var totalGoals int64
	var completedGoals int64

	err := r.db.Model(&models.PaymentGoal{}).
		Where("user_id = ?", userID).
		Count(&totalGoals).Error

	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.PaymentGoal{}).
		Where("user_id = ?", userID).
		Where("status_id = ?", models.PaymentGoalStatusCompleted).
		Count(&completedGoals).Error

	return totalGoals, completedGoals, err
}
