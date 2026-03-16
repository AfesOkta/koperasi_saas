package dto

type AccountCreateRequest struct {
	Code          string `json:"code" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Type          string `json:"type" validate:"required,oneof=asset liability equity revenue expense"`
	NormalBalance string `json:"normal_balance" validate:"required,oneof=debit credit"`
	ParentID      *uint  `json:"parent_id" validate:"omitempty"`
	Description   string `json:"description" validate:"omitempty"`
}

type AccountResponse struct {
	ID            uint   `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	NormalBalance string `json:"normal_balance"`
	ParentID      *uint  `json:"parent_id,omitempty"`
	IsActive      bool   `json:"is_active"`
	Description   string `json:"description"`
}

type JournalEntryLineRequest struct {
	AccountID   uint    `json:"account_id" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	Debit       float64 `json:"debit" validate:"min=0"`
	Credit      float64 `json:"credit" validate:"min=0"`
}

type JournalEntryCreateRequest struct {
	Date        string                    `json:"date" validate:"required,datetime=2006-01-02"`
	Description string                    `json:"description" validate:"required"`
	Lines       []JournalEntryLineRequest `json:"lines" validate:"required,min=2,dive"` // At least 2 lines for double entry
}

type JournalEntryResponse struct {
	ID              uint                       `json:"id"`
	ReferenceNumber string                     `json:"reference_number"`
	Date            string                     `json:"date"`
	Description     string                     `json:"description"`
	Status          string                     `json:"status"`
	CreatedAt       string                     `json:"created_at"`
	Lines           []JournalEntryLineResponse `json:"lines"`
}

type JournalEntryLineResponse struct {
	ID          uint    `json:"id"`
	AccountID   uint    `json:"account_id"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}
