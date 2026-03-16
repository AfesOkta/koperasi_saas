package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

type POSShift struct {
	model.TenantModel
	CashierID    uint       `json:"cashier_id" gorm:"not null;index"`
	StartTime    time.Time  `json:"start_time" gorm:"not null"`
	EndTime      *time.Time `json:"end_time"`
	StartBalance float64    `json:"start_balance" gorm:"type:decimal(15,2);not null"`
	EndBalance   float64    `json:"end_balance" gorm:"type:decimal(15,2)"` // Expected balance
	ActualCash   float64    `json:"actual_cash" gorm:"type:decimal(15,2)"` // Counted by cashier
	Difference   float64    `json:"difference" gorm:"type:decimal(15,2)"`
	Notes        string     `json:"notes" gorm:"type:text"`
	Status       string     `json:"status" gorm:"size:20;default:'open'"` // open, closed
}
