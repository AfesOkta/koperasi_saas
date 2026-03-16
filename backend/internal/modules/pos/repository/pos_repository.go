package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/pos/model"
	"gorm.io/gorm"
)

type POSRepository interface {
	OpenShift(ctx context.Context, shift *model.POSShift) error
	CloseShift(ctx context.Context, id uint, shift *model.POSShift) error
	GetActiveShift(ctx context.Context, cashierID uint) (*model.POSShift, error)
}

type posRepository struct {
	db *gorm.DB
}

func NewPOSRepository(db *gorm.DB) POSRepository {
	return &posRepository{db: db}
}

func (r *posRepository) OpenShift(ctx context.Context, shift *model.POSShift) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

func (r *posRepository) CloseShift(ctx context.Context, id uint, shift *model.POSShift) error {
	return r.db.WithContext(ctx).Model(&model.POSShift{}).Where("id = ?", id).Updates(shift).Error
}

func (r *posRepository) GetActiveShift(ctx context.Context, cashierID uint) (*model.POSShift, error) {
	var shift model.POSShift
	err := r.db.WithContext(ctx).Where("cashier_id = ? AND status = 'open'", cashierID).First(&shift).Error
	return &shift, err
}
