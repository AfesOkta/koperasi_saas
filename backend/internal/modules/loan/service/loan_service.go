package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/loan/dto"
	"github.com/koperasi-gresik/backend/internal/modules/loan/model"
	"github.com/koperasi-gresik/backend/internal/modules/loan/repository"
	"github.com/koperasi-gresik/backend/internal/shared/event"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type LoanService interface {
	CreateProduct(ctx context.Context, orgID uint, req dto.LoanProductCreateRequest) (*dto.LoanProductResponse, error)
	ListProducts(ctx context.Context, orgID uint) ([]dto.LoanProductResponse, error)

	ApplyForLoan(ctx context.Context, orgID uint, req dto.LoanApplicationRequest) (*dto.LoanResponse, error)
	GetLoanByID(ctx context.Context, orgID, id uint) (*dto.LoanResponse, error)

	ApproveLoan(ctx context.Context, orgID, loanID uint) error
	RecordPayment(ctx context.Context, orgID, loanID uint, req dto.LoanPaymentRequest) (*dto.LoanPaymentResponse, error)
}

type loanService struct {
	repo      repository.LoanRepository
	publisher event.Publisher
}

func NewLoanService(repo repository.LoanRepository, publisher event.Publisher) LoanService {
	return &loanService{repo: repo, publisher: publisher}
}

func (s *loanService) CreateProduct(ctx context.Context, orgID uint, req dto.LoanProductCreateRequest) (*dto.LoanProductResponse, error) {
	product := &model.LoanProduct{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		InterestRate: req.InterestRate,
		InterestType: req.InterestType,
		MaxAmount:    req.MaxAmount,
		MaxTerm:      req.MaxTerm,
	}
	product.OrganizationID = orgID

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.mapProductToResponse(product), nil
}

