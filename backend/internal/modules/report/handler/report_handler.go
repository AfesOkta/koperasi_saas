package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/report/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetBalanceSheet(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	res, err := h.service.GetBalanceSheet(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *ReportHandler) GetProfitLoss(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	res, err := h.service.GetProfitLoss(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *ReportHandler) GetSummary(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	res, err := h.service.GetSummary(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func RegisterRoutes(router fiber.Router, handler *ReportHandler, middlewares ...fiber.Handler) {
	group := router.Group("/reports", middlewares...)
	group.Get("/balance-sheet", handler.GetBalanceSheet)
	group.Get("/profit-loss", handler.GetProfitLoss)
	group.Get("/summary", handler.GetSummary)
}
