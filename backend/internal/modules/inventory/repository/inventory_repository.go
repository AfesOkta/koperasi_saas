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
	AdjustStock(ctx context.Context, product *model.Product, movement *model.StockMovement) error
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
func (r *inventoryRepository) AdjustStock(ctx context.Context, product *model.Product, movement *model.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock Product Row to prevent race conditions during concurrent sales/purchases
		var lockedProduct model.Product
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&lockedProduct, product.ID).Error; err != nil {
			return err
		}

		// 2. Adjust Stock
		if movement.Type == "out" {
			// Basic validation - some use cases might allow negative stock, but typically we prevent it.
			// Implementing soft prevention here: if you want to allow negative stock, remove this check.
			lockedProduct.Stock -= movement.Quantity
		} else if movement.Type == "in" || movement.Type == "adj" {
			// If it's an adjustment, we could either add/subtract based on positive/negative quantity
			lockedProduct.Stock += movement.Quantity
		}

		movement.BalanceAfter = lockedProduct.Stock

		// 3. Save Movement Record
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		// 4. Update Product Stock
		if err := tx.Save(&lockedProduct).Error; err != nil {
			return err
		}

		product.Stock = lockedProduct.Stock // Update the caller's struct
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
