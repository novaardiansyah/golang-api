package models

type PaymentAccount struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (PaymentAccount) TableName() string {
	return "payment_accounts"
}
