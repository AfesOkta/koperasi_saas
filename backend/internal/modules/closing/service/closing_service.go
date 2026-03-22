package service

import (
	"context"
	"fmt"
	"time"

	accountingDto "github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	accountingService "github.com/koperasi-gresik/backend/internal/modules/accounting/service"
	"github.com/koperasi-gresik/backend/internal/modules/closing/model"
	"github.com/koperasi-gresik/backend/internal/modules/closing/repository"
	loanRepo "github.com/koperasi-gresik/backend/internal/modules/loan/repository"
	orgRepo "github.com/koperasi-gresik/backend/internal/modules/organization/repository"
	savingRepoModel "github.com/koperasi-gresik/backend/internal/modules/savings/model"
	savingRepo "github.com/koperasi-gresik/backend/internal/modules/savings/repository"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/sirupsen/logrus"
)

type ClosingService interface {
	ProcessAllTenantsEOD(ctx context.Context, date string) error
	ProcessAllTenantsEOM(ctx context.Context, month, year int) error
	ProcessEOD(ctx context.Context, orgID uint, date string) error
	ProcessEOM(ctx context.Context, orgID uint, month, year int) error
}

type closingService struct {
	repo       repository.ClosingRepository
	loanRepo   loanRepo.LoanRepository
	savingRepo savingRepo.SavingRepository
	accSvc     accountingService.AccountingService
	orgRepo    orgRepo.OrganizationRepository
}

func NewClosingService(
	repo repository.ClosingRepository,
	loanRepo loanRepo.LoanRepository,
	savingRepo savingRepo.SavingRepository,
	accSvc accountingService.AccountingService,
	orgRepo orgRepo.OrganizationRepository,
) ClosingService {
	return &closingService{
		repo:       repo,
		loanRepo:   loanRepo,
		savingRepo: savingRepo,
		accSvc:     accSvc,
		orgRepo:    orgRepo,
	}
}

func (s *closingService) ProcessAllTenantsEOD(ctx context.Context, date string) error {
	orgs, _, err := s.orgRepo.List(ctx, pagination.Params{Limit: 1000})
	if err != nil {
		return err
	}

	for _, org := range orgs {
		if err := s.ProcessEOD(ctx, org.ID, date); err != nil {
			logrus.Errorf("[ClosingService] EOD failed for Org %d: %v", org.ID, err)
		}
	}
	return nil
}

func (s *closingService) ProcessAllTenantsEOM(ctx context.Context, month, year int) error {
	orgs, _, err := s.orgRepo.List(ctx, pagination.Params{Limit: 1000})
	if err != nil {
		return err
	}

	for _, org := range orgs {
		if err := s.ProcessEOM(ctx, org.ID, month, year); err != nil {
			logrus.Errorf("[ClosingService] EOM failed for Org %d: %v", org.ID, err)
		}
	}
	return nil
}

func (s *closingService) ProcessEOD(ctx context.Context, orgID uint, date string) error {
	// 1. Idempotency Check
	existing, err := s.repo.GetLogByDate(ctx, orgID, date, "EOD")
	if err == nil && (existing.Status == "SUCCESS" || existing.Status == "RUNNING") {
		return nil // Already processed or in progress
	}

	// 2. Create/Update Log
	logEntry := &model.ClosingLog{
		Date:      date,
		Type:      "EOD",
		Status:    "RUNNING",
		StartedAt: pointerTime(time.Now()),
	}
	logEntry.OrganizationID = orgID
	if existing != nil && existing.ID > 0 {
		logEntry.ID = existing.ID
		_ = s.repo.UpdateLog(ctx, logEntry)
	} else {
		_ = s.repo.CreateLog(ctx, logEntry)
	}

	defer func() {
		now := time.Now()
		logEntry.FinishedAt = &now
		_ = s.repo.UpdateLog(ctx, logEntry)
	}()

	// 3. Process Loans (Overdue & Penalties)
	loans, err := s.loanRepo.ListLoansByStatus(ctx, orgID, []string{"active", "overdue"})
	if err != nil {
		logEntry.Status = "FAILED"
		logEntry.ErrorMessage = err.Error()
		return err
	}

	today, _ := time.Parse("2006-01-02", date)
	for _, loan := range loans {
		loanUpdated := false
		var totalPenalty float64

		for i, sched := range loan.Schedules {
			if (sched.Status == "unpaid" || sched.Status == "partial") {
				dueDate, _ := time.Parse("2006-01-02", sched.DueDate)
				if dueDate.Before(today) {
					// Mark as overdue
					if sched.Status != "overdue" {
						loan.Schedules[i].Status = "overdue"
						loan.Status = "overdue"
						loanUpdated = true
					}

					// Calculate Penalty (Denda)
					// Note: For MVP, we need the Product to get penalty config
					product, _ := s.loanRepo.GetProductByID(ctx, orgID, loan.LoanProductID)
					if product != nil && product.PenaltyRate > 0 {
						overdueDays := int(today.Sub(dueDate).Hours() / 24)
						if overdueDays > product.PenaltyGraceDays {
							penalty := 0.0
							if product.PenaltyType == "percentage" {
								penalty = loan.PrincipalAmount * product.PenaltyRate
							} else {
								penalty = product.PenaltyRate // Flat
							}
							
							if product.PenaltyCap > 0 && penalty > product.PenaltyCap {
								penalty = product.PenaltyCap
							}
							totalPenalty += penalty
						}
					}
				}
			}
		}

		if loanUpdated {
			_ = s.loanRepo.UpdateLoanStatus(ctx, &loan)
		}

		// Post Penalty Journal
		if totalPenalty > 0 {
			_, _ = s.accSvc.CreateJournalEntryIdempotent(ctx, orgID, fmt.Sprintf("EOD-PEN-%d-%s", loan.ID, date), accountingDto.JournalEntryCreateRequest{
				Date:         date,
				Description:  fmt.Sprintf("Late fee penalty for Loan %s", loan.LoanNumber),
				SourceModule: "closing",
				Lines: []accountingDto.JournalEntryLineRequest{
					{AccountCode: "1201", Debit: totalPenalty, Description: "Loan Receivable Increase (Penalty)"},
					{AccountCode: "4201", Credit: totalPenalty, Description: "Penalty Income"},
				},
			})
		}
	}

	logEntry.Status = "SUCCESS"
	return nil
}

