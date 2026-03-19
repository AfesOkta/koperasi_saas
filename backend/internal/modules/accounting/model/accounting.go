package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

// Account represents a Chart of Accounts (CoA) entry.
type Account struct {
	model.TenantModel
	Code          string `json:"code" gorm:"not null;uniqueIndex:idx_acc_code_org"` // e.g., 1000, 2000, 4000
	Name          string `json:"name" gorm:"not null"`
	Type          string `json:"type" gorm:"not null"`           // Asset, Liability, Equity, Revenue, Expense
	NormalBalance string `json:"normal_balance" gorm:"not null"` // debit, credit
	ParentID      *uint  `json:"parent_id" gorm:"index"`
	IsActive      bool   `json:"is_active" gorm:"default:true"`
	Description   string `json:"description"`
}

// JournalEntry represents a double-entry accounting transaction block.
type JournalEntry struct {
	model.TenantModel
	ReferenceNumber string             `json:"reference_number" gorm:"not null;uniqueIndex"`
	IdempotencyKey  string             `json:"idempotency_key" gorm:"index;uniqueIndex:idx_idempotency_org"` // For duplicate prevention
	Date            time.Time          `json:"date" gorm:"type:date;not null"`
	Description     string             `json:"description"`
	Status          string             `json:"status" gorm:"default:'posted'"` // drafted, posted, voided
	SourceModule    string             `json:"source_module"`                  // savings, loan, cash, sales, purchasing
	SourceReference string             `json:"source_reference"`               // Original transaction ID
	ReversedEntryID *uint              `json:"reversed_entry_id" gorm:"index"` // Links to reversed entry
	ReversalReason  string             `json:"reversal_reason"`                // Reason for reversal
	PostedAt        *time.Time         `json:"posted_at"`                      // When entry was posted
	PostedBy        *uint              `json:"posted_by" gorm:"index"`
	Lines           []JournalEntryLine `json:"lines" gorm:"foreignKey:JournalEntryID"`
}

// JournalEntryLine represents a single debit or credit line in a JournalEntry.
type JournalEntryLine struct {
	model.TenantModel
	JournalEntryID uint    `json:"journal_entry_id" gorm:"not null;index"`
	AccountID      uint    `json:"account_id" gorm:"not null;index"`
	AccountCode    string  `json:"account_code"` // Denormalized for audit
	Description    string  `json:"description"`
	Debit          float64 `json:"debit" gorm:"type:decimal(15,2);default:0"`
	Credit         float64 `json:"credit" gorm:"type:decimal(15,2);default:0"`
	PartnerID      *uint   `json:"partner_id" gorm:"index"` // Member, supplier, etc.
	PartnerType    string  `json:"partner_type"`            // member, supplier, employee
}
