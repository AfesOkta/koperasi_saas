package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/supplier/dto"
	"github.com/koperasi-gresik/backend/internal/modules/supplier/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type SupplierHandler struct {
	service service.SupplierService
}

func NewSupplierHandler(service service.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.SupplierCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	supplier, err := h.service.Create(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, supplier, "Supplier created successfully")
}

func (h *SupplierHandler) Get(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid supplier ID")
	}

	supplier, err := h.service.GetByID(c.Context(), orgID, uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, supplier)
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid supplier ID")
	}

	var req dto.SupplierUpdateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	supplier, err := h.service.Update(c.Context(), orgID, uint(id), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, supplier, "Supplier updated successfully")
}

func (h *SupplierHandler) List(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	params := pagination.Parse(c)

	suppliers, total, err := h.service.List(c.Context(), orgID, params)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, suppliers, response.Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: params.TotalPages(total),
	})
}

// RegisterRoutes registers the supplier routes.
func RegisterRoutes(router fiber.Router, handler *SupplierHandler, middlewares ...fiber.Handler) {
	group := router.Group("/suppliers", middlewares...)

	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.Get)
	group.Put("/:id", handler.Update)
}
