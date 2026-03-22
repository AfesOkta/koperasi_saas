package service

import (
	"context"
	"strings"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/report/dto"
	"gorm.io/gorm"
)

type Aggregator struct {
	db *gorm.DB
}

func NewAggregator(db *gorm.DB) *Aggregator {
	return &Aggregator{db: db}
}

// GetAccountBalances calculates opening, period, and ending balances for all accounts.
// Generates a map of AccountNodes keyed by Account ID.
func (a *Aggregator) GetAccountBalances(ctx context.Context, orgID uint, dateFrom, dateTo time.Time) (map[uint]*dto.AccountNode, error) {
	// 1. Fetch all accounts
	var accounts []struct {
		ID            uint
		ParentID      *uint
		Code          string
		Name          string
		Type          string
		NormalBalance string
	}
	if err := a.db.WithContext(ctx).Table("accounts").
		Select("id, parent_id, code, name, type, normal_balance").
		Where("organization_id = ? AND is_active = ?", orgID, true).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// Initialize the map
	nodeMap := make(map[uint]*dto.AccountNode)
	for _, acc := range accounts {
		nodeMap[acc.ID] = &dto.AccountNode{
			ID:            acc.ID,
			ParentID:      acc.ParentID,
			Code:          acc.Code,
			Name:          acc.Name,
			Type:          dto.AccountType(strings.ToLower(acc.Type)),
			NormalBalance: dto.NormalBalance(strings.ToLower(acc.NormalBalance)),
			Children:      make([]*dto.AccountNode, 0),
		}
	}

	// 2. Fetch Aggregation
	startOfYear := time.Date(dateFrom.Year(), 1, 1, 0, 0, 0, 0, dateFrom.Location())

	type row struct {
		AccountID      uint
		AllPriorDebit  float64
		AllPriorCredit float64
		YtdPriorDebit  float64
		YtdPriorCredit float64
		PeriodDebit    float64
		PeriodCredit   float64
	}
	var rows []row

	query := `
		SELECT 
			l.account_id,
			SUM(CASE WHEN je.date < ? THEN l.debit ELSE 0 END) as all_prior_debit,
			SUM(CASE WHEN je.date < ? THEN l.credit ELSE 0 END) as all_prior_credit,
			SUM(CASE WHEN je.date >= ? AND je.date < ? THEN l.debit ELSE 0 END) as ytd_prior_debit,
			SUM(CASE WHEN je.date >= ? AND je.date < ? THEN l.credit ELSE 0 END) as ytd_prior_credit,
			SUM(CASE WHEN je.date >= ? AND je.date <= ? THEN l.debit ELSE 0 END) as period_debit,
			SUM(CASE WHEN je.date >= ? AND je.date <= ? THEN l.credit ELSE 0 END) as period_credit
		FROM journal_entry_lines l
		JOIN journal_entries je ON l.journal_entry_id = je.id
		WHERE je.organization_id = ? AND je.status = 'posted' AND je.deleted_at IS NULL
		GROUP BY l.account_id
	`
	if err := a.db.WithContext(ctx).Raw(query,
		dateFrom, dateFrom,
		startOfYear, dateFrom, startOfYear, dateFrom,
		dateFrom, dateTo,
		dateFrom, dateTo,
		orgID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// 3. Assign and Calculate Balances
	for _, r := range rows {
		node, ok := nodeMap[r.AccountID]
		if !ok {
			continue // Account might have been deleted/inactive
		}

		node.PeriodDebit = r.PeriodDebit
		node.PeriodCredit = r.PeriodCredit

		// Determine Opening Balance
		isBalanceSheet := node.Type == dto.TypeAsset || node.Type == dto.TypeLiability || node.Type == dto.TypeEquity
		var obDebit, obCredit float64

		if isBalanceSheet {
			obDebit = r.AllPriorDebit
			obCredit = r.AllPriorCredit
		} else {
			// P&L accounts only carry over YTD before dateFrom
			obDebit = r.YtdPriorDebit
			obCredit = r.YtdPriorCredit
		}

		if node.NormalBalance == dto.NormalDebit {
			node.OpeningBalance = obDebit - obCredit
			node.PeriodMovement = r.PeriodDebit - r.PeriodCredit
		} else {
			node.OpeningBalance = obCredit - obDebit
			node.PeriodMovement = r.PeriodCredit - r.PeriodDebit
		}
		node.EndingBalance = node.OpeningBalance + node.PeriodMovement
	}

	// 4. Build Hierarchy & Rollup (Bottom-up)
	return a.buildAndRollupTree(nodeMap)
}

func (a *Aggregator) buildAndRollupTree(nodeMap map[uint]*dto.AccountNode) (map[uint]*dto.AccountNode, error) {
	// Link children to parents
	var roots []*dto.AccountNode
	for _, node := range nodeMap {
		if node.ParentID != nil {
			if parent, ok := nodeMap[*node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				// Parent not found, treat as root
				roots = append(roots, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	// Recursive Rollup (updates parents based on children)
	for _, root := range roots {
		a.rollupNode(root)
	}

	return nodeMap, nil
}

func (a *Aggregator) rollupNode(node *dto.AccountNode) {
	// If it has children, its balances are purely the sum of its children's balances.
	// (Assuming parent accounts don't receive direct journals. If they do, this adds to them).
	
	// We only add children to parent if parent is the same type/normal balance. 
	// In SAK ETAP, they should be the same.
	for _, child := range node.Children {
		a.rollupNode(child)
		
		// If normal balances match, we just add.
		// If child is "contra" (e.g. Accum Depr is Credit but Parent Fixed Asset is Debit), 
		// we subtract.
		if node.NormalBalance == child.NormalBalance {
			node.OpeningBalance += child.OpeningBalance
			node.PeriodDebit += child.PeriodDebit
			node.PeriodCredit += child.PeriodCredit
			node.PeriodMovement += child.PeriodMovement
			node.EndingBalance += child.EndingBalance
		} else {
			node.OpeningBalance -= child.OpeningBalance
			node.PeriodDebit += child.PeriodDebit
			node.PeriodCredit += child.PeriodCredit
			node.PeriodMovement -= child.PeriodMovement
			node.EndingBalance -= child.EndingBalance
		}
	}
}
