package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/cash/dto"
	"github.com/koperasi-gresik/backend/internal/modules/cash/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type CashHandler struct {
	service service.CashService
}

func NewCashHandler(service service.CashService) *CashHandler {
	return &CashHandler{service: service}
}

func (h *CashHandler) CreateRegister(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	var req dto.CashRegisterCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	register, err := h.service.CreateRegister(c.Context(), orgID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, register, "Cash register created successfully")
}

func (h *CashHandler) ListRegisters(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	registers, err := h.service.ListRegisters(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, registers)
}

func (h *CashHandler) RecordTransaction(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	registerID, err := c.ParamsInt("registerId")
	if err != nil {
		return response.BadRequest(c, "Invalid register ID")
	}

	var req dto.CashTransactionRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	txn, err := h.service.RecordTransaction(c.Context(), orgID, uint(registerID), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, txn, "Cash transaction recorded successfully")
}

func (h *CashHandler) ListTransactions(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	registerID, err := c.ParamsInt("registerId")
	if err != nil {
		return response.BadRequest(c, "Invalid register ID")
	}

	params := pagination.Parse(c)
	txns, total, err := h.service.ListTransactionsByRegister(c.Context(), orgID, uint(registerID), params)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, txns, response.Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: params.TotalPages(total),
	})
}

// RegisterRoutes registers the cash routes.
func RegisterRoutes(router fiber.Router, handler *CashHandler, middlewares ...fiber.Handler) {
	group := router.Group("/cash", middlewares...)

	// Registers
	group.Post("/registers", handler.CreateRegister)
	group.Get("/registers", handler.ListRegisters)

	// Transactions
	group.Post("/registers/:registerId/transactions", handler.RecordTransaction)
	group.Get("/registers/:registerId/transactions", handler.ListTransactions)
}
