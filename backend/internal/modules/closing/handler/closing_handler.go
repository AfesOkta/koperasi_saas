package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/closing/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type ClosingHandler struct {
	svc service.ClosingService
}

func NewClosingHandler(svc service.ClosingService) *ClosingHandler {
	return &ClosingHandler{svc: svc}
}

// ProcessEOD manually triggers the EOD process for the current tenant.
func (h *ClosingHandler) ProcessEOD(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	date := c.Query("date", time.Now().Format("2006-01-02"))

	if err := h.svc.ProcessEOD(c.Context(), orgID, date); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, fiber.Map{"message": "EOD processed successfully", "date": date})
}

// ProcessEOM manually triggers the EOM process for the current tenant.
func (h *ClosingHandler) ProcessEOM(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(time.Now().Month()))))
	year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))

	if err := h.svc.ProcessEOM(c.Context(), orgID, month, year); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, fiber.Map{"message": "EOM processed and period closed successfully", "month": month, "year": year})
}

// RegisterRoutes registers the closing routes.
func RegisterRoutes(router fiber.Router, handler *ClosingHandler, rbacMiddleware fiber.Handler) {
	group := router.Group("/closing", rbacMiddleware)
	group.Post("/eod", handler.ProcessEOD)
	group.Post("/eom", handler.ProcessEOM)
}
