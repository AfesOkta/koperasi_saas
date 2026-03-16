package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/pos/model"
	"github.com/koperasi-gresik/backend/internal/modules/pos/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type POSHandler struct {
	repo repository.POSRepository
}

func NewPOSHandler(repo repository.POSRepository) *POSHandler {
	return &POSHandler{repo: repo}
}

func (h *POSHandler) OpenShift(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	cashierID := middleware.GetUserID(c)

	var req struct {
		StartBalance float64 `json:"start_balance"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	shift := &model.POSShift{
		CashierID:    cashierID,
		StartTime:    time.Now(),
		StartBalance: req.StartBalance,
		Status:       "open",
	}
	shift.OrganizationID = orgID

	if err := h.repo.OpenShift(c.Context(), shift); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, shift)
}

func (h *POSHandler) GetActiveShift(c *fiber.Ctx) error {
	cashierID := middleware.GetUserID(c)
	shift, err := h.repo.GetActiveShift(c.Context(), cashierID)
	if err != nil {
		return response.NotFound(c, "No active shift found")
	}
	return response.Success(c, shift)
}

// RegisterRoutes registers the POS routes.
func RegisterRoutes(router fiber.Router, handler *POSHandler, middlewares ...fiber.Handler) {
	group := router.Group("/pos", middlewares...)
	group.Post("/shifts/open", handler.OpenShift)
	group.Get("/shifts/active", handler.GetActiveShift)
}
