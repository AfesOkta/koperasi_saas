package service

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/report/dto"
	"gorm.io/gorm"
)

type ReportService interface {
	GetBalanceSheet(ctx context.Context, orgID uint) (*dto.BalanceSheetResponse, error)
	GetProfitLoss(ctx context.Context, orgID uint) (*dto.ProfitLossResponse, error)
	GetSummary(ctx context.Context, orgID uint) (*dto.MemberFinancialSummary, error)
}

type reportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportService{db: db}
}

func (s *reportService) GetBalanceSheet(ctx context.Context, orgID uint) (*dto.BalanceSheetResponse, error) {
	var results []struct {
		Code    string
		Name    string
		Type    string
		Balance float64
	}

	err := s.db.WithContext(ctx).Table("accounts").
		Select("code, name, type, balance").
		Where("organization_id = ? AND type IN ('asset', 'liability', 'equity')", orgID).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	res := &dto.BalanceSheetResponse{}
	for _, r := range results {
		item := dto.AccountBalance{AccountCode: r.Code, AccountName: r.Name, Balance: r.Balance}
		switch r.Type {
		case "asset":
			res.Assets = append(res.Assets, item)
			res.TotalAssets += r.Balance
		case "liability":
			res.Liabilities = append(res.Liabilities, item)
			res.TotalLia += r.Balance
		case "equity":
			res.Equity = append(res.Equity, item)
			res.TotalEquity += r.Balance
		}
	}
	return res, nil
}

func (s *reportService) GetProfitLoss(ctx context.Context, orgID uint) (*dto.ProfitLossResponse, error) {
	var results []struct {
		Code    string
		Name    string
		Type    string
		Balance float64
	}

	err := s.db.WithContext(ctx).Table("accounts").
		Select("code, name, type, balance").
		Where("organization_id = ? AND type IN ('income', 'expense')", orgID).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	res := &dto.ProfitLossResponse{}
	for _, r := range results {
		item := dto.AccountBalance{AccountCode: r.Code, AccountName: r.Name, Balance: r.Balance}
		if r.Type == "income" {
			res.Incomes = append(res.Incomes, item)
			res.TotalIncome += r.Balance
		} else {
			res.Expenses = append(res.Expenses, item)
			res.TotalExpenses += r.Balance
		}
	}
	res.NetProfit = res.TotalIncome - res.TotalExpenses
	return res, nil
}

func (s *reportService) GetSummary(ctx context.Context, orgID uint) (*dto.MemberFinancialSummary, error) {
	var summary dto.MemberFinancialSummary

	s.db.Table("members").Where("organization_id = ?", orgID).Count(&summary.ActiveMembers)
	s.db.Table("saving_accounts").Where("organization_id = ?", orgID).Select("SUM(balance)").Scan(&summary.TotalSavings)
	s.db.Table("loans").Where("organization_id = ? AND status = 'active'", orgID).Select("SUM(outstanding)").Scan(&summary.TotalLoanBalance)

	return &summary, nil
}
