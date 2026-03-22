package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/report/dto"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReportService interface {
	GetBalanceSheet(ctx context.Context, orgID uint, dateFrom, dateTo time.Time) (*dto.BalanceSheetResponse, error)
	GetProfitLoss(ctx context.Context, orgID uint, dateFrom, dateTo time.Time) (*dto.ProfitLossResponse, error)
	GetSummary(ctx context.Context, orgID uint) (*dto.MemberFinancialSummary, error)
	GetDashboardKPIs(ctx context.Context, orgID uint) (*dto.DashboardKPIs, error)
	
	InvalidateCache(ctx context.Context, orgID uint) error
}

type reportService struct {
	db         *gorm.DB
	rdb        *redis.Client
	aggregator *Aggregator
}

func NewReportService(db *gorm.DB, rdb *redis.Client) ReportService {
	return &reportService{
		db:         db,
		rdb:        rdb,
		aggregator: NewAggregator(db),
	}
}

func (s *reportService) GetBalanceSheet(ctx context.Context, orgID uint, dateFrom, dateTo time.Time) (*dto.BalanceSheetResponse, error) {
	// 1. Get Cache Version
	version, _ := s.rdb.Get(ctx, fmt.Sprintf("report:v:%d", orgID)).Result()
	if version == "" {
		version = "1"
	}

	// 2. Try Redis cache (1 Hour TTL)
	cacheKey := fmt.Sprintf("report:bs:%d:%s:%s:%s", orgID, version, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))
	if val, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		var cachedRes dto.BalanceSheetResponse
		if err := json.Unmarshal([]byte(val), &cachedRes); err == nil {
			return &cachedRes, nil
		}
	}

	// 3. Compute
	nodeMap, err := s.aggregator.GetAccountBalances(ctx, orgID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	res := &dto.BalanceSheetResponse{
		Assets:      make([]*dto.AccountNode, 0),
		Liabilities: make([]*dto.AccountNode, 0),
		Equity:      make([]*dto.AccountNode, 0),
	}

	var totalAssets, totalLia, totalEquity, retainedEarnings float64

	// Gather top-level accounts and compute retained earnings sweeps
	for _, node := range nodeMap {
		// Only collect top-level nodes for the hierarchy
		if node.ParentID == nil {
			switch node.Type {
			case dto.TypeAsset:
				res.Assets = append(res.Assets, node)
				totalAssets += node.EndingBalance
			case dto.TypeLiability:
				res.Liabilities = append(res.Liabilities, node)
				totalLia += node.EndingBalance
			case dto.TypeEquity:
				res.Equity = append(res.Equity, node)
				totalEquity += node.EndingBalance
			}
		}

		// Calculate Retained Earnings: accumulate P&L balances regardless of hierarchy
		// Unclosed prior periods + YTD profit
		if node.ParentID == nil && (node.Type == dto.TypeRevenue || node.Type == dto.TypeExpense) {
			if node.NormalBalance == dto.NormalCredit {
				retainedEarnings += node.EndingBalance
			} else {
				retainedEarnings -= node.EndingBalance
			}
		}
	}

	// Append Retained Earnings dynamically to Equity to balance the Neraca
	if retainedEarnings != 0 {
		retainedEquityNode := &dto.AccountNode{
			Code:          "auto-RE",
			Name:          "Laba Ditahan / Tahun Berjalan (Auto)",
			Type:          dto.TypeEquity,
			NormalBalance: dto.NormalCredit,
			EndingBalance: retainedEarnings,
		}
		res.Equity = append(res.Equity, retainedEquityNode)
		totalEquity += retainedEarnings
	}

	res.TotalAssets = totalAssets
	res.TotalLia = totalLia
	res.TotalEquity = totalEquity
    
    // 3. Cache the result
    if jsonStr, err := json.Marshal(res); err == nil {
        s.rdb.Set(ctx, cacheKey, jsonStr, 1*time.Hour)
    }

	return res, nil
}

func (s *reportService) GetProfitLoss(ctx context.Context, orgID uint, dateFrom, dateTo time.Time) (*dto.ProfitLossResponse, error) {
	// 1. Get Cache Version
	version, _ := s.rdb.Get(ctx, fmt.Sprintf("report:v:%d", orgID)).Result()
	if version == "" {
		version = "1"
	}

	// 2. Try Redis cache (1 Hour TTL)
	cacheKey := fmt.Sprintf("report:pl:%d:%s:%s:%s", orgID, version, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))
	if val, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		var cachedRes dto.ProfitLossResponse
		if err := json.Unmarshal([]byte(val), &cachedRes); err == nil {
			return &cachedRes, nil
		}
	}

	// 3. Compute
	nodeMap, err := s.aggregator.GetAccountBalances(ctx, orgID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	res := &dto.ProfitLossResponse{
		Revenues: make([]*dto.AccountNode, 0),
		Expenses: make([]*dto.AccountNode, 0),
	}

	var totalRevenue, totalExpenses float64

	for _, node := range nodeMap {
		if node.ParentID == nil {
			if node.Type == dto.TypeRevenue {
				res.Revenues = append(res.Revenues, node)
				// P&L only cares about PeriodMovement (the activity within dateFrom -> dateTo)
				totalRevenue += node.PeriodMovement
			} else if node.Type == dto.TypeExpense {
				res.Expenses = append(res.Expenses, node)
				totalExpenses += node.PeriodMovement
			}
		}
	}

	res.TotalRevenue = totalRevenue
	res.TotalExpenses = totalExpenses
	res.NetProfit = totalRevenue - totalExpenses

    // 3. Cache the result
    if jsonStr, err := json.Marshal(res); err == nil {
        s.rdb.Set(ctx, cacheKey, jsonStr, 1*time.Hour)
    }

	return res, nil
}

