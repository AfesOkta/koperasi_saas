package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// Category represents a product category (e.g., Electronics, Groceries)
type Category struct {
	model.TenantModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
}

// Product represents an item for sale in the cooperative
type Product struct {
	model.TenantModel
	CategoryID  uint    `json:"category_id" gorm:"index"`
	SKU         string  `json:"sku" gorm:"not null;uniqueIndex:idx_product_sku_org"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"type:decimal(15,2);not null"`
	CostPrice   float64 `json:"cost_price" gorm:"type:decimal(15,2);not null"` // What the coop paid
	Stock       int     `json:"stock" gorm:"not null;default:0"`
	MinStock    int     `json:"min_stock" gorm:"not null;default:5"` // Threshold for reorder alerts
	Unit        string  `json:"unit" gorm:"not null;default:'pcs'"`  // pcs, kg, box
	Status      string  `json:"status" gorm:"default:'active'"`      // active, inactive
}

// StockMovement tracks all changes to a product's stock.
type StockMovement struct {
	model.TenantModel
	ProductID       uint   `json:"product_id" gorm:"not null;index"`
	ReferenceNumber string `json:"reference_number" gorm:"not null;index"`
	Type            string `json:"type" gorm:"not null"` // in (purchase), out (sale), adj (adjustment)
	Quantity        int    `json:"quantity" gorm:"not null"`
	BalanceAfter    int    `json:"balance_after" gorm:"not null"`
	Notes           string `json:"notes"`
	RelatedEntity   string `json:"related_entity"` // e.g. "purchase_order", "sales_invoice"
	RelatedEntityID *uint  `json:"related_entity_id"`
}