func (s *loanService) ListProducts(ctx context.Context, orgID uint) ([]dto.LoanProductResponse, error) {
	products, err := s.repo.ListProducts(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.LoanProductResponse
	for _, p := range products {
		res = append(res, *s.mapProductToResponse(&p))
	}
	return res, nil
}

func (s *loanService) ApplyForLoan(ctx context.Context, orgID uint, req dto.LoanApplicationRequest) (*dto.LoanResponse, error) {
	product, err := s.repo.GetProductByID(ctx, orgID, req.LoanProductID)
	if err != nil {
		return nil, errors.New("loan product not found")
	}

	if req.PrincipalAmount > product.MaxAmount {
		return nil, errors.New("requested amount exceeds product maximum")
	}
	if req.TermMonths > product.MaxTerm {
		return nil, errors.New("requested term exceeds product maximum")
	}

	loan := &model.Loan{
		MemberID:        req.MemberID,
		LoanProductID:   req.LoanProductID,
		LoanNumber:      utils.GenerateCode("LOAN"),
		PrincipalAmount: req.PrincipalAmount,
		InterestRate:    product.InterestRate,
		TermMonths:      req.TermMonths,
		Status:          "pending",
	}
	loan.OrganizationID = orgID

	// Generate schedules (Flat vs Reducing Balance)
	var schedules []model.LoanSchedule
	var totalInterest, expectedTotal float64

	if product.InterestType == "reducing" || product.InterestType == "effective" {
		r := (product.InterestRate / 100) / 12
		n := float64(req.TermMonths)
		pmt := (req.PrincipalAmount * r * math.Pow(1+r, n)) / (math.Pow(1+r, n) - 1)
		if math.IsNaN(pmt) || math.IsInf(pmt, 0) {
			pmt = req.PrincipalAmount / n // Fallback to 0 interest
		}

		expectedTotal = pmt * n
		totalInterest = expectedTotal - req.PrincipalAmount
		
		loan.TotalInterest = totalInterest
		loan.ExpectedTotal = expectedTotal
		loan.Outstanding = expectedTotal
		
		balance := req.PrincipalAmount
		now := time.Now()
		for i := 1; i <= req.TermMonths; i++ {
			interestForMonth := balance * r
			principalForMonth := pmt - interestForMonth
			
			// Adjust last month rounding
			if i == req.TermMonths {
				principalForMonth = balance
				pmt = principalForMonth + interestForMonth
			}

			dueDate := now.AddDate(0, i, 0).Format("2006-01-02")
			schedule := model.LoanSchedule{
				Period:          i,
				DueDate:         dueDate,
				PrincipalAmount: principalForMonth,
				InterestAmount:  interestForMonth,
				TotalAmount:     pmt,
				Status:          "unpaid",
			}
			schedule.OrganizationID = orgID
			schedules = append(schedules, schedule)
			
			balance -= principalForMonth
		}
	} else {
		// Flat interest logic
		totalInterest = (req.PrincipalAmount * (product.InterestRate / 100)) * float64(req.TermMonths)
		loan.TotalInterest = totalInterest
		loan.ExpectedTotal = req.PrincipalAmount + totalInterest
		loan.Outstanding = loan.ExpectedTotal

		principalPerMonth := req.PrincipalAmount / float64(req.TermMonths)
		interestPerMonth := totalInterest / float64(req.TermMonths)

		now := time.Now()
		for i := 1; i <= req.TermMonths; i++ {
			dueDate := now.AddDate(0, i, 0).Format("2006-01-02")
			schedule := model.LoanSchedule{
				Period:          i,
				DueDate:         dueDate,
				PrincipalAmount: principalPerMonth,
				InterestAmount:  interestPerMonth,
				TotalAmount:     principalPerMonth + interestPerMonth,
				Status:          "unpaid",
			}
			schedule.OrganizationID = orgID
			schedules = append(schedules, schedule)
		}
	}

	loan.Schedules = schedules

	if err := s.repo.CreateLoan(ctx, loan); err != nil {
		return nil, err
	}

	return s.mapLoanToResponse(loan), nil
}

func (s *loanService) ApproveLoan(ctx context.Context, orgID, loanID uint) error {
	loan, err := s.repo.GetLoanByID(ctx, orgID, loanID)
	if err != nil {
		return err
	}

	if loan.Status != "pending" {
		return errors.New("only pending loans can be approved")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	loan.Status = "active"
	loan.ApprovedAt = &now
	loan.DisbursedAt = &now // For MVP, assume immediate disbursement

	if err := s.repo.UpdateLoanStatus(ctx, loan); err != nil {
		return err
	}

	payload := event.LoanTransactionPayload{
		MemberID:      loan.MemberID,
		Amount:        loan.PrincipalAmount,
		Description:   fmt.Sprintf("Disbursement for Loan %s", loan.LoanNumber),
	}
	evt := event.Event{
		Type:           event.EventLoanDisbursed,
		AggregateID:    loan.ID,
		OrganizationID: orgID,
		Payload:        payload,
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, evt)
	}

	return nil
}

func (s *loanService) RecordPayment(ctx context.Context, orgID, loanID uint, req dto.LoanPaymentRequest) (*dto.LoanPaymentResponse, error) {
	loan, err := s.repo.GetLoanByID(ctx, orgID, loanID)
	if err != nil {
		return nil, err
	}

	if loan.Status != "active" {
		return nil, errors.New("cannot make payments to an inactive loan")
	}

	if req.Amount > loan.Outstanding {
		return nil, errors.New("payment amount exceeds outstanding balance")
	}

	payment := &model.LoanPayment{
		LoanID:          loan.ID,
		ReferenceNumber: utils.GenerateCode("LPAY"),
		Amount:          req.Amount,
		PaymentDate:     time.Now().Format(time.RFC3339),
		Description:     req.Description,
	}
	payment.OrganizationID = orgID

	// Simple logic to distribute payment across unpaid schedules
	remainingPayment := req.Amount
	var schedulesToUpdate []model.LoanSchedule

	for i, schedule := range loan.Schedules {
		if remainingPayment <= 0 {
			break
		}
		if schedule.Status == "paid" {
			continue
		}

		unpaidAmount := schedule.TotalAmount - schedule.PaidAmount
		if remainingPayment >= unpaidAmount {
			schedule.PaidAmount = schedule.TotalAmount
			schedule.Status = "paid"
			remainingPayment -= unpaidAmount
		} else {
			schedule.PaidAmount += remainingPayment
			schedule.Status = "partial"
			remainingPayment = 0
		}
		loan.Schedules[i] = schedule
		schedulesToUpdate = append(schedulesToUpdate, schedule)
	}

	loan.Outstanding -= req.Amount

	if err := s.repo.RecordPayment(ctx, loan, payment, schedulesToUpdate); err != nil {
		return nil, err
	}

	// Prorate principal and interest for the payment
	principalRatio := loan.PrincipalAmount / loan.ExpectedTotal
	principalPart := req.Amount * principalRatio
	interestPart := req.Amount - principalPart

	payload := event.LoanTransactionPayload{
		MemberID:      loan.MemberID,
		Amount:        req.Amount,
		PrincipalPart: principalPart,
		InterestPart:  interestPart,
		Description:   req.Description,
	}
	evt := event.Event{
		Type:           event.EventLoanInstallmentPaid,
		AggregateID:    loan.ID,
		OrganizationID: orgID,
		Payload:        payload,
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, evt)
	}

	return s.mapPaymentToResponse(payment), nil
}

