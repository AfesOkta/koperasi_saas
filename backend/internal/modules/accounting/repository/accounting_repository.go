package repository

import (
	"context"
	"fmt"

	"github.com/koperasi-gresik/backend/internal/modules/accounting/dto"
	"github.com/koperasi-gresik/backend/internal/modules/accounting/model"
	"gorm.io/gorm"
)

type AccountingRepository interface {
	// Chart of Accounts
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccountByID(ctx context.Context, orgID, id uint) (*model.Account, error)
	GetAccountByCode(ctx context.Context, orgID uint, code string) (*model.Account, error)
	ListAccounts(ctx context.Context, orgID uint) ([]model.Account, error)
	ListAccountsByType(ctx context.Context, orgID uint, accountType string) ([]model.Account, error)

	// Journal Entries
	CreateJournalEntry(ctx context.Context, entry *model.JournalEntry) error
	GetJournalEntryByID(ctx context.Context, orgID, id uint) (*model.JournalEntry, error)
	GetJournalEntryByIdempotencyKey(ctx context.Context, orgID uint, key string) (*model.JournalEntry, error)
	ListJournalEntries(ctx context.Context, orgID uint) ([]model.JournalEntry, error)
	ListJournalEntriesFiltered(ctx context.Context, orgID uint, filter dto.JournalEntryFilter) ([]model.JournalEntry, error)
	UpdateJournalEntry(ctx context.Context, entry *model.JournalEntry) error
}

type accountingRepository struct {
	db *gorm.DB
}

func NewAccountingRepository(db *gorm.DB) AccountingRepository {
	return &accountingRepository{db: db}
}

func (r *accountingRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountingRepository) GetAccountByID(ctx context.Context, orgID, id uint) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&account, id).Error
	return &account, err
}

func (r *accountingRepository) GetAccountByCode(ctx context.Context, orgID uint, code string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Where("organization_id = ? AND code = ?", orgID, code).First(&account).Error
	return &account, err
}

func (r *accountingRepository) ListAccounts(ctx context.Context, orgID uint) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Order("code ASC").Find(&accounts).Error
	return accounts, err
}

func (r *accountingRepository) ListAccountsByType(ctx context.Context, orgID uint, accountType string) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.WithContext(ctx).Where("organization_id = ? AND type = ?", orgID, accountType).Order("code ASC").Find(&accounts).Error
	return accounts, err
}

// CreateJournalEntry inserts the entry and its lines within a single DB transaction.
func (r *accountingRepository) CreateJournalEntry(ctx context.Context, entry *model.JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(entry).Error
	})
}

func (r *accountingRepository) GetJournalEntryByID(ctx context.Context, orgID, id uint) (*model.JournalEntry, error) {
	var entry model.JournalEntry
	err := r.db.WithContext(ctx).Preload("Lines").Where("organization_id = ?", orgID).First(&entry, id).Error
	return &entry, err
}

// GetJournalEntryByIdempotencyKey retrieves a journal entry by its idempotency key
func (r *accountingRepository) GetJournalEntryByIdempotencyKey(ctx context.Context, orgID uint, key string) (*model.JournalEntry, error) {
	var entry model.JournalEntry
	err := r.db.WithContext(ctx).Preload("Lines").Where("organization_id = ? AND idempotency_key = ?", orgID, key).First(&entry).Error
	return &entry, err
}

func (r *accountingRepository) ListJournalEntries(ctx context.Context, orgID uint) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry
	err := r.db.WithContext(ctx).Preload("Lines").Where("organization_id = ?", orgID).Order("date DESC, created_at DESC").Find(&entries).Error
	return entries, err
}

// ListJournalEntriesFiltered retrieves journal entries with filters
func (r *accountingRepository) ListJournalEntriesFiltered(ctx context.Context, orgID uint, filter dto.JournalEntryFilter) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry

	query := r.db.WithContext(ctx).Preload("Lines").Where("organization_id = ?", orgID)

	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate)
	}
	if filter.SourceModule != "" {
		query = query.Where("source_module = ?", filter.SourceModule)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ReferenceNum != "" {
		query = query.Where("reference_number LIKE ?", fmt.Sprintf("%%%s%%", filter.ReferenceNum))
	}

	err := query.Order("date DESC, created_at DESC").Find(&entries).Error
	return entries, err
}

// UpdateJournalEntry updates an existing journal entry
func (r *accountingRepository) UpdateJournalEntry(ctx context.Context, entry *model.JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(entry).Error; err != nil {
			return err
		}
		// Update lines
		if len(entry.Lines) > 0 {
			// Delete existing lines
			if err := tx.Where("journal_entry_id = ?", entry.ID).Delete(&model.JournalEntryLine{}).Error; err != nil {
				return err
			}
			// Create new lines
			for i := range entry.Lines {
				entry.Lines[i].JournalEntryID = entry.ID
			}
			if err := tx.Create(&entry.Lines).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
