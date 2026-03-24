package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// LoanProduct defines the type of loan offered by the cooperative.
type LoanProduct struct {
	model.TenantModel
	Code         string  `json:"code" gorm:"not null;uniqueIndex:idx_lp_code_org"`
	Name         string  `json:"name" gorm:"not null"`
	Description  string  `json:"description"`
	InterestRate float64 `json:"interest_rate" gorm:"type:decimal(5,2);not null"` // e.g. 1.5% per month
	InterestType string  `json:"interest_type" gorm:"not null"`                   // flat, declining
	MaxAmount    float64 `json:"max_amount" gorm:"type:decimal(15,2);not null"`
	MaxTerm      int     `json:"max_term" gorm:"not null"`       // Maximum period in months
	Status       string  `json:"status" gorm:"default:'active'"` // active, inactive

	// Penalty Configuration (Denda)
	PenaltyType      string  `json:"penalty_type" gorm:"size:20;default:'percentage'"` // flat, percentage
	PenaltyRate      float64 `json:"penalty_rate" gorm:"type:decimal(15,4);default:0"` // 0.001 for 0.1% or 1000 for flat
	PenaltyGraceDays int     `json:"penalty_grace_days" gorm:"default:0"`
	PenaltyCap       float64 `json:"penalty_cap" gorm:"type:decimal(15,2);default:0"`
}

// Loan represents a specific loan application and its current state.
type Loan struct {
	model.TenantModel
	MemberID        uint    `json:"member_id" gorm:"not null;index"`
	LoanProductID   uint    `json:"loan_product_id" gorm:"not null;index"`
	LoanNumber      string  `json:"loan_number" gorm:"not null;uniqueIndex"`
	PrincipalAmount float64 `json:"principal_amount" gorm:"type:decimal(15,2);not null"` // Amount requested/approved
	InterestRate    float64 `json:"interest_rate" gorm:"type:decimal(5,2);not null"`     // Locked-in rate at approval
	TermMonths      int     `json:"term_months" gorm:"not null"`                         // Number of months

	TotalInterest float64 `json:"total_interest" gorm:"type:decimal(15,2);default:0"` // Calculated total interest
	ExpectedTotal float64 `json:"expected_total" gorm:"type:decimal(15,2);default:0"` // Principal + Total Interest
	Outstanding   float64 `json:"outstanding" gorm:"type:decimal(15,2);default:0"`    // Remaining balance

	Purpose            string `json:"purpose" gorm:"type:text"`                               // Tujuan pinjaman
	DisbursementMethod string `json:"disbursement_method" gorm:"size:20;default:'transfer'"` // transfer, cash

	Status      string  `json:"status" gorm:"default:'pending'"` // pending, staff_approved, supervisor_approved, active, rejected, paid, defaulted
	ApprovedAt  *string `json:"approved_at" gorm:"type:timestamp"`
	ApprovedBy  *uint   `json:"approved_by" gorm:"index"`
	RejectedAt  *string `json:"rejected_at" gorm:"type:timestamp"`
	RejectedBy  *uint   `json:"rejected_by" gorm:"index"`
	DisbursedAt *string `json:"disbursed_at" gorm:"type:timestamp"`

	Schedules    []LoanSchedule    `json:"schedules,omitempty" gorm:"foreignKey:LoanID"`
	Payments     []LoanPayment     `json:"payments,omitempty" gorm:"foreignKey:LoanID"`
	Collaterals  []LoanCollateral  `json:"collaterals,omitempty" gorm:"foreignKey:LoanID"`
	ApprovalLogs []ApprovalLog     `json:"approval_logs,omitempty" gorm:"foreignKey:LoanID"`
}

// LoanSchedule defines the expected monthly repayments.
type LoanSchedule struct {
	model.TenantModel
	LoanID          uint    `json:"loan_id" gorm:"not null;index"`
	Period          int     `json:"period" gorm:"not null"` // Month 1, 2, 3...
	DueDate         string  `json:"due_date" gorm:"type:date;not null"`
	PrincipalAmount float64 `json:"principal_amount" gorm:"type:decimal(15,2);not null"`
	InterestAmount  float64 `json:"interest_amount" gorm:"type:decimal(15,2);not null"`
	TotalAmount     float64 `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	PaidAmount      float64 `json:"paid_amount" gorm:"type:decimal(15,2);default:0"`
	Status          string  `json:"status" gorm:"default:'unpaid'"` // unpaid, partial, paid, overdue
}

// LoanPayment tracks the actual payments received.
type LoanPayment struct {
	model.TenantModel
	LoanID          uint    `json:"loan_id" gorm:"not null;index"`
	ReferenceNumber string  `json:"reference_number" gorm:"not null;uniqueIndex"`
	Amount          float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	PaymentDate     string  `json:"payment_date" gorm:"type:timestamp;not null"`
	Description     string  `json:"description"`
}
