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

type POSOrder struct {
	model.TenantModel
	ShiftID         uint           `json:"shift_id" gorm:"not null;index"`
	ReferenceNumber string         `json:"reference_number" gorm:"not null;index"`
	TotalAmount     float64        `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	TaxAmount       float64        `json:"tax_amount" gorm:"type:decimal(15,2);not null"`
	DiscountAmount  float64        `json:"discount_amount" gorm:"type:decimal(15,2);not null"`
	FinalAmount     float64        `json:"final_amount" gorm:"type:decimal(15,2);not null"`
	PaymentMethod   string         `json:"payment_method" gorm:"not null"` // cash, transfer, qris
	Status          string         `json:"status" gorm:"default:'pending'"` // pending, completed, cancelled
	Notes           string         `json:"notes"`
	Items           []POSOrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Shift           POSShift       `json:"shift,omitempty" gorm:"foreignKey:ShiftID"`
}

type POSOrderItem struct {
	model.TenantModel
	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Subtotal  float64 `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	KDSStatus string  `json:"kds_status" gorm:"default:'pending'"` // pending, preparing, ready, served
	Notes     string  `json:"notes"`
	Order     *POSOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
