package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

// Warehouse represents a storage location
type Warehouse struct {
	model.TenantModel
	Code        string `json:"code" gorm:"not null;uniqueIndex:idx_warehouse_code_org"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Address     string `json:"address"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}

// WarehouseItem tracks stock of a product in a specific warehouse
type WarehouseItem struct {
	model.TenantModel
	WarehouseID  uint    `json:"warehouse_id" gorm:"index;not null"`
	ProductID    uint    `json:"product_id" gorm:"index;not null"`
	Quantity     int     `json:"quantity" gorm:"not null;default:0"`
	MinStock     int     `json:"min_stock" gorm:"not null;default:5"`
	ReorderPoint int     `json:"reorder_point" gorm:"not null;default:10"`
	Warehouse    Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Product      Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

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
	WarehouseID     uint   `json:"warehouse_id" gorm:"not null;index"`
	ReferenceNumber string `json:"reference_number" gorm:"not null;index"`
	Type            string `json:"type" gorm:"not null"` // in (purchase), out (sale), adj (adjustment), transfer_in, transfer_out, transit
	Quantity        int    `json:"quantity" gorm:"not null"`
	BalanceAfter    int    `json:"balance_after" gorm:"not null"`
	Notes           string `json:"notes"`
	RelatedEntity   string `json:"related_entity"` // e.g. "purchase_order", "sales_invoice", "stock_transfer"
	RelatedEntityID *uint  `json:"related_entity_id"`
}

// StockTransfer represents stock moving between warehouses
type StockTransfer struct {
	model.TenantModel
	ProductID            uint   `json:"product_id" gorm:"not null;index"`
	Quantity             int    `json:"quantity" gorm:"not null"`
	ReferenceNumber      string `json:"reference_number" gorm:"not null;uniqueIndex:idx_transfer_ref_org"`
	FromWarehouseID      uint   `json:"from_warehouse_id" gorm:"not null;index"`
	ToWarehouseID        uint   `json:"to_warehouse_id" gorm:"not null;index"`
	Status               string `json:"status" gorm:"default:'pending'"` // pending, shipped, received, cancelled
	Notes                string `json:"notes"`
	ShippedAt            *time.Time `json:"shipped_at"`
	ReceivedAt           *time.Time `json:"received_at"`
	Product              Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	FromWarehouse        Warehouse `json:"from_warehouse,omitempty" gorm:"foreignKey:FromWarehouseID"`
	ToWarehouse          Warehouse `json:"to_warehouse,omitempty" gorm:"foreignKey:ToWarehouseID"`
}
