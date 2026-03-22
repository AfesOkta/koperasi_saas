package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/member/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

type MobileHandler struct {
	service service.MobileService
}

func NewMobileHandler(service service.MobileService) *MobileHandler {
	return &MobileHandler{service: service}
}

func (h *MobileHandler) GetDashboard(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	
	res, err := h.service.GetMemberDashboard(c.Context(), claims.OrganizationID, claims.UserID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	
	return response.Success(c, res, "Mobile dashboard loaded")
}

func RegisterMobileRoutes(router fiber.Router, handler *MobileHandler, middlewares ...fiber.Handler) {
	group := router.Group("/mobile", middlewares...)
	group.Get("/me", handler.GetDashboard)
}
