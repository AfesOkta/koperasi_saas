package service

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/cash/dto"
	"github.com/koperasi-gresik/backend/internal/modules/cash/model"
	"github.com/koperasi-gresik/backend/internal/modules/cash/repository"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type CashService interface {
	CreateRegister(ctx context.Context, orgID uint, req dto.CashRegisterCreateRequest) (*dto.CashRegisterResponse, error)
	ListRegisters(ctx context.Context, orgID uint) ([]dto.CashRegisterResponse, error)

	RecordTransaction(ctx context.Context, orgID, registerID uint, req dto.CashTransactionRequest) (*dto.CashTransactionResponse, error)
	ListTransactionsByRegister(ctx context.Context, orgID, registerID uint, params pagination.Params) ([]dto.CashTransactionResponse, int64, error)
}

type cashService struct {
	repo repository.CashRepository
}

func NewCashService(repo repository.CashRepository) CashService {
	return &cashService{repo: repo}
}

func (s *cashService) CreateRegister(ctx context.Context, orgID uint, req dto.CashRegisterCreateRequest) (*dto.CashRegisterResponse, error) {
	register := &model.CashRegister{
		Name:        req.Name,
		Type:        req.Type,
		AccountID:   req.AccountID,
		Description: req.Description,
	}
	register.OrganizationID = orgID

	if err := s.repo.CreateRegister(ctx, register); err != nil {
		return nil, err
	}

	return s.mapRegisterToResponse(register), nil
}

func (s *cashService) ListRegisters(ctx context.Context, orgID uint) ([]dto.CashRegisterResponse, error) {
	registers, err := s.repo.ListRegisters(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.CashRegisterResponse
	for _, r := range registers {
		res = append(res, *s.mapRegisterToResponse(&r))
	}
	return res, nil
}

func (s *cashService) RecordTransaction(ctx context.Context, orgID, registerID uint, req dto.CashTransactionRequest) (*dto.CashTransactionResponse, error) {
	register, err := s.repo.GetRegisterByID(ctx, orgID, registerID)
	if err != nil {
		return nil, err
	}

	txn := &model.CashTransaction{
		CashRegisterID:  register.ID,
		ReferenceNumber: utils.GenerateCode("CSH"),
		Type:            req.Type,
		Amount:          req.Amount,
		Category:        req.Category,
		Description:     req.Description,
		RelatedEntity:   req.RelatedEntity,
		RelatedEntityID: req.RelatedEntityID,
	}
	txn.OrganizationID = orgID

	// Note: We don't implement double entry here, this merely bumps the cash amount.
	// However, a true system triggers a domain event for Accounting to listen to, which generates a Journal Entry automatically.

	if err := s.repo.ExecuteTransaction(ctx, register, txn); err != nil {
		return nil, err
	}

	return s.mapTransactionToResponse(txn), nil
}

func (s *cashService) ListTransactionsByRegister(ctx context.Context, orgID, registerID uint, params pagination.Params) ([]dto.CashTransactionResponse, int64, error) {
	txns, total, err := s.repo.ListTransactionsByRegister(ctx, orgID, registerID, params)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.CashTransactionResponse
	for _, t := range txns {
		res = append(res, *s.mapTransactionToResponse(&t))
	}
	return res, total, nil
}

func (s *cashService) mapRegisterToResponse(r *model.CashRegister) *dto.CashRegisterResponse {
	return &dto.CashRegisterResponse{
		ID:          r.ID,
		Name:        r.Name,
		Type:        r.Type,
		Balance:     r.Balance,
		Status:      r.Status,
		AccountID:   r.AccountID,
		Description: r.Description,
	}
}

func (s *cashService) mapTransactionToResponse(t *model.CashTransaction) *dto.CashTransactionResponse {
	return &dto.CashTransactionResponse{
		ID:              t.ID,
		CashRegisterID:  t.CashRegisterID,
		ReferenceNumber: t.ReferenceNumber,
		Type:            t.Type,
		Amount:          t.Amount,
		BalanceAfter:    t.BalanceAfter,
		Category:        t.Category,
		Description:     t.Description,
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
