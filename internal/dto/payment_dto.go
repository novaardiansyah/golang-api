package dto

import (
	"encoding/json"
	"time"
)

type StorePaymentRequest struct {
	Name               *string `json:"name"`
	Amount             *int64  `json:"amount"`
	TypeID             uint    `json:"type_id"`
	Date               string  `json:"date"`
	PaymentAccountID   uint    `json:"payment_account_id"`
	PaymentAccountToID *uint   `json:"payment_account_to_id"`
	HasItems           bool    `json:"has_items"`
	IsDraft            bool    `json:"is_draft"`
	IsScheduled        bool    `json:"is_scheduled"`
}

type PaymentLogProperties struct {
	ID                 uint            `json:"id"`
	UserID             uint            `json:"user_id"`
	Code               string          `json:"code"`
	Name               *string         `json:"name"`
	Date               time.Time       `json:"date"`
	Amount             *int64          `json:"amount"`
	HasItems           bool            `json:"has_items"`
	IsScheduled        bool            `json:"is_scheduled"`
	IsDraft            bool            `json:"is_draft"`
	Attachments        json.RawMessage `json:"-"`
	TypeID             uint            `json:"type_id"`
	PaymentAccountID   uint            `json:"payment_account_id"`
	PaymentAccountToID *uint           `json:"payment_account_to_id"`
}
