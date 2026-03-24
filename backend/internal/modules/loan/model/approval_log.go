package model

import iammodel "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// ApprovalLog tracks each step of the multi-level loan approval process.
// Basic plan: Staff -> Supervisor (2 levels)
// Professional/Enterprise: Staff -> Supervisor -> Manager (3 levels)
type ApprovalLog struct {
	iammodel.TenantModel
	LoanID     uint   `json:"loan_id" gorm:"not null;index"`
	ApproverID uint   `json:"approver_id" gorm:"not null;index"`
	Role       string `json:"role" gorm:"size:20;not null"`   // staff, supervisor, manager
	Action     string `json:"action" gorm:"size:20;not null"` // approve, reject
	Notes      string `json:"notes" gorm:"type:text"`
}
