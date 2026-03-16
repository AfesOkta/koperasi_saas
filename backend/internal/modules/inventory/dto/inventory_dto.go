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
	Type            string `json:"type" validate:"required,oneof=in out adj"`
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
