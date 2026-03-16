package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// Order represents a sales transaction in the POS.
type Order struct {
	model.TenantModel
	MemberID      *uint          `json:"member_id" gorm:"index"`               // Optional, can be a non-member (guest)
	OrderID       string         `json:"order_id" gorm:"not null;uniqueIndex"` // Order Reference Number
	TotalAmount   float64        `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Discount      float64        `json:"discount" gorm:"type:decimal(15,2);default:0"`
	TaxAmount     float64        `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	FinalAmount   float64        `json:"final_amount" gorm:"type:decimal(15,2);not null"`
	PaymentStatus string         `json:"payment_status" gorm:"default:'unpaid'"` // unpaid, paid, partial
	Status        string         `json:"status" gorm:"default:'completed'"`      // completed, voided
	CashierID     uint           `json:"cashier_id" gorm:"not null;index"`       // User ID who processed the sale
	Items         []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
	Payments      []OrderPayment `json:"payments" gorm:"foreignKey:OrderID"`
}

// OrderItem represents a specific product in an Order.
type OrderItem struct {
	model.TenantModel
	OrderID     uint    `json:"order_id" gorm:"not null;index"`
	ProductID   uint    `json:"product_id" gorm:"not null;index"`
	Quantity    int     `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Subtotal    float64 `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	Discount    float64 `json:"discount" gorm:"type:decimal(15,2);default:0"`
	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(15,2);not null"`
}

// OrderPayment tracks payments made for an Order (supports split payments).
type OrderPayment struct {
	model.TenantModel
	OrderID        uint    `json:"order_id" gorm:"not null;index"`
	PaymentMethod  string  `json:"payment_method" gorm:"not null"` // cash, savings, transfer
	Amount         float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	ReferenceToken string  `json:"reference_token"` // e.g., Savings Account ID or Transfer Ref
	Status         string  `json:"status" gorm:"default:'completed'"`
}
