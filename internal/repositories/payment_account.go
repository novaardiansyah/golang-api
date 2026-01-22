/*
 * Project Name: repositories
 * File: payment_account.go
 * Created Date: Thursday January 22nd 2026
 * 
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 * 
 * Copyright (c) 2025-2026 Nova Ardiansyah, Org
 */

package repositories

import (
	"golang-api/internal/models"

	"gorm.io/gorm"
)

type PaymentAccountRepository struct {
  db *gorm.DB
}

func NewPaymentAccountRepository(db *gorm.DB) *PaymentAccountRepository {
  return &PaymentAccountRepository{db: db}
}

func (r *PaymentAccountRepository) Count(userID uint) (int64, error) {
  var count int64
  err := r.db.Model(&models.PaymentAccount{}).Where("user_id = ?", userID).Count(&count).Error
  return count, err
}

func (r *PaymentAccountRepository) FindAllPaginated(userID uint, page, limit int) ([]models.PaymentAccount, error) {
  var paymentAccounts []models.PaymentAccount
  
  offset := (page - 1) * limit

  err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&paymentAccounts).Error

  return paymentAccounts, err
}