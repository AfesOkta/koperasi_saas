package dto

type SavingProductCreateRequest struct {
	Code          string  `json:"code" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"omitempty"`
	IsWithdrawble bool    `json:"is_withdrawable"`
	InterestRate  float64 `json:"interest_rate" validate:"omitempty,min=0,max=100"`
}

type SavingProductResponse struct {
	ID            uint    `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	IsWithdrawble bool    `json:"is_withdrawable"`
	InterestRate  float64 `json:"interest_rate"`
	CreatedAt     string  `json:"created_at"`
}

type SavingAccountResponse struct {
	ID              uint    `json:"id"`
	MemberID        uint    `json:"member_id"`
	SavingProductID uint    `json:"saving_product_id"`
	AccountNumber   string  `json:"account_number"`
	Balance         float64 `json:"balance"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

type SavingTransactionRequest struct {
	MemberID        uint    `json:"member_id" validate:"required"`
	SavingProductID uint    `json:"saving_product_id" validate:"required"`
	Type            string  `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Description     string  `json:"description" validate:"omitempty"`
}

type SavingTransactionResponse struct {
	ID              uint    `json:"id"`
	SavingAccountID uint    `json:"saving_account_id"`
	ReferenceNumber string  `json:"reference_number"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	BalanceAfter    float64 `json:"balance_after"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}
