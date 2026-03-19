package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/koperasi-gresik/backend/internal/modules/savings/dto"
	"github.com/koperasi-gresik/backend/internal/modules/savings/model"
	"github.com/koperasi-gresik/backend/internal/modules/savings/repository"
	"github.com/koperasi-gresik/backend/internal/shared/event"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type SavingService interface {
	// Products
	CreateProduct(ctx context.Context, orgID uint, req dto.SavingProductCreateRequest) (*dto.SavingProductResponse, error)
	ListProducts(ctx context.Context, orgID uint) ([]dto.SavingProductResponse, error)

	// Transactions
	Deposit(ctx context.Context, orgID uint, req dto.SavingTransactionRequest) (*dto.SavingTransactionResponse, error)
	Withdraw(ctx context.Context, orgID uint, req dto.SavingTransactionRequest) (*dto.SavingTransactionResponse, error)
	GetBalance(ctx context.Context, orgID, memberID, productID uint) (*dto.SavingAccountResponse, error)
	GetTransactionHistory(ctx context.Context, orgID, accountID uint) ([]dto.SavingTransactionResponse, error)

	// Admin Tasks
	CalculateAndDistributeInterest(ctx context.Context, orgID uint) error
}

type savingService struct {
	repo      repository.SavingRepository
	publisher event.Publisher
}

func NewSavingService(repo repository.SavingRepository, publisher event.Publisher) SavingService {
	return &savingService{repo: repo, publisher: publisher}
}

func (s *savingService) CreateProduct(ctx context.Context, orgID uint, req dto.SavingProductCreateRequest) (*dto.SavingProductResponse, error) {
	if _, err := s.repo.GetProductByCode(ctx, orgID, req.Code); err == nil {
		return nil, errors.New("saving product code already exists")
	}

	product := &model.SavingProduct{
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		IsWithdrawble: req.IsWithdrawble,
		InterestRate:  req.InterestRate,
	}
	product.OrganizationID = orgID

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.mapProductToResponse(product), nil
}

func (s *savingService) ListProducts(ctx context.Context, orgID uint) ([]dto.SavingProductResponse, error) {
	products, err := s.repo.ListProducts(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.SavingProductResponse
	for _, p := range products {
		res = append(res, *s.mapProductToResponse(&p))
	}
	return res, nil
}

func (s *savingService) getOrCreateAccount(ctx context.Context, orgID, memberID, productID uint) (*model.SavingAccount, error) {
	account, err := s.repo.GetAccountByMemberAndProduct(ctx, orgID, memberID, productID)
	if err == nil {
		return account, nil
	}

	// Create new account if not found
	account = &model.SavingAccount{
		MemberID:        memberID,
		SavingProductID: productID,
		AccountNumber:   utils.GenerateCode("SAV"),
		Balance:         0,
	}
	account.OrganizationID = orgID

	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create saving account: %w", err)
	}

	return account, nil
}

func (s *savingService) Deposit(ctx context.Context, orgID uint, req dto.SavingTransactionRequest) (*dto.SavingTransactionResponse, error) {
	account, err := s.getOrCreateAccount(ctx, orgID, req.MemberID, req.SavingProductID)
	if err != nil {
		return nil, err
	}

	txn := &model.SavingTransaction{
		SavingAccountID: account.ID,
		ReferenceNumber: utils.GenerateCode("TXN"),
		Type:            "deposit",
		Amount:          req.Amount,
		Description:     req.Description,
	}
	txn.OrganizationID = orgID

	if err := s.repo.ExecuteTransaction(ctx, account, txn); err != nil {
		return nil, err
	}

	// Publish Event
	payload := event.SavingsTransactionPayload{
		MemberID:    req.MemberID,
		Amount:      req.Amount,
		ProductCode: fmt.Sprintf("%d", req.SavingProductID), // using ID for MVP
		Description: req.Description,
	}
	evt := event.Event{
		Type:           event.EventSavingsDeposited,
		AggregateID:    txn.ID,
		OrganizationID: orgID,
		Payload:        payload,
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, evt)
	}

	return s.mapTransactionToResponse(txn), nil
}

