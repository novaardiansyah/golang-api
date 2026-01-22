/*
 * Project Name: models
 * File: payment_account.go
 * Created Date: Sunday December 28th 2025
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/golang-api/blob/main/LICENSE
 *
 * Copyright (c) 2025-2026 Nova Ardiansyah, Org
 */

package models

import (
	"golang-api/pkg/utils"

	"gorm.io/gorm"
)

type PaymentAccount struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"-"`
	Name      string    `json:"name"`
	Deposit   int64     `json:"deposit"`
	Logo      string    `json:"logo"`
	Formatted Formatted `gorm:"-" json:"formatted"`
}

func (PaymentAccount) TableName() string {
	return "payment_accounts"
}

type Formatted struct {
	Deposit string `json:"deposit"`
	Logo    string `json:"logo"`
}

func (p *PaymentAccount) AfterFind(tx *gorm.DB) (err error) {
	var logo string

	if p.Logo != "" {
		logo = utils.GetExternalUrl("main", "storage/"+p.Logo)
	}

	p.Formatted = Formatted{
		Deposit: utils.FormatRupiah(p.Deposit),
		Logo:    logo,
	}

	return
}
