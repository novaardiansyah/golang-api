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

type PaymentItemSummaryResponse struct {
	PaymentID       uint   `json:"payment_id"`
	PaymentCode     string `json:"payment_code"`
	TotalItems      int64  `json:"total_items"`
	TotalQty        int64  `json:"total_qty"`
	TotalAmount     int64  `json:"total_amount"`
	FormattedAmount string `json:"formatted_amount"`
}

type PaymentItemAttachedResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	TypeID         uint      `json:"type_id"`
	Type           string    `json:"type"`
	Code           string    `json:"code"`
	Price          int64     `json:"price"`
	Quantity       int       `json:"quantity"`
	Total          int64     `json:"total"`
	FormattedPrice string    `json:"formatted_price"`
	FormattedTotal string    `json:"formatted_total"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaymentItemSummary struct {
	PaymentID   uint  `gorm:"column:payment_id" json:"payment_id"`
	TotalItems  int64 `gorm:"column:total_items" json:"total_items"`
	TotalQty    int64 `gorm:"column:total_qty" json:"total_qty"`
	TotalAmount int64 `gorm:"column:total_amount" json:"total_amount"`
}
