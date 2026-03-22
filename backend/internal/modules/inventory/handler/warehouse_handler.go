package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/dto"
	"github.com/koperasi-gresik/backend/internal/modules/inventory/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type WarehouseHandler struct {
	service service.WarehouseService
}

func NewWarehouseHandler(service service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) CreateWarehouse(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.WarehouseCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.CreateWarehouse(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, res, "Warehouse created successfully")
}

func (h *WarehouseHandler) ListWarehouses(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	res, err := h.service.ListWarehouses(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, res)
}

func (h *WarehouseHandler) InitiateTransfer(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.TransferCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.InitiateTransfer(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, res, "Transfer initiated successfully")
}

func (h *WarehouseHandler) ShipTransfer(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, _ := c.ParamsInt("id")

	if err := h.service.ShipTransfer(c.Context(), orgID, uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Transfer shipped successfully")
}

func (h *WarehouseHandler) ReceiveTransfer(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, _ := c.ParamsInt("id")

	if err := h.service.ReceiveTransfer(c.Context(), orgID, uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Transfer received successfully")
}

func RegisterWarehouseRoutes(router fiber.Router, handler *WarehouseHandler, middlewares ...fiber.Handler) {
	group := router.Group("/inventory", middlewares...)

	group.Post("/warehouses", handler.CreateWarehouse)
	group.Get("/warehouses", handler.ListWarehouses)
	
	group.Post("/transfers", handler.InitiateTransfer)
	group.Post("/transfers/:id/ship", handler.ShipTransfer)
	group.Post("/transfers/:id/receive", handler.ReceiveTransfer)
}
