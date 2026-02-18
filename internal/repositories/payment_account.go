/*
 * Project Name: repositories
 * File: payment_account.go
 * Created Date: Thursday January 22nd 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package repositories

import (
	"encoding/json"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/pkg/utils"

	"gorm.io/gorm"
)

type PaymentAccountRepository struct {
	activityLogRepository *ActivityLogRepository
	db                    *gorm.DB
}

func NewPaymentAccountRepository(db *gorm.DB) *PaymentAccountRepository {
	return &PaymentAccountRepository{
		activityLogRepository: NewActivityLogRepository(db),
		db:                    db,
	}
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

func (r *PaymentAccountRepository) Update(tx *gorm.DB, userId uint, paymentAccount *models.PaymentAccount, prevPaymentAccount *models.PaymentAccount) (*models.PaymentAccount, error) {
	var err error
	before := prevPaymentAccount

	if prevPaymentAccount == nil {
		before, err = r.SelectByID(paymentAccount.ID, []string{"id", "user_id", "name", "deposit"})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Where("id = ?", paymentAccount.ID).Updates(paymentAccount).Error; err != nil {
		return nil, err
	}

	prevDifference := before.Deposit - paymentAccount.Deposit
	decreaseDeposit := paymentAccount.Deposit < before.Deposit

	if decreaseDeposit {
		prevDifference = -prevDifference
	}

	difference := -prevDifference

	logPrevProps := dto.PaymentAccountLogProperties{
		ID:         before.ID,
		UserID:     before.UserID,
		Name:       before.Name,
		Deposit:    before.Deposit,
		Difference: &prevDifference,
	}

	logProps := dto.PaymentAccountLogProperties{
		ID:         paymentAccount.ID,
		UserID:     paymentAccount.UserID,
		Name:       paymentAccount.Name,
		Deposit:    paymentAccount.Deposit,
		Difference: &difference,
	}

	properties, _ := json.Marshal(logProps)
	prevProperties, _ := json.Marshal(logPrevProps)

	err = r.activityLogRepository.Store(&models.ActivityLog{
		Event:          "Updated",
		LogName:        "Resource",
		Description:    "Payment Account Updated by Nova Ardiansyah (Hardcode)",
		SubjectType:    utils.String("App\\Models\\PaymentAccount"),
		SubjectID:      &paymentAccount.ID,
		CauserType:     "App\\Models\\User",
		CauserID:       userId,
		PrevProperties: (*json.RawMessage)(&prevProperties),
		Properties:     properties,
	})

	if err != nil {
		return nil, err
	}

	return paymentAccount, nil
}

func (r *PaymentAccountRepository) FindByID(id uint) (*models.PaymentAccount, error) {
	var paymentAccount models.PaymentAccount
	err := r.db.First(&paymentAccount, id).Error
	return &paymentAccount, err
}

// ! saya ingin select bukan * tapi field sesuai dto aja
func (r *PaymentAccountRepository) SelectByID(id uint, fields []string) (*models.PaymentAccount, error) {
	var paymentAccount models.PaymentAccount
	err := r.db.Select(fields).First(&paymentAccount, id).Error
	return &paymentAccount, err
}
