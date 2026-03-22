package dto

type OrderCreateRequest struct {
	ShiftID       uint               `json:"shift_id" validate:"required"`
	PaymentMethod string             `json:"payment_method" validate:"required,oneof=cash transfer qris"`
	Notes         string             `json:"notes"`
	Items         []OrderItemRequest `json:"items" validate:"required,dive,required"`
}

type OrderItemRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Notes     string `json:"notes"`
}

type OrderResponse struct {
	ID              uint                `json:"id"`
	ReferenceNumber string              `json:"reference_number"`
	TotalAmount     float64             `json:"total_amount"`
	FinalAmount     float64             `json:"final_amount"`
	PaymentMethod   string              `json:"payment_method"`
	Status          string              `json:"status"`
	CreatedAt       string              `json:"created_at"`
	Items           []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	KDSStatus   string  `json:"kds_status"`
	Notes       string  `json:"notes"`
}

type KDSItemResponse struct {
	ID             uint   `json:"id"`
	OrderID        uint   `json:"order_id"`
	OrderReference string `json:"order_reference"`
	ProductName    string `json:"product_name"`
	Quantity       int    `json:"quantity"`
	KDSStatus      string `json:"kds_status"`
	Notes          string `json:"notes"`
	CreatedAt      string `json:"created_at"`
}

type ShiftResponse struct {
	ID           uint    `json:"id"`
	CashierID    uint    `json:"cashier_id"`
	StartTime    string  `json:"start_time"`
	Status       string  `json:"status"`
	StartBalance float64 `json:"start_balance"`
}
