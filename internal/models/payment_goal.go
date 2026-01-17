package models

import (
	"golang-api/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type PaymentGoalStatus struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

const (
	PaymentGoalStatusOngoing   = 1
	PaymentGoalStatusOverdue   = 2
	PaymentGoalStatusCompleted = 3
)

func (PaymentGoalStatus) TableName() string {
	return "payment_goal_statuses"
}

type PaymentGoal struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"user_id"`
	StatusID        uint      `json:"status_id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     *string   `json:"description"`
	Amount          int64     `json:"amount"`
	TargetAmount    int64     `json:"target_amount"`
	ProgressPercent int       `json:"progress_percent"`
	StartDate       DateOnly  `json:"start_date"`
	TargetDate      DateOnly  `json:"target_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Status *PaymentGoalStatus `gorm:"foreignKey:StatusID" json:"-"`

	StatusName string        `gorm:"-" json:"status"`
	Formatted  FormattedGoal `gorm:"-" json:"formatted"`
}

type FormattedGoal struct {
	Amount       string `json:"amount"`
	TargetAmount string `json:"target_amount"`
	Progress     string `json:"progress"`
	StartDate    string `json:"start_date"`
	TargetDate   string `json:"target_date"`
}

func (PaymentGoal) TableName() string {
	return "payment_goals"
}

func (p *PaymentGoal) AfterFind(tx *gorm.DB) (err error) {
	if p.Status != nil {
		p.StatusName = p.Status.Name
	}

	p.Formatted = FormattedGoal{
		Amount:       utils.FormatRupiah(p.Amount),
		TargetAmount: utils.FormatRupiah(p.TargetAmount),
		Progress:     utils.FormatPercent(p.ProgressPercent),
		StartDate:    utils.FormatDateID(time.Time(p.StartDate), "02/01/2006"),
		TargetDate:   utils.FormatDateID(time.Time(p.TargetDate), "02/01/2006"),
	}

	return
}
