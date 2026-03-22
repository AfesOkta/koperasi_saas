package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/inventory/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Categories
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, orgID, id uint) (*model.Category, error)
	ListCategories(ctx context.Context, orgID uint) ([]model.Category, error)

	// Products
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProductByID(ctx context.Context, orgID, id uint) (*model.Product, error)
	GetProductBySKU(ctx context.Context, orgID uint, sku string) (*model.Product, error)
	ListProducts(ctx context.Context, orgID uint, params pagination.Params) ([]model.Product, int64, error)
	UpdateProduct(ctx context.Context, product *model.Product) error

	// Stock Operations (Atomic)
	AdjustStock(ctx context.Context, warehouseID, productID uint, movement *model.StockMovement) error
	GetStockMovements(ctx context.Context, orgID, productID uint) ([]model.StockMovement, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// -- Categories
func (r *inventoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *inventoryRepository) GetCategoryByID(ctx context.Context, orgID, id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&category, id).Error
	return &category, err
}

func (r *inventoryRepository) ListCategories(ctx context.Context, orgID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&categories).Error
	return categories, err
}

// -- Products
func (r *inventoryRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *inventoryRepository) GetProductByID(ctx context.Context, orgID, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&product, id).Error
	return &product, err
}

func (r *inventoryRepository) GetProductBySKU(ctx context.Context, orgID uint, sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Where("organization_id = ? AND sku = ?", orgID, sku).First(&product).Error
	return &product, err
}

func (r *inventoryRepository) ListProducts(ctx context.Context, orgID uint, params pagination.Params) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(params.Scope()).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *inventoryRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// -- Stock Operations
func (r *inventoryRepository) AdjustStock(ctx context.Context, warehouseID, productID uint, movement *model.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock Product Row (Total Stock)
		var product model.Product
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&product, productID).Error; err != nil {
			return err
		}

		// 2. Lock/Find Warehouse Item
		var item model.WarehouseItem
		err := tx.Clauses(gorm.Expr("FOR UPDATE")).
			Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
			First(&item).Error
		
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			item = model.WarehouseItem{
				WarehouseID: warehouseID,
				ProductID:   productID,
				Quantity:    0,
			}
			item.OrganizationID = product.OrganizationID
		}

		// 3. Adjust Quantities
		qtyChange := movement.Quantity
		if movement.Type == "out" || movement.Type == "transfer_out" {
			qtyChange = -movement.Quantity
		}

		product.Stock += qtyChange
		item.Quantity += qtyChange

		movement.BalanceAfter = item.Quantity
		movement.WarehouseID = warehouseID

		// 4. Save Changes
		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *inventoryRepository) GetStockMovements(ctx context.Context, orgID, productID uint) ([]model.StockMovement, error) {
	var movements []model.StockMovement
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND product_id = ?", orgID, productID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}