func (s *closingService) ProcessEOM(ctx context.Context, orgID uint, month, year int) error {
	// Usually interest is posted on the last day of the month or 1st of next.
	// We'll use the last day of the month being closed.
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	date := lastDay.Format("2006-01-02")

	// 1. Idempotency Check
	existing, err := s.repo.GetLogByDate(ctx, orgID, date, "EOM")
	if err == nil && (existing.Status == "SUCCESS" || existing.Status == "RUNNING") {
		return nil
	}

	// 2. Log entry
	logEntry := &model.ClosingLog{
		Date:      date,
		Type:      "EOM",
		Status:    "RUNNING",
		StartedAt: pointerTime(time.Now()),
	}
	logEntry.OrganizationID = orgID
	_ = s.repo.CreateLog(ctx, logEntry)

	defer func() {
		now := time.Now()
		logEntry.FinishedAt = &now
		_ = s.repo.UpdateLog(ctx, logEntry)
	}()

	// 3. Process Savings Interest
	accounts, err := s.savingRepo.ListAccounts(ctx, orgID)
	if err != nil {
		logEntry.Status = "FAILED"
		logEntry.ErrorMessage = err.Error()
		return err
	}

	for _, acc := range accounts {
		// Get Product for interest rate
		product, _ := s.savingRepo.GetProductByID(ctx, orgID, acc.SavingProductID)
		if product == nil || product.InterestRate <= 0 {
			continue
		}

		// Calculate Min Balance
		minVal, err := s.savingRepo.GetMinBalance(ctx, acc.ID, month, year)
		if err != nil || minVal <= 0 {
			continue
		}

		interest := minVal * (product.InterestRate / 100 / 12)
		if interest > 0 {
			// Actually perform transaction in Saving module? 
			// We should probably call SavingService or just use Repository if we want to bypass validation.
			// Best is to call SavingService to ensure proper event emission if any.
			
			// For Step 3, we'll use repos directly to avoid circular dependency if service was complex
			// But wait, SavingRepository has ExecuteTransaction.
			txn := &savingRepoModel.SavingTransaction{
				SavingAccountID: acc.ID,
				ReferenceNumber: fmt.Sprintf("INT-%d-%d%02d", acc.ID, year, month),
				Type:            "interest",
				Amount:          interest,
				Description:     fmt.Sprintf("Monthly interest for %02d/%d", month, year),
			}
			txn.OrganizationID = orgID
			
			// Note: savingRepoModel needs to be imported or use full path.
			// I'll fix imports.
			
			if err := s.savingRepo.ExecuteTransaction(ctx, &acc, txn); err == nil {
				// Post Journal
				_, _ = s.accSvc.CreateJournalEntryIdempotent(ctx, orgID, fmt.Sprintf("EOM-INT-%d-%02d%d", acc.ID, month, year), accountingDto.JournalEntryCreateRequest{
					Date:         date,
					Description:  fmt.Sprintf("Monthly interest payout for account %s", acc.AccountNumber),
					SourceModule: "closing",
					Lines: []accountingDto.JournalEntryLineRequest{
						{AccountCode: "5201", Debit: interest, Description: "Interest Expense"},
						{AccountCode: "2101", Credit: interest, Description: "Savings Liability Increase (Interest)"},
					},
				})
			}
		}
	}

	// 4. Close Period
	_ = s.repo.ClosePeriod(ctx, &model.ClosedPeriod{
		Month:    month,
		Year:     year,
		ClosedAt: time.Now().Format(time.RFC3339),
		ClosedBy: 0, // System
	})

	logEntry.Status = "SUCCESS"
	return nil
}

func pointerTime(t time.Time) *time.Time {
	return &t
}
