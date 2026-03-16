package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/shu/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type SHUHandler struct {
	repo repository.SHURepository
}

func NewSHUHandler(repo repository.SHURepository) *SHUHandler {
	return &SHUHandler{repo: repo}
}

func (h *SHUHandler) ListConfigs(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	configs, err := h.repo.ListConfigs(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, configs)
}

// RegisterRoutes registers the SHU routes.
func RegisterRoutes(router fiber.Router, handler *SHUHandler, middlewares ...fiber.Handler) {
	group := router.Group("/shu", middlewares...)
	group.Get("/configs", handler.ListConfigs)
}
