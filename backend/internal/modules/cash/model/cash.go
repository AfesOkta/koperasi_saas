package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// CashRegister represents a physical or logical cash drawer/bank account.
type CashRegister struct {
	model.TenantModel
	Name        string  `json:"name" gorm:"not null"` // e.g., "Main Vault", "Teller 1", "BCA Account"
	Type        string  `json:"type" gorm:"not null"` // cash, bank, e-wallet
	Balance     float64 `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Status      string  `json:"status" gorm:"default:'active'"` // active, locked
	AccountID   *uint   `json:"account_id" gorm:"index"`        // Link to Accounting Chart of Accounts
	Description string  `json:"description"`
}

// CashTransaction represents money moving in or out of a Cash Register.
type CashTransaction struct {
	model.TenantModel
	CashRegisterID  uint    `json:"cash_register_id" gorm:"not null;index"`
	ReferenceNumber string  `json:"reference_number" gorm:"not null;uniqueIndex"`
	Type            string  `json:"type" gorm:"not null"` // in, out, transfer
	Amount          float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	BalanceAfter    float64 `json:"balance_after" gorm:"type:decimal(15,2);not null"`
	Category        string  `json:"category" gorm:"not null"` // e.g., fees, petty_cash, deposit
	Description     string  `json:"description"`
	RelatedEntity   string  `json:"related_entity"` // e.g. "savings_txn", "loan_payment"
	RelatedEntityID *uint   `json:"related_entity_id"`
}
