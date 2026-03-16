package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// -- Categories
func (h *InventoryHandler) CreateCategory(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.CategoryCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	category, err := h.service.CreateCategory(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, category, "Category created successfully")
}

func (h *InventoryHandler) ListCategories(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	categories, err := h.service.ListCategories(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, categories)
}

// -- Products
func (h *InventoryHandler) CreateProduct(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.ProductCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	product, err := h.service.CreateProduct(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, product, "Product created successfully")
}

func (h *InventoryHandler) ListProducts(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	params := pagination.Parse(c)

	products, total, err := h.service.ListProducts(c.Context(), orgID, params)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, products, response.Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: params.TotalPages(total),
	})
}

// -- Stock Adjustments
func (h *InventoryHandler) AdjustStock(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	productID, err := c.ParamsInt("productId")
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	var req dto.StockMovementRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	movement, err := h.service.AdjustStock(c.Context(), orgID, uint(productID), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, movement, "Stock updated successfully")
}

func (h *InventoryHandler) GetStockHistory(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	productID, err := c.ParamsInt("productId")
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	movements, err := h.service.GetStockHistory(c.Context(), orgID, uint(productID))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, movements)
}

// RegisterRoutes registers the inventory routes.
func RegisterRoutes(router fiber.Router, handler *InventoryHandler, middlewares ...fiber.Handler) {
	group := router.Group("/inventory", middlewares...)
	// Products and config

	// Categories
	group.Post("/categories", handler.CreateCategory)
	group.Get("/categories", handler.ListCategories)

	// Products
	group.Post("/products", handler.CreateProduct)
	group.Get("/products", handler.ListProducts)

	// Stock Interactions
	group.Post("/products/:productId/stock", handler.AdjustStock)
	group.Get("/products/:productId/stock-history", handler.GetStockHistory)
}
