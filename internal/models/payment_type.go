package models

type PaymentType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (PaymentType) TableName() string {
	return "payment_types"
}
