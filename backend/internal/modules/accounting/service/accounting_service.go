package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/model"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type AccountingService interface {
	// CoA
	CreateAccount(ctx context.Context, orgID uint, req dto.AccountCreateRequest) (*dto.AccountResponse, error)
	ListAccounts(ctx context.Context, orgID uint) ([]dto.AccountResponse, error)
	ListAccountsByType(ctx context.Context, orgID uint, accountType string) ([]dto.AccountResponse, error)
	GetAccountByCode(ctx context.Context, orgID uint, code string) (*dto.AccountResponse, error)
	GetAccountByID(ctx context.Context, orgID, id uint) (*dto.AccountResponse, error)

	// Journal Entries - Basic
	CreateJournalEntry(ctx context.Context, orgID uint, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error)
	ListJournalEntries(ctx context.Context, orgID uint) ([]dto.JournalEntryResponse, error)
	GetJournalEntryByID(ctx context.Context, orgID, id uint) (*dto.JournalEntryResponse, error)

	// Journal Entries - Enhanced
	CreateJournalEntryIdempotent(ctx context.Context, orgID uint, idempotencyKey string, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error)
	ListJournalEntriesFiltered(ctx context.Context, orgID uint, filter dto.JournalEntryFilter) ([]dto.JournalEntryResponse, error)
	ReverseJournalEntry(ctx context.Context, orgID, journalEntryID uint, reason string) (*dto.JournalEntryResponse, error)
}

// Journal entry status constants
const (
	JournalStatusDrafted = "drafted"
	JournalStatusPosted  = "posted"
	JournalStatusVoided  = "voided"
)

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

func (s *accountingService) ListAccountsByType(ctx context.Context, orgID uint, accountType string) ([]dto.AccountResponse, error) {
	accounts, err := s.repo.ListAccountsByType(ctx, orgID, accountType)
	if err != nil {
		return nil, err
	}

	var res []dto.AccountResponse
	for _, a := range accounts {
		res = append(res, *s.mapAccountToResponse(&a))
	}
	return res, nil
}

func (s *accountingService) GetAccountByCode(ctx context.Context, orgID uint, code string) (*dto.AccountResponse, error) {
	account, err := s.repo.GetAccountByCode(ctx, orgID, code)
	if err != nil {
		return nil, errors.New("account not found")
	}
	return s.mapAccountToResponse(account), nil
}

func (s *accountingService) GetAccountByID(ctx context.Context, orgID, id uint) (*dto.AccountResponse, error) {
	account, err := s.repo.GetAccountByID(ctx, orgID, id)
	if err != nil {
		return nil, errors.New("account not found")
	}
	return s.mapAccountToResponse(account), nil
}

// CreateJournalEntry - Enhanced with account code resolution
func (s *accountingService) CreateJournalEntry(ctx context.Context, orgID uint, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error) {
	return s.createJournalEntry(ctx, orgID, "", req)
}

// CreateJournalEntryIdempotent - Creates journal entry with idempotency check
func (s *accountingService) CreateJournalEntryIdempotent(ctx context.Context, orgID uint, idempotencyKey string, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error) {
	// Validate idempotency key length if provided
	if idempotencyKey != "" && len(idempotencyKey) > 255 {
		return nil, errors.New("idempotency key exceeds maximum length of 255 characters")
	}

	// Validate source module is provided when auto-generating idempotency key
	if idempotencyKey == "" && req.SourceModule == "" {
		return nil, errors.New("source_module is required for journal entry creation")
	}

	// Generate idempotency key if not provided
	if idempotencyKey == "" {
		// Calculate total amount for uniqueness
		var totalAmount float64
		for _, line := range req.Lines {
			totalAmount += line.Debit
			break // Just get first non-zero amount
		}
		idempotencyKey = s.generateIdempotencyKey(orgID, req.SourceModule, req.SourceReference, totalAmount)
	}

	// Check if entry already exists
	existing, err := s.repo.GetJournalEntryByIdempotencyKey(ctx, orgID, idempotencyKey)
	if err == nil && existing != nil {
		// Return existing entry - idempotent behavior
		return s.mapJournalEntryToResponse(existing), nil
	}

	// Create new entry
	return s.createJournalEntry(ctx, orgID, idempotencyKey, req)
}

// createJournalEntry - Core implementation for creating journal entries
func (s *accountingService) createJournalEntry(ctx context.Context, orgID uint, idempotencyKey string, req dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error) {
	// Parse date first - fail early before expensive operations
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Validate double-entry balance
	var totalDebit, totalCredit float64
	var lines []model.JournalEntryLine
	// Note: accountCache is per-request and has no TTL.
	// For production with concurrent account modifications, consider using a shared cache (e.g., Redis).
	accountCache := make(map[string]*model.Account)

	for _, reqLine := range req.Lines {
		var account *model.Account
		var err error

		// Resolve account - either by code or by ID
		if reqLine.AccountCode != "" {
			// Check cache first
			if cached, ok := accountCache[reqLine.AccountCode]; ok {
				account = cached
			} else {
				account, err = s.repo.GetAccountByCode(ctx, orgID, reqLine.AccountCode)
				if err != nil {
					return nil, fmt.Errorf("invalid account code %s in journal entry lines", reqLine.AccountCode)
				}
				accountCache[reqLine.AccountCode] = account
			}
		} else if reqLine.AccountID != nil {
			account, err = s.repo.GetAccountByID(ctx, orgID, *reqLine.AccountID)
			if err != nil {
				return nil, errors.New("invalid account ID in journal entry lines")
			}
		} else {
			return nil, errors.New("either account_code or account_id must be provided")
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
			AccountID:   account.ID,
			AccountCode: account.Code,
			Description: reqLine.Description,
			Debit:       reqLine.Debit,
			Credit:      reqLine.Credit,
			PartnerID:   reqLine.PartnerID,
			PartnerType: reqLine.PartnerType,
		}
		line.OrganizationID = orgID
		lines = append(lines, line)
	}

	if totalDebit != totalCredit {
		return nil, errors.New("journal entry debit and credit totals must balance")
	}

	entry := &model.JournalEntry{
		ReferenceNumber: utils.GenerateCode("JE"),
		IdempotencyKey:  idempotencyKey,
		Date:            parsedDate,
		Description:     req.Description,
		SourceModule:    req.SourceModule,
		SourceReference: req.SourceReference,
		Status:          JournalStatusPosted,
		Lines:           lines,
	}
	entry.OrganizationID = orgID

	now := time.Now()
	entry.PostedAt = &now

	if err := s.repo.CreateJournalEntry(ctx, entry); err != nil {
		return nil, err
	}

	return s.mapJournalEntryToResponse(entry), nil
}

