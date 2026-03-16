package dto

type BalanceSheetResponse struct {
	Assets      []AccountBalance `json:"assets"`
	Liabilities []AccountBalance `json:"liabilities"`
	Equity      []AccountBalance `json:"equity"`
	TotalAssets float64          `json:"total_assets"`
	TotalLia    float64          `json:"total_liabilities"`
	TotalEquity float64          `json:"total_equity"`
}

type ProfitLossResponse struct {
	Incomes       []AccountBalance `json:"incomes"`
	Expenses      []AccountBalance `json:"expenses"`
	TotalIncome   float64          `json:"total_income"`
	TotalExpenses float64          `json:"total_expenses"`
	NetProfit     float64          `json:"net_profit"`
}

type AccountBalance struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
}

type MemberFinancialSummary struct {
	TotalSavings     float64 `json:"total_savings"`
	TotalLoanBalance float64 `json:"total_loan_balance"`
	ActiveMembers    int64   `json:"active_members"`
}

type InventorySummary struct {
	TotalValue       float64 `json:"total_value"`
	LowStockProducts int64   `json:"low_stock_products"`
}
