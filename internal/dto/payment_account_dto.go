package dto

type PaymentAccountLogProperties struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"-"`
	Name       string `json:"name"`
	Deposit    int64  `json:"deposit"`
	Difference *int64 `json:"difference"`
	Logo       string `json:"-"`
}
