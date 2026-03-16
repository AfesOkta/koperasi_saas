package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/billing/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type BillingHandler struct {
	repo repository.BillingRepository
}

func NewBillingHandler(repo repository.BillingRepository) *BillingHandler {
	return &BillingHandler{repo: repo}
}

func (h *BillingHandler) GetSubscription(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	sub, err := h.repo.GetSubscription(c.Context(), orgID)
	if err != nil {
		return response.NotFound(c, "No active subscription found")
	}
	return response.Success(c, sub)
}

// RegisterRoutes registers the billing routes.
func RegisterRoutes(router fiber.Router, handler *BillingHandler, middlewares ...fiber.Handler) {
	group := router.Group("/billing", middlewares...)
	group.Get("/subscription", handler.GetSubscription)
}