func (s *loanService) GetLoanByID(ctx context.Context, orgID, id uint) (*dto.LoanResponse, error) {
	loan, err := s.repo.GetLoanByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return s.mapLoanToResponse(loan), nil
}

// Mappers
func (s *loanService) mapProductToResponse(p *model.LoanProduct) *dto.LoanProductResponse {
	return &dto.LoanProductResponse{
		ID:           p.ID,
		Code:         p.Code,
		Name:         p.Name,
		Description:  p.Description,
		InterestRate: p.InterestRate,
		InterestType: p.InterestType,
		MaxAmount:    p.MaxAmount,
		MaxTerm:      p.MaxTerm,
		Status:       p.Status,
	}
}

func (s *loanService) mapLoanToResponse(l *model.Loan) *dto.LoanResponse {
	res := &dto.LoanResponse{
		ID:              l.ID,
		MemberID:        l.MemberID,
		LoanProductID:   l.LoanProductID,
		LoanNumber:      l.LoanNumber,
		PrincipalAmount: l.PrincipalAmount,
		InterestRate:    l.InterestRate,
		TermMonths:      l.TermMonths,
		TotalInterest:   l.TotalInterest,
		ExpectedTotal:   l.ExpectedTotal,
		Outstanding:     l.Outstanding,
		Status:          l.Status,
		CreatedAt:       l.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, sch := range l.Schedules {
		res.Schedules = append(res.Schedules, dto.LoanScheduleResponse{
			ID:              sch.ID,
			Period:          sch.Period,
			DueDate:         sch.DueDate,
			PrincipalAmount: sch.PrincipalAmount,
			InterestAmount:  sch.InterestAmount,
			TotalAmount:     sch.TotalAmount,
			Status:          sch.Status,
		})
	}

	for _, pay := range l.Payments {
		res.Payments = append(res.Payments, *s.mapPaymentToResponse(&pay))
	}

	return res
}

func (s *loanService) mapPaymentToResponse(p *model.LoanPayment) *dto.LoanPaymentResponse {
	return &dto.LoanPaymentResponse{
		ID:              p.ID,
		LoanID:          p.LoanID,
		ReferenceNumber: p.ReferenceNumber,
		Amount:          p.Amount,
		PaymentDate:     p.PaymentDate,
		Description:     p.Description,
	}
}
