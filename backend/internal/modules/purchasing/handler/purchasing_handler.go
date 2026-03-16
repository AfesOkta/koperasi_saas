package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/purchasing/dto"
	"github.com/koperasi-gresik/backend/internal/modules/purchasing/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type PurchasingHandler struct {
	service service.PurchasingService
}

func NewPurchasingHandler(service service.PurchasingService) *PurchasingHandler {
	return &PurchasingHandler{service: service}
}

func (h *PurchasingHandler) CreatePO(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.PurchaseOrderCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	po, err := h.service.CreatePO(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, po, "Purchase order created successfully")
}

func (h *PurchasingHandler) GetPO(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid PO ID")
	}

	po, err := h.service.GetPO(c.Context(), orgID, uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, po)
}

func (h *PurchasingHandler) ReceivePO(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid PO ID")
	}

	if err := h.service.ReceivePO(c.Context(), orgID, uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "PO items received into inventory")
}

func (h *PurchasingHandler) ListPOs(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	pos, err := h.service.ListPOs(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, pos)
}

// RegisterRoutes registers the purchasing routes.
func RegisterRoutes(router fiber.Router, handler *PurchasingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/purchasing", middlewares...)

	group.Post("/orders", handler.CreatePO)
	group.Get("/orders", handler.ListPOs)
	group.Get("/orders/:id", handler.GetPO)
	group.Post("/orders/:id/receive", handler.ReceivePO)
}
