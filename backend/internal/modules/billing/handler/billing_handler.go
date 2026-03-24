package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/billing/dto"
	"github.com/koperasi-gresik/backend/internal/modules/billing/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type BillingHandler struct {
	service service.BillingService
}

func NewBillingHandler(service service.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) GetSubscription(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	sub, err := h.service.GetSubscription(c.Context(), orgID)
	if err != nil {
		return response.NotFound(c, "No active subscription found")
	}
	return response.Success(c, sub)
}

// Subscription Plan CRUD (Superadmin only)

func (h *BillingHandler) ListPlans(c *fiber.Ctx) error {
	plans, err := h.service.ListPlans(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, plans)
}

func (h *BillingHandler) GetPlan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	plan, err := h.service.GetPlanByID(c.Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Plan not found")
	}
	return response.Success(c, plan)
}

func (h *BillingHandler) CreatePlan(c *fiber.Ctx) error {
	var req dto.SubscriptionPlanRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	plan, err := h.service.CreatePlan(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, plan, "Subscription plan created successfully")
}

func (h *BillingHandler) UpdatePlan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.SubscriptionPlanRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	plan, err := h.service.UpdatePlan(c.Context(), uint(id), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, plan, "Subscription plan updated successfully")
}

func (h *BillingHandler) DeletePlan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeletePlan(c.Context(), uint(id)); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, nil, "Subscription plan deleted successfully")
}

// RegisterRoutes registers the billing routes.
func RegisterRoutes(router fiber.Router, handler *BillingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/billing", middlewares...)

	// Public/Organization level
	group.Get("/subscription", handler.GetSubscription)
	group.Get("/plans", handler.ListPlans) // Everyone can see plans

	// Superadmin level
	adminGroup := group.Group("/admin", middleware.SuperAdminGuard())
	adminGroup.Post("/plans", handler.CreatePlan)
	adminGroup.Get("/plans/:id", handler.GetPlan)
	adminGroup.Put("/plans/:id", handler.UpdatePlan)
	adminGroup.Delete("/plans/:id", handler.DeletePlan)
}
