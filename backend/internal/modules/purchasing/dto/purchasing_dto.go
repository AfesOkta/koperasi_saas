package dto

type PurchaseOrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	CostPrice float64 `json:"cost_price" validate:"required,gt=0"`
}

type PurchaseOrderCreateRequest struct {
	SupplierID uint                       `json:"supplier_id" validate:"required"`
	Discount   float64                    `json:"discount" validate:"min=0"`
	Notes      string                     `json:"notes" validate:"omitempty"`
	Items      []PurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type PurchaseOrderResponse struct {
	ID            uint                        `json:"id"`
	PONumber      string                      `json:"po_number"`
	SupplierID    uint                        `json:"supplier_id"`
	TotalAmount   float64                     `json:"total_amount"`
	FinalAmount   float64                     `json:"final_amount"`
	PaymentStatus string                      `json:"payment_status"`
	Status        string                      `json:"status"`
	CreatedAt     string                      `json:"created_at"`
	Items         []PurchaseOrderItemResponse `json:"items"`
}

type PurchaseOrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	CostPrice float64 `json:"cost_price"`
	Subtotal  float64 `json:"subtotal"`
}

type PurchasePaymentRequest struct {
	Method         string  `json:"method" validate:"required,oneof=cash transfer"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	ReferenceToken string  `json:"reference_token" validate:"omitempty"`
}
