package models

type PaymentType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (PaymentType) TableName() string {
	return "payment_types"
}

const (
	PaymentTypeExpense    uint = 1
	PaymentTypeIncome     uint = 2
	PaymentTypeTransfer   uint = 3
	PaymentTypeWithdrawal uint = 4
)
