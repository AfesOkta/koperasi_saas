package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/notification/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type NotificationHandler struct {
	repo repository.NotificationRepository
}

func NewNotificationHandler(repo repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	userID := middleware.GetUserID(c)
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	ns, total, err := h.repo.ListByUser(c.Context(), orgID, userID, limit, offset)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, fiber.Map{"items": ns, "total": total})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.repo.MarkAsRead(c.Context(), uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, nil)
}

// RegisterRoutes registers the notification routes.
func RegisterRoutes(router fiber.Router, handler *NotificationHandler, middlewares ...fiber.Handler) {
	group := router.Group("/notifications", middlewares...)
	group.Get("/", handler.ListNotifications)
	group.Put("/:id/read", handler.MarkAsRead)
}
