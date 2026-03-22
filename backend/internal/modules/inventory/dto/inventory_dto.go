package dto

type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WarehouseResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	IsActive    bool   `json:"is_active"`
}

type WarehouseCreateRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	Address     string `json:"address" validate:"omitempty"`
}

type ProductCreateRequest struct {
	CategoryID  uint    `json:"category_id" validate:"omitempty"`
	SKU         string  `json:"sku" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,min=0"`
	CostPrice   float64 `json:"cost_price" validate:"required,min=0"`
	MinStock    int     `json:"min_stock" validate:"min=0"`
	Unit        string  `json:"unit" validate:"required"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	CategoryID  uint    `json:"category_id,omitempty"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CostPrice   float64 `json:"cost_price"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
	Unit        string  `json:"unit"`
	Status      string  `json:"status"`
}

type StockMovementRequest struct {
	WarehouseID     uint   `json:"warehouse_id"` // Optional, defaults to main
	Type            string `json:"type" validate:"required,oneof=in out adj transfer_in transfer_out transit"`
	Quantity        int    `json:"quantity" validate:"required"`
	Notes           string `json:"notes" validate:"omitempty"`
	RelatedEntity   string `json:"related_entity" validate:"omitempty"`
	RelatedEntityID *uint  `json:"related_entity_id" validate:"omitempty"`
}

type StockMovementResponse struct {
	ID              uint   `json:"id"`
	ProductID       uint   `json:"product_id"`
	ReferenceNumber string `json:"reference_number"`
	Type            string `json:"type"`
	Quantity        int    `json:"quantity"`
	BalanceAfter    int    `json:"balance_after"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"created_at"`
}

type TransferCreateRequest struct {
	ProductID       uint   `json:"product_id" validate:"required"`
	FromWarehouseID uint   `json:"from_warehouse_id" validate:"required"`
	ToWarehouseID   uint   `json:"to_warehouse_id" validate:"required"`
	Quantity        int    `json:"quantity" validate:"required,gt=0"`
	Notes           string `json:"notes" validate:"omitempty"`
}

type TransferResponse struct {
	ID              uint              `json:"id"`
	ReferenceNumber string            `json:"reference_number"`
	FromWarehouse   WarehouseResponse `json:"from_warehouse"`
	ToWarehouse     WarehouseResponse `json:"to_warehouse"`
	Status          string            `json:"status"`
	Notes           string            `json:"notes"`
	CreatedAt       string            `json:"created_at"`
}