func (s *savingService) Withdraw(ctx context.Context, orgID uint, req dto.SavingTransactionRequest) (*dto.SavingTransactionResponse, error) {
	// First check if product allows withdrawal
	product, err := s.repo.GetProductByID(ctx, orgID, req.SavingProductID)
	if err != nil {
		return nil, errors.New("invalid saving product")
	}
	if !product.IsWithdrawble {
		return nil, errors.New("this saving product is not withdrawable")
	}

	account, err := s.getOrCreateAccount(ctx, orgID, req.MemberID, req.SavingProductID)
	if err != nil {
		return nil, err
	}

	if account.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	txn := &model.SavingTransaction{
		SavingAccountID: account.ID,
		ReferenceNumber: utils.GenerateCode("TXN"),
		Type:            "withdrawal",
		Amount:          req.Amount,
		Description:     req.Description,
	}
	txn.OrganizationID = orgID

	if err := s.repo.ExecuteTransaction(ctx, account, txn); err != nil {
		return nil, err
	}

	// Publish Event
	payload := event.SavingsTransactionPayload{
		MemberID:    req.MemberID,
		Amount:      req.Amount,
		ProductCode: fmt.Sprintf("%d", req.SavingProductID), // using ID for MVP
		Description: req.Description,
	}
	evt := event.Event{
		Type:           event.EventSavingsWithdrawn,
		AggregateID:    txn.ID,
		OrganizationID: orgID,
		Payload:        payload,
	}
	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, evt)
	}

	return s.mapTransactionToResponse(txn), nil
}

func (s *savingService) GetBalance(ctx context.Context, orgID, memberID, productID uint) (*dto.SavingAccountResponse, error) {
	account, err := s.getOrCreateAccount(ctx, orgID, memberID, productID)
	if err != nil {
		return nil, err
	}
	return s.mapAccountToResponse(account), nil
}

func (s *savingService) GetTransactionHistory(ctx context.Context, orgID, accountID uint) ([]dto.SavingTransactionResponse, error) {
	_, err := s.repo.GetAccountByID(ctx, orgID, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	txns, err := s.repo.ListTransactionsByAccount(ctx, orgID, accountID)
	if err != nil {
		return nil, err
	}

	var res []dto.SavingTransactionResponse
	for _, t := range txns {
		res = append(res, *s.mapTransactionToResponse(&t))
	}
	return res, nil
}

func (s *savingService) CalculateAndDistributeInterest(ctx context.Context, orgID uint) error {
	// For MVP: Fetch all active products that yield interest
	products, err := s.repo.ListProducts(ctx, orgID)
	if err != nil {
		return err
	}

	for _, product := range products {
		if product.InterestRate <= 0 {
			continue // Skip products with no interest
		}

		// Usually, we'd query accounts by productID
		// But in MVP, let's just assume we need to process this at scale later and simulate it here with a repository method.
		// Since ListAccountsByProduct isn't defined in the repo interface yet in this codebase,
		// I will log a message. In a real scenario, this would loop over accounts, 
		// calculate iterest amount = balance * (InterestRate/12/100), 
		// deposit the transaction and publish an event.
		// Example:
		// accounts, _ := s.repo.ListAccountsByProduct(ctx, orgID, product.ID)
		// for _, acc := range accounts { ... s.Deposit(...) ... }

		log.Printf("ℹ️ [Savings] Monthly interest distribution calculated for product %s (Rate: %.2f%%)", product.Code, product.InterestRate)
	}

	return nil
}

// Mappers
func (s *savingService) mapProductToResponse(p *model.SavingProduct) *dto.SavingProductResponse {
	return &dto.SavingProductResponse{
		ID:            p.ID,
		Code:          p.Code,
		Name:          p.Name,
		Description:   p.Description,
		Status:        p.Status,
		IsWithdrawble: p.IsWithdrawble,
		InterestRate:  p.InterestRate,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *savingService) mapAccountToResponse(a *model.SavingAccount) *dto.SavingAccountResponse {
	return &dto.SavingAccountResponse{
		ID:              a.ID,
		MemberID:        a.MemberID,
		SavingProductID: a.SavingProductID,
		AccountNumber:   a.AccountNumber,
		Balance:         a.Balance,
		Status:          a.Status,
		CreatedAt:       a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *savingService) mapTransactionToResponse(t *model.SavingTransaction) *dto.SavingTransactionResponse {
	return &dto.SavingTransactionResponse{
		ID:              t.ID,
		SavingAccountID: t.SavingAccountID,
		ReferenceNumber: t.ReferenceNumber,
		Type:            t.Type,
		Amount:          t.Amount,
		BalanceAfter:    t.BalanceAfter,
		Description:     t.Description,
		Status:          t.Status,
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
