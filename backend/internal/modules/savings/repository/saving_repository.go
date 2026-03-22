package repository

import (
	"context"
	"errors"
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/savings/model"
	"gorm.io/gorm"
)

type SavingRepository interface {
	// Products
	CreateProduct(ctx context.Context, product *model.SavingProduct) error
	GetProductByID(ctx context.Context, orgID, id uint) (*model.SavingProduct, error)
	GetProductByCode(ctx context.Context, orgID uint, code string) (*model.SavingProduct, error)
	ListProducts(ctx context.Context, orgID uint) ([]model.SavingProduct, error)

	// Accounts
	CreateAccount(ctx context.Context, account *model.SavingAccount) error
	GetAccountByMemberAndProduct(ctx context.Context, orgID, memberID, productID uint) (*model.SavingAccount, error)
	GetAccountByID(ctx context.Context, orgID, id uint) (*model.SavingAccount, error)

	// Transactions
	ExecuteTransaction(ctx context.Context, account *model.SavingAccount, transaction *model.SavingTransaction) error
	ListTransactionsByAccount(ctx context.Context, orgID, accountID uint) ([]model.SavingTransaction, error)

	// Utils
	ListAccounts(ctx context.Context, orgID uint) ([]model.SavingAccount, error)
	ListAccountsByMember(ctx context.Context, orgID, memberID uint) ([]model.SavingAccount, error)
	GetMinBalance(ctx context.Context, accountID uint, month int, year int) (float64, error)
	GetDB() *gorm.DB
}

type savingRepository struct {
	db *gorm.DB
}

func NewSavingRepository(db *gorm.DB) SavingRepository {
	return &savingRepository{db: db}
}

func (r *savingRepository) CreateProduct(ctx context.Context, product *model.SavingProduct) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *savingRepository) GetProductByID(ctx context.Context, orgID, id uint) (*model.SavingProduct, error) {
	var product model.SavingProduct
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&product, id).Error
	return &product, err
}

func (r *savingRepository) GetProductByCode(ctx context.Context, orgID uint, code string) (*model.SavingProduct, error) {
	var product model.SavingProduct
	err := r.db.WithContext(ctx).Where("organization_id = ? AND code = ?", orgID, code).First(&product).Error
	return &product, err
}

func (r *savingRepository) ListProducts(ctx context.Context, orgID uint) ([]model.SavingProduct, error) {
	var products []model.SavingProduct
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&products).Error
	return products, err
}

func (r *savingRepository) CreateAccount(ctx context.Context, account *model.SavingAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *savingRepository) GetAccountByMemberAndProduct(ctx context.Context, orgID, memberID, productID uint) (*model.SavingAccount, error) {
	var account model.SavingAccount
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND member_id = ? AND saving_product_id = ?", orgID, memberID, productID).
		First(&account).Error
	return &account, err
}

func (r *savingRepository) GetAccountByID(ctx context.Context, orgID, id uint) (*model.SavingAccount, error) {
	var account model.SavingAccount
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		First(&account, id).Error
	return &account, err
}

// ExecuteTransaction handles atomic deposits/withdrawals using database transactions
func (r *savingRepository) ExecuteTransaction(ctx context.Context, account *model.SavingAccount, txn *model.SavingTransaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the account row for update
		var lockedAccount model.SavingAccount
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&lockedAccount, account.ID).Error; err != nil {
			return err
		}

		// Recalculate based on transaction type
		if txn.Type == "deposit" {
			lockedAccount.Balance += txn.Amount
		} else if txn.Type == "withdrawal" {
			if lockedAccount.Balance < txn.Amount {
				return errors.New("insufficient balance")
			}
			lockedAccount.Balance -= txn.Amount
		} else {
			return errors.New("unsupported transaction type")
		}

		// Update balance after
		txn.BalanceAfter = lockedAccount.Balance

		// Save account balance
		if err := tx.Save(&lockedAccount).Error; err != nil {
			return err
		}

		// Save transaction record
		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// Update the referenced account object for the caller
		account.Balance = lockedAccount.Balance
		return nil
	})
}

func (r *savingRepository) ListTransactionsByAccount(ctx context.Context, orgID, accountID uint) ([]model.SavingTransaction, error) {
	var txns []model.SavingTransaction
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND saving_account_id = ?", orgID, accountID).
		Order("created_at DESC").
		Find(&txns).Error
	return txns, err
}

func (r *savingRepository) ListAccounts(ctx context.Context, orgID uint) ([]model.SavingAccount, error) {
	var accounts []model.SavingAccount
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&accounts).Error
	return accounts, err
}

func (r *savingRepository) ListAccountsByMember(ctx context.Context, orgID, memberID uint) ([]model.SavingAccount, error) {
	var accounts []model.SavingAccount
	err := r.db.WithContext(ctx).Where("organization_id = ? AND member_id = ?", orgID, memberID).Find(&accounts).Error
	return accounts, err
}

func (r *savingRepository) GetMinBalance(ctx context.Context, accountID uint, month int, year int) (float64, error) {
	// Formula: The minimum 'balance_after' during the month. 
	// If no transactions in the month, current balance is likely unchanged from previous month.
	// But we need to handle the case where the account was opened mid-month.
	
	var minVal struct {
		MinBalance float64
	}
	
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	err := r.db.WithContext(ctx).Table("saving_transactions").
		Select("MIN(balance_after) as min_balance").
		Where("saving_account_id = ? AND created_at >= ? AND created_at < ?", accountID, startDate, endDate).
		Scan(&minVal).Error
	
	if err != nil {
		return 0, err
	}

	// If no transactions, get the latest balance before the month started
	if minVal.MinBalance == 0 {
		var lastTxn model.SavingTransaction
		err = r.db.WithContext(ctx).
			Where("saving_account_id = ? AND created_at < ?", accountID, startDate).
			Order("created_at DESC").
			First(&lastTxn).Error
		if err == nil {
			return lastTxn.BalanceAfter, nil
		}
		// If still nothing (brand new account), get current balance
		var acc model.SavingAccount
		r.db.WithContext(ctx).First(&acc, accountID)
		return acc.Balance, nil
	}

	return minVal.MinBalance, nil
}

func (r *savingRepository) GetDB() *gorm.DB {
	return r.db
}
