package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/accounting/model"
	"gorm.io/gorm"
)

type AccountingRepository interface {
	// Chart of Accounts
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccountByID(ctx context.Context, orgID, id uint) (*model.Account, error)
	GetAccountByCode(ctx context.Context, orgID uint, code string) (*model.Account, error)
	ListAccounts(ctx context.Context, orgID uint) ([]model.Account, error)

	// Journal Entries
	CreateJournalEntry(ctx context.Context, entry *model.JournalEntry) error
	GetJournalEntryByID(ctx context.Context, orgID, id uint) (*model.JournalEntry, error)
	ListJournalEntries(ctx context.Context, orgID uint) ([]model.JournalEntry, error)
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

func (r *accountingRepository) ListJournalEntries(ctx context.Context, orgID uint) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry
	err := r.db.WithContext(ctx).Preload("Lines").Where("organization_id = ?", orgID).Order("date DESC").Find(&entries).Error
	return entries, err
}
