package service

import (
	"context"
	"errors"

	"github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/model"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type AccountingService interface {
	// CoA
	CreateAccount(ctx context.Context, orgID uint, req dto.AccountCreateRequest) (*dto.AccountResponse, error)
	ListAccounts(ctx context.Context, orgID uint) ([]dto.AccountResponse, error)

	// Journal Entries
	CreateJournalEntry(ctx context.Context, orgID uint, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error)
	ListJournalEntries(ctx context.Context, orgID uint) ([]dto.JournalEntryResponse, error)
	GetJournalEntryByID(ctx context.Context, orgID, id uint) (*dto.JournalEntryResponse, error)
}

type accountingService struct {
	repo repository.AccountingRepository
}

func NewAccountingService(repo repository.AccountingRepository) AccountingService {
	return &accountingService{repo: repo}
}

func (s *accountingService) CreateAccount(ctx context.Context, orgID uint, req dto.AccountCreateRequest) (*dto.AccountResponse, error) {
	if _, err := s.repo.GetAccountByCode(ctx, orgID, req.Code); err == nil {
		return nil, errors.New("account code already exists")
	}

	account := &model.Account{
		Code:          req.Code,
		Name:          req.Name,
		Type:          req.Type,
		NormalBalance: req.NormalBalance,
		ParentID:      req.ParentID,
		Description:   req.Description,
	}
	account.OrganizationID = orgID

	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return nil, err
	}

	return s.mapAccountToResponse(account), nil
}

func (s *accountingService) ListAccounts(ctx context.Context, orgID uint) ([]dto.AccountResponse, error) {
	accounts, err := s.repo.ListAccounts(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.AccountResponse
	for _, a := range accounts {
		res = append(res, *s.mapAccountToResponse(&a))
	}
	return res, nil
}

func (s *accountingService) CreateJournalEntry(ctx context.Context, orgID uint, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error) {
	// Validate double-entry balance
	var totalDebit, totalCredit float64
	var lines []model.JournalEntryLine

	for _, reqLine := range req.Lines {
		// Ensure account exists
		if _, err := s.repo.GetAccountByID(ctx, orgID, reqLine.AccountID); err != nil {
			return nil, errors.New("invalid account ID in journal entry lines")
		}

		if reqLine.Debit == 0 && reqLine.Credit == 0 {
			return nil, errors.New("journal entry line must have either debit or credit")
		}
		if reqLine.Debit > 0 && reqLine.Credit > 0 {
			return nil, errors.New("journal entry line cannot have both debit and credit")
		}

		totalDebit += reqLine.Debit
		totalCredit += reqLine.Credit

		line := model.JournalEntryLine{
			AccountID:   reqLine.AccountID,
			Description: reqLine.Description,
			Debit:       reqLine.Debit,
			Credit:      reqLine.Credit,
		}
		line.OrganizationID = orgID
		lines = append(lines, line)
	}

	if totalDebit != totalCredit {
		return nil, errors.New("journal entry debit and credit totals must balance")
	}

	entry := &model.JournalEntry{
		ReferenceNumber: utils.GenerateCode("JE"),
		Date:            req.Date,
		Description:     req.Description,
		Lines:           lines,
	}
	entry.OrganizationID = orgID

	if err := s.repo.CreateJournalEntry(ctx, entry); err != nil {
		return nil, err
	}

	return s.mapJournalEntryToResponse(entry), nil
}

func (s *accountingService) ListJournalEntries(ctx context.Context, orgID uint) ([]dto.JournalEntryResponse, error) {
	entries, err := s.repo.ListJournalEntries(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.JournalEntryResponse
	for _, e := range entries {
		res = append(res, *s.mapJournalEntryToResponse(&e))
	}
	return res, nil
}

func (s *accountingService) GetJournalEntryByID(ctx context.Context, orgID, id uint) (*dto.JournalEntryResponse, error) {
	entry, err := s.repo.GetJournalEntryByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return s.mapJournalEntryToResponse(entry), nil
}

func (s *accountingService) mapAccountToResponse(a *model.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:            a.ID,
		Code:          a.Code,
		Name:          a.Name,
		Type:          a.Type,
		NormalBalance: a.NormalBalance,
		ParentID:      a.ParentID,
		IsActive:      a.IsActive,
		Description:   a.Description,
	}
}

func (s *accountingService) mapJournalEntryToResponse(e *model.JournalEntry) *dto.JournalEntryResponse {
	res := &dto.JournalEntryResponse{
		ID:              e.ID,
		ReferenceNumber: e.ReferenceNumber,
		Date:            e.Date,
		Description:     e.Description,
		Status:          e.Status,
		CreatedAt:       e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, l := range e.Lines {
		res.Lines = append(res.Lines, dto.JournalEntryLineResponse{
			ID:          l.ID,
			AccountID:   l.AccountID,
			Description: l.Description,
			Debit:       l.Debit,
			Credit:      l.Credit,
		})
	}
	return res
}
