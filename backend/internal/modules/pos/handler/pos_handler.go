package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/pos/dto"
	"github.com/koperasi-gresik/backend/internal/modules/pos/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type POSHandler struct {
	service service.POSService
}

func NewPOSHandler(service service.POSService) *POSHandler {
	return &POSHandler{service: service}
}

func (h *POSHandler) OpenShift(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	cashierID := middleware.GetUserID(c)

	var req struct {
		StartBalance float64 `json:"start_balance" validate:"required,min=0"`
	}
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.OpenShift(c.Context(), orgID, cashierID, req.StartBalance)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, res, "Shift opened successfully")
}

func (h *POSHandler) GetActiveShift(c *fiber.Ctx) error {
	cashierID := middleware.GetUserID(c)

	res, err := h.service.GetActiveShift(c.Context(), cashierID)
	if err != nil {
		return response.NotFound(c, "No active shift found")
	}

	return response.Success(c, res)
}

func (h *POSHandler) CreateOrder(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.OrderCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.CreateOrder(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, res, "Order completed successfully")
}

func (h *POSHandler) GetKDSItems(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	status := c.Query("status")
	
	var statuses []string
	if status != "" {
		statuses = append(statuses, status)
	}

	res, err := h.service.GetKDSItems(c.Context(), orgID, statuses)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, res)
}

func (h *POSHandler) UpdateKDSStatus(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, _ := c.ParamsInt("id")
	
	var req struct {
		Status string `json:"status" validate:"required,oneof=pending preparing ready served"`
	}
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	if err := h.service.UpdateKDSStatus(c.Context(), orgID, uint(id), req.Status); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "KDS status updated successfully")
}

func (h *POSHandler) GenerateReceipt(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	id, _ := c.ParamsInt("id")

	res, err := h.service.GenerateReceipt(c.Context(), orgID, uint(id))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, res)
}

func RegisterRoutes(router fiber.Router, handler *POSHandler, middlewares ...fiber.Handler) {
	group := router.Group("/pos", middlewares...)

	group.Post("/shifts/open", handler.OpenShift)
	group.Get("/shifts/active", handler.GetActiveShift)
	
	group.Post("/orders", handler.CreateOrder)
	
	group.Get("/kds", handler.GetKDSItems)
	group.Patch("/kds/items/:id/status", handler.UpdateKDSStatus)

	group.Get("/orders/:id/receipt", handler.GenerateReceipt)
}
