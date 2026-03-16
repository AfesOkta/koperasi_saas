package dto

type CashRegisterCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=cash bank e-wallet"`
	AccountID   *uint  `json:"account_id" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
}

type CashRegisterResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	AccountID   *uint   `json:"account_id,omitempty"`
	Description string  `json:"description"`
}

type CashTransactionRequest struct {
	Type            string  `json:"type" validate:"required,oneof=in out transfer"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Category        string  `json:"category" validate:"required"`
	Description     string  `json:"description" validate:"omitempty"`
	RelatedEntity   string  `json:"related_entity" validate:"omitempty"`
	RelatedEntityID *uint   `json:"related_entity_id" validate:"omitempty"`
}

type CashTransactionResponse struct {
	ID              uint    `json:"id"`
	CashRegisterID  uint    `json:"cash_register_id"`
	ReferenceNumber string  `json:"reference_number"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	BalanceAfter    float64 `json:"balance_after"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
}
