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

	// Orders
	CreateOrder(ctx context.Context, order *model.POSOrder) error
	GetOrderByID(ctx context.Context, orgID, orderID uint) (*model.POSOrder, error)
	ListOrdersByShift(ctx context.Context, orgID, shiftID uint) ([]model.POSOrder, error)

	// KDS
	UpdateOrderItemKDSStatus(ctx context.Context, orgID, itemID uint, status string) error
	GetKDSItems(ctx context.Context, orgID uint, statuses []string) ([]model.POSOrderItem, error)
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

// -- Orders
func (r *posRepository) CreateOrder(ctx context.Context, order *model.POSOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *posRepository) GetOrderByID(ctx context.Context, orgID, orderID uint) (*model.POSOrder, error) {
	var order model.POSOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("organization_id = ?", orgID).
		First(&order, orderID).Error
	return &order, err
}

func (r *posRepository) ListOrdersByShift(ctx context.Context, orgID, shiftID uint) ([]model.POSOrder, error) {
	var orders []model.POSOrder
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND shift_id = ?", orgID, shiftID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// -- KDS
func (r *posRepository) UpdateOrderItemKDSStatus(ctx context.Context, orgID, itemID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.POSOrderItem{}).
		Where("organization_id = ? AND id = ?", orgID, itemID).
		Update("kds_status", status).Error
}

func (r *posRepository) GetKDSItems(ctx context.Context, orgID uint, statuses []string) ([]model.POSOrderItem, error) {
	var items []model.POSOrderItem
	query := r.db.WithContext(ctx).
		Preload("Order").
		Where("organization_id = ?", orgID)

	if len(statuses) > 0 {
		query = query.Where("kds_status IN ?", statuses)
	}

	err := query.Order("created_at ASC").Find(&items).Error
	return items, err
}
