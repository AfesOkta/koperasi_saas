package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/organization/dto"
	"github.com/koperasi-gresik/backend/internal/modules/organization/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type OrganizationHandler struct {
	service service.OrganizationService
}

func NewOrganizationHandler(service service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var req dto.OrganizationCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err // validator already formatted the fiber response
	}

	org, err := h.service.Create(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, org, "Organization created successfully")
}

func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid organization ID")
	}

	org, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, org)
}

func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	params := pagination.Parse(c)

	orgs, total, err := h.service.List(c.Context(), params)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, orgs, response.Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: params.TotalPages(total),
	})
}

func (h *OrganizationHandler) UpdateSettings(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid organization ID")
	}

	var req struct {
		Settings map[string]interface{} `json:"settings" validate:"required"`
	}
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	org, err := h.service.UpdateSettings(c.Context(), uint(id), req.Settings)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, org, "Organization settings updated successfully")
}

func (h *OrganizationHandler) Onboard(c *fiber.Ctx) error {
	var req dto.OnboardingRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	onboarding, err := h.service.Onboard(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, onboarding, "Public onboarding successful")
}

// RegisterPublicRoutes registers the public organization routes.
func RegisterPublicRoutes(router fiber.Router, handler *OrganizationHandler) {
	group := router.Group("/organizations")
	group.Post("/onboard", handler.Onboard)
}

// RegisterRoutes registers the organization routes.
func RegisterRoutes(router fiber.Router, handler *OrganizationHandler, middlewares ...fiber.Handler) {
	group := router.Group("/organizations", middlewares...)

	// Platform Admin Only (Org ID 1) using SuperAdminGuard
	group.Get("/", middleware.SuperAdminGuard(), handler.List)
	group.Post("/", middleware.SuperAdminGuard(), handler.Create)
	group.Patch("/:id/settings", middleware.SuperAdminGuard(), handler.UpdateSettings)

	// Admin of the organization can view their own details
	group.Get("/:id", handler.Get)
}
