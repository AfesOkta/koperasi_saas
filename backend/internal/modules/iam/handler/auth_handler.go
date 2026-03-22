package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/iam/dto"
	"github.com/koperasi-gresik/backend/internal/modules/iam/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"github.com/koperasi-gresik/backend/internal/shared/validator"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.Login(c.Context(), req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, res, "Login successful")
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	res, err := h.service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, res, "Token refreshed")
}

func (h *AuthHandler) RegisterDeviceToken(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	var req dto.DeviceTokenRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	err := h.service.RegisterDeviceToken(c.Context(), claims.OrganizationID, claims.UserID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Device token registered")
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.UserCreateRequest
	if err := validator.Validate(c, &req); err != nil {
		return err
	}

	res, err := h.service.Register(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, res, "Registration successful")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Dummy response for MVP - extract from claims
	claims := c.Locals("claims").(*middleware.Claims)

	return response.Success(c, fiber.Map{
		"id":              claims.UserID,
		"organization_id": claims.OrganizationID,
		"email":           claims.Email,
		"role_id":         claims.RoleID,
	})
}

// RegisterPublicRoutes registers the public auth routes.
func RegisterPublicRoutes(router fiber.Router, handler *AuthHandler) {
	group := router.Group("/auth")
	group.Post("/login", handler.Login)
	group.Post("/register", handler.Register) // Public for MVP testing purpose
	group.Post("/refresh", handler.Refresh)
}

// RegisterRoutes registers the protected auth routes.
func RegisterRoutes(router fiber.Router, handler *AuthHandler, middlewares ...fiber.Handler) {
	group := router.Group("/auth", middlewares...)
	group.Get("/me", handler.Me)
	group.Post("/device-token", handler.RegisterDeviceToken)
}
