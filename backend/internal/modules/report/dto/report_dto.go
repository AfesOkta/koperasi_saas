package dto

type AccountType string

const (
	TypeAsset     AccountType = "asset"
	TypeLiability AccountType = "liability"
	TypeEquity    AccountType = "equity"
	TypeRevenue   AccountType = "revenue"
	TypeExpense   AccountType = "expense"
)

type NormalBalance string

const (
	NormalDebit  NormalBalance = "debit"
	NormalCredit NormalBalance = "credit"
)

// AccountNode represents an account in the Chart of Accounts tree.
type AccountNode struct {
	ID             uint           `json:"id"`
	ParentID       *uint          `json:"parent_id,omitempty"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	Type           AccountType    `json:"type"`
	NormalBalance  NormalBalance  `json:"normal_balance"`
	OpeningBalance float64        `json:"opening_balance"` // Balance just before date_from
	PeriodDebit    float64        `json:"period_debit"`    // Debit movements during the period
	PeriodCredit   float64        `json:"period_credit"`   // Credit movements during the period
	PeriodMovement float64        `json:"period_movement"` // Net movement during period based on normal balance
	EndingBalance  float64        `json:"ending_balance"`  // Balance as of date_to
	Children       []*AccountNode `json:"children,omitempty"`
}

// BalanceSheetResponse is the API response for Neraca.
type BalanceSheetResponse struct {
	Assets      []*AccountNode `json:"assets"`
	Liabilities []*AccountNode `json:"liabilities"`
	Equity      []*AccountNode `json:"equity"`
	TotalAssets float64        `json:"total_assets"`
	TotalLia    float64        `json:"total_liabilities"`
	TotalEquity float64        `json:"total_equity"`
}

// ProfitLossResponse is the API response for Laba Rugi.
type ProfitLossResponse struct {
	Revenues      []*AccountNode `json:"revenues"`
	Expenses      []*AccountNode `json:"expenses"`
	TotalRevenue  float64        `json:"total_revenue"`
	TotalExpenses float64        `json:"total_expenses"`
	NetProfit     float64        `json:"net_profit"`
}

type TrialBalanceResponse struct {
	Accounts []*AccountNode `json:"accounts"`
}

type DashboardKPIs struct {
	TotalSavings     float64 `json:"total_savings"`
	SavingsGrowthMoM float64 `json:"savings_growth_mom"` // Not MVP, set to 0
	
	TotalLoanBalance float64 `json:"total_loan_balance"`
	LoanDisbursedMTD float64 `json:"loan_disbursed_mtd"`
	PAR              float64 `json:"par_percentage"` // Portfolio at Risk
	NPL              float64 `json:"npl_percentage"` // Non-Performing Loans
	
	ActiveMembers    int64   `json:"active_members"`
	NewMembersMTD    int64   `json:"new_members_mtd"`
	
	NetProfitMTD     float64 `json:"net_profit_mtd"`
}

type MemberFinancialSummary struct {
	TotalSavings     float64 `json:"total_savings"`
	TotalLoanBalance float64 `json:"total_loan_balance"`
	ActiveMembers    int64   `json:"active_members"`
}
