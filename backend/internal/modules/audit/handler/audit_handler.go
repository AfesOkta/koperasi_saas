package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/audit/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type AuditHandler struct {
	repo repository.AuditRepository
}

func NewAuditHandler(repo repository.AuditRepository) *AuditHandler {
	return &AuditHandler{repo: repo}
}

func (h *AuditHandler) ListLogs(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	logs, total, err := h.repo.ListLogs(c.Context(), orgID, limit, offset)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, fiber.Map{
		"items": logs,
		"total": total,
	})
}

// RegisterRoutes registers the audit routes.
func RegisterRoutes(router fiber.Router, handler *AuditHandler, middlewares ...fiber.Handler) {
	group := router.Group("/audit", middlewares...)
	group.Get("/logs", handler.ListLogs)
}
