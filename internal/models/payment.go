package models

import (
	"encoding/json"
	"fmt"
	"golang-api/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))), nil
}

type Payment struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	UserID             uint            `json:"user_id"`
	Code               string          `json:"code"`
	Name               *string         `json:"name"`
	Date               time.Time       `json:"date"`
	Amount             *int64          `json:"amount"`
	HasItems           bool            `json:"has_items"`
	IsScheduled        bool            `json:"is_scheduled"`
	IsDraft            bool            `json:"is_draft"`
	Attachments        json.RawMessage `gorm:"type:json" json:"-"`
	TypeID             uint            `json:"type_id"`
	PaymentAccountID   uint            `json:"payment_account_id"`
	PaymentAccountToID *uint           `json:"payment_account_to_id"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`

	PaymentType      *PaymentType    `gorm:"foreignKey:TypeID" json:"-"`
	PaymentAccount   *PaymentAccount `gorm:"foreignKey:PaymentAccountID" json:"-"`
	PaymentAccountTo *PaymentAccount `gorm:"foreignKey:PaymentAccountToID" json:"-"`

	Type               string       `gorm:"-" json:"type"`
	FormattedAmount    string       `gorm:"-" json:"formatted_amount"`
	FormattedDate      string       `gorm:"-" json:"formatted_date"`
	FormattedUpdatedAt string       `gorm:"-" json:"formatted_updated_at"`
	AttachmentsCount   int          `gorm:"-" json:"attachments_count"`
	ItemsCount         int          `gorm:"->" json:"items_count"`
	Account            *AccountInfo `gorm:"-" json:"account"`
	AccountTo          *AccountInfo `gorm:"-" json:"account_to"`
}

type AccountInfo struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) GetAttachmentsCount() int {
	if len(p.Attachments) == 0 {
		return 0
	}

	var attachments []interface{}
	if err := json.Unmarshal(p.Attachments, &attachments); err != nil {
		return 0
	}
	return len(attachments)
}

func (p *Payment) AfterFind(tx *gorm.DB) (err error) {
	amount := int64(0)

	if p.Amount != nil {
		amount = *p.Amount
	}

	p.FormattedUpdatedAt = utils.FormatDateID(p.UpdatedAt, "Monday, 2 Jan 2006, 15.04 WIB")
	p.FormattedAmount = utils.FormatRupiah(amount)
	p.FormattedDate = utils.FormatDateID(time.Time(p.Date), "Monday, 2 Jan 2006")
	p.AttachmentsCount = p.GetAttachmentsCount()

	if p.PaymentType != nil {
		p.Type = strings.ToLower(p.PaymentType.Name)
	}

	if p.PaymentAccount != nil {
		p.Account = &AccountInfo{ID: &p.PaymentAccount.ID, Name: &p.PaymentAccount.Name}
	} else {
		p.Account = &AccountInfo{}
	}

	if p.PaymentAccountTo != nil {
		p.AccountTo = &AccountInfo{ID: &p.PaymentAccountTo.ID, Name: &p.PaymentAccountTo.Name}
	} else {
		p.AccountTo = &AccountInfo{}
	}

	return
}
