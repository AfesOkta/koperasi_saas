package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/supplier/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *model.Supplier) error
	GetByID(ctx context.Context, orgID, id uint) (*model.Supplier, error)
	GetByCode(ctx context.Context, orgID uint, code string) (*model.Supplier, error)
	Update(ctx context.Context, supplier *model.Supplier) error
	List(ctx context.Context, orgID uint, params pagination.Params) ([]model.Supplier, int64, error)
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *model.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

func (r *supplierRepository) GetByID(ctx context.Context, orgID, id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&supplier, id).Error
	return &supplier, err
}

func (r *supplierRepository) GetByCode(ctx context.Context, orgID uint, code string) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.WithContext(ctx).Where("organization_id = ? AND code = ?", orgID, code).First(&supplier).Error
	return &supplier, err
}

func (r *supplierRepository) Update(ctx context.Context, supplier *model.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

func (r *supplierRepository) List(ctx context.Context, orgID uint, params pagination.Params) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Supplier{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(params.Scope()).Order("name ASC").Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
