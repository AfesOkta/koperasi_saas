package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/sales/dto"
	"github.com/koperasi-gresik/backend/internal/modules/sales/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type SalesHandler struct {
	service service.SalesService
}

func NewSalesHandler(service service.SalesService) *SalesHandler {
	return &SalesHandler{service: service}
}

func (h *SalesHandler) CreateOrder(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	cashierID := middleware.GetUserID(c)

	var req dto.OrderCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	order, err := h.service.CreateOrder(c.Context(), orgID, cashierID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, order, "Sales order created successfully")
}

func (h *SalesHandler) GetOrder(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid order ID")
	}

	order, err := h.service.GetOrder(c.Context(), orgID, uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, order)
}

func (h *SalesHandler) ListOrders(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	orders, err := h.service.ListOrders(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, orders)
}

// RegisterRoutes registers the sales routes.
func RegisterRoutes(router fiber.Router, handler *SalesHandler, middlewares ...fiber.Handler) {
	group := router.Group("/sales", middlewares...)

	group.Post("/orders", handler.CreateOrder)
	group.Get("/orders", handler.ListOrders)
	group.Get("/orders/:id", handler.GetOrder)
}
