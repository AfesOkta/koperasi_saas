package dto

type OrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Discount  float64 `json:"discount" validate:"min=0"`
}

type OrderPaymentRequest struct {
	Method         string  `json:"method" validate:"required,oneof=cash savings transfer"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	ReferenceToken string  `json:"reference_token" validate:"omitempty"`
}

type OrderCreateRequest struct {
	MemberID *uint                 `json:"member_id" validate:"omitempty"`
	Discount float64               `json:"discount" validate:"min=0"`
	Items    []OrderItemRequest    `json:"items" validate:"required,min=1,dive"`
	Payments []OrderPaymentRequest `json:"payments" validate:"required,min=1,dive"`
}

type OrderResponse struct {
	ID            uint                   `json:"id"`
	OrderID       string                 `json:"order_id"`
	MemberID      *uint                  `json:"member_id,omitempty"`
	TotalAmount   float64                `json:"total_amount"`
	Discount      float64                `json:"discount"`
	TaxAmount     float64                `json:"tax_amount"`
	FinalAmount   float64                `json:"final_amount"`
	PaymentStatus string                 `json:"payment_status"`
	Status        string                 `json:"status"`
	CashierID     uint                   `json:"cashier_id"`
	CreatedAt     string                 `json:"created_at"`
	Items         []OrderItemResponse    `json:"items"`
	Payments      []OrderPaymentResponse `json:"payments"`
}

type OrderItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	Discount    float64 `json:"discount"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderPaymentResponse struct {
	ID            uint    `json:"id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}
