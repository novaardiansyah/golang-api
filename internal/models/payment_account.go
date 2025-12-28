package models

type PaymentAccount struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	Deposit int64  `json:"deposit"`
}

func (PaymentAccount) TableName() string {
	return "payment_accounts"
}
