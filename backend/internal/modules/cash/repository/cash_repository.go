package repository

import (
	"context"
	"errors"

	"github.com/koperasi-gresik/backend/internal/modules/cash/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type CashRepository interface {
	// Cash Registers
	CreateRegister(ctx context.Context, register *model.CashRegister) error
	GetRegisterByID(ctx context.Context, orgID, id uint) (*model.CashRegister, error)
	ListRegisters(ctx context.Context, orgID uint) ([]model.CashRegister, error)

	// Transactions
	ExecuteTransaction(ctx context.Context, register *model.CashRegister, txn *model.CashTransaction) error
	ListTransactionsByRegister(ctx context.Context, orgID, registerID uint, params pagination.Params) ([]model.CashTransaction, int64, error)
}

type cashRepository struct {
	db *gorm.DB
}

func NewCashRepository(db *gorm.DB) CashRepository {
	return &cashRepository{db: db}
}

func (r *cashRepository) CreateRegister(ctx context.Context, register *model.CashRegister) error {
	return r.db.WithContext(ctx).Create(register).Error
}

func (r *cashRepository) GetRegisterByID(ctx context.Context, orgID, id uint) (*model.CashRegister, error) {
	var register model.CashRegister
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&register, id).Error
	return &register, err
}

func (r *cashRepository) ListRegisters(ctx context.Context, orgID uint) ([]model.CashRegister, error) {
	var registers []model.CashRegister
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&registers).Error
	return registers, err
}

// ExecuteTransaction atomically locks the register, updates its balance, and inserts the transaction record.
func (r *cashRepository) ExecuteTransaction(ctx context.Context, register *model.CashRegister, txn *model.CashTransaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedReg model.CashRegister
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&lockedReg, register.ID).Error; err != nil {
			return err
		}

		if txn.Type == "in" {
			lockedReg.Balance += txn.Amount
		} else if txn.Type == "out" {
			if lockedReg.Balance < txn.Amount {
				return errors.New("insufficient balance in cash register")
			}
			lockedReg.Balance -= txn.Amount
		} else {
			return errors.New("unsupported cash transaction type")
		}

		txn.BalanceAfter = lockedReg.Balance

		if err := tx.Save(&lockedReg).Error; err != nil {
			return err
		}

		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		register.Balance = lockedReg.Balance
		return nil
	})
}

func (r *cashRepository) ListTransactionsByRegister(ctx context.Context, orgID, registerID uint, params pagination.Params) ([]model.CashTransaction, int64, error) {
	var txns []model.CashTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&model.CashTransaction{}).
		Where("organization_id = ? AND cash_register_id = ?", orgID, registerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(params.Scope()).Order("created_at DESC").Find(&txns).Error; err != nil {
		return nil, 0, err
	}

	return txns, total, nil
}
