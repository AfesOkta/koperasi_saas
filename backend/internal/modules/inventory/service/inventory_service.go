package service

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/model"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/repository"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type InventoryService interface {
	// Categories
	CreateCategory(ctx context.Context, orgID uint, req dto.CategoryCreateRequest) (*dto.CategoryResponse, error)
	ListCategories(ctx context.Context, orgID uint) ([]dto.CategoryResponse, error)

	// Products
	CreateProduct(ctx context.Context, orgID uint, req dto.ProductCreateRequest) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, orgID uint, params pagination.Params) ([]dto.ProductResponse, int64, error)
	GetProduct(ctx context.Context, orgID, productID uint) (*dto.ProductResponse, error)

	// Stock Operations
	AdjustStock(ctx context.Context, orgID, productID uint, req dto.StockMovementRequest) (*dto.StockMovementResponse, error)
	GetStockHistory(ctx context.Context, orgID, productID uint) ([]dto.StockMovementResponse, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

// Categories
func (s *inventoryService) CreateCategory(ctx context.Context, orgID uint, req dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	category.OrganizationID = orgID

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return s.mapCategoryToResponse(category), nil
}

func (s *inventoryService) ListCategories(ctx context.Context, orgID uint) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.ListCategories(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.CategoryResponse
	for _, c := range categories {
		res = append(res, *s.mapCategoryToResponse(&c))
	}
	return res, nil
}

// Products
func (s *inventoryService) CreateProduct(ctx context.Context, orgID uint, req dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	product := &model.Product{
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CostPrice:   req.CostPrice,
		MinStock:    req.MinStock,
		Unit:        req.Unit,
		Stock:       0, // Initialize explicitly to zero
	}
	product.OrganizationID = orgID

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.mapProductToResponse(product), nil
}

func (s *inventoryService) ListProducts(ctx context.Context, orgID uint, params pagination.Params) ([]dto.ProductResponse, int64, error) {
	products, total, err := s.repo.ListProducts(ctx, orgID, params)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.ProductResponse
	for _, p := range products {
		res = append(res, *s.mapProductToResponse(&p))
	}
	return res, total, nil
}

func (s *inventoryService) GetProduct(ctx context.Context, orgID, productID uint) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProductByID(ctx, orgID, productID)
	if err != nil {
		return nil, err
	}
	return s.mapProductToResponse(product), nil
}

// Stock Operations
func (s *inventoryService) AdjustStock(ctx context.Context, orgID, productID uint, req dto.StockMovementRequest) (*dto.StockMovementResponse, error) {
	product, err := s.repo.GetProductByID(ctx, orgID, productID)
	if err != nil {
		return nil, err
	}

	movement := &model.StockMovement{
		ProductID:       product.ID,
		ReferenceNumber: utils.GenerateCode("STK"),
		Type:            req.Type,
		Quantity:        req.Quantity,
		Notes:           req.Notes,
		RelatedEntity:   req.RelatedEntity,
		RelatedEntityID: req.RelatedEntityID,
	}
	movement.OrganizationID = orgID

	// NOTE: Quantity direction manipulation handled securely at atomic query level inside repo
	if err := s.repo.AdjustStock(ctx, product, movement); err != nil {
		return nil, err
	}

	return s.mapStockMovementToResponse(movement), nil
}

func (s *inventoryService) GetStockHistory(ctx context.Context, orgID, productID uint) ([]dto.StockMovementResponse, error) {
	// Verify product exists first to be safe
	if _, err := s.repo.GetProductByID(ctx, orgID, productID); err != nil {
		return nil, err
	}

	movements, err := s.repo.GetStockMovements(ctx, orgID, productID)
	if err != nil {
		return nil, err
	}

	var res []dto.StockMovementResponse
	for _, m := range movements {
		res = append(res, *s.mapStockMovementToResponse(&m))
	}
	return res, nil
}

// Mappers
func (s *inventoryService) mapCategoryToResponse(c *model.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}

func (s *inventoryService) mapProductToResponse(p *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CostPrice:   p.CostPrice,
		Stock:       p.Stock,
		MinStock:    p.MinStock,
		Unit:        p.Unit,
		Status:      p.Status,
	}
}

func (s *inventoryService) mapStockMovementToResponse(m *model.StockMovement) *dto.StockMovementResponse {
	return &dto.StockMovementResponse{
		ID:              m.ID,
		ProductID:       m.ProductID,
		ReferenceNumber: m.ReferenceNumber,
		Type:            m.Type,
		Quantity:        m.Quantity,
		BalanceAfter:    m.BalanceAfter,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