func (s *reportService) GetSummary(ctx context.Context, orgID uint) (*dto.MemberFinancialSummary, error) {
	// 1. Try Redis cache (30 Min TTL)
	cacheKey := fmt.Sprintf("report:summary:%d", orgID)
	if val, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		var cachedRes dto.MemberFinancialSummary
		if err := json.Unmarshal([]byte(val), &cachedRes); err == nil {
			return &cachedRes, nil
		}
	}

	var summary dto.MemberFinancialSummary

	s.db.Table("members").Where("organization_id = ? AND deleted_at IS NULL AND status = 'active'", orgID).Count(&summary.ActiveMembers)
	s.db.Table("saving_accounts").Where("organization_id = ? AND deleted_at IS NULL AND status = 'active'", orgID).Select("COALESCE(SUM(balance), 0)").Scan(&summary.TotalSavings)
	s.db.Table("loans").Where("organization_id = ? AND status = 'active' AND deleted_at IS NULL", orgID).Select("COALESCE(SUM(outstanding), 0)").Scan(&summary.TotalLoanBalance)

    // 2. Cache the result
    if jsonStr, err := json.Marshal(summary); err == nil {
        s.rdb.Set(ctx, cacheKey, jsonStr, 30*time.Minute)
    }

	return &summary, nil
}

func (s *reportService) InvalidateCache(ctx context.Context, orgID uint) error {
	return s.rdb.Incr(ctx, fmt.Sprintf("report:v:%d", orgID)).Err()
}

func (s *reportService) GetDashboardKPIs(ctx context.Context, orgID uint) (*dto.DashboardKPIs, error) {
	// 1. Try Redis cache (5 min TTL)
	cacheKey := fmt.Sprintf("report:dashboard:org:%d", orgID)
	if val, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		var kpis dto.DashboardKPIs
		if err := json.Unmarshal([]byte(val), &kpis); err == nil {
			return &kpis, nil
		}
	}

	// 2. Compute from DB
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nplDate := now.AddDate(0, 0, -90).Format("2006-01-02")
	todayStr := now.Format("2006-01-02")

	kpis := &dto.DashboardKPIs{}

	// Active Members
	s.db.WithContext(ctx).Table("members").Where("organization_id = ? AND deleted_at IS NULL AND status = 'active'", orgID).Count(&kpis.ActiveMembers)
	s.db.WithContext(ctx).Table("members").Where("organization_id = ? AND deleted_at IS NULL AND created_at >= ?", orgID, firstOfMonth).Count(&kpis.NewMembersMTD)

	// Savings
	s.db.WithContext(ctx).Table("saving_accounts").Where("organization_id = ? AND deleted_at IS NULL AND status = 'active'", orgID).Select("COALESCE(SUM(balance), 0)").Scan(&kpis.TotalSavings)

	// Loans
	s.db.WithContext(ctx).Table("loans").Where("organization_id = ? AND status = 'active' AND deleted_at IS NULL", orgID).Select("COALESCE(SUM(outstanding), 0)").Scan(&kpis.TotalLoanBalance)

	// Disbursed MTD
	s.db.WithContext(ctx).Table("loans").
		Where("organization_id = ? AND status IN ('active', 'paid') AND deleted_at IS NULL", orgID).
		Where("disbursed_at >= ? AND disbursed_at <= ?", firstOfMonth.Format("2006-01-02")+" 00:00:00", now.Format("2006-01-02")+" 23:59:59").
		Select("COALESCE(SUM(principal_amount), 0)").Scan(&kpis.LoanDisbursedMTD)

	if kpis.TotalLoanBalance > 0 {
		var parTotal, nplTotal float64

		// PAR: Outstanding of active loans with ANY schedule unpaid/partial AND due_date < today
		s.db.WithContext(ctx).Raw(`
			SELECT COALESCE(SUM(outstanding), 0) FROM loans 
			WHERE organization_id = ? AND status = 'active' AND deleted_at IS NULL
			AND id IN (
				SELECT loan_id FROM loan_schedules 
				WHERE status IN ('unpaid', 'partial', 'overdue') AND due_date < ? AND deleted_at IS NULL
			)`, orgID, todayStr).Scan(&parTotal)

		// NPL: Outstanding of active loans with ANY schedule unpaid/partial AND due_date < 90 days ago
		s.db.WithContext(ctx).Raw(`
			SELECT COALESCE(SUM(outstanding), 0) FROM loans 
			WHERE organization_id = ? AND status = 'active' AND deleted_at IS NULL
			AND id IN (
				SELECT loan_id FROM loan_schedules 
				WHERE status IN ('unpaid', 'partial', 'overdue') AND due_date < ? AND deleted_at IS NULL
			)`, orgID, nplDate).Scan(&nplTotal)

		kpis.PAR = (parTotal / kpis.TotalLoanBalance) * 100
		kpis.NPL = (nplTotal / kpis.TotalLoanBalance) * 100
	}

	// Net Profit MTD
	plResponse, _ := s.GetProfitLoss(ctx, orgID, firstOfMonth, now)
	if plResponse != nil {
		kpis.NetProfitMTD = plResponse.NetProfit
	}

	// 3. Cache the result
	if jsonStr, err := json.Marshal(kpis); err == nil {
		s.rdb.Set(ctx, cacheKey, jsonStr, 5*time.Minute)
	}

	return kpis, nil
}
