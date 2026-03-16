package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/sales/model"
	"gorm.io/gorm"
)

type SalesRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrderByID(ctx context.Context, orgID, id uint) (*model.Order, error)
	ListOrders(ctx context.Context, orgID uint) ([]model.Order, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderID uint, status string) error
}

type salesRepository struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save Order Header, Items, and Payments in one go
		return tx.Create(order).Error
	})
}

func (r *salesRepository) GetOrderByID(ctx context.Context, orgID, id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Payments").
		Where("organization_id = ?", orgID).
		First(&order, id).Error
	return &order, err
}

func (r *salesRepository) ListOrders(ctx context.Context, orgID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *salesRepository) UpdateOrderPaymentStatus(ctx context.Context, orderID uint, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Update("payment_status", status).Error
}
