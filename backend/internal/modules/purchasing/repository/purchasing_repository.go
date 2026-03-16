package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/purchasing/model"
	"gorm.io/gorm"
)

type PurchasingRepository interface {
	CreatePO(ctx context.Context, po *model.PurchaseOrder) error
	GetPOByID(ctx context.Context, orgID, id uint) (*model.PurchaseOrder, error)
	ListPOs(ctx context.Context, orgID uint) ([]model.PurchaseOrder, error)
	UpdatePOStatus(ctx context.Context, poID uint, status string) error
	AddPayment(ctx context.Context, payment *model.PurchasePayment) error
}

type purchasingRepository struct {
	db *gorm.DB
}

func NewPurchasingRepository(db *gorm.DB) PurchasingRepository {
	return &purchasingRepository{db: db}
}

func (r *purchasingRepository) CreatePO(ctx context.Context, po *model.PurchaseOrder) error {
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *purchasingRepository) GetPOByID(ctx context.Context, orgID, id uint) (*model.PurchaseOrder, error) {
	var po model.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Payments").
		Where("organization_id = ?", orgID).
		First(&po, id).Error
	return &po, err
}

func (r *purchasingRepository) ListPOs(ctx context.Context, orgID uint) ([]model.PurchaseOrder, error) {
	var pos []model.PurchaseOrder
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&pos).Error
	return pos, err
}

func (r *purchasingRepository) UpdatePOStatus(ctx context.Context, poID uint, status string) error {
	return r.db.Model(&model.PurchaseOrder{}).Where("id = ?", poID).Update("status", status).Error
}

func (r *purchasingRepository) AddPayment(ctx context.Context, payment *model.PurchasePayment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}
