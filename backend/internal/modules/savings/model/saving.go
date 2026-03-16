package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// SavingProduct represents the different types of savings (Principal, Mandatory, Voluntary, etc.)
type SavingProduct struct {
	model.TenantModel
	Code          string  `json:"code" gorm:"not null;uniqueIndex:idx_sp_code_org"` // PRINCIPAL, MANDATORY, VOLUNTARY, etc.
	Name          string  `json:"name" gorm:"not null"`
	Description   string  `json:"description"`
	Status        string  `json:"status" gorm:"default:'active'"` // active, inactive
	IsWithdrawble bool    `json:"is_withdrawable" gorm:"default:false"`
	InterestRate  float64 `json:"interest_rate" gorm:"type:decimal(5,2);default:0"` // e.g. 5.00 for 5%
}

// SavingAccount is the actual holding per member per product.
type SavingAccount struct {
	model.TenantModel
	MemberID        uint    `json:"member_id" gorm:"not null;index"`
	SavingProductID uint    `json:"saving_product_id" gorm:"not null;index"`
	AccountNumber   string  `json:"account_number" gorm:"not null;uniqueIndex"`
	Balance         float64 `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Status          string  `json:"status" gorm:"default:'active'"` // active, closed, frozen
}

// SavingTransaction represents a deposit or withdrawal.
type SavingTransaction struct {
	model.TenantModel
	SavingAccountID uint    `json:"saving_account_id" gorm:"not null;index"`
	ReferenceNumber string  `json:"reference_number" gorm:"not null;uniqueIndex"`
	Type            string  `json:"type" gorm:"not null"` // deposit, withdrawal, interest, fee
	Amount          float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	BalanceAfter    float64 `json:"balance_after" gorm:"type:decimal(15,2);not null"`
	Description     string  `json:"description"`
	Status          string  `json:"status" gorm:"default:'completed'"` // pending, completed, failed, reversed
}
