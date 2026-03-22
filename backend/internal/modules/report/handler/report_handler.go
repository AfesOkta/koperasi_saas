package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/report/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type ReportHandler struct {
	service  service.ReportService
	exporter *service.Exporter
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{
		service:  svc,
		exporter: service.NewExporter(),
	}
}

func (h *ReportHandler) GetBalanceSheet(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	
	// Default to current month if not provided
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	
	dateFrom := firstOfMonth
	if dateFromStr != "" {
		if d, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = d
		}
	}
	
	dateTo := now
	if dateToStr != "" {
		if d, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
		}
	}

	res, err := h.service.GetBalanceSheet(c.Context(), orgID, dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *ReportHandler) GetProfitLoss(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)

	// Default to current month if not provided
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	dateFrom := firstOfMonth
	if dateFromStr != "" {
		if d, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = d
		}
	}

	dateTo := now
	if dateToStr != "" {
		if d, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
		}
	}

	res, err := h.service.GetProfitLoss(c.Context(), orgID, dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *ReportHandler) GetDashboard(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	res, err := h.service.GetDashboardKPIs(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *ReportHandler) ExportBalanceSheet(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	dateFrom := firstOfMonth
	if dateFromStr != "" {
		if d, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = d
		}
	}
	dateTo := now
	if dateToStr != "" {
		if d, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
		}
	}

	res, err := h.service.GetBalanceSheet(c.Context(), orgID, dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	format := c.Query("format", "excel")
	if format != "excel" {
		return response.BadRequest(c, "Unsupported format: "+format)
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", `attachment; filename="BalanceSheet.xlsx"`)

	return h.exporter.ExportBalanceSheetExcel(res, c.Response().BodyWriter())
}

func (h *ReportHandler) ExportProfitLoss(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	dateFrom := firstOfMonth
	if dateFromStr != "" {
		if d, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = d
		}
	}
	dateTo := now
	if dateToStr != "" {
		if d, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
		}
	}

	res, err := h.service.GetProfitLoss(c.Context(), orgID, dateFrom, dateTo)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	format := c.Query("format", "excel")
	if format != "excel" {
		return response.BadRequest(c, "Unsupported format: "+format)
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", `attachment; filename="ProfitLoss.xlsx"`)

	return h.exporter.ExportProfitLossExcel(res, c.Response().BodyWriter())
}

func RegisterRoutes(router fiber.Router, handler *ReportHandler, middlewares ...fiber.Handler) {
	group := router.Group("/reports", middlewares...)
	group.Get("/balance-sheet", handler.GetBalanceSheet)
	group.Get("/balance-sheet/export", handler.ExportBalanceSheet)
	group.Get("/profit-loss", handler.GetProfitLoss)
	group.Get("/profit-loss/export", handler.ExportProfitLoss)
	group.Get("/dashboard", handler.GetDashboard)
}
