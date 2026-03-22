package handler

import (
	"github.com/gofiber/fiber/v2"
	iamService "github.com/koperasi-gresik/backend/internal/modules/iam/service"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// RoleHandler handles RBAC management HTTP requests.
type RoleHandler struct {
	roleService iamService.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleService iamService.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// RegisterRoleRoutes registers all RBAC management routes.
func RegisterRoleRoutes(router fiber.Router, h *RoleHandler, authMid, tenantMid fiber.Handler, cache *middleware.PermissionCache) {
	r := router.Group("/iam", authMid, tenantMid)

	// Read-only
	r.Get("/permissions", h.ListPermissions)
	r.Get("/roles", middleware.RequirePermission("role:read", cache), h.ListRoles)
	r.Get("/roles/:id", middleware.RequirePermission("role:read", cache), h.GetRole)

	// Role management (requires role:manage)
	r.Post("/roles", middleware.RequirePermission("role:manage", cache), h.CreateRole)
	r.Put("/roles/:id", middleware.RequirePermission("role:manage", cache), h.UpdateRole)
	r.Delete("/roles/:id", middleware.RequirePermission("role:manage", cache), h.DeleteRole)

	// Permission assignment
	r.Post("/roles/:id/permissions", middleware.RequirePermission("role:manage", cache), h.AssignPermissions)
	r.Delete("/roles/:id/permissions", middleware.RequirePermission("role:manage", cache), h.RemovePermissions)

	// User-role assignment (requires user:update)
	r.Put("/users/:id/role", middleware.RequirePermission("user:update", cache), h.AssignRoleToUser)
}

// ListPermissions returns all system permissions.
func (h *RoleHandler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.roleService.ListPermissions(c.Context())
	if err != nil {
		return response.InternalError(c, "Failed to list permissions")
	}
	return response.Success(c, perms)
}

// ListRoles returns all roles for the organization.
func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roles, err := h.roleService.ListRoles(c.Context(), orgID)
	if err != nil {
		return response.InternalError(c, "Failed to list roles")
	}
	return response.Success(c, roles)
}

// GetRole returns a single role by ID.
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roleID, err := c.ParamsInt("id")
	if err != nil || roleID <= 0 {
		return response.BadRequest(c, "Invalid role ID")
	}
	role, err := h.roleService.GetRole(c.Context(), orgID, uint(roleID))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}
	return response.Success(c, role)
}

// CreateRole creates a new custom role.
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return response.BadRequest(c, "Name is required")
	}
	role, err := h.roleService.CreateRole(c.Context(), orgID, req.Name, req.Description)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, role)
}

// UpdateRole updates a custom role.
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roleID, err := c.ParamsInt("id")
	if err != nil || roleID <= 0 {
		return response.BadRequest(c, "Invalid role ID")
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	role, err := h.roleService.UpdateRole(c.Context(), orgID, uint(roleID), req.Name, req.Description)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, role)
}

// DeleteRole deletes a custom (non-system) role.
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roleID, err := c.ParamsInt("id")
	if err != nil || roleID <= 0 {
		return response.BadRequest(c, "Invalid role ID")
	}
	if err := h.roleService.DeleteRole(c.Context(), orgID, uint(roleID)); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, nil)
}

// AssignPermissions adds permissions to a role.
func (h *RoleHandler) AssignPermissions(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roleID, err := c.ParamsInt("id")
	if err != nil || roleID <= 0 {
		return response.BadRequest(c, "Invalid role ID")
	}
	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.BodyParser(&req); err != nil || len(req.Permissions) == 0 {
		return response.BadRequest(c, "Permissions list is required")
	}
	if err := h.roleService.AssignPermissions(c.Context(), orgID, uint(roleID), req.Permissions); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "Permissions assigned"})
}

// RemovePermissions removes permissions from a role.
func (h *RoleHandler) RemovePermissions(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	roleID, err := c.ParamsInt("id")
	if err != nil || roleID <= 0 {
		return response.BadRequest(c, "Invalid role ID")
	}
	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.BodyParser(&req); err != nil || len(req.Permissions) == 0 {
		return response.BadRequest(c, "Permissions list is required")
	}
	if err := h.roleService.RemovePermissions(c.Context(), orgID, uint(roleID), req.Permissions); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "Permissions removed"})
}

// AssignRoleToUser changes the role assigned to a user (single role per user).
func (h *RoleHandler) AssignRoleToUser(c *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(c)
	userID, err := c.ParamsInt("id")
	if err != nil || userID <= 0 {
		return response.BadRequest(c, "Invalid user ID")
	}
	var req struct {
		RoleID uint `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.RoleID == 0 {
		return response.BadRequest(c, "role_id is required")
	}
	if err := h.roleService.AssignRoleToUser(c.Context(), orgID, uint(userID), req.RoleID); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "Role assigned"})
}
