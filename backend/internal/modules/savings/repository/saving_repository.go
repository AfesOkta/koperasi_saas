package repository

import (
	"context"
	"errors"

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
