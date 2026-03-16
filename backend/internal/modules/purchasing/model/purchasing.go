package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

// PurchaseOrder represents a purchase transaction from a supplier.
type PurchaseOrder struct {
	model.TenantModel
	SupplierID    uint                `json:"supplier_id" gorm:"not null;index"`
	PONumber      string              `json:"po_number" gorm:"not null;uniqueIndex"`
	TotalAmount   float64             `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Discount      float64             `json:"discount" gorm:"type:decimal(15,2);default:0"`
	TaxAmount     float64             `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	FinalAmount   float64             `json:"final_amount" gorm:"type:decimal(15,2);not null"`
	PaymentStatus string              `json:"payment_status" gorm:"default:'unpaid'"` // unpaid, paid, partial
	Status        string              `json:"status" gorm:"default:'pending'"`        // pending, ordered, received, cancelled
	Notes         string              `json:"notes"`
	ReceivedAt    *time.Time          `json:"received_at"`
	Items         []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
	Payments      []PurchasePayment   `json:"payments" gorm:"foreignKey:PurchaseOrderID"`
}

// PurchaseOrderItem represents an item in a Purchase Order.
type PurchaseOrderItem struct {
	model.TenantModel
	PurchaseOrderID uint    `json:"purchase_order_id" gorm:"not null;index"`
	ProductID       uint    `json:"product_id" gorm:"not null;index"`
	Quantity        int     `json:"quantity" gorm:"not null"`
	CostPrice       float64 `json:"cost_price" gorm:"type:decimal(15,2);not null"`
	Subtotal        float64 `json:"subtotal" gorm:"type:decimal(15,2);not null"`
}

// PurchasePayment tracks payments to suppliers.
type PurchasePayment struct {
	model.TenantModel
	PurchaseOrderID uint      `json:"purchase_order_id" gorm:"not null;index"`
	PaymentMethod   string    `json:"payment_method" gorm:"not null"` // cash, transfer
	Amount          float64   `json:"amount" gorm:"type:decimal(15,2);not null"`
	PaymentDate     time.Time `json:"payment_date" gorm:"not null"`
	ReferenceToken  string    `json:"reference_token"`
	Status          string    `json:"status" gorm:"default:'completed'"`
}
