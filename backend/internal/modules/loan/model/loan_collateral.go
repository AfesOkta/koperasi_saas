package model

import iammodel "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// LoanCollateral stores collateral information for loans.
// Only available for Professional/Enterprise plans.
type LoanCollateral struct {
	iammodel.TenantModel
	LoanID      uint   `json:"loan_id" gorm:"not null;index"`
	Type        string `json:"type" gorm:"size:50;not null"`  // e.g., BPKB, Sertifikat Tanah, Sertifikat Deposito
	Description string `json:"description" gorm:"type:text"`
	DocumentURL string `json:"document_url" gorm:"size:500"`
}
