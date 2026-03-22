package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/inventory/model"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetWarehouseByID(ctx context.Context, orgID, id uint) (*model.Warehouse, error)
	GetWarehouseByCode(ctx context.Context, orgID uint, code string) (*model.Warehouse, error)
	ListWarehouses(ctx context.Context, orgID uint) ([]model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error

	// Stock Operations per Warehouse
	GetWarehouseItem(ctx context.Context, orgID, warehouseID, productID uint) (*model.WarehouseItem, error)
	GetLowStockItems(ctx context.Context, orgID uint) ([]model.WarehouseItem, error)
	UpsertWarehouseItem(ctx context.Context, item *model.WarehouseItem) error
	
	// Atomic Stock Adjustment with Movement
	AdjustStock(ctx context.Context, warehouseID, productID uint, movement *model.StockMovement) error
	
	// Transfers
	CreateTransfer(ctx context.Context, transfer *model.StockTransfer) error
	GetTransferByID(ctx context.Context, orgID, id uint) (*model.StockTransfer, error)
	UpdateTransfer(ctx context.Context, transfer *model.StockTransfer) error
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *warehouseRepository) GetWarehouseByID(ctx context.Context, orgID, id uint) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&warehouse, id).Error
	return &warehouse, err
}

func (r *warehouseRepository) GetWarehouseByCode(ctx context.Context, orgID uint, code string) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	err := r.db.WithContext(ctx).Where("organization_id = ? AND code = ?", orgID, code).First(&warehouse).Error
	return &warehouse, err
}

func (r *warehouseRepository) ListWarehouses(ctx context.Context, orgID uint) ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&warehouses).Error
	return warehouses, err
}

func (r *warehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *warehouseRepository) GetWarehouseItem(ctx context.Context, orgID, warehouseID, productID uint) (*model.WarehouseItem, error) {
	var item model.WarehouseItem
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND warehouse_id = ? AND product_id = ?", orgID, warehouseID, productID).
		First(&item).Error
	return &item, err
}

func (r *warehouseRepository) GetLowStockItems(ctx context.Context, orgID uint) ([]model.WarehouseItem, error) {
	var items []model.WarehouseItem
	err := r.db.WithContext(ctx).
		Preload("Warehouse").
		Where("organization_id = ? AND quantity < min_stock", orgID).
		Find(&items).Error
	return items, err
}

func (r *warehouseRepository) UpsertWarehouseItem(ctx context.Context, item *model.WarehouseItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *warehouseRepository) AdjustStock(ctx context.Context, warehouseID, productID uint, movement *model.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock WarehouseItem for safety
		var item model.WarehouseItem
		err := tx.Clauses(gorm.Expr("FOR UPDATE")).
			Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
			First(&item).Error
		
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// If not found, create new warehouse item (orgID from movement)
		if err == gorm.ErrRecordNotFound {
			item = model.WarehouseItem{
				WarehouseID:  warehouseID,
				ProductID:    productID,
				Quantity:     0,
				MinStock:     5,
				ReorderPoint: 10,
			}
			item.OrganizationID = movement.OrganizationID
		}

		// 2. Adjust Quantity
		if movement.Type == "out" || movement.Type == "transfer_out" {
			item.Quantity -= movement.Quantity
		} else {
			item.Quantity += movement.Quantity
		}

		movement.BalanceAfter = item.Quantity
		movement.WarehouseID = warehouseID

		// 3. Save Record
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *warehouseRepository) CreateTransfer(ctx context.Context, transfer *model.StockTransfer) error {
	return r.db.WithContext(ctx).Create(transfer).Error
}

func (r *warehouseRepository) GetTransferByID(ctx context.Context, orgID, id uint) (*model.StockTransfer, error) {
	var transfer model.StockTransfer
	err := r.db.WithContext(ctx).
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Where("organization_id = ?", orgID).
		First(&transfer, id).Error
	return &transfer, err
}

func (r *warehouseRepository) UpdateTransfer(ctx context.Context, transfer *model.StockTransfer) error {
	return r.db.WithContext(ctx).Save(transfer).Error
}
