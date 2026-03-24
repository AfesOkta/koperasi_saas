package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/loan/model"
	"gorm.io/gorm"
)

type LoanRepository interface {
	// Products
	CreateProduct(ctx context.Context, product *model.LoanProduct) error
	GetProductByID(ctx context.Context, orgID, id uint) (*model.LoanProduct, error)
	ListProducts(ctx context.Context, orgID uint) ([]model.LoanProduct, error)

	// Loans
	CreateLoan(ctx context.Context, loan *model.Loan) error
	GetLoanByID(ctx context.Context, orgID, id uint) (*model.Loan, error)
	ListLoansByMember(ctx context.Context, orgID, memberID uint) ([]model.Loan, error)
	UpdateLoanStatus(ctx context.Context, loan *model.Loan) error
	ListLoansByStatus(ctx context.Context, orgID uint, statuses []string) ([]model.Loan, error)
	HasOutstandingLoan(ctx context.Context, orgID, memberID uint) (bool, error)

	// Payments
	RecordPayment(ctx context.Context, loan *model.Loan, payment *model.LoanPayment, schedulesToUpdate []model.LoanSchedule) error

	// Collaterals
	CreateCollateral(ctx context.Context, coll *model.LoanCollateral) error
	GetCollateralsByLoanID(ctx context.Context, loanID uint) ([]model.LoanCollateral, error)

	// Approval Logs
	CreateApprovalLog(ctx context.Context, log *model.ApprovalLog) error
	GetApprovalLogsByLoanID(ctx context.Context, loanID uint) ([]model.ApprovalLog, error)

	// Utils
	GetDB() *gorm.DB
}

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *loanRepository) CreateProduct(ctx context.Context, product *model.LoanProduct) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *loanRepository) GetProductByID(ctx context.Context, orgID, id uint) (*model.LoanProduct, error) {
	var product model.LoanProduct
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&product, id).Error
	return &product, err
}

func (r *loanRepository) ListProducts(ctx context.Context, orgID uint) ([]model.LoanProduct, error) {
	var products []model.LoanProduct
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&products).Error
	return products, err
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan *model.Loan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(loan).Error
	})
}

func (r *loanRepository) GetLoanByID(ctx context.Context, orgID, id uint) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.WithContext(ctx).
		Preload("Schedules").
		Preload("Payments").
		Preload("Collaterals").
		Preload("ApprovalLogs").
		Where("organization_id = ?", orgID).
		First(&loan, id).Error
	return &loan, err
}

func (r *loanRepository) ListLoansByMember(ctx context.Context, orgID, memberID uint) ([]model.Loan, error) {
	var loans []model.Loan
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND member_id = ?", orgID, memberID).
		Order("created_at DESC").
		Find(&loans).Error
	return loans, err
}

func (r *loanRepository) UpdateLoanStatus(ctx context.Context, loan *model.Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

func (r *loanRepository) ListLoansByStatus(ctx context.Context, orgID uint, statuses []string) ([]model.Loan, error) {
	var loans []model.Loan
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND status IN ?", orgID, statuses).
		Preload("Schedules").
		Find(&loans).Error
	return loans, err
}

func (r *loanRepository) HasOutstandingLoan(ctx context.Context, orgID, memberID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Loan{}).
		Where("organization_id = ? AND member_id = ? AND outstanding > 0 AND status NOT IN ?", orgID, memberID, []string{"paid", "defaulted", "rejected"}).
		Count(&count).Error
	return count > 0, err
}

func (r *loanRepository) RecordPayment(ctx context.Context, loan *model.Loan, payment *model.LoanPayment, schedulesToUpdate []model.LoanSchedule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock loan
		var lockedLoan model.Loan
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&lockedLoan, loan.ID).Error; err != nil {
			return err
		}

		// 2. Save payment
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 3. Update schedules
		for _, s := range schedulesToUpdate {
			if err := tx.Save(&s).Error; err != nil {
				return err
			}
		}

		// 4. Update loan balances
		lockedLoan.Outstanding = loan.Outstanding
		if lockedLoan.Outstanding <= 0 {
			lockedLoan.Status = "paid"
		}
		if err := tx.Save(&lockedLoan).Error; err != nil {
			return err
		}

		return nil
	})
}

// Collateral methods
func (r *loanRepository) CreateCollateral(ctx context.Context, coll *model.LoanCollateral) error {
	return r.db.WithContext(ctx).Create(coll).Error
}

func (r *loanRepository) GetCollateralsByLoanID(ctx context.Context, loanID uint) ([]model.LoanCollateral, error) {
	var cols []model.LoanCollateral
	err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Find(&cols).Error
	return cols, err
}

// Approval Log methods
func (r *loanRepository) CreateApprovalLog(ctx context.Context, log *model.ApprovalLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *loanRepository) GetApprovalLogsByLoanID(ctx context.Context, loanID uint) ([]model.ApprovalLog, error) {
	var logs []model.ApprovalLog
	err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Order("created_at ASC").Find(&logs).Error
	return logs, err
}

