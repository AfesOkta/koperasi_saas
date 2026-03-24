package dto

type LoanProductCreateRequest struct {
	Code         string  `json:"code" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"omitempty"`
	InterestRate float64 `json:"interest_rate" validate:"required,min=0,max=100"`
	InterestType string  `json:"interest_type" validate:"required,oneof=flat declining"`
	MaxAmount    float64 `json:"max_amount" validate:"required,gt=0"`
	MaxTerm      int     `json:"max_term" validate:"required,gt=0"`
}

type LoanProductResponse struct {
	ID           uint    `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	InterestRate float64 `json:"interest_rate"`
	InterestType string  `json:"interest_type"`
	MaxAmount    float64 `json:"max_amount"`
	MaxTerm      int     `json:"max_term"`
	Status       string  `json:"status"`
}

// CollateralRequest is used for submitting collateral with a loan application.
type CollateralRequest struct {
	Type        string `json:"type" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	DocumentURL string `json:"document_url" validate:"omitempty"`
}

type LoanApplicationRequest struct {
	MemberID           uint               `json:"member_id" validate:"required"`
	LoanProductID      uint               `json:"loan_product_id" validate:"required"`
	PrincipalAmount    float64            `json:"principal_amount" validate:"required,gt=0"`
	TermMonths         int                `json:"term_months" validate:"required,gt=0"`
	Purpose            string             `json:"purpose" validate:"omitempty"`
	DisbursementMethod string             `json:"disbursement_method" validate:"omitempty,oneof=transfer cash"`
	Collateral         *CollateralRequest `json:"collateral,omitempty"`
}

// ApprovalRequest is used by staff/supervisor/manager to approve or reject a loan.
type ApprovalRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Notes  string `json:"notes" validate:"omitempty"`
}

type LoanResponse struct {
	ID                 uint                   `json:"id"`
	MemberID           uint                   `json:"member_id"`
	LoanProductID      uint                   `json:"loan_product_id"`
	LoanNumber         string                 `json:"loan_number"`
	PrincipalAmount    float64                `json:"principal_amount"`
	InterestRate       float64                `json:"interest_rate"`
	TermMonths         int                    `json:"term_months"`
	TotalInterest      float64                `json:"total_interest"`
	ExpectedTotal      float64                `json:"expected_total"`
	Outstanding        float64                `json:"outstanding"`
	Purpose            string                 `json:"purpose"`
	DisbursementMethod string                 `json:"disbursement_method"`
	Status             string                 `json:"status"`
	CreatedAt          string                 `json:"created_at"`
	Schedules          []LoanScheduleResponse `json:"schedules,omitempty"`
	Payments           []LoanPaymentResponse  `json:"payments,omitempty"`
	Collaterals        []CollateralResponse   `json:"collaterals,omitempty"`
	ApprovalLogs       []ApprovalLogResponse  `json:"approval_logs,omitempty"`
}

type LoanScheduleResponse struct {
	ID              uint    `json:"id"`
	Period          int     `json:"period"`
	DueDate         string  `json:"due_date"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestAmount  float64 `json:"interest_amount"`
	TotalAmount     float64 `json:"total_amount"`
	Status          string  `json:"status"`
}

type LoanPaymentRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description" validate:"omitempty"`
}

type LoanPaymentResponse struct {
	ID              uint    `json:"id"`
	LoanID          uint    `json:"loan_id"`
	ReferenceNumber string  `json:"reference_number"`
	Amount          float64 `json:"amount"`
	PaymentDate     string  `json:"payment_date"`
	Description     string  `json:"description"`
}

type CollateralResponse struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	DocumentURL string `json:"document_url"`
}

type ApprovalLogResponse struct {
	ID         uint   `json:"id"`
	ApproverID uint   `json:"approver_id"`
	Role       string `json:"role"`
	Action     string `json:"action"`
	Notes      string `json:"notes"`
	CreatedAt  string `json:"created_at"`
}

