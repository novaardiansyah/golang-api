package dto

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