// ReverseJournalEntry - Creates an offsetting journal entry
func (s *accountingService) ReverseJournalEntry(ctx context.Context, orgID, journalEntryID uint, reason string) (*dto.JournalEntryResponse, error) {
	// Get the original entry
	original, err := s.repo.GetJournalEntryByID(ctx, orgID, journalEntryID)
	if err != nil {
		return nil, errors.New("journal entry not found")
	}

	if original.Status == JournalStatusVoided {
		return nil, errors.New("journal entry is already voided")
	}

	// Generate reversal idempotency key
	reversalIdempotencyKey := fmt.Sprintf("reversal.%d.%s", original.ID, original.IdempotencyKey)

	// Check if reversal already exists
	existing, err := s.repo.GetJournalEntryByIdempotencyKey(ctx, orgID, reversalIdempotencyKey)
	if err == nil && existing != nil {
		return s.mapJournalEntryToResponse(existing), nil
	}

	// Create reversal lines (swap debits and credits)
	var reversalLines []model.JournalEntryLine
	for _, line := range original.Lines {
		reversalLine := model.JournalEntryLine{
			AccountID:   line.AccountID,
			AccountCode: line.AccountCode,
			Description: fmt.Sprintf("REVERSAL - %s", line.Description),
			Debit:       line.Credit, // Swap debit/credit
			Credit:      line.Debit,
			PartnerID:   line.PartnerID,
			PartnerType: line.PartnerType,
		}
		reversalLine.OrganizationID = orgID
		reversalLines = append(reversalLines, reversalLine)
	}

	now := time.Now()
	reversalEntry := &model.JournalEntry{
		ReferenceNumber: utils.GenerateCode("JE"),
		IdempotencyKey:  reversalIdempotencyKey,
		Date:            now,
		Description:     fmt.Sprintf("REVERSAL: %s", original.Description),
		SourceModule:    "accounting",
		SourceReference: original.ReferenceNumber,
		Status:          JournalStatusPosted,
		ReversedEntryID: &original.ID,
		ReversalReason:  reason,
		Lines:           reversalLines,
	}
	reversalEntry.OrganizationID = orgID
	reversalEntry.PostedAt = &now

	if err := s.repo.CreateJournalEntry(ctx, reversalEntry); err != nil {
		return nil, err
	}

	// Mark original entry as voided
	original.Status = JournalStatusVoided
	if err := s.repo.UpdateJournalEntry(ctx, original); err != nil {
		// Log error but don't fail - reversal was created successfully
		log.Printf("Failed to mark original entry as voided: %v", err)
	}

	return s.mapJournalEntryToResponse(reversalEntry), nil
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

// ListJournalEntriesFiltered - List journal entries with filters
func (s *accountingService) ListJournalEntriesFiltered(ctx context.Context, orgID uint, filter dto.JournalEntryFilter) ([]dto.JournalEntryResponse, error) {
	entries, err := s.repo.ListJournalEntriesFiltered(ctx, orgID, filter)
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

// generateIdempotencyKey creates a unique key for idempotency
func (s *accountingService) generateIdempotencyKey(orgID uint, sourceModule, sourceReference string, totalAmount float64) string {
	// Include orgID, sourceModule, sourceReference, and amount hash to avoid collisions
	// Use underscore replacement to avoid issues with dotted source references
	amountHash := fmt.Sprintf("%.0f", totalAmount*100) // Convert to cents to avoid float issues
	safeRef := strings.ReplaceAll(sourceReference, ".", "_")
	return fmt.Sprintf("%d.%s.%s.%s", orgID, sourceModule, safeRef, amountHash)
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
	var totalDebit, totalCredit float64

	res := &dto.JournalEntryResponse{
		ID:              e.ID,
		ReferenceNumber: e.ReferenceNumber,
		IdempotencyKey:  e.IdempotencyKey,
		Date:            e.Date.Format("2006-01-02"),
		Description:     e.Description,
		Status:          e.Status,
		SourceModule:    e.SourceModule,
		SourceReference: e.SourceReference,
		ReversedEntryID: e.ReversedEntryID,
		ReversalReason:  e.ReversalReason,
		CreatedAt:       e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, l := range e.Lines {
		totalDebit += l.Debit
		totalCredit += l.Credit

		res.Lines = append(res.Lines, dto.JournalEntryLineResponse{
			ID:          l.ID,
			AccountID:   l.AccountID,
			AccountCode: l.AccountCode,
			Description: l.Description,
			Debit:       l.Debit,
			Credit:      l.Credit,
			PartnerID:   l.PartnerID,
			PartnerType: l.PartnerType,
		})
	}

	res.TotalDebit = totalDebit
	res.TotalCredit = totalCredit

	return res
}
