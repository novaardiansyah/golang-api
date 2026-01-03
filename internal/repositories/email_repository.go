package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type EmailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) GetByUID(uid string) (*models.Email, error) {
	var email models.Email
	err := r.db.Preload("Files").Where("uid = ?", uid).First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}
