package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/shu/dto"
	"github.com/koperasi-gresik/backend/internal/modules/shu/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type SHUHandler struct {
	service service.SHUService
}

func NewSHUHandler(service service.SHUService) *SHUHandler {
	return &SHUHandler{service: service}
}

func (h *SHUHandler) CreateConfig(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	var req dto.SHUConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	res, err := h.service.CreateConfig(c.Context(), orgID, req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, res)
}

func (h *SHUHandler) ListConfigs(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	configs, err := h.service.ListConfigs(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, configs)
}

func (h *SHUHandler) CalculateSHU(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	configIDStr := c.Params("id")
	configID, _ := strconv.ParseUint(configIDStr, 10, 32)
	
	if err := h.service.CalculateSHU(c.Context(), orgID, uint(configID)); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "SHU calculation completed successfully"})
}

func (h *SHUHandler) GetMemberDistributions(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	memberIDStr := c.Params("memberId")
	memberID, _ := strconv.ParseUint(memberIDStr, 10, 32)
	
	dists, err := h.service.GetMemberDistributions(c.Context(), orgID, uint(memberID))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, dists)
}

// RegisterRoutes registers the SHU routes.
func RegisterRoutes(router fiber.Router, handler *SHUHandler, middlewares ...fiber.Handler) {
	group := router.Group("/shu", middlewares...)
	group.Post("/configs", handler.CreateConfig)
	group.Get("/configs", handler.ListConfigs)
	group.Post("/configs/:id/calculate", handler.CalculateSHU)
	group.Get("/members/:memberId/distributions", handler.GetMemberDistributions)
}
