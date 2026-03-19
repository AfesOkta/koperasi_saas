package dto

import "time"

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

// JournalEntryLineRequest - Enhanced with account code support
type JournalEntryLineRequest struct {
	AccountCode string  `json:"account_code"` // Primary: use code to lookup account
	AccountID   *uint   `json:"account_id"`   // Alternative: direct ID
	Description string  `json:"description"`
	Debit       float64 `json:"debit" validate:"min=0"`
	Credit      float64 `json:"credit" validate:"min=0"`
	PartnerID   *uint   `json:"partner_id"`
	PartnerType string  `json:"partner_type" validate:"omitempty,oneof=member supplier employee"`
}

// JournalEntryCreateRequest - Enhanced with idempotency and source tracking
type JournalEntryCreateRequest struct {
	Date            string                    `json:"date" validate:"required,datetime=2006-01-02"`
	Description     string                    `json:"description" validate:"required"`
	SourceModule    string                    `json:"source_module" validate:"required"`
	SourceReference string                    `json:"source_reference"`
	IdempotencyKey  string                    `json:"idempotency_key" validate:"max=255"`
	Lines           []JournalEntryLineRequest `json:"lines" validate:"required,min=2,dive"`
}

// JournalEntryResponse - Enhanced with additional fields
type JournalEntryResponse struct {
	ID              uint                       `json:"id"`
	ReferenceNumber string                     `json:"reference_number"`
	IdempotencyKey  string                     `json:"idempotency_key,omitempty"`
	Date            string                     `json:"date"`
	Description     string                     `json:"description"`
	Status          string                     `json:"status"`
	SourceModule    string                     `json:"source_module,omitempty"`
	SourceReference string                     `json:"source_reference,omitempty"`
	ReversedEntryID *uint                      `json:"reversed_entry_id,omitempty"`
	ReversalReason  string                     `json:"reversal_reason,omitempty"`
	TotalDebit      float64                    `json:"total_debit"`
	TotalCredit     float64                    `json:"total_credit"`
	CreatedAt       string                     `json:"created_at"`
	Lines           []JournalEntryLineResponse `json:"lines"`
}

// JournalEntryLineResponse - Enhanced with account code and partner info
type JournalEntryLineResponse struct {
	ID          uint    `json:"id"`
	AccountID   uint    `json:"account_id"`
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name,omitempty"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	PartnerID   *uint   `json:"partner_id,omitempty"`
	PartnerType string  `json:"partner_type,omitempty"`
}

// JournalEntryFilter - For filtering journal entries
type JournalEntryFilter struct {
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	SourceModule string     `json:"source_module"`
	Status       string     `json:"status"`
	AccountID    uint       `json:"account_id"`
	ReferenceNum string     `json:"reference_number"`
}

// ReverseJournalEntryRequest - For reversing a journal entry
type ReverseJournalEntryRequest struct {
	Reason string `json:"reason" validate:"required"`
}
