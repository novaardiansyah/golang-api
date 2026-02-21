/*
 * Project Name: repositories
 * File: payment_type_repository.go
 * Created Date: Saturday February 21st 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentTypeRepository struct {
	db *gorm.DB
}

func NewPaymentTypeRepository(db *gorm.DB) *PaymentTypeRepository {
	return &PaymentTypeRepository{db: db}
}

func (r *PaymentTypeRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.PaymentType{}).Count(&count).Error
	return count, err
}

func (r *PaymentTypeRepository) FindAllPaginated(page, limit int) ([]models.PaymentType, error) {
	var paymentTypes []models.PaymentType

	offset := (page - 1) * limit

	err := r.db.Offset(offset).Limit(limit).Find(&paymentTypes).Error

	return paymentTypes, err
}
